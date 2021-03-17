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
					Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
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
