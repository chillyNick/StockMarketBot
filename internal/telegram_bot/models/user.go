package models

const (
	UserStateMenu            = "menu"
	UserStateAddStock        = "add_stock"
	UserStateRemoveStock     = "remove_stock"
	UserStateDiff            = "diff"
	UserStateAddNotification = "add_notification"
)

type User struct {
	Id           int64
	ChatId       int64
	ServerUserId int32
	State        string
}
