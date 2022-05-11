package core

import (
	"errors"
	"strings"
	"sync"

	tele "gopkg.in/telebot.v3"
)

var IsReceiveNovel = false

type Novel struct {
	id            string
	title         string
	updateDate    string
	contentTitles []struct {
		id    string
		title string
	}
}

func novelHandler(b *tele.Bot) {
	b.Handle(tele.OnText, func(c tele.Context) error {
		// 判断是否等待接收novel
		if IsReceiveNovel {
			IsReceiveNovel = false
			var text = c.Text()
			err := subscribeNovels(strings.Split(text, ","))
			if err != nil {
				return errors.New("订阅失败! error: " + err.Error())
			}

			return c.Reply("订阅成功!")
		}

		return nil
	})

	b.Handle("/subnovel", func(c tele.Context) error {
		IsReceiveNovel = true
		return c.Reply("🤖: 请发送小说ID, 如果有多个用逗号隔开, 例如: 1234,2234,3234")
	})
}

func subscribeNovels(ids []string) error {
	wg := &sync.WaitGroup{}
	limiter := make(chan bool, 10) // 限制并发数

	for _, id := range ids {
		wg.Add(1)
		limiter <- true
		go subscribeNovel(id, limiter, wg)
	}

	wg.Wait()

	return nil
}

func subscribeNovel(id string, limiter chan bool, wg *sync.WaitGroup) error {
	defer wg.Done()

	var wg2 sync.WaitGroup
	var err error
	bodyc := ""

	wg2.Add(1)
	go func() {
		body, e := request(BASE_URL + "/ajax/novel/series/" + id + "/content_titles")
		if e != nil {
			err = e
		}
		bodyc = body
		wg2.Done()
	}()

	body, err := request(BASE_URL + "/ajax/novel/series/" + id)

	wg2.Wait()

	if err != nil {
		return err
	}

	println("================")
	println(bodyc)
	println(body)

	<-limiter

	return nil
}
