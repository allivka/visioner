package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/allivka/visioner/unit_gate/sci"
	"go.bug.st/serial"
)

const (
	serverPort = 8080
)

func attemptPorts(names []string, mode *serial.Mode) (serial.Port, []error) {

	var errors []error

	for _, name := range names {
		port, err := serial.Open(name, mode)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		return port, nil

	}

	return nil, errors

}

func serve(port serial.Port) func() {

	ctx, cancel := context.WithCancel(context.Background())

	in := make(chan sci.Behavior, 10)

	receiver := sci.RunSCI(ctx, sci.SCIConfig{
		PerWriteTimeout: 1 * time.Second,
		ReadTimeout:     1 * time.Second,
		Port:            port,
		QueueSize:       10,
		In:              in,
	})

	sci.MaintainChannel(ctx, receiver, 9, 200*time.Millisecond)

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", serverPort),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,

		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer request.Body.Close()

			body, err := io.ReadAll(request.Body)

			if err != nil {
				slog.Warn(fmt.Sprintf("Failed reading request body: %v", err))
				return
			}

			switch request.Method {
			case http.MethodPut:
				fallthrough
			case http.MethodPost:
				behavior, err := sci.ValidateBehaviorBuffer(*bytes.NewBuffer(body))

				if err != nil {
					slog.Warn(fmt.Sprintf("Invalid behavior packet received: %v", err))
					writer.WriteHeader(http.StatusBadRequest)
					writer.Write([]byte(err.Error()))
				}

				in <- behavior
				writer.WriteHeader(http.StatusOK)

			case http.MethodGet:
				angle := <-receiver

				err := binary.Write(writer, binary.BigEndian, angle)

				if err != nil {
					slog.Warn(fmt.Sprintf("Failed writing angle to response: %v", err))
					writer.WriteHeader(http.StatusInternalServerError)
					writer.Write([]byte(err.Error()))
				}

				writer.WriteHeader(http.StatusOK)

			default:
				writer.WriteHeader(http.StatusNotImplemented)
			}
		}),
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatalf("Failed starting http server at port %v: %v", serverPort, err)
			return
		}
	}()

	return func() {
		cancel()
		err := server.Close()

		if err != nil {
			slog.Warn(fmt.Sprintf("Failed successfully closing http server at port %v: %v", serverPort, err))
		}
	}
}

func main() {

	if len(os.Args) > 2 {

		port, errors := attemptPorts(os.Args[1:], &serial.Mode{BaudRate: 115200})

		if errors != nil {
			slog.Error("Failed to open any of the specified ports:")
			for i, err := range errors {
				slog.Error(fmt.Sprintf("%v) Port '%v' opened with error: %v", i, os.Args[i+1], err))
			}
			return
		}

		slog.Info("Successfully opened serial port")
		serve(port)
		return
	}

	names, err := serial.GetPortsList()

	if err != nil {
		log.Fatal(err)
		return
	}

	if len(names) == 0 {
		log.Fatal("No ports specified and no ports available for connection")
	}

	slog.Info(fmt.Sprintf("Got available serial ports on the machine, will use first successful connection: %v", names))

	port, errors := attemptPorts(names, &serial.Mode{BaudRate: 115200})

	if errors != nil {
		slog.Error("Failed to open any of the available ports:")
		for i, err := range errors {
			slog.Error(fmt.Sprintf("%v) Port '%v' opened with error: %v", i, os.Args[i+1], err))
		}
		return
	}

	slog.Info("Successfully opened serial port")
	serve(port)

}
