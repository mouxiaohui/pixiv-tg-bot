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
}

type Novel struct {
	Id            string         `json:"id"`
	Title         string         `json:"title"`
	UpdateDate    string         `json:"updateDate"`
	ContentTitles []ContentTitle `json:"contentTitles"`
}

type ContentTitle struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

// 机器人小说相关功能
func novelHandler(b *tele.Bot) {
	b.Handle(tele.OnText, func(c tele.Context) error {
		// 判断是否等待接收novel
		if IsReceiveNovel {
			IsReceiveNovel = false
			var text = c.Text()
			var reply string

			sd := subscribeNovels(strings.Split(text, ","))
			if len(sd.Success) > 0 {
				reply += "订阅成功:\n"
				for _, s := range sd.Success {
					reply += (">  " + s + "\n")
				}
			}
			if len(sd.Failure) > 0 {
				reply += "订阅失败:\n"
				for _, f := range sd.Failure {
					reply += (">  " + f + "\n")
				}
			}

			return c.Reply(reply)
		}

		return nil
	})

	b.Handle("/subnovel", func(c tele.Context) error {
		IsReceiveNovel = true
		return c.Reply("请发送小说ID, 如果有多个用逗号隔开, 例如: 1234,2234,3234")
	})

	b.Handle("/showsubnovel", func(c tele.Context) error {
		ns, err := queryAllNovel()
		if err != nil {
			return c.Reply("查询失败! Error: " + err.Error())
		}

		r := "查询结果:\n"
		for i, n := range ns {
			r += fmt.Sprintf("%d. %s\n", i, n.Title)
			r += fmt.Sprintf("    > 地址: %s%s\n", BASE_URL+"/novel/series/", n.Id)
		}

		return c.Reply(r)
	})
}

// 订阅小说列表
func subscribeNovels(ids []string) SubscribeDetails {
	sd := SubscribeDetails{}
	ch := make(chan ChanResult, len(ids))

	for _, id := range ids {
		go func(id string) {
			n := subscribeNovel(id, ch)
			if n.Id != "" {
				err := saveNovel(n)
				if err != nil {
					sd.Success = removeArrVal(sd.Success, id)
					sd.Failure = append(sd.Failure, id)
				}
			}
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

// 持久化小说
func saveNovel(n Novel) error {
	stmt, err := DB.Prepare("INSERT INTO novels(id, title, update_date) values(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(n.Id, n.Title, n.UpdateDate)
	if err != nil {
		return err
	}

	stmt2, err := DB.Prepare("INSERT INTO content_titles(id, title, novel_id) values(?,?,?)")
	if err != nil {
		return err
	}

	for _, ct := range n.ContentTitles {
		_, err = stmt2.Exec(ct.Id, ct.Title, n.Id)
		if err != nil {
			break
		}
	}

	return nil
}

// 查询所有小说
func queryAllNovel() ([]Novel, error) {
	var ns []Novel
	rows, err := DB.Query("select * from novels")
	if err != nil {
		return ns, err
	}

	for rows.Next() {
		var id string
		var title string
		var updateDate string

		err := rows.Scan(&id, &title, &updateDate)
		if err != nil {
			return ns, err
		}

		ns = append(ns, Novel{Id: id, Title: title, UpdateDate: updateDate})
	}

	return ns, nil
}

// 查询小说所有的内容
func queryNovelContent(id string) ([]ContentTitle, error) {
	var cts []ContentTitle
	r, err := DB.Query("select * from content_titles where novel_id = " + id)
	if err != nil {
		return cts, err
	}

	for r.Next() {
		var id string
		var title string
		var novel_id string
		err = r.Scan(&id, &title, &novel_id)
		if err != nil {
			return cts, err
		}
		cts = append(cts, ContentTitle{Id: id, Title: title})
	}

	return cts, nil
}
