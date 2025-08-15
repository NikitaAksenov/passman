package commands

import (
	"fmt"

	"github.com/NikitaAksenov/passman/cmd/passman-cli/app"
	"github.com/spf13/cobra"
)

func InfoCommand(app *app.CliApp) *cobra.Command {
	command := cobra.Command{
		Use:   "info target",
		Short: "Returns info about entry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]

			// Get entry info from storage
			targetInfo, err := app.Storage.GetTargetInfo(target)
			if err != nil {
				fmt.Printf("getting from storage failed: %s", err.Error())
				return
			}

			fmt.Printf("Target info\n%s\n", targetInfo)
		},
	}

	return &command
}
