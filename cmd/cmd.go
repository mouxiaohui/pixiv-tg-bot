package cmd

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	BOT_TOKEN     string
	PROXY_ADDRESS string
	DB_PATH       string
)

func init() {
	app := &cli.App{
		Version: "1.0",
		Name:    "pixiv telegram bot",
		Action: func(c *cli.Context) error {
			initDBPath()
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
		&cli.StringFlag{
			Name:        "dbPath",
			Aliases:     []string{"d"},
			Usage:       "Sqlite3的数据库路径(默认为: './database/pixiv.db')",
			Destination: &DB_PATH,
			Required:    false,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// 初始化数据库路径
func initDBPath() {
	if DB_PATH == "" {
		DB_PATH = "./database/pixiv.db"
	}

	_, err := os.Stat(DB_PATH)
	if err != nil {
		panic(err)
	}
	if !os.IsExist(err) {
		f, err := os.Create(DB_PATH)
		defer f.Close()
		if err != nil {
			panic(err)
		}
	}
}
