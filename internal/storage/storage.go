package storage

import (
	"errors"
	"fmt"
)

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
	GetTargetInfo(target string) (*TargetInfo, error)
}

type TargetInfo struct {
	Target      string
	Created     string
	LastUpdated string
	LastRead    string
}

func (info TargetInfo) String() string {
	return fmt.Sprintf("Target:      %s\nCreated:     %s\nLastUpdated: %s\nLastRead:    %s",
		info.Target, info.Created, info.LastUpdated, info.LastRead)
}
