package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rehydrate1/shorty/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveURL(ctx context.Context, alias, urlToSave string) error {
	query := `INSERT INTO links (short_key, original_url) VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, alias, urlToSave)
	return err
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	query := `SELECT original_url FROM links WHERE short_key=$1`

	var longURL string
	if err := s.db.QueryRowContext(ctx, query, alias).Scan(&longURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", err
	}

	return longURL, nil
}