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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Visioner struct {
	address string
}

func(v Visioner) getAngle() (float32, error) {
	resp, err := http.Get("http://" + v.address)
	
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

func rotateLine(origin *fyne.Position, final *fyne.Position, angle float32) {
	vec := fyne.NewPos(final.X - origin.X, final.Y - origin.Y)
	
	sin, cos := math.Sincos(float64(angle / 180 * math.Pi))
	
	*final = fyne.NewPos(float32(cos * float64(vec.X) - sin * float64(vec.Y)) + origin.X, float32(sin * float64(vec.X) + cos * float64(vec.Y)) + origin.Y)
}

type Arrow struct {
	l1, l2, l3, l4, l5, l6 *canvas.Line
	
	pos1, pos2, pos3, pos4, pos5 fyne.Position
	
	storage *fyne.Container
}

func(a *Arrow) Construct() *Arrow {
	
	makeOne := func(l **canvas.Line) {
		*l = canvas.NewLine(color.Black)
		(*l).StrokeWidth = 5
		(*l).StrokeColor = color.Black
	}

	makeOne(&a.l1)
	makeOne(&a.l2)
	makeOne(&a.l3)
	makeOne(&a.l4)
	makeOne(&a.l5)
	makeOne(&a.l6)
	
	a.storage = container.NewStack(
		container.NewWithoutLayout(a.l1),
		container.NewWithoutLayout(a.l2),
		container.NewWithoutLayout(a.l3),
		container.NewWithoutLayout(a.l4),
		container.NewWithoutLayout(a.l5),
		container.NewWithoutLayout(a.l6),
	)
	
	return a
}

func (a Arrow) render(pos fyne.Position, origin fyne.Position, angle float32, scale float32) {
	a.pos1 = fyne.NewPos(pos.X - 20 * scale, pos.Y)
	a.pos2 = fyne.NewPos(pos.X + 20 * scale, pos.Y)
	a.pos3 = fyne.NewPos(pos.X - 20 * scale, pos.Y - 100 * scale)
	a.pos4 = fyne.NewPos(pos.X + 20 * scale, pos.Y - 100 * scale)
	a.pos5 = fyne.NewPos(pos.X, pos.Y - 125 * scale)
	
	rotateLine(&origin, &a.pos1, angle)
	rotateLine(&origin, &a.pos2, angle)
	rotateLine(&origin, &a.pos3, angle)
	rotateLine(&origin, &a.pos4, angle)
	rotateLine(&origin, &a.pos5, angle)
	
	a.l1.Position1 = a.pos1
	a.l1.Position2 = a.pos2
	
	a.l2.Position1 = a.pos1
	a.l2.Position2 = a.pos3
	
	a.l3.Position1 = a.pos2
	a.l3.Position2 = a.pos4
	
	a.l4.Position1 = a.pos3
	a.l4.Position2 = a.pos4
	
	a.l5.Position1 = a.pos3
	a.l5.Position2 = a.pos5
	
	a.l6.Position1 = a.pos4
	a.l6.Position2 = a.pos5
		
	a.l1.Refresh()
	a.l2.Refresh()
	a.l3.Refresh()
	a.l4.Refresh()
	a.l5.Refresh()
	a.l6.Refresh()
	
}

func main() {
	
	var (
		visioner Visioner
		angleChan chan float32 = make(chan float32)
	)
	
	application := app.New()
	application.SetIcon(canvas.NewImageFromFile("robot.png").Resource)
	
	loginWindow := application.NewWindow("Login")
	
	label := widget.NewLabel("Enter visioner device address please:")
	input := widget.NewEntry()
	
	circle :=  canvas.NewCircle(color.Black)
	circle.StrokeWidth = 5
	circle.Resize(fyne.NewSize(400, 400))
	circle.StrokeColor = color.Black
	circle.FillColor = color.White
	
	arrow := (&Arrow{}).Construct()
	
	robotContainer := container.NewStack(
		container.NewWithoutLayout(circle),
		arrow.storage,
	)
	
	angleText := canvas.NewText("Angle: 0", color.Black)
	angleText.TextSize = 32
	
	window := application.NewWindow("Visioner controller")
	
	window.SetContent(container.NewStack(
		canvas.NewRectangle(color.White),
		container.NewVBox(
			container.NewCenter(angleText),
			robotContainer,
		),
	))
	window.Resize(fyne.NewSize(500, 500))
	
	
	content := container.NewVBox(label, input, widget.NewButton("Submit", func() {
		
		visioner.address = input.Text
		_, err := visioner.getAngle()
		
		if err != nil && input.Text != "pass" {
			dialog.ShowInformation("Error, Invalid visioner device address", err.Error(), loginWindow)
			return
		}
		
		loginWindow.Close()
		
		go func() {
			ticker := time.NewTicker(time.Second / 60)
			s := fyne.NewSize(0, 0)
			pos := fyne.NewPos(0, 0)
			var angle float32
			
			go func(){
				for angle = range angleChan {}
			}()
			
			
			for range ticker.C {
				
				fyne.Do(func() {
					s = fyne.NewSize(robotContainer.Size().Width - 100, window.Canvas().Size().Height - label.Size().Height - 100)
					circle.Resize(fyne.NewSize(min(s.Height, s.Width), min(s.Height, s.Width)))
					pos = fyne.NewPos(s.Width / 2 - circle.Size().Width / 2 + 50, 50)
					circle.Move(pos)
					circle.Refresh()
					
					arrow.render(
						fyne.NewPos(pos.X + circle.Size().Width / 2, pos.Y + 100 * circle.Size().Height / 100 / 3),
						fyne.NewPos(pos.X + circle.Size().Width / 2, pos.Y + circle.Size().Height / 2),
						angle,
						circle.Size().Height / 100 / 3,
					)
				})
			}
		}()
		
		time.Sleep(time.Microsecond * 100)
		
		go func(){
			ticker := time.NewTicker(time.Second / 10)
			var (
				err error
				angle float32
			)
			
			_ = err
			
			for range ticker.C {
				fyne.Do(func() {
					// angle, err = visioner.getAngle()
					
					// if err != nil {
					// 	slog.Warn(err.Error())
					// 	return
					// }
					
					angle++
					
					angleChan <- angle
					
					angleText.Text = "Angle: " + fmt.Sprint(angle) + "°"
					angleText.Refresh()
				})
			}
		}()
		
		window.Show()
	}))
	
	loginWindow.SetContent(content)
	
	loginWindow.Show()
	
	application.Run()
}
