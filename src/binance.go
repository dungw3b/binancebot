package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	//"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

/*var (
	validate *validator.Validate
)*/

// Constants
const (
	APIErrorTimeout = 1 * time.Minute
	APIPing         = "/api/v3/ping"
	APITime         = "/api/v3/time"
	APITickerPrice  = "/api/v3/ticker/price"
)

// Global Vars
var (
	client = resty.New()
)

// SymbolPrice ..
type SymbolPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

func searchCoin(data []SymbolPrice) []SymbolPrice {
	result := []SymbolPrice{}

	for _, item := range data {
		if strings.Contains(ListSymbols, item.Symbol) {
			result = append(result, item)
		}
	}
	return result
}

func callAPIGet(ctx *cli.Context, api string) (string, error) {
	log.Info().Msg("---Start API Endpoint " + BinanceAPIEndPoint)

	if ctx.Bool("verbose") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Msg("GET " + api)

	client.SetTimeout(APIErrorTimeout)

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		ForceContentType("application/json").
		Get(BinanceAPIEndPoint + api)

	if err != nil {
		log.Error().Msg(err.Error())
		log.Error().Msg(resp.String())
		return "", err
	}

	if resp.StatusCode() != 200 {
		log.Error().Msg("Status Code: " + strconv.Itoa(resp.StatusCode()))
		log.Debug().Msg("Status               : " + resp.Status())
		log.Debug().Msg("Proto                : " + resp.Proto())
		log.Debug().Msg("Time                 : " + resp.Time().String())
		log.Debug().Msg("Received At          : " + resp.ReceivedAt().String())
		log.Debug().Msg("x-mbx-uuid           : " + resp.Header().Get("x-mbx-uuid"))
		log.Debug().Msg("x-mbx-used-weight    : " + resp.Header().Get("x-mbx-used-weight"))
		log.Debug().Msg("x-mbx-used-weight-1m : " + resp.Header().Get("x-mbx-used-weight-1m"))
		log.Debug().Msg("Body                 : " + resp.String())
		return "", err
	}

	if ctx.Bool("verbose") {
		log.Debug().Msg("Status Code: " + strconv.Itoa(resp.StatusCode()))
		log.Debug().Msg("Status               : " + resp.Status())
		log.Debug().Msg("Proto                : " + resp.Proto())
		log.Debug().Msg("Time                 : " + resp.Time().String())
		log.Debug().Msg("Received At          : " + resp.ReceivedAt().String())
		log.Debug().Msg("x-mbx-uuid           : " + resp.Header().Get("x-mbx-uuid"))
		log.Debug().Msg("x-mbx-used-weight    : " + resp.Header().Get("x-mbx-used-weight"))
		log.Debug().Msg("x-mbx-used-weight-1m : " + resp.Header().Get("x-mbx-used-weight-1m"))
		log.Debug().Msg("Body                 : " + resp.String())
	}

	return resp.String(), nil
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