package repository

import (
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	"golang.org/x/net/context"
)

type Repository interface {
	CreateUser(ctx context.Context, id int64, chatId int64, serverUserId int32) error
	GetUser(ctx context.Context, id int64) (*models.User, error)
	GetUserByServerUserId(ctx context.Context, serverUserId int32) (*models.User, error)
	UpdateUserState(ctx context.Context, id int64, state string) error
}
