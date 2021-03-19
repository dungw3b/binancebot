package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func searchCoin(data []SymbolPrice) []SymbolPrice {
	result := []SymbolPrice{}

	if ListSymbols == "BTCUSDT,ETHUSDT" {
		log.Info().Msg("Query default list symbols \"" + ListSymbols + "\"")
	} else {
		log.Info().Msg("Query list symbols \"" + ListSymbols + "\"")
	}

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

func GetTimestampUTC() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func HMACSHA256(str string, key string) string {
	kbyte := []byte(key)
	sig := hmac.New(sha256.New, kbyte)
	sig.Write([]byte(str))

	return hex.EncodeToString(sig.Sum(nil))
}
