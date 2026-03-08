package sci

//SCI = Serial Communication Interface

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"log/slog"
	"math"
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
			value float32
			counter int
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
				n, err := c.Port.Read(buffer[counter:])

				if err != nil {
					slog.Warn(fmt.Sprintf("Failed receiving data over serial port: %v", err))
					continue
				}
				counter += n
				
				if counter == 8 {
					counter = 0
					value = math.Float32frombits(binary.LittleEndian.Uint32(buffer[:8]))
					
					serialReceiver <- value
				}

				
			}
		}
	}(ctx)

	return serialReceiver
}
