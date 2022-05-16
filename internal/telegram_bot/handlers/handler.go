package handlers

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
)

const brokenMessage = "Something broken please try again later"
const StockNotFoundMessage = "Stock with such name not found"

func HandleUpdate(s *telegram_bot.Server, update tgbotapi.Update) *tgbotapi.MessageConfig {
	if update.Message == nil {
		return nil
	}

	var msgCnf tgbotapi.MessageConfig
	if update.Message.IsCommand() {
		msgCnf = handleCommand(s, update.Message)
	} else {
		msgCnf = handleText(s, update.Message)
	}

	return &msgCnf
}

func handleCommand(s *telegram_bot.Server, msg *tgbotapi.Message) tgbotapi.MessageConfig {
	if msg.Command() == "start" {
		return handleStartCommand(s, msg)
	}

	u, ok, text := getUser(s.Repo, msg.From.ID)
	if !ok {
		return tgbotapi.NewMessage(msg.Chat.ID, text)
	}

	switch msg.Command() {
	case "show":
		return handleShowCommand(s, msg, u)
	case "diff":
		return handleDiffCommand(s, msg, u)
	case "add_stock":
		return handleAddStockCommand(s, msg, u)
	case "remove_stock":
		return handleRemoveStockCommand(s, msg, u)
	case "add_notification":
		return handleAddNotificationCommand(s, msg, u)
	}

	return tgbotapi.NewMessage(msg.Chat.ID, "Unknown command")
}

func handleText(s *telegram_bot.Server, msg *tgbotapi.Message) tgbotapi.MessageConfig {
	u, ok, text := getUser(s.Repo, msg.From.ID)
	if !ok {
		return tgbotapi.NewMessage(msg.Chat.ID, text)
	}

	switch u.State {
	case models.UserStateDiff:
		return handleDiffText(s, msg, u)
	case models.UserStateAddStock:
		return handleAddStockText(s, msg, u)
	case models.UserStateRemoveStock:
		return handleRemoveStockText(s, msg, u)
	case models.UserStateAddNotification:
		return handleAddNotificationText(s, msg, u)
	}

	return tgbotapi.NewMessage(msg.Chat.ID, "Type command")
}

func getUser(repo telegram_bot.Repository, id int64) (u *models.User, ok bool, msg string) {
	u, err := repo.GetUser(context.Background(), id)
	if err == nil {
		return u, true, ""
	}

	var text string
	if errors.Is(err, db.ErrNotFound) {
		text = "Firstly type /start"
	} else {
		logger.Error.Printf("Failed to get user by id:%v err: %v\n", id, err)
		text = brokenMessage
	}

	return nil, false, text
}
