package main

import (
	"log"
	"net/http"
	"time"

	"pixiv-tg-bot/cmd"

	"golang.org/x/net/proxy"
	tele "gopkg.in/telebot.v3"
)

func run() error {
	// 设置代理
	dialer, err := proxy.SOCKS5("tcp", cmd.PROXY_ADDRESS, &proxy.Auth{}, proxy.Direct)
	if err != nil {
		return err
	}

	pref := tele.Settings{
		Token:  cmd.BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		Client: &http.Client{Transport: &http.Transport{Dial: dialer.Dial}},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return err
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("欢迎使用 Pixiv Bot! 😘")
	})

	b.Start()

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
