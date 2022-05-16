package models

const (
	NotificationTypeUp   = "up"
	NotificationTypeDown = "down"
)

type Notification struct {
	Id        int32
	StockName string
	UserId    int32
	Threshold float64
	Type      string
}
