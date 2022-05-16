package stock_market

import (
	"context"
	"encoding/json"
	"github.com/piquette/finance-go/quote"
	"github.com/streadway/amqp"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/models"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/queue"
	"log"
	"time"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const sleepTime = time.Hour

func TrackNotification(repo repository, url string) {
	conn, err := amqp.Dial(url)
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	q, err := amqpChannel.QueueDeclare("notification", true, false, false, false, nil)
	handleError(err, "Could not declare `notification` queue")

	for {
		ids, err := repo.GetUserIdsWithNotifications(context.Background())
		if err != nil {
			log.Println(err)
			time.Sleep(sleepTime)

			continue
		}

		for _, id := range ids {
			ntfs, err := repo.GetNotifications(context.Background(), id)
			if err != nil {
				println(err)
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
				quote := quotes.Quote()
				n := ntfByStockName[quote.Symbol]
				if (n.Type == models.NotificationTypeUp && quote.Bid < n.Threshold) ||
					(n.Type == models.NotificationTypeDown && quote.Bid > n.Threshold) {
					continue
				}

				body, err := json.Marshal(queue.Notification{
					StockName:  n.StockName,
					UserId:     n.UserId,
					Threshold:  n.Threshold,
					StockPrice: quote.Bid,
					EventTime:  time.Now(),
				})
				if err != nil {
					log.Println(err)
					continue
				}

				err = amqpChannel.Publish("", q.Name, false, false, amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         body,
				})

				if err != nil {
					log.Printf("Error publishing message: %s", err)
					continue
				}

				err = repo.RemoveNotification(context.Background(), n.Id)
				if err != nil {
					log.Println(err)
				}
			}

			if err = quotes.Err(); err != nil {
				log.Println(err)
			}

		}

		time.Sleep(sleepTime)
	}
}
