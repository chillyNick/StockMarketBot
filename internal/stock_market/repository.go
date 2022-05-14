package stock_market

import (
	"context"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/models"
)

type repository interface {
	CreateUser(ctx context.Context) (int32, error)
	GetStock(ctx context.Context, userId int32, name string) (*models.Stock, error)
	GetStocks(ctx context.Context, userId int32) ([]models.Stock, error)
	UpdateStock(ctx context.Context, id, amount int32, price float64) error
	AddStock(ctx context.Context, userId int32, name string, amount int32, price float64) error
	RemoveStock(ctx context.Context, id int32) error
}
