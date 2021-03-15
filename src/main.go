package main

import (
	"os"
	"sort"
	"time"

	plog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/tsdb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	// BinanceAPIEndPoint ...
	BinanceAPIEndPoint = "https://api.binance.com"
)

// Global Vars
var (
	ListSymbols    = "BTCUSDT,ETHUSDT"
	RuleMonitoring string
	TSDatabase     *tsdb.DB
)

func main() {
	var err error

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	// open tsdb
	log.Info().Msg("Opening TSDB ./data")
	TSDatabase, err = tsdb.Open(
		"./data",
		plog.NewLogfmtLogger(os.Stdout),
		prometheus.NewRegistry(),
		tsdb.DefaultOptions,
	)
	if err != nil {
		log.Error().Msg("Can not open database")
		log.Panic().Msg(err.Error())
	}

	app := &cli.App{
		Name:  "Binance Client Bot",
		Usage: "@dungw3b",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "apikey",
				Usage:    "Binance API Key (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "secretkey",
				Usage:    "Binance Secret Key (required)",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "verbose",
				Usage:    "Set debug logging",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "ping",
				Usage:  "Test Binance API server connectivity",
				Action: BPing,
			},
			{
				Name:   "time",
				Usage:  "Check Binance API server time",
				Action: BTime,
			},
			{
				Name:   "price",
				Usage:  "Latest price for a symbol",
				Action: BTickerPrice,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "symbols",
						Usage:       "List of symbol",
						Required:    false,
						DefaultText: "BTCUSDT,ETHUSDT",
						Value:       "BTCUSDT,ETHUSDT",
						Destination: &ListSymbols,
					},
				},
			},
			{
				Name:   "daemon",
				Usage:  "Monitoring best coin for traders",
				Action: BDaemon,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "rule",
						Usage:       "Rule for monitoring (cnitgay, culi)",
						Required:    true,
						DefaultText: "culi",
						Value:       "culi",
						Destination: &RuleMonitoring,
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err = app.Run(os.Args)
	if err != nil {
		log.Info().Msg("")
		log.Error().Msg(err.Error())
	}
}
