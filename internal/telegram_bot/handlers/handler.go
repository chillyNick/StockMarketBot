package handlers

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/repository"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"golang.org/x/net/context"
)

type Handler struct {
	repo       repository.Repository
	grpcClient pb.StockMarketServiceClient
}

const brokenMessage = "Something broken please try again later"
const StockNotFoundMessage = "Stock with such name not found"

func New(repo repository.Repository, grpcClient pb.StockMarketServiceClient) *Handler {
	return &Handler{
		repo:       repo,
		grpcClient: grpcClient,
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) *tgbotapi.MessageConfig {
	if update.Message == nil {
		return nil
	}

	var msgCnf tgbotapi.MessageConfig
	if update.Message.IsCommand() {
		msgCnf = h.handleCommand(update.Message)
	} else {
		msgCnf = h.handleText(update.Message)
	}

	return &msgCnf
}

func (h *Handler) handleCommand(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	if msg.Command() == "start" {
		return h.handleStartCommand(msg)
	}

	u, ok, text := getUser(h.repo, msg.From.ID)
	if !ok {
		return tgbotapi.NewMessage(msg.Chat.ID, text)
	}

	switch msg.Command() {
	case "show":
		return h.handleShowCommand(msg, u)
	case "diff":
		return h.handleDiffCommand(msg, u)
	case "add_stock":
		return h.handleAddStockCommand(msg, u)
	case "remove_stock":
		return h.handleRemoveStockCommand(msg, u)
	case "add_notification":
		return h.handleAddNotificationCommand(msg, u)
	}

	return tgbotapi.NewMessage(msg.Chat.ID, "Unknown command")
}

func (h *Handler) handleText(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	u, ok, text := getUser(h.repo, msg.From.ID)
	if !ok {
		return tgbotapi.NewMessage(msg.Chat.ID, text)
	}

	switch u.State {
	case models.UserStateDiff:
		return h.handleDiffText(msg, u)
	case models.UserStateAddStock:
		return h.handleAddStockText(msg, u)
	case models.UserStateRemoveStock:
		return h.handleRemoveStockText(msg, u)
	case models.UserStateAddNotification:
		return h.handleAddNotificationText(msg, u)
	}

	return tgbotapi.NewMessage(msg.Chat.ID, "Type command")
}

func getUser(repo repository.Repository, id int64) (u *models.User, ok bool, msg string) {
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
