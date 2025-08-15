package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/NikitaAksenov/passman/internal/directory"
	"github.com/NikitaAksenov/passman/internal/storage"
	"github.com/NikitaAksenov/passman/internal/storage/sqlite"
	"golang.design/x/clipboard"
)

var appConfiguration = "dev"

const (
	AppName = "passman"

	storageDir  = "storage"
	storageName = "storage.db"
)

type App struct {
	Storage storage.Storage
}

func New() (*App, error) {
	switch appConfiguration {
	case "dev":
		fmt.Println("dev configuration")
	}

	// Check and create app dir
	appDir, err := GetAppDir(appConfiguration)
	if err != nil {
		return nil, fmt.Errorf("failed to get app dir: %s", err.Error())
	}

	// Check and create storage dir
	storageDir := filepath.Join(appDir, storageDir)
	err = directory.CheckAndCreateDir(storageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to check storage dir %s", err.Error())
	}
	storagePath := filepath.Join(storageDir, storageName)

	// Init storage
	var storage storage.Storage
	storage, err = sqlite.New(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %s", err.Error())
	}

	// Init clipboard
	err = clipboard.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to init clipboard: %s", err.Error())
	}

	return &App{
		Storage: storage,
	}, nil
}

func GetAppDir(configuration string) (string, error) {
	switch configuration {
	case "dev":
		return "./", nil
	case "prod":
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config dir")
		}

		appDir := filepath.Join(userConfigDir, AppName)
		err = directory.CheckAndCreateDir(appDir)
		if err != nil {
			return "", fmt.Errorf("failed to check app dir: %s", err.Error())
		}

		return appDir, nil
	default:
		return "", fmt.Errorf("unknown configuration %s", configuration)
	}
}
