package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (da *desktopApp) NewMainWindow() fyne.Window {
	mainWindow := da.FyneApp.NewWindow("passman")

	da.UpdateEntries()

	mainWindow.SetMaster()
	mainWindow.Resize(fyne.NewSize(400.0, 600.0))
	mainWindow.SetFixedSize(true)

	mainWindow_EntryList := widget.NewList(
		func() int {
			return len(da.entries)
		},
		func() fyne.CanvasObject {
			return NewEntryLabel("none")
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			entryContainer := co.(*fyne.Container)

			label := entryContainer.Objects[0].(*widget.Label)

			buttons := entryContainer.Objects[1].(*fyne.Container)
			showButton := buttons.Objects[0].(*widget.Button)
			copyButton := buttons.Objects[1].(*widget.Button)
			deleteButton := buttons.Objects[2].(*widget.Button)

			// for i, button := range buttons.Objects {
			// 	fmt.Println(i, button.(*widget.Button).Text)
			// }

			target := da.entries[lii].target

			label.SetText(target)
			showButton.SetText("Show")

			showButton.OnTapped = func() {
				switch showButton.Text {
				case "Show":
					showButton.SetText("Hide")

					pass, _ := da.PassmanApp.GetPassword(target, da.key)
					label.SetText(pass)
				case "Hide":
					showButton.SetText("Show")

					label.SetText(da.entries[lii].target)
				}
			}

			copyButton.OnTapped = func() {
				pass, _ := da.PassmanApp.GetPassword(target, da.key)
				da.PassmanApp.SendToClipboard(pass)
			}

			deleteButton.OnTapped = func() {
				_, err := da.PassmanApp.Storage.DeleteTarget(target)
				if err != nil {
					fmt.Println(err)
					return
				}

				da.UpdateEntries()

				mainWindow.Canvas().Content().Refresh()
			}
		},
	)

	mainWindow_AddEntryButton := widget.NewButton("Add entry", func() {
		addEntryWindow := da.NewAddEntryWindow()

		mainWindow.Hide()

		addEntryWindow.Show()
	})

	mainWindow_EnterKeyContainer := da.NewEnterKeyContainer()

	mainWindow.SetContent(container.NewBorder(
		mainWindow_EnterKeyContainer, mainWindow_AddEntryButton, nil, nil, mainWindow_EntryList,
	))

	return mainWindow
}

func (da *desktopApp) NewEnterKeyContainer() *fyne.Container {
	enterKeyEntry := widget.NewEntry()
	enterKeyEntry.Password = true
	enterKeyEntry.PlaceHolder = "Enter key..."
	enterKeyEntry.OnChanged = func(s string) {
		da.key = enterKeyEntry.Text
	}

	entryKeyContainer := container.NewBorder(
		nil, nil,
		widget.NewLabel("Key:"),
		nil,
		enterKeyEntry,
	)

	return entryKeyContainer
}

func (da *desktopApp) UpdateEntries() error {
	da.entries = nil

	targets, err := da.PassmanApp.Storage.GetTargets(20, 0)
	if err != nil {
		return fmt.Errorf("failed getting targets from storage: %s", err.Error())
	}

	for _, v := range targets {
		da.entries = append(da.entries, NewEntry(v))
	}

	return nil
}

func NewEntryLabel(text string) *fyne.Container {
	entry := widget.NewLabel(text)
	showButton := widget.NewButton("Show", func() {})
	copyButton := widget.NewButton("Copy", func() {})
	deleteButton := widget.NewButton("Delete", func() {})

	return container.NewBorder(
		nil, nil, nil, container.NewHBox(showButton, copyButton, deleteButton), entry,
	)
}
