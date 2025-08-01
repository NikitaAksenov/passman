package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/NikitaAksenov/passman/internal/encrypt"
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

func AddCommand(app *App) *cobra.Command {
	command := cobra.Command{
		Use:   "add [target]",
		Short: "Adds new target and it's password",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]

			// Read password silently
			fmt.Print("Enter password: ")
			bytePass, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Printf("error during reading password: %s", err.Error())
				return
			}
			fmt.Println()

			// Read key silently
			var key, keyRepeat []byte
			for keysEqual := false; !keysEqual; {
				fmt.Print("Enter key: ")
				key, err = term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					fmt.Printf("error during reading key: %s", err.Error())
					return
				}
				fmt.Println()

				fmt.Print("Enter key again: ")
				keyRepeat, err = term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					fmt.Printf("error during reading key again: %s", err.Error())
					return
				}
				fmt.Println()

				keysEqual = string(key) == string(keyRepeat)

				if !keysEqual {
					fmt.Println("First and second key are different, but must be equal. Please, try again...")
				}
			}

			// Resize key to 16 bytes
			resizedKey := encrypt.ResizeKey(key)

			// Encrypt password
			encryptedPass, err := encrypt.EncryptString(resizedKey, string(bytePass))
			if err != nil {
				fmt.Printf("encryption failed: %s", err.Error())
				return
			}

			// -- Add encrypted password to storage
			_, err = app.Storage.AddPass(target, encryptedPass)
			if err != nil {
				fmt.Printf("adding to storage failed: %s", err.Error())
				return
			}
		},
	}

	return &command
}

func GetCommand(app *App) *cobra.Command {
	command := cobra.Command{
		Use:   "get [target]",
		Short: "Returns [target]'s password",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]

			// Read key silently
			fmt.Print("Enter key: ")
			key, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Printf("error during reading key: %s", err.Error())
				return
			}
			fmt.Println()

			// Resize key to 16 bytes
			resizedKey := encrypt.ResizeKey([]byte(key))

			// Get encrypted password from storage
			encryptedPass, err := app.Storage.GetPass(target)
			if err != nil {
				fmt.Printf("getting from storage failed: %s", err.Error())
				return
			}

			// Decrypt password
			pass, err := encrypt.DecryptString(resizedKey, encryptedPass)
			if err != nil {
				fmt.Printf("decryption failed: %s", err.Error())
				return
			}

			fmt.Println("Password:", pass)
		},
	}

	return &command
}

func ListCommand(app *App) *cobra.Command {
	command := cobra.Command{
		Use:   "list",
		Short: "Lists existing targets",
		Run: func(cmd *cobra.Command, args []string) {
			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				fmt.Println("failed to get \"limit\" value")
				return
			}

			offset, err := cmd.Flags().GetInt("offset")
			if err != nil {
				fmt.Println("failed to get \"offset\" value")
				return
			}

			targets, err := app.Storage.GetTargets(limit, offset)
			if err != nil {
				fmt.Printf("failed getting targets from storage: %s", err.Error())
				return
			}

			for i, v := range targets {
				fmt.Printf("#%d %s\n", i+1, v)
			}
		},
	}

	command.Flags().IntP("limit", "l", 20, "limits amount of targets listed")
	command.Flags().IntP("offset", "o", 0, "list targets starting from passed value")

	return &command
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
