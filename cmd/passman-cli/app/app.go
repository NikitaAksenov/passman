package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/NikitaAksenov/passman/internal/directory"
	"github.com/NikitaAksenov/passman/internal/storage"
	"github.com/NikitaAksenov/passman/internal/storage/sqlite"
)

const (
	appName     = "passman"
	storageDir  = "storage"
	storageName = "storage.db"
)

type App struct {
	Storage storage.Storage
}

func New() (*App, error) {
	// - Check app directories and create them if needed
	// -- Get user config dir
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config dir")
	}

	// -- Check and create app dir
	appDir := filepath.Join(userConfigDir, appName)
	err = directory.CheckAndCreateDir(appDir)
	if err != nil {
		return nil, fmt.Errorf("failed to check app dir: %s", err.Error())
	}

	// -- Check and create storage dir
	storageDir := filepath.Join(appDir, storageDir)
	err = directory.CheckAndCreateDir(storageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to check storage dir %s", err.Error())
	}
	storagePath := filepath.Join(storageDir, storageName)

	// - Init storage
	var storage storage.Storage
	storage, err = sqlite.New(storagePath)
	if err != nil {
		return nil, fmt.Errorf("ailed to init storage: %s", err.Error())
	}

	return &App{
		Storage: storage,
	}, nil
}
