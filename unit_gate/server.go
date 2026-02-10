package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"
	"unsafe"

	"github.com/allivka/visioner/unit_gate/sci"
	"go.bug.st/serial"
)

const (
	serverPort = 8080
)

func serveGate(port serial.Port) func() {

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
				behavior, err := sci.ValidateBehaviorBuffer(body)

				if err != nil {
					slog.Warn(fmt.Sprintf("Invalid behavior packet received: %v", err))
					writer.WriteHeader(http.StatusBadRequest)
					writer.Write([]byte(err.Error()))
				}

				in <- behavior
				writer.WriteHeader(http.StatusOK)

			case http.MethodGet:
				angle := <-receiver

				writer.Write(unsafe.Slice((*byte)(unsafe.Pointer(&angle)), unsafe.Sizeof(angle)))

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