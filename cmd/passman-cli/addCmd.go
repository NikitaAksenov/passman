package main

import (
	"fmt"
	"syscall"

	"github.com/NikitaAksenov/passman/internal/encrypt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

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
