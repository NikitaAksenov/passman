package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	passman_app "github.com/NikitaAksenov/passman/internal/app"
	"github.com/NikitaAksenov/passman/internal/encrypt"
)

type Entry struct {
	target string
}

type desktopApp struct {
	PassmanApp *passman_app.App
	FyneApp    fyne.App

	key     string
	entries []Entry

	mainWindow     fyne.Window
	keyEnterWindow fyne.Window
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

func (da *desktopApp) InitWindows() {
	da.mainWindow = da.NewMainWindow()
	//da.keyEnterWindow = da.NewKeyEnteredWindow()
}

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
			copyButton := buttons.Objects[0].(*widget.Button)
			deleteButton := buttons.Objects[1].(*widget.Button)
			label.SetText(da.entries[lii].target)
			copyButton.OnTapped = func() {
				target := da.entries[lii].target
				pass, _ := da.PassmanApp.GetPassword(target, da.key)
				da.PassmanApp.SendToClipboard(pass)
			}
			deleteButton.OnTapped = func() {
				target := label.Text
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

func (da *desktopApp) NewAddEntryWindow() fyne.Window {
	addEntryWindow := da.FyneApp.NewWindow("New entry")

	addEntryWindow.Resize(fyne.NewSize(300.0, 150.0))

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

	addEntryWindow.SetContent(container.NewBorder(
		nil, addEntryWindow_EnterButton, nil, nil, container.NewVBox(addEntryWindow_Entry, addEntryWindow_FirstPass, addEntryWindow_SecondPass),
	))

	return addEntryWindow
}

func (da *desktopApp) NewKeyEnteredWindow() fyne.Window {
	keyEnterWindow := da.FyneApp.NewWindow("Enter key")

	keyEnterWindow.Resize(fyne.NewSize(300.0, 50.0))

	keyEnterWindow_Entry := widget.NewEntry()
	keyEnterWindow_EnterButton := widget.NewButton("Enter key", func() {
		da.key = keyEnterWindow_Entry.Text

		da.mainWindow.Show()
		keyEnterWindow.Close()
	})

	keyEnterWindow_Entry.Password = true
	keyEnterWindow_Entry.OnChanged = func(s string) {
		keyEnterWindow_EnterButton.Disable()

		if len(s) > 2 {
			keyEnterWindow_EnterButton.Enable()
		}
	}

	keyEnterWindow_EnterButton.Disable()

	keyEnterWindow.SetContent(container.NewBorder(
		nil, nil, nil, keyEnterWindow_EnterButton, keyEnterWindow_Entry,
	))

	return keyEnterWindow
}

func (da *desktopApp) Run() {
	da.mainWindow.Show()
	da.FyneApp.Run()
}

func NewEntryLabel(text string) *fyne.Container {
	entry := widget.NewLabel(text)
	copyButton := widget.NewButton("Copy", func() {})
	deleteButton := widget.NewButton("Delete", func() {})

	return container.NewBorder(
		nil, nil, nil, container.NewHBox(copyButton, deleteButton), entry,
	)
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
