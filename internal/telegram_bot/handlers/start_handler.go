package handlers

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

var greetingMessage = "Welcome to the stock market bot."

func (h *Handler) handleStartCommand(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	_, err := h.repo.GetUser(context.Background(), msg.From.ID)
	if err == nil {
		return tgbotapi.NewMessage(msg.Chat.ID, greetingMessage)
	}

	if !errors.Is(err, db.ErrNotFound) {
		logger.Error.Printf("Failed to get user by id:%v err: %v\n", msg.From.ID, err)

		return tgbotapi.NewMessage(msg.Chat.ID, brokenMessage)
	}

	id, err := h.grpcClient.CreateUser(context.Background(), &emptypb.Empty{})
	if err != nil {
		logger.Error.Printf("Failed to create user at server side %v\n", err)

		return tgbotapi.NewMessage(msg.Chat.ID, brokenMessage)

	}

	err = h.repo.CreateUser(context.Background(), msg.From.ID, msg.Chat.ID, id.Id)
	if err != nil {
		logger.Error.Printf("Failed to create user %v\n", err)

		return tgbotapi.NewMessage(msg.Chat.ID, brokenMessage)
	}

	return tgbotapi.NewMessage(msg.Chat.ID, greetingMessage)
}
