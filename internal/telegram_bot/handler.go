package telegram_bot

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"log"
	"strconv"
	"strings"
)

const brokenMessage = "Something broken please try again later"

func (s *server) handle(update tgbotapi.Update, client pb.StockMarketServiceClient) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		s.handleCommand(update.Message)
	} else {
		s.handleText(update.Message)
	}
}

func (s *server) handleCommand(msg *tgbotapi.Message) {
	var text string
	switch msg.Command() {
	case "start":
		text = s.handleStartCommand(msg)
	case "help":
		text = "/add_stock - Add to the stock portfolio\n/remove_stock - Remove from the stock portfolio"
	}
	if text != "" {
		s.send(tgbotapi.NewMessage(msg.Chat.ID, text))

		return
	}

	_, ok := s.getUser(msg.From.ID, msg.Chat.ID)
	if !ok {
		return
	}

	switch msg.Command() {
	case "show":
		s.send(tgbotapi.NewMessage(msg.Chat.ID, "Show current stocks"))
	case "diff":
		var periodKeyboard = tgbotapi.NewOneTimeReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("hour"),
				tgbotapi.NewKeyboardButton("day"),
				tgbotapi.NewKeyboardButton("week"),
				tgbotapi.NewKeyboardButton("all"),
			),
		)
		var nMsg tgbotapi.MessageConfig
		err := s.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateDiff)
		if err != nil {
			log.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
			nMsg = tgbotapi.NewMessage(msg.From.ID, brokenMessage)
		} else {
			nMsg = tgbotapi.NewMessage(msg.From.ID, "Choose a period")
			nMsg.ReplyMarkup = periodKeyboard
		}
		s.send(nMsg)
	case "add_stock":
		err := s.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateAddStock)
		if err != nil {
			log.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
			text = brokenMessage
		} else {
			text = "Send a message in the next format: stockName amount"
		}

		s.send(tgbotapi.NewMessage(msg.Chat.ID, text))
	case "remove_stock":
		s.send(tgbotapi.NewMessage(msg.Chat.ID, ""))
	case "add_notification":
		s.send(tgbotapi.NewMessage(msg.Chat.ID, ""))
	case "remove_notification":
		s.send(tgbotapi.NewMessage(msg.Chat.ID, ""))
	}
}

func (s *server) handleText(msg *tgbotapi.Message) {
	u, ok := s.getUser(msg.From.ID, msg.Chat.ID)
	if !ok {
		return
	}

	var err error
	var text string
	switch u.State {
	case models.UserStateAddStock:
		text, err = handleAddStockText(msg)
		if err != nil {
			log.Printf("something broken %v\n", err)
		} else {
			err = s.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateMenu)
			if err != nil {
				log.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
			}
		}
		s.send(tgbotapi.NewMessage(msg.Chat.ID, text))
	case models.UserStateRemoveStock:
	case models.UserStateDiff:
	case models.UserStateAddNotification:
	case models.UserStateRemoveNotification:

	default:
		s.send(tgbotapi.NewMessage(msg.Chat.ID, "NOTHING TO SEND!!!"))
	}

}

func (s *server) getUser(id, chatId int64) (*models.User, bool) {
	u, err := s.repo.GetUser(context.Background(), id)
	if err == nil {
		return u, true
	}

	var text string
	if errors.Is(err, ErrNotFound) {
		text = "Firstly type /start"
	} else {
		log.Printf("Failed to get user by id:%v err: %v\n", id, err)
		text = brokenMessage
	}

	msg := tgbotapi.NewMessage(chatId, text)
	s.send(msg)

	return nil, false
}

func (s *server) handleStartCommand(msg *tgbotapi.Message) string {
	//todo firstly create user at server and check if user already exist in db

	err := s.repo.CreateUser(context.Background(), msg.From.ID, msg.Chat.ID, 0)
	if err != nil {
		log.Printf("Failed to create user %v\n", err)

		return brokenMessage
	}

	return "Welcome to the stock market bot. To see all available command type /help"
}

func handleAddStockText(msg *tgbotapi.Message) (string, error) {
	splitMsg := strings.Split(msg.Text, " ")
	if len(splitMsg) != 2 {
		return "Type in the next format stockName amount", errors.New("")
	}

	amount, err := strconv.Atoi(splitMsg[1])
	if err != nil || amount <= 0 {
		return "amount must be a positive number", errors.New("")
	}

	//if stock == "" {
	//	return fmt.Sprintf("Couldn't found stock with name %v", splitMsg[1])
	//}
	//
	//stockRes, err := client.FindStock(context.Background(), &pb.StockName{Name: splitMsg[1]})
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("successfully client call %v\n", stockRes.GetName())
	//}

	return fmt.Sprintf("%v %v was added into portfolio", amount, splitMsg[0]), nil
}
