package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/NikitaAksenov/passgen/pkg/passgen"
	"github.com/NikitaAksenov/passman/internal/encrypt"
)

func (da *desktopApp) NewAddEntryWindow() fyne.Window {
	addEntryWindow := da.FyneApp.NewWindow("New entry")

	addEntryWindow.Resize(fyne.NewSize(300.0, 150.0))
	addEntryWindow.SetFixedSize(true)

	addEntryWindow_Entry := widget.NewEntry()
	addEntryWindow_Entry.PlaceHolder = "Target"

	addEntryWindow_FirstPass := widget.NewEntry()
	addEntryWindow_FirstPass.Password = true
	addEntryWindow_FirstPass.PlaceHolder = "Enter password..."

	addEntryWindow_SecondPass := widget.NewEntry()
	addEntryWindow_SecondPass.Password = true
	addEntryWindow_SecondPass.PlaceHolder = "Re-enter password..."

	addEntryWindow_EnterButton := widget.NewButton("Enter target", func() {
		if addEntryWindow_FirstPass.Text != addEntryWindow_SecondPass.Text {
			fmt.Println("Pass's must match")
			return
		}

		target := addEntryWindow_Entry.Text
		password := addEntryWindow_FirstPass.Text

		// Resize key to 16 bytes
		resizedKey := encrypt.ResizeKey([]byte(da.key))

		// Encrypt password
		encryptedPass, err := encrypt.EncryptString(resizedKey, password)
		if err != nil {
			fmt.Printf("encryption failed: %s", err.Error())
			return
		}

		// Add encrypted password to storage
		_, err = da.PassmanApp.Storage.AddPass(target, encryptedPass)
		if err != nil {
			fmt.Printf("adding to storage failed: %s", err.Error())
			return
		}

		da.UpdateEntries()

		addEntryWindow.Close()

		da.mainWindow.RequestFocus()
		da.mainWindow.Canvas().Content().Refresh()
	})
	addEntryWindow_EnterButton.Disable()

	enterButtonValidate := func(s string) {
		addEntryWindow_EnterButton.Enable()
		addEntryWindow_EnterButton.SetText("Enter")

		if len(addEntryWindow_Entry.Text) < 3 {
			addEntryWindow_EnterButton.Disable()
			addEntryWindow_EnterButton.SetText("Target too short")
			return
		}

		if len(addEntryWindow_FirstPass.Text) < 3 {
			addEntryWindow_EnterButton.Disable()
			addEntryWindow_EnterButton.SetText("Password too short")
			return
		}

		if addEntryWindow_FirstPass.Text != addEntryWindow_SecondPass.Text {
			addEntryWindow_EnterButton.Disable()
			addEntryWindow_EnterButton.SetText("Passwords don't match")
			return
		}
	}
	addEntryWindow_Entry.OnChanged = enterButtonValidate
	addEntryWindow_FirstPass.OnChanged = enterButtonValidate
	addEntryWindow_SecondPass.OnChanged = enterButtonValidate

	addEntryWindow_GenerateButton := widget.NewButton("Generate", func() {
		// Generate password
		generatedPassword, _ := passgen.Generate(10)

		addEntryWindow_FirstPass.SetText(generatedPassword)
		addEntryWindow_SecondPass.SetText(generatedPassword)
	})

	passwordsContainer := container.NewBorder(
		nil, nil, nil, addEntryWindow_GenerateButton, container.NewVBox(addEntryWindow_FirstPass, addEntryWindow_SecondPass),
	)

	addEntryWindow.SetContent(container.NewBorder(
		nil, addEntryWindow_EnterButton, nil, nil, container.NewVBox(addEntryWindow_Entry, passwordsContainer),
	))

	return addEntryWindow
}
