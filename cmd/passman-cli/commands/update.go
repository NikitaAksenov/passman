package commands

import (
	"fmt"
	"syscall"

	"github.com/NikitaAksenov/passgen/pkg/passgen"
	"github.com/NikitaAksenov/passman/cmd/passman-cli/app"
	"github.com/NikitaAksenov/passman/internal/encrypt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func UpdateCommand(app *app.App) *cobra.Command {
	command := cobra.Command{
		Use:     "update target",
		Aliases: []string{"upd"},
		Short:   "Update target's password",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]

			norepeat, err := cmd.Flags().GetBool("norepeat")
			if err != nil {
				fmt.Println("failed to get \"norepeat\" value")
				return
			}

			generate, err := cmd.Flags().GetBool("generate")
			if err != nil {
				fmt.Println("failed to get \"generate\" value")
				return
			}

			var password string
			if generate {
				// Generate password
				password, err = passgen.Generate(10)
				if err != nil {
					fmt.Printf("failed to generate password: %s", err.Error())
					return
				}
			} else {
				// Read password silently
				fmt.Print("Enter password: ")
				bytePass, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					fmt.Printf("failed to read password: %s", err.Error())
					return
				}
				fmt.Println()
				password = string(bytePass)
			}

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

				// If norepeat flag is set then don't prompt to enter key again
				if norepeat {
					break
				}

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
			encryptedPass, err := encrypt.EncryptString(resizedKey, password)
			if err != nil {
				fmt.Printf("encryption failed: %s", err.Error())
				return
			}

			// Update password in storage
			rowsAffected, err := app.Storage.UpdatePassword(target, encryptedPass)
			if err != nil {
				fmt.Printf("failed to update password: %s", err.Error())
				return
			}

			// Check if target was updated
			if rowsAffected == 0 {
				fmt.Printf("Failed to find target [%s]\n", target)
			} else {
				fmt.Printf("Successfully updated target [%s]\n", target)
			}
		},
	}

	command.Flags().BoolP("norepeat", "n", false, "if set then user won't be prompt to enter key twice")
	command.Flags().BoolP("generate", "g", false, "if set then password will be generated automatically")

	return &command
}
