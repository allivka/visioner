package main

import (
	"context"
	"fmt"
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

func writeData(request *http.Request) error {
	
	return nil
}

func returnData(writer http.ResponseWriter) error {
	
	return nil
}

func serve(port serial.Port) func() {
	
	ctx, cancel := context.WithCancel(context.Background())
	
	in := make(chan sci.Behavior, 10)
	
	sci.RunSCI(ctx, sci.SCIConfig{
		PerWriteTimeout: 1 * time.Second,
		ReadTimeout: 1 * time.Second,
		Port: port,
		QueueSize: 10,
		In: in,
	})
	
	server := http.Server{
		Addr: fmt.Sprintf(":%v", serverPort),
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 10,
		
		Handler: http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
			defer request.Body.Close()
			
			switch request.Method {
				case http.MethodPut: fallthrough
				case http.MethodPost:
					writeData(request)
				case http.MethodGet:
					returnData(writer)
				default:
					writer.WriteHeader(http.StatusBadRequest)
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
	
	return func () {
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
				slog.Error(fmt.Sprintf("%v) Port '%v' opened with error: %v", i, os.Args[i + 1], err))
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
			slog.Error(fmt.Sprintf("%v) Port '%v' opened with error: %v", i, os.Args[i + 1], err))
		}
		return
	}
	
	slog.Info("Successfully opened serial port")
	serve(port)
	
	
}