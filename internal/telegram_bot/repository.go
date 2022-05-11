package telegram_bot

import (
	"context"
	"errors"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
)

var ErrNotFound = errors.New("user not found")

type repository interface {
	CreateUser(ctx context.Context, id int64, chatId int64, serverUserId int64) error
	GetUser(ctx context.Context, id int64) (*models.User, error)
	UpdateUserState(ctx context.Context, id int64, state string) error
}
