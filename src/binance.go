package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	//"github.com/go-playground/validator"
	plog "github.com/go-kit/kit/log"
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/tsdb"
	"github.com/prometheus/tsdb/labels"
	"github.com/roylee0704/gron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

/*var (
	validate *validator.Validate
)*/

type Metric struct {
	Series labels.Labels

	Timestamp int64
	Value     float64
}

type Query struct {
	Promql  string `json:"promql"`
	MinTime int64  `json:"mint"`
	MaxTime int64  `json:"maxt"`
}

type Point struct {
	T int64   `json:"t"`
	V float64 `json:"v"`
}

type Series struct {
	Labels labels.Labels `json:"labels"`
	Points []Point       `json:"points"`
}

// Constants
const (
	APIErrorTimeout = 1 * time.Minute
	APIPing         = "/api/v3/ping"
	APITime         = "/api/v3/time"
	APITickerPrice  = "/api/v3/ticker/price"
	APIOrderTest    = "/api/v3/order/test"
	APIOrder        = "/api/v3/order"
)

// Global Vars
var (
	TSDatabase *tsdb.DB
	client     = resty.New()
)

// SymbolPrice ..
type SymbolPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

// BPing /api/v3/ping
func BPing(ctx *cli.Context) error {
	_, err := callAPIGet(ctx, APIPing)

	if err != nil {
		return err
	}

	log.Info().Msg("ServerResponse: PONG")
	return nil
}

// BTime /api/v3/time
func BTime(ctx *cli.Context) error {
	result := struct {
		ServerTime int `json:"serverTime"`
	}{}

	response, err := callAPIGet(ctx, APITime)

	if err != nil {
		return err
	}

	json.Unmarshal([]byte(response), &result)

	log.Info().Msg("ServerTime: " + strconv.Itoa(result.ServerTime))
	return nil
}

// BTickerPrice /api/v3/ticker/price
func BTickerPrice(ctx *cli.Context) error {

	result := []struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}{}

	response, err := callAPIGet(ctx, APITickerPrice)

	if err != nil {
		return err
	}

	json.Unmarshal([]byte(response), &result)
	data := []SymbolPrice{}
	for _, item := range result {
		price, _ := strconv.ParseFloat(item.Price, 64)
		data = append(data, SymbolPrice{
			Symbol: item.Symbol,
			Price:  price,
		})
	}

	search := searchCoin(data)

	output, _ := json.Marshal(search)
	log.Info().Msg(string(output))

	return nil
}

func BStart(ctx *cli.Context) error {
	var err error

	// open tsdb
	log.Info().Msg("Loading TSDB ./data")
	TSDatabase, err = tsdb.Open(
		"./data",
		plog.NewLogfmtLogger(os.Stdout),
		prometheus.NewRegistry(),
		tsdb.DefaultOptions,
	)
	if err != nil {
		log.Error().Msg("Can not open database")
		log.Panic().Msg(err.Error())
		return err
	}

	c := gron.New()
	log.Info().Msg("Interval query API " + strconv.Itoa(IntervalQuery) + " seconds")
	c.AddFunc(gron.Every(time.Duration(IntervalQuery)*time.Second), func() {
		result := []struct {
			Symbol string `json:"symbol"`
			Price  string `json:"price"`
		}{}

		response, err := callAPIGet(ctx, APITickerPrice)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		app := TSDatabase.Appender()
		json.Unmarshal([]byte(response), &result)
		for _, item := range result {
			if strings.Contains(item.Symbol, "USDT") {
				price, _ := strconv.ParseFloat(item.Price, 64)
				data2 := Metric{
					Series: labels.Labels{
						labels.Label{
							Name:  "symbol",
							Value: item.Symbol,
						},
					},
					Timestamp: GetTimestampUTC(),
					Value:     price,
				}
				_, err := app.Add(data2.Series, data2.Timestamp, data2.Value)
				if err != nil {
					log.Error().Msg(err.Error())
					continue
				}
				log.Info().Msg("Inserted symbol " + item.Symbol + ":" + item.Price + "(timestamp " + strconv.Itoa(int(data2.Timestamp)) + ")")
			}
		}
		if err := app.Commit(); err != nil {
			log.Error().Msg(err.Error())
		}
	})

	c.Start()
	defer c.Stop()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	return nil
}

func QueryDB() {

	/*q := Query{
		Promql:  "BTCUSDT",
		MinTime: 1615964266544,
		MaxTime: 2615964266544,
	}
	matchedSer := make([]Series, 0)
	querier := TSDatabase.Querier(q.MinTime, q.MaxTime)

	matchers, err := PromQLToMatchers(q.Promql)

	defer querier.Close()*/
}

