package main

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"log/slog"
	"math"
	"net/http"
	"time"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Visioner struct {
	url string
}

func(v Visioner) getAngle() (float32, error) {
	resp, err := http.Get(v.url)
	
	if err != nil {
		slog.Warn(err.Error())
		return 0, err
	}
	
	buffer, err := io.ReadAll(resp.Body)
	
	if err != nil {
		slog.Warn(err.Error())
		return 0, err
	}
	
	return math.Float32frombits(binary.LittleEndian.Uint32(buffer)), nil
}

func main() {
	
	var visioner Visioner
	
	application := app.New()
	application.SetIcon(canvas.NewImageFromFile("Icon.png").Resource)
	
	loginWindow := application.NewWindow("Login")
	
	label := widget.NewLabel("Enter visioner device URL please:")
	input := widget.NewEntry()
	// input.Validator = func(s string) error {
		
	// 	visioner.url = s
	// 	_, err := visioner.getAngle()
		
	// 	return err
	// }
	
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
		
		_, err := visioner.getAngle()
		
		if err != nil {
			dialog.ShowInformation("Error, Invalid visioner URL", err.Error(), loginWindow)
			return
		}
		
		loginWindow.Close()
		
		window.Show()
	}))
	
	loginWindow.SetContent(content)
	
	loginWindow.Show()
	
	go func(){
		ticker := time.NewTicker(time.Second / 10)
		
		for range ticker.C {
			fyne.Do(func() {
				// angleText.Text = "Angle: " + fmt.Sprint(visioner.getAngle())
				angleText.Refresh()
			})
		}
	}()
	
	application.Run()
}
