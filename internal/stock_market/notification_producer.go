package stock_market

import (
	"context"
	"encoding/json"
	"github.com/piquette/finance-go/quote"
	"github.com/streadway/amqp"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/models"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/logger"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/queue"
	"time"
)

const sleepTime = time.Hour

func TrackNotification(repo Repository, url string) {
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

	for {
		ids, err := repo.GetUserIdsWithNotifications(context.Background())
		if err != nil {
			logger.Error.Printf("Failed to get user ids with notifications: %s", err)
			time.Sleep(sleepTime)

			continue
		}

		for _, id := range ids {
			ntfs, err := repo.GetNotifications(context.Background(), id)
			if err != nil {
				logger.Error.Printf("Failed to get notifications of user %s: %s", id, err)
				continue
			}

			ntfByStockName := make(map[string]models.Notification, len(ntfs))
			stockNames := make([]string, 0, len(ntfs))
			for _, n := range ntfs {
				ntfByStockName[n.StockName] = n
				stockNames = append(stockNames, n.StockName)
			}

			quotes := quote.List(stockNames)
			for quotes.Next() {
				bid := quotes.Quote().Bid
				n := ntfByStockName[quotes.Quote().Symbol]
				if (n.Type == models.NotificationTypeUp && bid < n.Threshold) ||
					(n.Type == models.NotificationTypeDown && bid > n.Threshold) {
					continue
				}

				body, err := json.Marshal(queue.Notification{
					StockName:  n.StockName,
					UserId:     n.UserId,
					Threshold:  n.Threshold,
					StockPrice: bid,
					EventTime:  time.Now(),
				})
				if err != nil {
					logger.Error.Printf("Error to marshalized: %s", err)
					continue
				}

				err = amqpChannel.Publish("", q.Name, false, false, amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         body,
				})

				if err != nil {
					logger.Error.Printf("Error publishing message: %s", err)
					continue
				}

				err = repo.RemoveNotification(context.Background(), n.Id)
				if err != nil {
					logger.Error.Printf("Failed to remove notification %s: %s", n.Id, err)
				}
			}

			if err = quotes.Err(); err != nil {
				logger.Error.Printf("Failed to get quotes: %s", err)
			}

		}

		time.Sleep(sleepTime)
	}
}
