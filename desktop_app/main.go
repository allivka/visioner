package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/app"
)

func getAngle() float64 {
	
	return 0
}

func main() {
	
	
	application := app.New()
	application.SetIcon(canvas.NewImageFromFile("Icon.png").Resource)
	
	loginWindow := application.NewWindow("Login")
	
	label := widget.NewLabel("Enter visioner device IP address please:")
	input := widget.NewEntry()
	
	robotImage := canvas.NewImageFromFile("robot.png")
	robotImage.FillMode = canvas.ImageFillOriginal
	
	angleText := canvas.NewText("Angle: 0", color.Black)
	angleText.TextSize = 32
	
	window := application.NewWindow("Visioner controller")
	
	window.SetContent(container.NewStack(
		canvas.NewRectangle(color.White),
		container.NewCenter(
			container.NewVBox(
				container.NewCenter(angleText),
				robotImage,
			),
		),
	))
	
	content := container.NewVBox(label, input, widget.NewButton("Submit", func() {
		fmt.Println(input.Text)
		loginWindow.Close()
		
		window.Show()
	}))
	
	loginWindow.SetContent(content)
	
	loginWindow.Show()
	
	go func(){
		ticker := time.NewTicker(time.Second / 10)
		
		for range ticker.C {
			fyne.Do(func() {
				angleText.Text = "Angle: " + fmt.Sprint(getAngle())
				angleText.Refresh()
			})
		}
	}()
	
	application.Run()
}
