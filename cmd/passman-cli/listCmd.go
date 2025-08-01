package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
