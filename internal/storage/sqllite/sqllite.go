package sqllite

import (
	"database/sql"
	"errors"
	"fmt"
	"urlShotener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.sqllite.New"
	db, err := sql.Open("sqlite3", dbPath)
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
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

// Сохраняет новый URL и его Alias в базе данных
// Возвращает id записи
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqllite.SaveURL"
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqllite.GetURL"
	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("%s error Prepare: %w", op, err)
	}
	var recieved string
	err = stmt.QueryRow(alias).Scan(&recieved)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s error Query: %w", op, err)
	}
	return recieved, nil
}

func (s *Storage) DeleteURL(alias string) (int64, error) {
	const op = "storage.sqllite.DeleteURL"
	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias = ?`)
	if err != nil {
		return 0, fmt.Errorf("%s error Prepare: %w", op, err)
	}
	res, err := stmt.Exec(alias)
	if err != nil {
		return 0, fmt.Errorf("%s error Execute: %w", op, err)
	}
	rowsCount, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows", op, err)
	}
	return rowsCount, nil
}
