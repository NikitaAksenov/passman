package storage

import "errors"

var (
	ErrEmptyTarget    = errors.New("target is empty")
	ErrEmptyPassword  = errors.New("password is empty")
	ErrTargetExist    = errors.New("target already exists")
	ErrTargetNotFound = errors.New("target not found")
)

type Storage interface {
	AddPass(target string, pass string) (int64, error)
	GetPass(target string) (string, error)
	GetTargets(limit int, offset int) ([]string, error)
	DeleteTarget(target string) (int64, error)
	UpdatePassword(target string, pass string) (int64, error)
}
