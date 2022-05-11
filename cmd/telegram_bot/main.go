package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/repository"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
}

func main() {
	add := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_TG_BOT_DB"),
	)

	adp, err := db.New(context.Background(), add)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telegram_bot.New(os.Getenv("TELEGRAM_APITOKEN"), repository.New(adp), true)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
