package handlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
)

func handleShowCommand(s *telegram_bot.Server, msg *tgbotapi.Message, user *models.User) tgbotapi.MessageConfig {
	stocks, err := s.GrpcClient.GetStocks(context.Background(), &pb.UserId{Id: user.ServerUserId})
	if err != nil {
		logger.Error.Println(err)

		return tgbotapi.NewMessage(msg.Chat.ID, brokenMessage)
	}

	text := "Portfolio:\n"
	for _, stock := range stocks.GetStocks() {
		text += fmt.Sprintf("%v %v \n", stock.GetName(), stock.GetAmount())
	}

	return tgbotapi.NewMessage(msg.Chat.ID, text)
}
