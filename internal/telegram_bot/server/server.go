package server

import (
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
}

func StartAndServe() error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		return err
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		telegram_bot.Handle(bot, update)
	}

	return nil
}
