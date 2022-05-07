package main

import (
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/server"
	"log"
)

func main() {
	err := server.StartAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
