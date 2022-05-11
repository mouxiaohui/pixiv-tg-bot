package core

import (
	"io/ioutil"
	"net/http"
)

var HEADER = map[string]string{
	"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36 Edg/101.0.1210.39",
}
var CLIENT = &http.Client{}

func request(url string) (body string, err error) {
	body = ""
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	// 设置Header
	for k, v := range HEADER {
		req.Header.Add(k, v)
	}

	resp, err := CLIENT.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	body = string(b)

	return
}
