package queue

import "time"

type Notification struct {
	StockName  string
	UserId     int32
	Threshold  float64
	StockPrice float64
	EventTime  time.Time
}
