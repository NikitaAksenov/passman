package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	passman_app "github.com/NikitaAksenov/passman/internal/app"
)

type desktopApp struct {
	PassmanApp *passman_app.App
	FyneApp    *fyne.App

	key string
}

func main() {
	var dApp desktopApp
	pApp, _ := passman_app.New()
	fApp := app.New()
	dApp.PassmanApp = pApp
	dApp.FyneApp = &fApp

	mainWindow := fApp.NewWindow("passman")
	mainWindow.Resize(fyne.NewSize(500.0, 500.0))
	mainWindow.SetMaster()

	newEntryWindow := fApp.NewWindow("Add")
	newEntryWindow.Resize(fyne.NewSize(300.0, 150.0))
	newEntryTarget := widget.NewEntry()
	newEntryTarget.PlaceHolder = "Target"
	newEntryFirstPassword := widget.NewEntry()
	newEntryFirstPassword.PlaceHolder = "Enter password..."
	newEntryFirstPassword.Password = true
	newEntrySecondPassword := widget.NewEntry()
	newEntrySecondPassword.PlaceHolder = "Re-enter password..."
	newEntrySecondPassword.Password = true
	newEntryFields := container.NewVBox(newEntryTarget, newEntryFirstPassword, newEntrySecondPassword)
	newEntryButton := widget.NewButton("Add entry", func() {
		fmt.Println("Add entry clicked")
	})
	newEntryContainer := container.NewBorder(
		nil, newEntryButton, nil, nil, newEntryFields,
	)
	newEntryWindow.SetContent(newEntryContainer)

	addEntryButton := widget.NewButton("Add", func() {
		newEntryWindow.Show()
	})

	enterKeyWindow := fApp.NewWindow("Enter key")
	enterKeyWindow.Resize(fyne.NewSize(300.0, 50.0))
	enterKeyEntry := widget.NewEntry()
	enterKeyButton := widget.NewButton("Enter", func() {
		dApp.key = enterKeyEntry.Text

		entries, _ := pApp.Storage.GetTargets(20, 0)
		entryList := widget.NewList(
			func() int {
				return len(entries)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(lii widget.ListItemID, co fyne.CanvasObject) {
				entryLabel := co.(*widget.Label)
				entryLabel.SetText(entries[lii])
			},
		)

		mainContainer := container.NewBorder(
			nil, addEntryButton, nil, nil, entryList,
		)
		mainWindow.SetContent(mainContainer)

		mainWindow.Show()
		enterKeyWindow.Close()
	})
	enterKeyButton.Disable()
	enterKeyEntry.Password = true
	enterKeyEntry.OnChanged = func(s string) {
		enterKeyButton.Disable()
		if len(s) > 2 {
			enterKeyButton.Enable()
		}
	}
	enterKeyContainer := container.NewBorder(
		nil, nil, nil, enterKeyButton, enterKeyEntry,
	)
	enterKeyWindow.SetContent(enterKeyContainer)

	enterKeyWindow.ShowAndRun()
}
