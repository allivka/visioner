package sci

//SCI = Serial Communication Interface

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"log/slog"
	"time"

	"go.bug.st/serial"
)

type SCIConfig struct {
	PerWriteTimeout time.Duration
	ReadTimeout time.Duration
	Port serial.Port
	In <-chan Behavior
	QueueSize int
}

func RunSCI(ctx context.Context, c SCIConfig) (chan float64) {
	
	go func(ctx context.Context) {
		ctx2, cancel := context.WithTimeout(ctx, c.PerWriteTimeout)
		for { select {
			case value := <-c.In:
				go func(ctx context.Context) {
					buff, err := value.Serialize()
					
					if err != nil {
						slog.Warn(fmt.Sprintf("Failed serializing behavior package: %v", err))
						return
					}
					
					_, err = c.Port.Write(buff.Bytes())
					
					if err != nil {
						slog.Warn(fmt.Sprintf("Failed sending data over serial port: %v", err))
					}
				}(ctx2)
				
			case <-ctx.Done(): cancel(); slog.Info("Exiting serial writer goroutine"); return
		}}
	}(ctx)
	
	serialReceiver := make(chan float64, c.QueueSize)
	
	go func(ctx context.Context) {
		
		
		buffer := make([]byte, 8)
		
		var (
			err error
			value float64
		)
		
		err = c.Port.SetReadTimeout(c.ReadTimeout)
		
		if err != nil {
			log.Fatalf("Could not set read timeout for serial port: %v", err)
			return
		}
		
		for { select {
			
			case <-ctx.Done(): slog.Info("Exiting serial receiver goroutine"); return
			
			default:
				_, err = c.Port.Read(buffer)
				
				if err != nil {
					slog.Warn(fmt.Sprintf("Failed receiving data over serial port: %v", err))
					continue
				}
				
				err = binary.Read(bytes.NewBuffer(buffer), binary.BigEndian, value)
				
				if err != nil {
					slog.Warn(fmt.Sprintf("Failed deserializing data received over serial: %v", err))
					continue
				}
				
				serialReceiver <- value
		}}
	}(ctx)
	
	
	return serialReceiver
}