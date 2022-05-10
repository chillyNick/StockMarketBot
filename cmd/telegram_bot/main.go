package main

import (
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot"
	"log"
)

func main() {
	err := telegram_bot.StartAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
