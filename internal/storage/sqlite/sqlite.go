package sqlite

import (
	"accelolabs/cat-api/internal/storage"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveURL(targetUrl string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM url WHERE url = ?)", targetUrl).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if exists {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
	}

	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM url WHERE alias = ?)", alias).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if exists {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrAliasExists)
	}

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(targetUrl, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	var url string
	err := s.db.QueryRow("SELECT url FROM url WHERE alias = ?", alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) GetAlias(targetUrl string) (string, error) {
	const op = "storage.sqlite.GetURL"

	var alias string
	err := s.db.QueryRow("SELECT alias FROM url WHERE url = ?", targetUrl).Scan(&alias)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrAliasNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return alias, nil
}
