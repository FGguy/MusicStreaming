package main

import (
	"errors"
	"flag"
	"fmt"
	handlers "music-streaming/internal/adapter/handlers"
	"music-streaming/internal/adapter/repositories"
	"music-streaming/internal/core/config"
	"music-streaming/internal/core/services"
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

	// Load config file
	// Should be injected into application components that need it
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("Failed loading server configuration file. Error: %s", err)
	}

	// Setup Dependencies

	// Repositories
	// TODO: Use factory to choose repository implementations based on config
	InMemoryUserManagementRepository := repositories.NewInMemoryUserManagementRepository()
	InMemoryMediaBrowsingRepository := repositories.NewInMemoryMediaBrowsingRepository()

	// Services
	UserAuthenticationService := services.NewUserAuthenticationService(InMemoryUserManagementRepository)
	UserManagementService := services.NewUserManagementService(InMemoryUserManagementRepository)
	MediaBrowsingService := services.NewMediaBrowsingService(InMemoryMediaBrowsingRepository)
	MediaRetrievalService := services.NewMediaRetrievalService(InMemoryMediaBrowsingRepository)
	MediaScanningService := services.NewMediaScanningService(config)

	// Middleware
	UserAuthenticationMiddleware := handlers.NewUserManagementMiddleware(UserAuthenticationService)

	// Handlers
	UserManagementHandler := handlers.NewUserManagementHandler(UserManagementService)
	MediaBrowsingHandler := handlers.NewMediaBrowsingHandler(MediaBrowsingService)
	MediaRetrievalHandler := handlers.NewMediaRetrievalHandler(MediaRetrievalService)
	MediaScanningHandler := handlers.NewMediaScanningHandler(MediaScanningService)
	SystemHandler := handlers.NewSystemHandler()

	app := handlers.
		NewApplication().
		WithMiddleware(
			handlers.ValidateSubsonicQueryParameters,
			UserAuthenticationMiddleware.WithAuth,
		).
		WithHandlers(
			UserManagementHandler,
			MediaBrowsingHandler,
			MediaRetrievalHandler,
			MediaScanningHandler,
			SystemHandler,
		).
		RegisterHandlers()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", PORT),
		Handler:        app.Router,
		MaxHeaderBytes: 4 * 1024,
		ReadTimeout:    5 * time.Second,
	}

	log.Info().Msgf("Starting server at address :%d", PORT)
	serverError := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
