package pgx_repository

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/chillyNick/homework-2/internal/telegram_bot/models"
	"gitlab.ozon.dev/chillyNick/homework-2/pkg/db"
	"golang.org/x/net/context"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) CreateUser(ctx context.Context, id int64, chatId int64, serverUserId int32) error {
	const query = `
		INSERT INTO "user" (
			id,
			chat_id,
			server_user_id,
			state
		) VALUES (
			$1, $2, $3, $4
		)
	`

	_, err := r.pool.Exec(ctx, query,
		id,
		chatId,
		serverUserId,
		models.UserStateMenu,
	)

	return err
}

func (r *repository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	const query = `
		SELECT
			id,
			chat_id,
			server_user_id,
			state
		FROM "user"
		WHERE id = $1
	`

	u := models.User{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.Id,
		&u.ChatId,
		&u.ServerUserId,
		&u.State,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	return &u, nil
}

func (r *repository) GetUserByServerUserId(ctx context.Context, serverUserId int32) (*models.User, error) {
	const query = `
		SELECT
			id,
			chat_id,
			server_user_id,
			state
		FROM "user"
		WHERE server_user_id = $1
	`

	u := models.User{}

	err := r.pool.QueryRow(ctx, query, serverUserId).Scan(
		&u.Id,
		&u.ChatId,
		&u.ServerUserId,
		&u.State,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, db.ErrNotFound
	}

	return &u, nil
}

func (r *repository) UpdateUserState(ctx context.Context, id int64, state string) error {
	const query = `
		UPDATE "user"
		SET state = $2
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, state)

	return err
}
