package handlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

var addNotificationResponse = "Send a message in the next format: ticker amount"

func (h *Handler) handleAddNotificationCommand(msg *tgbotapi.Message, user *models.User) tgbotapi.MessageConfig {
	err := h.repo.UpdateUserState(context.Background(), user.Id, models.UserStateAddNotification)
	if err != nil {
		logger.Error.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
		return tgbotapi.NewMessage(msg.From.ID, brokenMessage)
	}

	return tgbotapi.NewMessage(msg.From.ID, addNotificationResponse)
}

func (h *Handler) handleAddNotificationText(msg *tgbotapi.Message, user *models.User) tgbotapi.MessageConfig {
	splitMsg := strings.Split(msg.Text, " ")
	if len(splitMsg) != 2 {
		return tgbotapi.NewMessage(msg.From.ID, addNotificationResponse)
	}

	threshold, err := strconv.ParseFloat(splitMsg[1], 64)
	if err != nil {
		return tgbotapi.NewMessage(msg.From.ID, "Threshold must be a positive decimal number")
	}

	_, err = h.grpcClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
		UserId:    &pb.UserId{Id: user.ServerUserId},
		StockName: splitMsg[0],
		Threshold: threshold,
	})

	if err != nil {
		text := brokenMessage
		if status.Code(err) == codes.NotFound {
			text = StockNotFoundMessage
		}

		return tgbotapi.NewMessage(msg.Chat.ID, text)
	}

	err = h.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateMenu)
	if err != nil {
		logger.Error.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)

		return tgbotapi.NewMessage(msg.Chat.ID, brokenMessage)
	}

	return tgbotapi.NewMessage(
		msg.Chat.ID,
		fmt.Sprintf("Notification for ticker %v by threshold %v was added", splitMsg[0], threshold),
	)
}
