package app

import (
	"database/sql"

	"github.com/NikitaAksenov/passman/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

var DBPath = "D:/Programs/passman/passwords.db"

type application struct {
	passwords *models.PasswordModel
}

func NewApplication() (*application, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return nil, err
	}

	app := application{passwords: &models.PasswordModel{DB: db}}

	return &app, nil
}

func (app *application) Close() {
	app.passwords.DB.Close()
}

func (app *application) Add(title, value string) error {
	return app.passwords.Add(title, value)
}

func (app *application) Get(title string) (string, error) {
	return app.passwords.Get(title)
}

func (app *application) List() ([]*models.Password, error) {
	return app.passwords.List()
}
