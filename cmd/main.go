package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Linus-Regander/Go-Microservice/cmd/config"
)

const serviceName = "go-microservice"

type Service struct {
	Config *config.Config
	Name   string
	Logger *log.Logger
}

func main() {
	//
	// Service Setup.
	//

	var err error

	service := &Service{
		Name: serviceName,
	}

	startupLogger := log.New(
		os.Stdout,
		serviceName,
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	serviceCtx, cancelCtx := context.WithCancel(context.Background())
	defer func() {
		if r := recover(); r != nil {
			startupLogger.Fatal(fmt.Sprintf("recovered panic: %v", r))
		}

		startupLogger.Print("service shutting down")

		cancelCtx()

		return
	}()

	cfg, err := config.Setup()
	if err != nil {
		startupLogger.Print(fmt.Errorf("setup config err: %w", err))

		return
	}

	service.Config = cfg
	service.Logger = startupLogger

	//
	// Service Startup.
	//

	srv := &http.Server{
		Addr: service.Config.Service.Port,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return
		}
	}()
	defer func() {
		shutdownCtx, shutdownCancelCtx := context.WithCancel(serviceCtx)
		defer shutdownCancelCtx()

		if err = srv.Shutdown(shutdownCtx); err != nil {
			service.Logger.Print(fmt.Errorf("Shutdown recieved error: %w", err))

			return
		} else {
			service.Logger.Print("Server shutdown successfully")

			return
		}
	}()

	fmt.Println("Service successfully started")

	<-shutdown()
}

func shutdown() <-chan os.Signal {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	return ch
}
