package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	passman_app "github.com/NikitaAksenov/passman/internal/app"
)

type Entry struct {
	target string
}

type desktopApp struct {
	PassmanApp *passman_app.App
	FyneApp    fyne.App

	key     string
	entries []Entry

	mainWindow fyne.Window
}

func NewEntry(s string) Entry {
	return Entry{target: s}
}

func New() (*desktopApp, error) {
	pApp, err := passman_app.New()
	if err != nil {
		return nil, err
	}

	fApp := app.New()

	return &desktopApp{
		PassmanApp: pApp,
		FyneApp:    fApp,
	}, nil
}

func (da *desktopApp) InitWindows() {
	da.mainWindow = da.NewMainWindow()
}

func (da *desktopApp) Run() {
	da.mainWindow.Show()
	da.FyneApp.Run()
}

func main() {
	da, err := New()
	if err != nil {
		fmt.Println(err)
		return
	}

	da.InitWindows()
	da.Run()
}
