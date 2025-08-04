package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/NikitaAksenov/passman/cmd/passman-cli/app"
	"github.com/NikitaAksenov/passman/cmd/passman-cli/commands"
)

var rootCmd = &cobra.Command{
	Use:   "passman",
	Short: "Simple CLI password manager",
	Long:  `Passman is a simple CLI tool for password management.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to passman! Use --help for usage.")
	},
}

func main() {
	App, err := app.New()
	if err != nil {
		log.Fatalf("failed init app: %s", err.Error())
	}

	rootCmd.AddCommand(commands.AddCommand(App))
	rootCmd.AddCommand(commands.GetCommand(App))
	rootCmd.AddCommand(commands.ListCommand(App))
	rootCmd.AddCommand(commands.DeleteCommand(App))
	rootCmd.AddCommand(commands.UpdateCommand(App))
	rootCmd.AddCommand(commands.InfoCommand(App))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
