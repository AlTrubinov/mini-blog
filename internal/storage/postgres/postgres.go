package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"mini-blog/internal/config"
	"mini-blog/internal/models/note"
	"mini-blog/internal/models/user"
	"mini-blog/pkg/apperror"
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
		return nil, fmt.Errorf("%w: unable to create connection pool", err)
	}
	if err = dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("%w: unable to ping connection pool", err)
	}

	return &Storage{db: dbPool}, nil
}

func (storage *Storage) Close() {
	storage.db.Close()
}

func (storage *Storage) SaveUser(ctx context.Context, username string, passwordHash string) (int64, error) {
	var u user.User

	err := storage.db.QueryRow(
		ctx,
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, created_at",
		username, passwordHash,
	).Scan(&u.Id, &u.Username, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("%w: timeout while save user", apperror.ErrTimeout)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%w: username already exists", apperror.ErrValidation)
		}
		return 0, fmt.Errorf("%w: failed to save user", apperror.ErrInternal)
	}

	return u.Id, nil
}

func (storage *Storage) GetUser(ctx context.Context, username string) (user.User, error) {
	var u user.User

	stmt := "SELECT id, username, created_at, password FROM users WHERE username = $1"
	err := storage.db.QueryRow(ctx, stmt, username).Scan(&u.Id, &u.Username, &u.CreatedAt, &u.Password)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return user.User{}, fmt.Errorf("%w: user not found", apperror.ErrNotFound)
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return user.User{}, fmt.Errorf("%w: timeout while retrieving user", apperror.ErrTimeout)
	case err != nil:
		return user.User{}, fmt.Errorf("%w: failed to get user", apperror.ErrInternal)
	default:
		return u, nil
	}
}

func (storage *Storage) CreateNote(ctx context.Context, userId int64, title string, content string) (int64, error) {
	var n note.Note

	err := storage.db.QueryRow(
		ctx,
		"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id",
		userId, title, content,
	).Scan(&n.Id)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("%w: context canceled or timed out", apperror.ErrTimeout)
		}
		return 0, fmt.Errorf("%w: failed to create note", apperror.ErrInternal)
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
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: context canceled or timed out", apperror.ErrTimeout)
		}
		return nil, fmt.Errorf("%w: failed to get notes list", apperror.ErrInternal)
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
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: rows iteration timeout", apperror.ErrTimeout)
		}
		return nil, fmt.Errorf("%w: rows iteration error", apperror.ErrInternal)
	}

	return resNotes, nil
}

func (storage *Storage) GetUserNote(ctx context.Context, userId int64, noteId int64) (note.Note, error) {
	var n note.Note

	stmt := "SELECT id,user_id,title,content,created_at,updated_at FROM notes WHERE user_id = $1 AND id = $2"

	err := storage.db.QueryRow(ctx, stmt, userId, noteId).Scan(&n.Id, &n.UserId, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return note.Note{}, fmt.Errorf("%w: note %d not found", apperror.ErrNotFound, noteId)
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return note.Note{}, fmt.Errorf("%w: timeout while retrieving note", apperror.ErrTimeout)
	case err != nil:
		return note.Note{}, fmt.Errorf("%w: failed to get note", apperror.ErrInternal)
	default:
		return n, nil
	}
}

func (storage *Storage) UpdateNote(ctx context.Context, userId int64, noteId int64, title string, content string) error {
	stmt := "UPDATE notes SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE user_id = $3 AND id = $4"

	res, err := storage.db.Exec(ctx, stmt, title, content, userId, noteId)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: context canceled or timed out", apperror.ErrTimeout)
		}
		return fmt.Errorf("%w: failed to update note", apperror.ErrInternal)
	}

	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("%w: note %d not found", apperror.ErrNotFound, noteId)
	}

	return nil
}

func (storage *Storage) DeleteNote(ctx context.Context, userId int64, noteId int64) error {
	stmt := "DELETE FROM notes WHERE user_id = $1 AND id = $2"

	res, err := storage.db.Exec(ctx, stmt, userId, noteId)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%w: context canceled or timed out", apperror.ErrTimeout)
		}
		return fmt.Errorf("%w: failed to delete note", apperror.ErrInternal)
	}

	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("%w: note %d not found", apperror.ErrNotFound, noteId)
	}

	return nil
}
