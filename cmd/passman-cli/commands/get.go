package commands

import (
	"fmt"
	"syscall"

	"github.com/NikitaAksenov/passman/cmd/passman-cli/app"
	"github.com/NikitaAksenov/passman/internal/encrypt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func GetCommand(app *app.App) *cobra.Command {
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
