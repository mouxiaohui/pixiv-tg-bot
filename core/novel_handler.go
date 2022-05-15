package core

import (
	"fmt"
	"strings"
	"sync"

	tele "gopkg.in/telebot.v3"
)

var IsReceiveNovel = false
var IsRemoveNovels = false

type ChanResult struct {
	Id  string
	Err error
}

type SubscribeDetails struct {
	Success []string
	Failure []string
}

type Novel struct {
	Id         string   `json:"id"`
	Title      string   `json:"title"`
	UpdateDate string   `json:"updateDate"`
	Content    []string `json:"content"`
}

type ContentTitle struct {
	Id string `json:"id"`
}

// 机器人小说相关功能
func novelHandler(b *tele.Bot) {
	b.Handle(tele.OnText, func(c tele.Context) error {
		// 判断是否等待接收novels
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

		// 判断是否等待接收移除novels
		if IsRemoveNovels {
			IsRemoveNovels = false
			var text = c.Text()
			err := removeNovels(strings.Split(text, ","))
			if err != nil {
				return c.Reply("移除失败, Error: ", err.Error())
			}

			return c.Reply("移除成功!")
		}

		return nil
	})

	// 订阅小说
	b.Handle("/subnovels", func(c tele.Context) error {
		IsReceiveNovel = true
		return c.Reply("请发送小说ID, 如果有多个用逗号隔开, 例如: 1234,2234,3234")
	})

	// 查看已经订阅的小说
	b.Handle("/showsubnovels", func(c tele.Context) error {
		ns, err := queryAllNovel()
		if err != nil {
			return c.Reply("查询失败! Error: " + err.Error())
		}

		r := "查询结果:\n"
		for i, n := range ns {
			r += fmt.Sprintf("%d. %s\n", i+1, n.Title)
			r += fmt.Sprintf("    > 地址: %s%s\n", BASE_URL+"/novel/series/", n.Id)
		}

		return c.Reply(r)
	})

	// 查看订阅的小说是否更新
	b.Handle("/checknovelupdate", func(c tele.Context) error {
		ns, err := queryAllNovel()
		if err != nil {
			return c.Reply("查询失败! Error: " + err.Error())
		}

		wg := &sync.WaitGroup{}
		var updateNovels []Novel
		for _, n := range ns {
			wg.Add(1)
			go func(n Novel) {
				updateNovel, err := checkNovelUpdate(n)
				if err == nil && updateNovel.Id != "" {
					updateNovels = append(updateNovels, updateNovel)
				}
				wg.Done()
			}(n)
		}

		wg.Wait()

		if len(updateNovels) == 0 {
			return c.Reply("暂无更新...")
		}

		reply := "查询到更新:\n"
		for _, n := range updateNovels {
			reply += (n.Title + "\n")
			for _, c := range n.Content {
				reply += fmt.Sprintf("    > %s%s\n", BASE_URL+"/novel/show.php?id=", c)
			}
		}

		return c.Reply(reply)
	})

	b.Handle("/removesubnovels", func(c tele.Context) error {
		IsRemoveNovels = true
		return c.Reply("请发送要移除订阅的小说ID(可以通过/showsubnovels查询), 如果有多个用逗号隔开, 例如: 1234,2234")
	})
}

// 查询小说是否更新，返回小说更新的内容
func checkNovelUpdate(n Novel) (Novel, error) {
	var res Novel

	respNovel, err := request[Novel](BASE_URL + "/ajax/novel/series/" + n.Id)
	if err != nil {
		return res, err
	}
	if respNovel.Body.UpdateDate != n.UpdateDate && respNovel.Body.Id != "" {
		respContent, err := request[[]ContentTitle](BASE_URL + "/ajax/novel/series/" + n.Id + "/content_titles")
		if err != nil {
			return res, err
		}

		var content []string
		for _, ct := range respContent.Body {
			content = append(content, ct.Id)
		}

		respNovel.Body.Content = content[len(n.Content):]
		res = respNovel.Body
	}

	return res, nil
}

// 订阅小说列表
func subscribeNovels(ids []string) SubscribeDetails {
	sd := SubscribeDetails{}
	ch := make(chan ChanResult, len(ids))

	for _, id := range ids {
		go func(id string) {
			n, err := queryNovel(id)
			if err != nil {
				println("ERROR: " + err.Error())
				return
			}

			if n.Id == id {
				ch <- ChanResult{Id: id, Err: nil}
				return
			}

			println("=====")

			n = subscribeNovel(id, ch)
			if n.Id != "" {
				err := saveNovel(n)
				if err != nil {
					println("ERROR: " + err.Error())
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
		resp, err := request[[]ContentTitle](BASE_URL + "/ajax/novel/series/" + id + "/content_titles")
		if err == nil {
			ctTitles = resp.Body
		}
	}()

	// 获取小说信息
	respNovel, err := request[Novel](BASE_URL + "/ajax/novel/series/" + id)

	wg2.Wait()

	if err != nil {
		chRes.Err = err
		ch <- chRes
		return n
	}

	for _, ct := range ctTitles {
		respNovel.Body.Content = append(respNovel.Body.Content, ct.Id)
	}

	ch <- chRes
	return respNovel.Body
}

// 持久化小说
func saveNovel(n Novel) error {
	stmt, err := DB.Prepare("INSERT INTO novels(id, title, update_date, content) values(?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(n.Id, n.Title, n.UpdateDate, arrToString(n.Content, ","))
	if err != nil {
		return err
	}

	return nil
}

// 查询小说
func queryNovel(novelId string) (Novel, error) {
	var n Novel
	row, err := DB.Query("select * from novels where id=" + novelId)
	if err != nil {
		return n, err
	}

	if row.Next() {
		var id string
		var title string
		var updateDate string
		var content string
		err = row.Scan(&id, &title, &updateDate, &content)
		if err != nil {
			return n, err
		}

		return Novel{
			Id:         id,
			Title:      title,
			UpdateDate: updateDate,
			Content:    strings.Split(content, ","),
		}, nil
	}

	return n, nil
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
		var content string

		err := rows.Scan(&id, &title, &updateDate, &content)
		if err != nil {
			return ns, err
		}

		ns = append(ns, Novel{
			Id:         id,
			Title:      title,
			UpdateDate: updateDate,
			Content:    strings.Split(content, ","),
		})
	}

	return ns, nil
}

// 移除小说
func removeNovels(ids []string) error {
	stmt := `DELETE FROM novels where id IN (`
	for i, id := range ids {
		stmt += id
		if i != len(ids)-1 {
			stmt += ","
		}
	}
	stmt += ");"

	_, err := DB.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
