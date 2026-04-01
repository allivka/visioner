package sci

//SCI = Serial Communication Interface

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"

	"go.bug.st/serial"
)

type SCIConfig struct {
	PerWriteTimeout time.Duration
	ReadTimeout     time.Duration
	Port            serial.Port
	In              <-chan Behavior
	QueueSize       int
}

func RunSCI(ctx context.Context, c SCIConfig) chan float32 {

	go func(ctx context.Context) {
		cancels := make([]func(), 1)
		for {
			select {
			case value := <-c.In:
				ct, cancel := context.WithTimeout(ctx, c.PerWriteTimeout)
				cancels = append(cancels, cancel)
				go func(ctx context.Context) {
					buff := value.Serialize()

					_, err := c.Port.Write(buff)

					if err != nil {
						slog.Warn(fmt.Sprintf("Failed sending data over serial port: %v", err))
					}
				}(ct)

			case <-ctx.Done():
				for _, cancel := range cancels {
					cancel()
				}
				slog.Info("Exiting serial writer goroutine")
				return
			}
		}
	}(ctx)

	serialReceiver := make(chan float32, c.QueueSize)

	go func(ctx context.Context) {

		buffer := make([]byte, 4)

		var (
			err   error
			value float64
		)

		err = c.Port.SetReadTimeout(c.ReadTimeout)

		if err != nil {
			log.Fatalf("Could not set read timeout for serial port: %v", err)
			return
		}

		for {
			select {

			case <-ctx.Done():
				slog.Info("Exiting serial receiver goroutine")
				return

			default:				
				buffer, _, err = bufio.NewReader(c.Port).ReadLine()

				if err != nil {
					slog.Warn(fmt.Sprintf("Failed receiving data over serial port: %v", err))
					continue
				}
				
				value, err = strconv.ParseFloat(string(buffer), 32)
				
				if err != nil {
					slog.Warn(err.Error())
					continue
				}
				
				// go fmt.Println(value)
				
				serialReceiver <- float32(value)

				
			}
		}
	}(ctx)

	return serialReceiver
}
