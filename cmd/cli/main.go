package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/NikitaAksenov/passman/internal/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Handle commands
	if len(os.Args) < 2 {
		fmt.Println("No command was passed")
		return
	}

	db, err := sql.Open("sqlite3", "passwords.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	app := app.NewApplication(db)

	command := os.Args[1]

	switch command {
	case "add":
		{
			if len(os.Args) < 3 {
				fmt.Println("No title was passed")
				return
			}

			title := os.Args[2]

			if len(os.Args) < 4 {
				fmt.Println("No value was passed")
				return
			}

			value := os.Args[3]

			err = app.Add(title, value)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case "get":
		{
			if len(os.Args) < 3 {
				fmt.Println("No title was passed")
				return
			}

			title := os.Args[2]

			value, err := app.Get(title)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println(value)
		}
	case "list":
		passwords, err := app.List()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, password := range passwords {
			fmt.Println(password.Title)
		}
	default:
		fmt.Println("Unknown command:", command)
		return
	}
}
