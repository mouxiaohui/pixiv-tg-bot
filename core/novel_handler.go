package core

import (
	"fmt"
	"strings"
	"sync"

	tele "gopkg.in/telebot.v3"
)

var IsReceiveNovel = false

type ChanResult struct {
	Id  string
	Err error
}

type SubscribeDetails struct {
	Success []string
	Failure []string
	Errors  []error
}

type Novel struct {
	Id            string         `json:"id"`
	Title         string         `json:"title"`
	UpdateDate    string         `json:"updateDate"`
	ContentTitles []ContentTitle `json:"contentTitles"`
}

type ContentTitle struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Available bool   `json:"available"`
}

func novelHandler(b *tele.Bot) {
	b.Handle(tele.OnText, func(c tele.Context) error {
		// 判断是否等待接收novel
		if IsReceiveNovel {
			IsReceiveNovel = false
			var text = c.Text()
			sd := subscribeNovels(strings.Split(text, ","))
			if len(sd.Errors) > 0 {
				for _, e := range sd.Errors {
					fmt.Println("error: ", e)
				}
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

// 订阅小说列表
func subscribeNovels(ids []string) SubscribeDetails {
	sd := SubscribeDetails{}
	ch := make(chan ChanResult, len(ids))

	for _, id := range ids {
		go func(id string) {
			n := subscribeNovel(id, ch)
			fmt.Println("Novel: ", n)
		}(id)
	}

	i := len(ids)
Loop:
	for {
		select {
		case res := <-ch:
			i--
			if res.Err != nil {
				sd.Failure = append(sd.Failure, res.Id)
				sd.Errors = append(sd.Errors, res.Err)
			} else {
				sd.Success = append(sd.Success, res.Id)
			}
			if i == 0 {
				break Loop
			}
		}
	}

	return sd
}

// 订阅小说
func subscribeNovel(id string, ch chan ChanResult) Novel {
	var err error
	wg2 := &sync.WaitGroup{}
	n := Novel{}
	ctTitles := []ContentTitle{}
	chRes := ChanResult{Id: id, Err: nil}

	wg2.Add(1)
	go func() {
		defer wg2.Done()
		// 获取小说章节
		res, err := request[[]ContentTitle](BASE_URL + "/ajax/novel/series/" + id + "/content_titles")
		if err == nil {
			ctTitles = res.Body
		}
	}()

	// 获取小说信息
	resNovel, err := request[Novel](BASE_URL + "/ajax/novel/series/" + id)

	wg2.Wait()

	if err != nil {
		chRes.Err = err
		ch <- chRes
		return n
	}

	resNovel.Body.ContentTitles = ctTitles

	ch <- chRes
	return resNovel.Body
}
