package app

import (
	"github.com/NikitaAksenov/passman/internal/app"
)

type CliApp struct {
	*app.App
}

func New() (*CliApp, error) {
	app, err := app.New()
	if err != nil {
		return nil, err
	}

	return &CliApp{
		app,
	}, nil
}
