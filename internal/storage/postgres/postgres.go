package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"mini-blog/internal/config"
	"mini-blog/internal/models/user"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, DBConf config.Database) (*Storage, error) {
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		DBConf.User,
		DBConf.Password,
		DBConf.Host,
		DBConf.Port,
		DBConf.Name,
	)
	dbPool, err := pgxpool.New(ctx, dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err = dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("unable to ping connection pool: %w", err)
	}

	return &Storage{db: dbPool}, nil
}

func (storage *Storage) Close() {
	storage.db.Close()
}

func (storage *Storage) SaveUser(ctx context.Context, username string) (int64, error) {
	var u user.User

	err := storage.db.QueryRow(ctx,
		"INSERT INTO users (username) VALUES ($1) RETURNING id, username, created_at",
		username,
	).Scan(&u.Id, &u.Username, &u.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("save url error, %w", err)
	}

	return u.Id, nil
}
