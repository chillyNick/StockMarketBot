package telegram_bot

import (
	"context"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
)

type repository interface {
	CreateUser(ctx context.Context, id int64, chatId int64, serverUserId int32) error
	GetUser(ctx context.Context, id int64) (*models.User, error)
	GetUserByServerUserId(ctx context.Context, serverUserId int32) (*models.User, error)
	UpdateUserState(ctx context.Context, id int64, state string) error
}
