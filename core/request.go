package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var HEADER = map[string]string{
	"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36 Edg/101.0.1210.39",
}
var CLIENT = &http.Client{}

type Result[T any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    T      `json:"body"`
}

func unmarshalJSON[T any](data []byte) (T, error) {
	var s T
	err := json.Unmarshal(data, &s)

	return s, err
}

func request[T any](url string) (Result[T], error) {
	r := Result[T]{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return r, err
	}

	// 设置Header
	for k, v := range HEADER {
		req.Header.Add(k, v)
	}

	resp, err := CLIENT.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}
