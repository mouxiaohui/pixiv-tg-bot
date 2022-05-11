package main

import (
	"log"

	"pixiv-tg-bot/core"
)

func main() {
	if err := core.Run(); err != nil {
		log.Fatal(err)
	}
}
