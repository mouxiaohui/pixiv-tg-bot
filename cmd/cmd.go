package cmd

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	BOT_TOKEN     string
	PROXY_ADDRESS string
)

func init() {
	app := &cli.App{
		Version: "1.0",
		Name:    "pixiv telegram bot",
		Action: func(c *cli.Context) error {
			println("Start...")
			return nil
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Aliases:     []string{"t"},
			Usage:       "机器人的 Token",
			Destination: &BOT_TOKEN,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "proxy",
			Aliases:     []string{"p"},
			Usage:       "代理地址, 比如(127.0.0.1:10808)",
			Destination: &PROXY_ADDRESS,
			Required:    false,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
