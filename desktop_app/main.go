package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/app"
)

func main() {
	
	
	application := app.New()
	application.SetIcon(canvas.NewImageFromFile("Icon.png").Resource)
	
	loginWindow := application.NewWindow("Login")
	
	label := widget.NewLabel("Enter visioner device IP address please:")
	input := widget.NewEntry()
	
	robotImage := canvas.NewImageFromFile("robot.png")
	robotImage.SetMinSize(fyne.NewSize(500, 500))
	
	content := container.NewVBox(label, input, widget.NewButton("Submit", func() {
		fmt.Println(input.Text)
		loginWindow.Close()
		
		window := application.NewWindow("Visioner controller")
		window.SetContent(container.New(layout.NewCenterLayout(),
			robotImage,
		))
		window.Show()
	}))

	loginWindow.SetContent(content)
	
	loginWindow.Show()
	
	application.Run()
}
