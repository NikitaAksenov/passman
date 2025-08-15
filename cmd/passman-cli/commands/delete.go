package commands

import (
	"fmt"

	"github.com/NikitaAksenov/passman/cmd/passman-cli/app"
	"github.com/spf13/cobra"
)

func DeleteCommand(app *app.CliApp) *cobra.Command {
	command := cobra.Command{
		Use:     "delete target",
		Short:   "Deletes entry for provided target",
		Aliases: []string{"del"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]

			rowsAffected, err := app.Storage.DeleteTarget(target)
			if err != nil {
				fmt.Printf("failed to delete target %s: %s", target, err.Error())
				return
			}

			if rowsAffected == 0 {
				fmt.Printf("Failed to find target [%s]\n", target)
			} else {
				fmt.Printf("Successfully deleted target [%s]\n", target)
			}
		},
	}

	return &command
}
