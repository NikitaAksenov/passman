package commands

import (
	"fmt"
	"syscall"

	"github.com/NikitaAksenov/passman/cmd/passman-cli/app"
	"github.com/NikitaAksenov/passman/internal/encrypt"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
	"golang.org/x/term"
)

func GetCommand(app *app.CliApp) *cobra.Command {
	command := cobra.Command{
		Use:   "get target",
		Short: "Sends target's password to the clipboard",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]

			showFlag, err := cmd.Flags().GetBool("show")
			if err != nil {
				fmt.Println("failed to get \"show\" flag value")
				return
			}

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

			// Send password to clipboard
			clipboard.Write(clipboard.FmtText, []byte(pass))
			fmt.Println("Password copied to clipboard")

			if showFlag {
				fmt.Println("Password:", pass)
			}
		},
	}

	command.Flags().BoolP("show", "s", false, "if set then password will be shown in the terminal")

	return &command
}
