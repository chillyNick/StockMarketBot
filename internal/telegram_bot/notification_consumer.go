package telegram_bot

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/streadway/amqp"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/queue"
	"log"
	"os"
	"time"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const sleepTime = time.Hour

func TrackNotification(s *server, url string) {
	conn, err := amqp.Dial(url)
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	q, err := amqpChannel.QueueDeclare("notification", true, false, false, false, nil)
	handleError(err, "Could not declare `notification` queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	log.Printf("Consumer ready, PID: %d", os.Getpid())
	for d := range messageChannel {
		log.Printf("Received a message: %s", d.Body)

		notification := queue.Notification{}
		err := json.Unmarshal(d.Body, &notification)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			continue
		}

		u, err := s.repo.GetUserByServerUserId(context.Background(), notification.UserId)
		if err != nil {
			log.Printf("Error to get user by userServerId: %v %s", notification.UserId, err)
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
		_, err = s.bot.Send(msgConfig)
		if err != nil {
			log.Printf("Failed to send tg message: %s", err)
			continue
		}

		if err := d.Ack(false); err != nil {
			log.Printf("Error acknowledging message : %s", err)
		} else {
			log.Printf("Acknowledged message")
		}
	}
}
