package storage

import "errors"

var (
	ErrTargetExist    = errors.New("target already exists")
	ErrTargetNotFound = errors.New("target not found")
)
