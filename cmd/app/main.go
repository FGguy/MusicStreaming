package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	handlers "music-streaming/internal/adapter/handlers"
	"music-streaming/internal/adapter/repositories"
	"music-streaming/internal/core/config"
	"music-streaming/internal/core/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const (
	PORT = 8080
)

var (
	logLevels = map[string]slog.Level{
		"info":  slog.LevelInfo,
		"debug": slog.LevelDebug,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)

func main() {
	//Setup Logging
	logLevelFlag := flag.String("loglevel", "info", "Used to set global logging level.")
	flag.Parse()

	logLevel, ok := logLevels[*logLevelFlag]
	if !ok {
		logLevel = slog.LevelInfo // Default to info level
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.NewLogLogger(handler, logLevel)
	jsonLogger := slog.New(handler)

	//Load .env
	if err := godotenv.Load(); err != nil {
		jsonLogger.Error("Error loading .env file", slog.String("error", err.Error()))
	}

	// Load config file
	// Should be injected into application components that need it
	config, err := config.LoadConfig()
	if err != nil {
		jsonLogger.Error("Failed loading server configuration file", slog.String("error", err.Error()))
	}

	dbURL, ok := os.LookupEnv("POSTGRES_CONNECTION_STRING")
	if !ok {
		jsonLogger.Error("POSTGRES_CONNECTION_STRING environment variable is not set")
	}

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		jsonLogger.Error("Failed to connect to database", slog.String("error", err.Error()))
	}
	defer db.Close(context.Background()) // nolint:errcheck

	// Test database connection
	if err := db.Ping(context.Background()); err != nil {
		jsonLogger.Error("Failed to ping database", slog.String("error", err.Error()))
	}
	jsonLogger.Info("Successfully connected to database")

	// Setup Redis Connection
	redisURL, ok := os.LookupEnv("REDIS_CONNECTION_STRING")
	if !ok {
		jsonLogger.Error("REDIS_CONNECTION_STRING environment variable is not set")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		jsonLogger.Error("Failed to connect to Redis", slog.String("error", err.Error()))
	}
	jsonLogger.Info("Successfully connected to Redis")
	defer redisClient.Close() // nolint:errcheck

	// Setup Dependencies
	// Repositories
	userManagementRepository := repositories.NewSQLUserManagementRepository(db, redisClient)
	mediaBrowsingRepository := repositories.NewSQLMediaBrowsingRepository(db)

	// Services
	userAuthenticationService := services.NewUserAuthenticationService(userManagementRepository, jsonLogger)
	userManagementService := services.NewUserManagementService(userManagementRepository, jsonLogger)
	mediaBrowsingService := services.NewMediaBrowsingService(mediaBrowsingRepository, jsonLogger)
	mediaRetrievalService := services.NewMediaRetrievalService(mediaBrowsingRepository, jsonLogger)
	mediaScanningService := services.NewMediaScanningService(mediaBrowsingRepository, config, jsonLogger)

	// Middleware
	userAuthenticationMiddleware := handlers.NewUserManagementMiddleware(userAuthenticationService, jsonLogger)

	// Handlers
	userManagementHandler := handlers.NewUserManagementHandler(userManagementService, jsonLogger)
	mediaBrowsingHandler := handlers.NewMediaBrowsingHandler(mediaBrowsingService, jsonLogger)
	mediaRetrievalHandler := handlers.NewMediaRetrievalHandler(mediaRetrievalService, jsonLogger)
	mediaScanningHandler := handlers.NewMediaScanningHandler(mediaScanningService, jsonLogger)
	systemHandler := handlers.NewSystemHandler(jsonLogger)

	app := handlers.
		NewApplication().
		WithMiddleware(
			handlers.ValidateSubsonicQueryParameters,
			userAuthenticationMiddleware.WithAuth,
		).
		WithHandlers(
			userManagementHandler,
			mediaBrowsingHandler,
			mediaRetrievalHandler,
			mediaScanningHandler,
			systemHandler,
		).
		RegisterHandlers()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", PORT),
		Handler:        app.Router,
		MaxHeaderBytes: 4 * 1024,
		ReadTimeout:    5 * time.Second,
		ErrorLog:       logger,
	}

	jsonLogger.Info("Starting server at address :%d", slog.Int("port", PORT))
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
		jsonLogger.Error("Server error", slog.String("error", err.Error()))
	case sig := <-stop:
		jsonLogger.Info("Shutting down server", slog.String("signal", sig.String()))
	}
}
