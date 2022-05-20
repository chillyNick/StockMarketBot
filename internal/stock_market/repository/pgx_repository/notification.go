package pgx_repository

import (
	"gitlab.ozon.dev/chillyNick/homework-2/internal/stock_market/models"
	"golang.org/x/net/context"
)

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
