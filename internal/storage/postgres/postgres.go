package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"mini-blog/internal/config"
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
