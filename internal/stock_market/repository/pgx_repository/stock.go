package pgx_repository

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/db"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/models"
	"golang.org/x/net/context"
)

func (r *repository) GetStock(ctx context.Context, userId int32, name string) (*models.Stock, error) {
	const query = `
		SELECT
			id,
			name,
			user_id,
			amount,
			price,
			created_at
		FROM stock
		WHERE user_id = $1 AND name = $2
	`

	s := models.Stock{}

	err := r.pool.QueryRow(ctx, query, userId, name).Scan(
		&s.Id,
		&s.Name,
		&s.UserId,
		&s.Amount,
		&s.Price,
		&s.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	return &s, nil
}

func (r *repository) GetStocks(ctx context.Context, userId int32) ([]models.Stock, error) {
	const query = `
		SELECT
			id,
			name,
			user_id,
			amount,
			price,
			created_at
		FROM stock
		WHERE user_id = $1
	`

	rows, err := r.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stocks := make([]models.Stock, 0)
	for rows.Next() {
		var s models.Stock
		err = rows.Scan(
			&s.Id,
			&s.Name,
			&s.UserId,
			&s.Amount,
			&s.Price,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		stocks = append(stocks, s)
	}

	return stocks, nil
}

func (r *repository) UpdateStock(ctx context.Context, id, amount int32, price float64) error {
	const query = `
		UPDATE stock
		set amount = $2, price = $3
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id, amount, price)

	return err
}

func (r *repository) AddStock(ctx context.Context, userId int32, name string, amount int32, price float64) error {
	const query = `
		INSERT INTO stock (
			name, user_id, amount, price
		) VALUES (
			$1, $2, $3, $4
		)
	`

	_, err := r.pool.Exec(ctx, query,
		name,
		userId,
		amount,
		price,
	)

	return err
}

func (r *repository) RemoveStock(ctx context.Context, id int32) error {
	const query = `
		DELETE FROM stock
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id)

	return err
}
