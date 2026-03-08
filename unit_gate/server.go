package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net"
	"net/http"
	"time"

	"github.com/allivka/visioner/unit_gate/sci"
	"go.bug.st/serial"
	"gocv.io/x/gocv"
)

const (
	serverPort = 80
	cameraStreamPort = ":8080"
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
				writer.WriteHeader(http.StatusOK)
				buffer := make([]byte, 8)
				binary.LittleEndian.PutUint32(buffer, math.Float32bits(angle))
				writer.Write(buffer)
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
	
	go func() {
		listener, err := net.Listen("tcp", cameraStreamPort)
		
		if err != nil {
			log.Fatal(err)
		}
		
		defer listener.Close()
		
		connection, err := listener.Accept()
		
		for err != nil {
			slog.Warn(err.Error())
			connection, err = listener.Accept()
		}
		
		defer connection.Close()
		
		writer := bufio.NewWriter(connection)
		_ = writer
		camera, err := gocv.OpenVideoCapture(0)
		
		if err != nil {
			slog.Warn(err.Error())
			return
		}
		
		defer camera.Close()
		
		camera.Set(gocv.VideoCaptureFrameWidth, 640)
		camera.Set(gocv.VideoCaptureFrameHeight, 480)
		
		mat := gocv.NewMat()
		
		defer mat.Close()
		
		var (
			ok bool
		)
		
		for {
			if ok = camera.Read(&mat); !ok || mat.Empty() {
				continue
			}
			
			buff, err := gocv.IMEncode(".jpg", mat)
			
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			
			data := buff.GetBytes()
			binary.Write(writer, binary.LittleEndian, uint32(len(data)))
			writer.Write(data)
			writer.Flush()
			
			buff.Close()
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
