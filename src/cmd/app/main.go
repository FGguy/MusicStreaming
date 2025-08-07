package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	controller "music-streaming/controller"
	"music-streaming/data"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PORT = 8080
)

func main() {
	SetupLogging()

	if err := godotenv.Load(); err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	dataLayer, err := data.New(context.Background())
	if err != nil {
		log.Fatal().Msgf("Failed initializing data layer. Error: %s", err)
	}
	defer dataLayer.Pg_pool.Close()

	config, err := controller.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("Failed loading server configuration file. Error: %s", err)
	}

	server := controller.NewServer(dataLayer, config)

	serverError := make(chan error, 1)
	go func() {
		if err = server.Router.Run(fmt.Sprintf(":%d", PORT)); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Info().Msgf("Server error: %v", err)
	case sig := <-stop:
		log.Info().Msgf("Received shutdown signal: %v", sig)
	}
}

func SetupLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logLevel := flag.String("loglevel", "info", "Used to set global logging level.")
	flag.Parse()

	zerolog.SetGlobalLevel(getLogLevel(logLevel))

	log.Info().Msgf("Log level set to %s", *logLevel)
}

func getLogLevel(level *string) zerolog.Level {
	switch *level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
