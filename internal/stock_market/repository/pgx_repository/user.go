package pgx_repository

import "context"

func (r *repository) CreateUser(ctx context.Context) (id int32, err error) {
	const query = `
		INSERT into "user" (id) values (default) returning id
	`
	err = r.pool.QueryRow(ctx, query).Scan(&id)

	return
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
