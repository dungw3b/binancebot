package main

import (
	"os"
	"sort"
	"time"

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
	RuleMonitoring = "culi"
	IntervalQuery  = 60
	APIKey         string
	APISecret      string
)

func main() {
	var err error

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	app := &cli.App{
		Name:  "Binance Client Bot",
		Usage: "@dungw3b",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "verbose",
				Usage:    "Set debug logging",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "price",
				Usage:  "Latest price for symbols",
				Action: BTickerPrice,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "symbols",
						Usage:       "List of symbols",
						Required:    false,
						DefaultText: "BTCUSDT,ETHUSDT",
						Value:       "BTCUSDT,ETHUSDT",
						Destination: &ListSymbols,
					},
				},
			},
			{
				Name:   "start",
				Usage:  "Monitoring coin symbols for rule",
				Action: BStart,
				Flags: []cli.Flag{
					/*&cli.StringFlag{
						Name:        "apikey",
						Usage:       "Binance API Key (required)",
						Required:    true,
						Destination: &APIKey,
					},
					&cli.StringFlag{
						Name:        "secretkey",
						Usage:       "Binance Secret Key (required)",
						Required:    true,
						Destination: &APISecret,
					},*/
					&cli.StringFlag{
						Name:        "rule",
						Usage:       "Rule for monitoring",
						Required:    false,
						DefaultText: "culi",
						Value:       "culi",
						Destination: &RuleMonitoring,
					},
					&cli.IntFlag{
						Name:        "interval",
						Usage:       "Interval time for query API in seconds",
						Required:    false,
						DefaultText: "60",
						Value:       60,
						Destination: &IntervalQuery,
					},
				},
			},
		},
	}

	/*app := &cli.App{
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
	}*/

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err = app.Run(os.Args)
	if err != nil {
		log.Info().Msg("")
		log.Error().Msg(err.Error())
	}
}
