package main

import (
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Update Time")

	message := widget.NewLabel("Welcome")
	button := widget.NewButton("Update", func() {
		formatted := time.Now().Format("Time: 03:04:05")
		message.SetText(formatted)
	})

	w.SetContent(container.NewVBox(message, button))
	w.ShowAndRun()
}