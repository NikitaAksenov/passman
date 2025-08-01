package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/NikitaAksenov/passman/internal/encrypt"
	"github.com/NikitaAksenov/passman/internal/storage/sqlite"
)

const (
	appName     = "passman"
	storageDir  = "storage"
	storageName = "storage.db"
)

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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatalf("failed to init storage: %s", err)
	}

	// - Handle user input
	command := os.Args[1]
	switch command {
	case "add":
		reader := bufio.NewReader(os.Stdin)

		// -- Read target
		fmt.Print("Enter target: ")
		target, err := reader.ReadString('\n')
		if len(target) == 0 {
			log.Fatal("target must not be empty")
		}
		if err != nil {
			log.Fatalf("error during reading target: %s", err.Error())
		}
		target = strings.TrimSpace(target)

		// -- Read password silently
		fmt.Print("Enter password: ")
		bytePass, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("error during reading password: %s", err.Error())
		}
		fmt.Println()

		// -- Read key silently
		var key, keyRepeat []byte
		for keysEqual := false; !keysEqual; {
			fmt.Print("Enter key: ")
			key, err = term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatalf("error during reading key: %s", err.Error())
			}
			fmt.Println()

			fmt.Print("Enter key again: ")
			keyRepeat, err = term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatalf("error during reading key again: %s", err.Error())
			}
			fmt.Println()

			keysEqual = string(key) == string(keyRepeat)

			if !keysEqual {
				fmt.Println("First and second key are different, but must be equal. Please, try again...")
			}
		}

		// fmt.Printf("[%s] [%s] [%s]\n", target, string(bytePass), string(key))

		// -- Resize key to 16 bytes
		resizedKey := encrypt.ResizeKey(key)

		// -- Encrypt password
		encryptedPass, err := encrypt.EncryptString(resizedKey, string(bytePass))
		if err != nil {
			log.Fatalf("encryption failed: %s", err.Error())
		}

		// -- Add encrypted password to storage
		_, err = storage.AddPass(target, encryptedPass)
		if err != nil {
			log.Fatalf("adding to storage failed: %s", err.Error())
		}
	case "get":
		reader := bufio.NewReader(os.Stdin)

		// -- Read target
		fmt.Print("Enter target: ")
		target, err := reader.ReadString('\n')
		if len(target) == 0 {
			log.Fatal("target must not be empty")
		}
		if err != nil {
			log.Fatalf("error during reading target: %s", err.Error())
		}
		target = strings.TrimSpace(target)

		// -- Read key silently
		fmt.Print("Enter key: ")
		key, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("error during reading key: %s", err.Error())
		}
		fmt.Println()

		// fmt.Printf("[%s] [%s]\n", target, string(key))

		// -- Resize key to 16 bytes
		resizedKey := encrypt.ResizeKey([]byte(key))

		// -- Get encrypted password from storage
		encryptedPass, err := storage.GetPass(target)
		if err != nil {
			log.Fatalf("getting from storage failed: %s", err.Error())
		}

		// -- Decrypt password
		pass, err := encrypt.DecryptString(resizedKey, encryptedPass)
		if err != nil {
			log.Fatalf("decryption failed: %s", err.Error())
		}

		fmt.Println("Password:", pass)
	case "list":
		var limit, offset int

		limit, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("failed reading limit: %s", err.Error())
		}
		offset, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("failed reading offset: %s", err.Error())
		}

		targets, err := storage.GetTargets(limit, offset)
		if err != nil {
			log.Fatalf("failed getting targets from storage: %s", err.Error())
		}

		fmt.Println(targets)
	default:
		fmt.Println("Unknonw command", command)
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