func BOrderTest(ctx *cli.Context) error {
	log.Info().Msg("---Start API Endpoint " + BinanceAPIEndPoint)

	if ctx.Bool("verbose") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Msg("GET " + APIOrderTest)

	// build parameters
	/*
		symbol=LTCBTC
		&side=BUY
		&type=LIMIT
		&timeInForce=GTC
		&quantity=1
		&price=0.1
		&recvWindow=5000
		&timestamp=1499827319559
	*/
	params := "symbol=" + OrderSymbol
	params = params + "&side=" + strings.ToUpper(OrderSide)
	params = params + "&type=" + strings.ToUpper(OrderType)
	params = params + "&timeInForce=GTC"
	params = params + "&quantity=" + strconv.FormatFloat(OrderQuantity, 'f', 8, 64)
	params = params + "&price=" + strconv.FormatFloat(OrderPrice, 'f', 8, 64)
	params = params + "&timestamp=" + strconv.FormatInt(GetTimestampUTC(), 10)

	signature := HMACSHA256(params, APISecret)
	params = params + "&signature=" + signature

	client.SetTimeout(APIErrorTimeout)

	resp, err := client.R().
		ForceContentType("application/json").
		SetHeader("X-MBX-APIKEY", APIKey).
		SetBody(params).
		Post(BinanceAPIEndPoint + APIOrderTest)

	if err != nil {
		log.Error().Msg(err.Error())
		log.Error().Msg(resp.String())
		return err
	}

	if resp.StatusCode() != 200 {
		log.Debug().Msg("Request Data         : " + params)
		log.Error().Msg("Status Code          : " + strconv.Itoa(resp.StatusCode()))
		log.Debug().Msg("Status               : " + resp.Status())
		log.Debug().Msg("Proto                : " + resp.Proto())
		log.Debug().Msg("Time                 : " + resp.Time().String())
		log.Debug().Msg("Received At          : " + resp.ReceivedAt().String())
		log.Debug().Msg("x-mbx-uuid           : " + resp.Header().Get("x-mbx-uuid"))
		log.Debug().Msg("x-mbx-used-weight    : " + resp.Header().Get("x-mbx-used-weight"))
		log.Debug().Msg("x-mbx-used-weight-1m : " + resp.Header().Get("x-mbx-used-weight-1m"))
		log.Debug().Msg("Body                 : " + resp.String())
		return err
	}

	if ctx.Bool("verbose") {
		log.Debug().Msg("Request Data         : " + params)
		log.Debug().Msg("Status Code          : " + strconv.Itoa(resp.StatusCode()))
		log.Debug().Msg("Status               : " + resp.Status())
		log.Debug().Msg("Proto                : " + resp.Proto())
		log.Debug().Msg("Time                 : " + resp.Time().String())
		log.Debug().Msg("Received At          : " + resp.ReceivedAt().String())
		log.Debug().Msg("x-mbx-uuid           : " + resp.Header().Get("x-mbx-uuid"))
		log.Debug().Msg("x-mbx-used-weight    : " + resp.Header().Get("x-mbx-used-weight"))
		log.Debug().Msg("x-mbx-used-weight-1m : " + resp.Header().Get("x-mbx-used-weight-1m"))
		log.Debug().Msg("Body                 : " + resp.String())
	}

	return nil
}

func BOrder(ctx *cli.Context) error {
	log.Info().Msg("---Start API Endpoint " + BinanceAPIEndPoint)

	if ctx.Bool("verbose") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Msg("GET " + APIOrder)

	// build parameters
	/*
		symbol=LTCBTC
		&side=BUY
		&type=LIMIT
		&timeInForce=GTC
		&quantity=1
		&price=0.1
		&recvWindow=5000
		&timestamp=1499827319559
	*/
	params := "symbol=" + OrderSymbol
	params = params + "&side=" + strings.ToUpper(OrderSide)
	params = params + "&type=" + strings.ToUpper(OrderType)
	params = params + "&timeInForce=GTC"
	params = params + "&quantity=" + strconv.FormatFloat(OrderQuantity, 'f', 8, 64)
	params = params + "&price=" + strconv.FormatFloat(OrderPrice, 'f', 8, 64)
	params = params + "&timestamp=" + strconv.FormatInt(GetTimestampUTC(), 10)

	signature := HMACSHA256(params, APISecret)
	params = params + "&signature=" + signature

	client.SetTimeout(APIErrorTimeout)

	resp, err := client.R().
		ForceContentType("application/json").
		SetHeader("X-MBX-APIKEY", APIKey).
		SetBody(params).
		Post(BinanceAPIEndPoint + APIOrder)

	if err != nil {
		log.Error().Msg(err.Error())
		log.Error().Msg(resp.String())
		return err
	}

	if resp.StatusCode() != 200 {
		log.Debug().Msg("Request Data         : " + params)
		log.Error().Msg("Status Code          : " + strconv.Itoa(resp.StatusCode()))
		log.Debug().Msg("Status               : " + resp.Status())
		log.Debug().Msg("Proto                : " + resp.Proto())
		log.Debug().Msg("Time                 : " + resp.Time().String())
		log.Debug().Msg("Received At          : " + resp.ReceivedAt().String())
		log.Debug().Msg("x-mbx-uuid           : " + resp.Header().Get("x-mbx-uuid"))
		log.Debug().Msg("x-mbx-used-weight    : " + resp.Header().Get("x-mbx-used-weight"))
		log.Debug().Msg("x-mbx-used-weight-1m : " + resp.Header().Get("x-mbx-used-weight-1m"))
		log.Debug().Msg("Body                 : " + resp.String())
		return err
	}

	if ctx.Bool("verbose") {
		log.Debug().Msg("Request Data         : " + params)
		log.Debug().Msg("Status Code          : " + strconv.Itoa(resp.StatusCode()))
		log.Debug().Msg("Status               : " + resp.Status())
		log.Debug().Msg("Proto                : " + resp.Proto())
		log.Debug().Msg("Time                 : " + resp.Time().String())
		log.Debug().Msg("Received At          : " + resp.ReceivedAt().String())
		log.Debug().Msg("x-mbx-uuid           : " + resp.Header().Get("x-mbx-uuid"))
		log.Debug().Msg("x-mbx-used-weight    : " + resp.Header().Get("x-mbx-used-weight"))
		log.Debug().Msg("x-mbx-used-weight-1m : " + resp.Header().Get("x-mbx-used-weight-1m"))
		log.Debug().Msg("Body                 : " + resp.String())
	}

	return nil
}
