package db

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/net/context"
)

func New(ctx context.Context, adr string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, adr)
}
