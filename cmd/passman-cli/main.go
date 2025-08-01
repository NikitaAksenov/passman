package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

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

var rootCmd = &cobra.Command{
	Use:   "passman",
	Short: "Simple CLI password manager",
	Long:  `Passman is a simple CLI tool for password management.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to passman! Use --help for usage.")
	},
}

func main() {
	// - Check app directories and create them if needed
	// -- Get user config dir
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("failed to get user config dir")
	}
	// fmt.Println("userConfigDir:", userConfigDir)

	// -- Check and create app dir
	appDir := filepath.Join(userConfigDir, appName)
	// fmt.Println("appDir:", appDir)
	err = CheckAndCreateDir(appDir)
	if err != nil {
		log.Fatal("failed to check app dir:", err.Error())
	}

	// -- Check and create storage dir
	storageDir := filepath.Join(appDir, storageDir)
	// fmt.Println("storageDir:", storageDir)
	err = CheckAndCreateDir(storageDir)
	if err != nil {
		log.Fatal("failed to check storage dir:", err.Error())
	}
	storagePath := filepath.Join(storageDir, storageName)
	// fmt.Println("storagePath:", storagePath)

	// - Init storage
	var storage storage.Storage
	storage, err = sqlite.New(storagePath)
	if err != nil {
		log.Fatalf("failed to init storage: %s", err)
	}

	App := App{
		Storage: storage,
	}

	rootCmd.AddCommand(AddCommand(&App))
	rootCmd.AddCommand(GetCommand(&App))
	rootCmd.AddCommand(ListCommand(&App))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CheckAndCreateDir(path string) error {
	// Check if path is empty
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	file, err := os.Stat(path)
	if os.IsNotExist(err) {
		// If dir does not exits - create one
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed creating dir: %s", err.Error())
		}

		return nil
	} else if err != nil {
		return fmt.Errorf("failed getting path info: %s", err.Error())
	}

	if !file.IsDir() {
		return fmt.Errorf("path is not dir")
	}

	return nil
}
