package app

import (
	"database/sql"

	"github.com/NikitaAksenov/passman/internal/models"
)

type application struct {
	passwords *models.PasswordModel
}

func NewApplication(db *sql.DB) *application {
	app := application{passwords: &models.PasswordModel{DB: db}}

	return &app
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
