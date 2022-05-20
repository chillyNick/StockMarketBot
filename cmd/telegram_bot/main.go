package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/db"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/logger"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/repository/pgx_repository"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Error.Fatalf("Error loading .env file %s", err)
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
		logger.Error.Fatalf("Db connect failed: %s", err)
	}
	defer adp.Close()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT")),
		opts...,
	)
	if err != nil {
		logger.Error.Fatalf("Grpc connect failed: %s", err)
	}
	defer conn.Close()

	bot, err := telegram_bot.New(os.Getenv("TELEGRAM_APITOKEN"), pgx_repository.New(adp), conn, false)
	if err != nil {
		logger.Error.Fatal(err)
	}

	go telegram_bot.TrackNotification(bot, os.Getenv("RABBITMQ_URL"))

	bot.Serve()
}
