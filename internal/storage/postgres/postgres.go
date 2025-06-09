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
		errMsg := fmt.Sprintf("unable to create connection pool: %v", err.Error())
		return nil, fmt.Errorf(errMsg)
	}
	if err = dbPool.Ping(ctx); err != nil {
		errMsg := fmt.Sprintf("unable to ping connection pool: %v", err.Error())
		dbPool.Close()
		return nil, fmt.Errorf(errMsg)
	}

	return &Storage{db: dbPool}, nil
}

func (storage *Storage) Close() {
	storage.db.Close()
}
