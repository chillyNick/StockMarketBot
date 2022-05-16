package handlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

var removeStockResponse = "Send a message in the next format: ticker amount"

func handleRemoveStockCommand(s *telegram_bot.Server, msg *tgbotapi.Message, user *models.User) tgbotapi.MessageConfig {
	err := s.Repo.UpdateUserState(context.Background(), user.Id, models.UserStateRemoveStock)
	if err != nil {
		logger.Error.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
		return tgbotapi.NewMessage(msg.From.ID, brokenMessage)
	}

	return tgbotapi.NewMessage(msg.From.ID, removeStockResponse)
}

func handleRemoveStockText(s *telegram_bot.Server, msg *tgbotapi.Message, user *models.User) tgbotapi.MessageConfig {
	splitMsg := strings.Split(msg.Text, " ")
	if len(splitMsg) != 2 {
		return tgbotapi.NewMessage(msg.From.ID, addStockResponse)
	}

	amount, err := strconv.Atoi(splitMsg[1])
	if err != nil || amount <= 0 {
		return tgbotapi.NewMessage(msg.From.ID, "Amount must be a positive number")
	}

	_, err = s.GrpcClient.RemoveStock(context.Background(), &pb.StockRequest{
		Stock: &pb.Stock{
			Name:   splitMsg[0],
			Amount: int32(amount),
		},
		UserId: &pb.UserId{Id: user.ServerUserId},
	})

	if err != nil {
		text := brokenMessage
		if status.Code(err) == codes.NotFound {
			text = StockNotFoundMessage + " or you don't have in portfolio"
		}

		return tgbotapi.NewMessage(msg.Chat.ID, text)
	}

	err = s.Repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateMenu)
	if err != nil {
		logger.Error.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)

		return tgbotapi.NewMessage(msg.Chat.ID, brokenMessage)
	}

	return tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("%v %v was remove from portfolio", amount, splitMsg[0]))
}
