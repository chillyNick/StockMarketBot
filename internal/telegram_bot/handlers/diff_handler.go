package handlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"strings"
)

var (
	periods        = []string{"hour", "day", "week", "all"}
	periodKeyboard tgbotapi.ReplyKeyboardMarkup
)

func init() {
	buttons := make([]tgbotapi.KeyboardButton, len(periods))
	for i, p := range periods {
		buttons[i] = tgbotapi.NewKeyboardButton(p)
	}

	periodKeyboard = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(buttons...),
	)
}

func (h *Handler) handleDiffCommand(msg *tgbotapi.Message, user *models.User) tgbotapi.MessageConfig {
	err := h.repo.UpdateUserState(context.Background(), user.Id, models.UserStateDiff)
	if err != nil {
		logger.Error.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)

		return tgbotapi.NewMessage(msg.From.ID, brokenMessage)
	}

	msgCnf := tgbotapi.NewMessage(msg.From.ID, "Choose a period")
	msgCnf.ReplyMarkup = periodKeyboard

	return msgCnf
}

func (h *Handler) handleDiffText(msg *tgbotapi.Message, user *models.User) tgbotapi.MessageConfig {
	if !validatePeriod(msg.Text) {
		return tgbotapi.NewMessage(msg.Chat.ID, "Incorrect period")
	}

	changes, err := h.grpcClient.GetPortfolioChanges(
		context.Background(),
		&pb.GetPortfolioChangesRequest{
			Period: pb.Period(pb.Period_value[strings.ToUpper(msg.Text)]),
			UserId: &pb.UserId{Id: user.ServerUserId},
		},
	)
	if err != nil {
		logger.Error.Printf("Failed to get portfolio changes: %v", err)

		return tgbotapi.NewMessage(msg.Chat.ID, brokenMessage)
	}

	err = h.repo.UpdateUserState(context.Background(), user.Id, models.UserStateMenu)
	if err != nil {
		logger.Error.Printf("Failed to update user state with id:%v err: %v\n", user.Id, err)

		return tgbotapi.NewMessage(msg.Chat.ID, brokenMessage)
	}

	text := fmt.Sprintf("Portfolio changes by %v:\n", msg.Text)
	for _, ch := range changes.GetStocks() {
		text += fmt.Sprintf("%v %v %v \n", ch.Stock.GetName(), ch.Stock.GetAmount(), ch.GetCurrentPrice()-ch.GetOldPrice())
	}

	msgCnf := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgCnf.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	return msgCnf
}

func validatePeriod(period string) (valid bool) {
	for _, p := range periods {
		if period == p {
			return true
		}
	}

	return true
}
