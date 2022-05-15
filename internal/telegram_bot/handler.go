package telegram_bot

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	pb "gitlab.ozon.dev/chillyNick/homework-2/pkg/api"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strconv"
	"strings"
)

const brokenMessage = "Something broken please try again later"

func (s *server) handle(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		s.handleCommand(update.Message)
	} else {
		s.handleText(update.Message)
	}
}

var periods = []string{"hour", "day", "week", "all"}

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

	u, ok := s.getUser(msg.From.ID, msg.Chat.ID)
	if !ok {
		s.send(tgbotapi.NewMessage(msg.Chat.ID, "Firstly type /start"))

		return
	}

	switch msg.Command() {
	case "show":
		stocks, err := s.grpcClient.GetStocks(context.Background(), &pb.UserId{Id: u.ServerUserId})
		if err != nil {
			log.Println(err)
			s.send(tgbotapi.NewMessage(msg.Chat.ID, brokenMessage))
		}
		text = "Portfolio:\n"
		for _, stock := range stocks.GetStocks() {
			text += fmt.Sprintf("%v %v \n", stock.GetName(), stock.GetAmount())
		}

		s.send(tgbotapi.NewMessage(msg.Chat.ID, text))
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
		err := s.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateRemoveStock)
		if err != nil {
			log.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
			text = brokenMessage
		} else {
			text = "Send a message in the next format: stockName amount"
		}

		s.send(tgbotapi.NewMessage(msg.Chat.ID, text))
	case "add_notification":
		err := s.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateAddNotification)
		if err != nil {
			log.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
			text = brokenMessage
		} else {
			text = "Send a message in the next format: stockName threshold"
		}

		s.send(tgbotapi.NewMessage(msg.Chat.ID, text))
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
		text, err = s.handleAddStockText(msg, u)
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
		text, err = s.handleRemoveStockText(msg, u)
		if err != nil {
			log.Printf("something broken %v\n", err)
		} else {
			err = s.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateMenu)
			if err != nil {
				log.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
			}
		}
		s.send(tgbotapi.NewMessage(msg.Chat.ID, text))
	case models.UserStateDiff:
		ok := false
		for _, p := range periods {
			if msg.Text == p {
				ok = true
				break
			}
		}

		if !ok {
			s.send(tgbotapi.NewMessage(msg.Chat.ID, "Incorrect period"))

			return
		}

		changes, err := s.grpcClient.GetPortfolioChanges(context.Background(), &pb.GetPortfolioChangesRequest{Period: pb.Period(pb.Period_value[strings.ToUpper(msg.Text)]), UserId: &pb.UserId{Id: u.ServerUserId}})
		if err != nil {
			s.send(tgbotapi.NewMessage(msg.Chat.ID, brokenMessage))

			return
		}

		err = s.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateMenu)
		if err != nil {
			log.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
			s.send(tgbotapi.NewMessage(msg.Chat.ID, brokenMessage))

			return
		}

		text = fmt.Sprintf("Portfolio changes by %v:\n", msg.Text)
		for _, ch := range changes.GetStocks() {
			text += fmt.Sprintf("%v %v %v \n", ch.Stock.GetName(), ch.Stock.GetAmount(), ch.GetOldPrice()-ch.GetCurrentPrice())
		}

		msgConfig := tgbotapi.NewMessage(msg.Chat.ID, text)
		msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		s.send(msgConfig)
	case models.UserStateAddNotification:
		text, err = s.handleAddNotificationText(msg, u)
		if err != nil {
			log.Printf("something broken %v\n", err)
		} else {
			err = s.repo.UpdateUserState(context.Background(), msg.From.ID, models.UserStateMenu)
			if err != nil {
				log.Printf("Failed to update user state with id:%v err: %v\n", msg.From.ID, err)
			}
		}
		s.send(tgbotapi.NewMessage(msg.Chat.ID, text))
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
	if errors.Is(err, db.ErrNotFound) {
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
	_, ok := s.getUser(msg.From.ID, msg.Chat.ID)
	if ok {
		return "Welcome to the stock market bot. To see all available command type /help"
	}

	id, err := s.grpcClient.CreateUser(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Printf("Failed to create user %v\n", err)

		return brokenMessage
	}

	err = s.repo.CreateUser(context.Background(), msg.From.ID, msg.Chat.ID, id.Id)
	if err != nil {
		log.Printf("Failed to create user %v\n", err)

		return brokenMessage
	}

	return "Welcome to the stock market bot. To see all available command type /help"
}

func (s *server) handleAddStockText(msg *tgbotapi.Message, u *models.User) (string, error) {
	splitMsg := strings.Split(msg.Text, " ")
	if len(splitMsg) != 2 {
		return "Type in the next format stockName amount", errors.New("")
	}

	amount, err := strconv.Atoi(splitMsg[1])
	if err != nil || amount <= 0 {
		return "amount must be a positive number", errors.New("")
	}

	_, err = s.grpcClient.AddStock(context.Background(), &pb.StockRequest{
		Name:   splitMsg[0],
		Amount: int32(amount),
		UserId: &pb.UserId{Id: u.ServerUserId},
	})

	if err == nil {
		return fmt.Sprintf("%v %v was added into portfolio", amount, splitMsg[0]), nil
	}

	if status.Code(err) == codes.NotFound {
		return "Stock with such name not found", err
	}

	return brokenMessage, err
}

func (s *server) handleRemoveStockText(msg *tgbotapi.Message, u *models.User) (string, error) {
	splitMsg := strings.Split(msg.Text, " ")
	if len(splitMsg) != 2 {
		return "Type in the next format stockName amount", errors.New("")
	}

	amount, err := strconv.Atoi(splitMsg[1])
	if err != nil || amount <= 0 {
		return "amount must be a positive number", errors.New("")
	}

	_, err = s.grpcClient.RemoveStock(context.Background(), &pb.StockRequest{
		Name:   splitMsg[0],
		Amount: int32(amount),
		UserId: &pb.UserId{Id: u.ServerUserId},
	})

	if err == nil {
		return fmt.Sprintf("%v %v was remove into portfolio", amount, splitMsg[0]), nil
	}

	if status.Code(err) == codes.NotFound {
		return "Stock with such name not found or you don't have in portfolio", err
	}

	return brokenMessage, err
}

func (s *server) handleAddNotificationText(msg *tgbotapi.Message, u *models.User) (string, error) {
	splitMsg := strings.Split(msg.Text, " ")
	if len(splitMsg) != 2 {
		return "Type in the next format stockName threshold", errors.New("")
	}

	threshold, err := strconv.ParseFloat(splitMsg[1], 64)
	if err != nil {
		return "Threshold must be a positive decimal number", errors.New("")
	}

	_, err = s.grpcClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
		UserId:    &pb.UserId{Id: u.ServerUserId},
		StockName: splitMsg[0],
		Threshold: threshold,
	})

	if err == nil {
		return fmt.Sprintf("Notification was added"), nil
	}

	if status.Code(err) == codes.NotFound {
		return "Stock with such name not found", err
	}

	return brokenMessage, err
}
