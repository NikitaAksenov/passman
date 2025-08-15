package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	FyneApp *fyne.App
}

func main() {
	a := app.New()

	mainWindow := a.NewWindow("passman")
	mainWindow.Resize(fyne.NewSize(500.0, 500.0))

	mainWindow.ShowAndRun()
}
