package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log/slog"
	"math"
	"net"
	"net/http"
	"time"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/veandco/go-sdl2/sdl"
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
	
	defer resp.Body.Close()
	
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
	
	streamContainer := canvas.NewImageFromImage(nil)
	streamContainer.FillMode = canvas.ImageFillContain
	
	streamWindow := application.NewWindow("camera stream")
	streamWindow.Resize(fyne.NewSize(640, 480))
	streamWindow.SetContent(streamContainer)
	streamWindow.SetMaster()
	
	window := application.NewWindow("Visioner controller")
	window.SetContent(container.NewStack(
		canvas.NewRectangle(color.White),
	
		container.NewVBox(
			container.NewCenter(angleText),
			robotContainer,
			
		),
	))
	window.Resize(fyne.NewSize(500, 500))
	window.SetMaster()
	
	content := container.NewVBox(label, input, widget.NewButton("Submit", func() {
		
		visioner.address = input.Text
		_, err := visioner.getAngle()
		
		if err != nil && input.Text != "pass" {
			dialog.ShowInformation("Error, Invalid visioner device address", err.Error(), loginWindow)
			return
		}
		
		loginWindow.Close()
		
		slog.Info("Starting main app")
		
		go func() {
			
			slog.Info("Starting visualization ")
			
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
			slog.Info("Starting angle data receiver")
			
			ticker := time.NewTicker(time.Second / 10)
			var (
				err error
				angle float32
			)
			
			_ = err
			
			for range ticker.C {
				fyne.Do(func() {
					angle, err = visioner.getAngle()
					
					if err != nil {
						slog.Warn(err.Error())
						return
					}
					
					// angle++
					
					angleChan <- angle
					
					angleText.Text = "Angle: " + fmt.Sprint(angle) + "°"
					angleText.Refresh()
				})
			}
		}()
		
		go func () {
			
			slog.Info("Starting video stream receiver")
			
			connection, err := net.Dial("tcp", visioner.address + ":8080")
			
			for err != nil {
				slog.Warn(err.Error())
				time.Sleep(1 * time.Second)
				connection, err = net.Dial("tcp", visioner.address + ":8080")
			}
			
			defer connection.Close()
			
			slog.Info("Established connection with video stream sender")
			
			reader := bufio.NewReader(connection)
			
			var (
				length uint32
				n int
				frame image.Image
			)
			
			
			for {
				err = binary.Read(reader, binary.LittleEndian, &length)
				
				if err != nil {
					slog.Warn(err.Error())
					continue
				}
				
				data := make([]byte, length)
				
				n, err = io.ReadFull(reader, data)
				
				if err != nil {
					slog.Warn(err.Error())
					continue
				}
				
				if n < int(length) {
					slog.Warn("The amount of read bytes is less then stated size of the frame")
				}
				
				frame, err = jpeg.Decode(bytes.NewReader(data))
				
				if err != nil {
					slog.Warn(err.Error())
					continue
				}
				
				fyne.Do(func ()  {
					streamContainer.Image = frame
					streamContainer.Refresh()
				})
				
				
			}
		}()
		
		go func() {
			if err := sdl.Init(sdl.INIT_GAMECONTROLLER); err != nil {
				log.Fatal(err)
			}
			defer sdl.Quit()
			
			var controller *sdl.GameController
			
			outer:
			for controller == nil {
				slog.Info("Attempting to find a controller")
				
				for i := 0; i < sdl.NumJoysticks(); i++ {
					if sdl.IsGameController(i) {
						controller = sdl.GameControllerOpen(i)
						if controller != nil {
							slog.Info("Found controller:" + controller.Name())
							break outer
						}
					}
				}
			}
			
			defer controller.Close()
			
			for event := sdl.PollEvent(); true; event = sdl.PollEvent() {
				if event == nil {
					sdl.Delay(10)
					continue
				}
				
				switch e := event.(type) {
				case *sdl.ControllerButtonEvent:
					switch e.State {
					case sdl.PRESSED:
						fmt.Print("Button pressed: ")
					case sdl.RELEASED:
						fmt.Print("Button release: ")
					}
					
					fmt.Println(e.Button)
				case *sdl.ControllerAxisEvent:
					v := float64(e.Value) /32768.0
					
					if math.Abs(v) < 0.3 {
						continue
					}
					
					fmt.Println(e.Which, e.Axis, v)
				
				}
				
				sdl.Delay(10)
			}
			
		}()
		
		window.Show()
		streamWindow.Show()
	}))
	
	loginWindow.SetContent(content)
	
	loginWindow.Show()
	
	application.Run()
}
