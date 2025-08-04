package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NikitaAksenov/passman/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

func New(storagePath string) (*SqliteStorage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS pass(
		id INTEGER PRIMARY KEY,
		target TEXT NOT NULL UNIQUE,
		pass TEXT NOT NULL,
		created TEXT NOT NULL,
		lastUpdate TEXT,
		lastRead TEXT);
	`)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &SqliteStorage{db: db}, nil
}

func (s *SqliteStorage) AddPass(target string, pass string) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO pass(target, pass, created) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	created, err := PrepareTime(time.Now().UTC())
	if err != nil {
		return 0, fmt.Errorf("failed to prepare time: %s", err)
	}

	res, err := stmt.Exec(target, pass, created)
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

func (s *SqliteStorage) GetPass(target string) (string, error) {
	stmt, err := s.db.Prepare("SELECT pass FROM pass WHERE target = ?")
	if err != nil {
		return "", err
	}

	var resultPass string

	err = stmt.QueryRow(target).Scan(&resultPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrTargetNotFound
		}

		return "", err
	}

	return resultPass, nil
}

func (s *SqliteStorage) GetTargets(limit int, offset int) ([]string, error) {
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

func (s *SqliteStorage) DeleteTarget(target string) (int64, error) {
	stmt, err := s.db.Prepare("DELETE FROM pass WHERE target = ?")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(target)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (s *SqliteStorage) UpdatePassword(target string, pass string) (int64, error) {
	if target == "" {
		return 0, storage.ErrEmptyTarget
	}

	if pass == "" {
		return 0, storage.ErrEmptyPassword
	}

	lastUpdate, err := PrepareTime(time.Now().UTC())
	if err != nil {
		return 0, fmt.Errorf("failed to prepare time: %s", err)
	}

	stmt, err := s.db.Prepare("UPDATE pass SET pass = ?, lastUpdate = ? WHERE target = ?")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(pass, lastUpdate, target)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s *SqliteStorage) GetTargetInfo(target string) (*storage.TargetInfo, error) {
	if target == "" {
		return nil, storage.ErrEmptyTarget
	}

	stmt, err := s.db.Prepare("SELECT created, lastUpdate, lastRead FROM pass WHERE target = ?")
	if err != nil {
		return nil, err
	}

	var created string
	var lastUpdated, lastRead sql.NullString

	err = stmt.QueryRow(target).Scan(&created, &lastUpdated, &lastRead)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrTargetNotFound
		}

		return nil, err
	}

	return &storage.TargetInfo{
		Target:      target,
		Created:     created,
		LastUpdated: lastUpdated.String,
		LastRead:    lastRead.String,
	}, nil
}

func PrepareTime(t time.Time) (string, error) {
	return t.Format(time.RFC822), nil
}
