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

	OrderSymbol   string
	OrderSide     string
	OrderType     string
	OrderQuantity float64
	OrderPrice    float64
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
			{
				Name:   "ordertest",
				Usage:  "Test new order creation",
				Action: BOrderTest,
				Flags: []cli.Flag{
					&cli.StringFlag{
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
					},
					&cli.StringFlag{
						Name:        "symbol",
						Usage:       "Symbol for trading",
						Required:    true,
						Destination: &OrderSymbol,
					},
					&cli.StringFlag{
						Name:        "side",
						Usage:       "BUY, SELL",
						Required:    true,
						Destination: &OrderSide,
					},
					&cli.StringFlag{
						Name:        "type",
						Usage:       "LIMIT, MARKET, STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, TAKE_PROFIT_LIMIT, LIMIT_MAKER",
						Required:    true,
						Destination: &OrderType,
					},
					&cli.Float64Flag{
						Name:        "quantity",
						Usage:       "Quantity of symbol",
						Required:    true,
						Destination: &OrderQuantity,
					},
					&cli.Float64Flag{
						Name:        "price",
						Usage:       "Price of symbol",
						Required:    true,
						Destination: &OrderPrice,
					},
				},
			},
			{
				Name:   "order",
				Usage:  "Make a new order",
				Action: BOrder,
				Flags: []cli.Flag{
					&cli.StringFlag{
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
					},
					&cli.StringFlag{
						Name:        "symbol",
						Usage:       "Symbol for trading",
						Required:    true,
						Destination: &OrderSymbol,
					},
					&cli.StringFlag{
						Name:        "side",
						Usage:       "BUY, SELL",
						Required:    true,
						Destination: &OrderSide,
					},
					&cli.StringFlag{
						Name:        "type",
						Usage:       "LIMIT, MARKET, STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, TAKE_PROFIT_LIMIT, LIMIT_MAKER",
						Required:    true,
						Destination: &OrderType,
					},
					&cli.Float64Flag{
						Name:        "quantity",
						Usage:       "Quantity of symbol",
						Required:    true,
						Destination: &OrderQuantity,
					},
					&cli.Float64Flag{
						Name:        "price",
						Usage:       "Price of symbol",
						Required:    true,
						Destination: &OrderPrice,
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
