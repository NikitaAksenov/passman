package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikitaAksenov/passman/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS pass(
		id INTEGER PRIMARY KEY,
		target TEXT NOT NULL UNIQUE,
		pass TEXT NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddPass(target string, pass string) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO pass(target, pass) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	res, err := stmt.Exec(target, pass)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%w", storage.ErrTargetExist)
		}

		return 0, fmt.Errorf("%w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func (s *Storage) GetPass(target string) (string, error) {
	stmt, err := s.db.Prepare("SELECT pass FROM pass WHERE target = ?")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	var resultPass string

	err = stmt.QueryRow(target).Scan(&resultPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrTargetNotFound
		}

		return "", fmt.Errorf("%w", err)
	}

	return resultPass, nil
}

func (s *Storage) GetTargets(limit int, offset int) ([]string, error) {
	stmt, err := s.db.Prepare("SELECT target FROM pass DESC LIMIT ? OFFSET ?")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []string

	for rows.Next() {
		var target string
		err = rows.Scan(&target)
		if err != nil {
			return nil, err
		}

		targets = append(targets, target)
	}

	return targets, nil
}
