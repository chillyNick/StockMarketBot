package telegram_bot

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/streadway/amqp"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/queue"
	"os"
)

func TrackNotification(s *Server, url string) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error.Fatalf("can't connect to AMQP: %s", err)
	}
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	if err != nil {
		logger.Error.Fatalf("can't create a amqpChannel: %s", err)
	}
	defer amqpChannel.Close()

	q := queue.CreateNotificationQueue(amqpChannel)

	err = amqpChannel.Qos(1, 0, false)
	if err != nil {
		logger.Error.Fatalf("could not configure QoS: %s", err)
	}

	messageChannel, err := amqpChannel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error.Fatalf("could not register consumer: %s", err)
	}

	logger.Info.Printf("Consumer ready, PID: %d", os.Getpid())
	for d := range messageChannel {
		logger.Info.Printf("Received a message: %s", d.Body)

		notification := queue.Notification{}
		err := json.Unmarshal(d.Body, &notification)
		if err != nil {
			logger.Error.Printf("Error decoding JSON: %s", err)
			continue
		}

		u, err := s.Repo.GetUserByServerUserId(context.Background(), notification.UserId)
		if err != nil {
			logger.Error.Printf("Error to get user by userServerId: %v %s", notification.UserId, err)
			continue
		}

		text := fmt.Sprintf(
			"%s was reached %v at %s, price %v",
			notification.StockName,
			notification.Threshold,
			notification.EventTime,
			notification.StockPrice,
		)

		msgConfig := tgbotapi.NewMessage(u.ChatId, text)
		if !s.send(msgConfig) {
			continue
		}

		if err := d.Ack(false); err != nil {
			logger.Error.Printf("Error acknowledging message : %s", err)
		} else {
			logger.Info.Printf("Acknowledged message")
		}
	}
}
