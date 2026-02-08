package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"go.bug.st/serial"
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

func openPort() serial.Port {
	if len(os.Args) > 2 {

		port, errors := attemptPorts(os.Args[1:], &serial.Mode{BaudRate: 115200})

		if errors != nil {
			slog.Error("Failed to open any of the specified ports:")
			for i, err := range errors {
				slog.Error(fmt.Sprintf("%v) Port '%v' opened with error: %v", i, os.Args[i+1], err))
			}
			os.Exit(1)
		}

		slog.Info("Successfully opened serial port")
		return port
	}

	names, err := serial.GetPortsList()

	if err != nil {
		log.Fatal(err)
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
		os.Exit(1)
	}

	slog.Info("Successfully opened serial port")
	
	return port
}

func main() {

	port := openPort()
	serveGate(port)

}
