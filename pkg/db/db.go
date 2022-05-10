package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func New(ctx context.Context, adr string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, adr)
}
