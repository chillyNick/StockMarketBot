package queue

import (
	"github.com/streadway/amqp"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/logger"
	"time"
)

type Notification struct {
	StockName  string
	UserId     int32
	Threshold  float64
	StockPrice float64
	EventTime  time.Time
}

func CreateNotificationQueue(channel *amqp.Channel) amqp.Queue {
	q, err := channel.QueueDeclare("notification", true, false, false, false, nil)
	if err != nil {
		logger.Error.Fatalf("could not declare `notification` queue: %s", err)
	}

	return q
}
