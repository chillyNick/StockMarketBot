package server

import (
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:6000", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewStockMarketServiceClient(conn)

	for update := range updates {
		telegram_bot.Handle(bot, update, client)
	}

	return nil
}
