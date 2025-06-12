package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mini-blog/internal/config"
	"mini-blog/internal/models/note"
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

	err := storage.db.QueryRow(
		ctx,
		"INSERT INTO users (username) VALUES ($1) RETURNING id, username, created_at",
		username,
	).Scan(&u.Id, &u.Username, &u.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("save url error, %w", err)
	}

	return u.Id, nil
}

func (storage *Storage) CreateNote(ctx context.Context, userId int64, title string, content string) (int64, error) {
	var n note.Note

	err := storage.db.QueryRow(
		ctx,
		"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id",
		userId, title, content,
	).Scan(&n.Id)
	if err != nil {
		return 0, fmt.Errorf("save note error, %w", err)
	}

	return n.Id, nil
}

func (storage *Storage) GetUserNotes(ctx context.Context, userId int64, limit int, offset int, order string) ([]note.Note, error) {
	var resNotes []note.Note

	stmt := fmt.Sprintf(
		"SELECT id,user_id,title,content,created_at,updated_at FROM notes WHERE user_id = $1 ORDER BY id %s LIMIT $2 OFFSET $3",
		order,
	)

	rows, err := storage.db.Query(ctx, stmt, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get notes list error, %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n note.Note
		err := rows.Scan(&n.Id, &n.UserId, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			continue
		}
		resNotes = append(resNotes, n)
	}

	if err := rows.Err(); err != nil {
		slog.Error("rows error, %w", err)
	}
	return resNotes, nil
}

func (storage *Storage) GetUserNote(ctx context.Context, userId int64, noteId int64) (note.Note, error) {
	var n note.Note

	stmt := "SELECT id,user_id,title,content,created_at,updated_at FROM notes WHERE user_id = $1 AND id = $2"

	err := storage.db.QueryRow(ctx, stmt, userId, noteId).Scan(&n.Id, &n.UserId, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return note.Note{}, err
	} else if err != nil {
		return note.Note{}, fmt.Errorf("get note error, %w", err)
	}
	return n, nil
}

func (storage *Storage) UpdateNote(ctx context.Context, userId int64, noteId int64, title string, content string) error {
	stmt := "UPDATE notes SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE user_id = $3 AND id = $4"

	res, err := storage.db.Exec(ctx, stmt, title, content, userId, noteId)
	if err != nil {
		return fmt.Errorf("update note error, %w", err)
	}

	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}
