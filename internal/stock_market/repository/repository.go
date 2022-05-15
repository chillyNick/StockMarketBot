package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/models"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) CreateUser(ctx context.Context) (int32, error) {
	const query = `
		INSERT into "user" (id) values (default) returning id
	`
	var id int32
	err := r.pool.QueryRow(ctx, query).Scan(&id)

	return id, err
}

func (r *repository) GetUserIdsWithNotifications(ctx context.Context) ([]int32, error) {
	const query = `
		SELECT u.id
		FROM "user" AS u
        JOIN notification AS n ON u.id = n.user_id
		GROUP BY u.id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int32, 0)
	for rows.Next() {
		var id int32
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

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

func (r *repository) AddNotification(ctx context.Context, userId int32, stockName string, threshold float64, nType string) error {
	const query = `
		INSERT INTO notification (
			stock_name, user_id, threshold, type
		) VALUES (
			$1, $2, $3, $4
		)
	`

	_, err := r.pool.Exec(ctx, query,
		stockName,
		userId,
		threshold,
		nType,
	)

	return err
}

func (r *repository) RemoveNotification(ctx context.Context, id int32) error {
	const query = `
		DELETE FROM notification
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id)

	return err
}

func (r *repository) GetNotifications(ctx context.Context, userId int32) ([]models.Notification, error) {
	const query = `
		SELECT
			id,
			stock_name,
			user_id,
			threshold,
			type
		FROM notification
		WHERE user_id = $1
	`

	rows, err := r.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ntfs := make([]models.Notification, 0)
	for rows.Next() {
		var n models.Notification
		err = rows.Scan(
			&n.Id,
			&n.StockName,
			&n.UserId,
			&n.Threshold,
			&n.Type,
		)
		if err != nil {
			return nil, err
		}

		ntfs = append(ntfs, n)
	}

	return ntfs, nil
}
