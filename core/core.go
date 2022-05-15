package core

import (
	"net/http"
	"time"

	"pixiv-tg-bot/cmd"

	"golang.org/x/net/proxy"
	tele "gopkg.in/telebot.v3"
)

const BASE_URL = "https://www.pixiv.net"
const HELP_MESSAGE = `
🤖功能列表
/start              快速开始
/help               查看帮助信息
/subnovels          订阅小说
/showsubnovels      查看已经订阅的小说
/checknovelupdate   查看订阅的小说是否更新
/removesubnovels    移除订阅的小说
`

func Run() error {
	pref := tele.Settings{
		Token:  cmd.BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	if cmd.PROXY_ADDRESS != "" {
		// 设置代理
		dialer, err := proxy.SOCKS5("tcp", cmd.PROXY_ADDRESS, &proxy.Auth{}, proxy.Direct)
		if err != nil {
			return err
		}
		CLIENT.Transport = &http.Transport{Dial: dialer.Dial}
		pref.Client = CLIENT
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return err
	}

	addFeatures(b)

	b.Start()

	return nil
}

func addFeatures(b *tele.Bot) {
	b.Handle("/start", func(c tele.Context) error {
		return c.Send("欢迎使用 Pixiv Bot! 😘")
	})

	b.Handle("/help", func(c tele.Context) error {
		return c.Send(HELP_MESSAGE)
	})

	novelHandler(b)
}
