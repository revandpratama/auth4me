package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/revandpratama/auth4me/config"
	"github.com/revandpratama/auth4me/internal/app"
	"github.com/rs/zerolog/log"
)

type Server struct {
	shutdownCh chan os.Signal
	errCh      chan error
}

func NewServer() *Server {
	return &Server{
		shutdownCh: make(chan os.Signal, 1),
		errCh:      make(chan error),
	}
}

func main() {

	server := NewServer()

	signal.Notify(server.shutdownCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := config.LoadConfig(); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	log.Info().Msg("initializing server...")

	apps, err := app.NewApp(
		app.WithDB(),
		app.WithRESTServer(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create app")
	}

	select {
	case shutdown := <-server.shutdownCh:
		log.Info().Msgf("gracefully shutting down the app: %v", shutdown)
		if err := apps.Stop(); err != nil {
			log.Error().Err(err).Msgf("failed to stop app cleanly, cause: %v", err)
		}
		log.Info().Msg("server stopped, goodbye!")
	case err := <-server.errCh:
		log.Error().Err(err).Msgf("failed to start server, cause: %v", err)
	}
}
