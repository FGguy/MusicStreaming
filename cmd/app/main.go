package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	controller "music-streaming/cmd/app/controller"
	"music-streaming/internal/data"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PORT = 8080
)

var (
	logLevels = map[string]zerolog.Level{
		"trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
		"info":  zerolog.InfoLevel,
	}
)

func main() {
	//Setup Logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logLevelFlag := flag.String("loglevel", "info", "Used to set global logging level.")
	flag.Parse()

	logLevel, ok := logLevels[*logLevelFlag]
	if !ok {
		logLevel = zerolog.InfoLevel
		log.Info().Msg("No log level passed or invalid value. Log level set to: info")
	} else {
		log.Info().Msgf("Log level set to: %s", *logLevelFlag)
	}
	zerolog.SetGlobalLevel(logLevel)

	//Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	//Initialize data layer, ie. connect to databases
	dataLayer, err := data.New(context.Background())
	if err != nil {
		log.Fatal().Msgf("Failed initializing data layer. Error: %s", err)
	}
	defer dataLayer.Pg_pool.Close()

	//Load config file
	config, err := controller.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("Failed loading server configuration file. Error: %s", err)
	}

	app := controller.NewApplication(dataLayer, config)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", PORT),
		Handler:        app.Router,
		MaxHeaderBytes: 4 * 1024,
		ReadTimeout:    5 * time.Second,
	}

	log.Info().Msgf("Starting server at address :%d", PORT)
	serverError := make(chan error, 1)
	go func() {
		if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
