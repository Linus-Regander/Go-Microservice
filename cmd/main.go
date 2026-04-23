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
	"github.com/Linus-Regander/Go-Microservice/internal/api/handler"
	userRepository "github.com/Linus-Regander/Go-Microservice/internal/api/repository/user"
	"github.com/Linus-Regander/Go-Microservice/internal/api/router"
	"github.com/Linus-Regander/Go-Microservice/internal/api/service"

	"github.com/go-chi/chi/v5"

	"github.com/sanity-io/litter"

	httpSwagger "github.com/swaggo/http-swagger"
)

const serviceName = "go-microservice"

type Service struct {
	Config *config.Config
	Name   string
	Logger *log.Logger
}

// @title User Service API
// @version 1.4
// @description Microservice for User Management
// @BasePath /
func main() {
	//
	// Service Setup.
	//

	var err error

	microService := &Service{
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

	microService.Config = cfg
	microService.Logger = startupLogger

	//
	// Service API.
	//

	mainRouter := chi.NewRouter()
	mainRouter.Get("/swagger/*", httpSwagger.WrapHandler)

	serviceRouter := router.New(handler.New(microService.Logger, service.New(microService.Logger, userRepository.New(microService.Logger, nil))))
	mainPath, userAPI := serviceRouter.SetupChi()

	mainRouter.Group(func(r chi.Router) {
		r.Route("/service", func(r chi.Router) {
			r.Mount(mainPath, userAPI)
		})
	})

	routeMap, err := routes(mainRouter)
	if err != nil {
		startupLogger.Print(fmt.Errorf("setup config err: %w", err))

		return
	}

	litter.Dump(routeMap)

	//
	// Service Startup.
	//

	srv := &http.Server{
		Addr: microService.Config.Service.Port,
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
			microService.Logger.Print(fmt.Errorf("Shutdown recieved error: %w", err))

			return
		} else {
			microService.Logger.Print("Server shutdown successfully")

			return
		}
	}()

	fmt.Println("Service successfully started")

	<-shutdown()
}

func routes(r *chi.Mux) (map[string]string, error) {
	var (
		err      error
		routeMap = make(map[string]string)
	)

	err = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routeMap[route] = method

		return nil
	})
	if err != nil {
		return nil, err
	}

	return routeMap, nil
}

func shutdown() <-chan os.Signal {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	return ch
}
