package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"music-streaming/internal/adapter/repositories"
	handlers "music-streaming/internal/adapter/handlers"
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
	"github.com/rs/zerolog/log"
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
		log.Fatal().Msg("Error loading .env file")
	}

	// Load config file
	// Should be injected into application components that need it
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("Failed loading server configuration file. Error: %s", err)
	}

	// Setup Database Connection
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		dbUser = "myuser"
	}
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		dbPassword = "mypassword"
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = "music_streaming"
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to database: %s", err)
	}
	defer db.Close(context.Background())

	// Test database connection
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal().Msgf("Failed to ping database: %s", err)
	}
	log.Info().Msg("Successfully connected to database")

	// Setup Redis Connection
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = ""
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal().Msgf("Failed to connect to Redis: %s", err)
	}
	log.Info().Msg("Successfully connected to Redis")
	defer redisClient.Close()

	// Setup Dependencies
	// Repositories
	UserManagementRepository := repositories.NewSQLUserManagementRepository(db, redisClient)
	MediaBrowsingRepository := repositories.NewSQLMediaBrowsingRepository(db)

	// Services
	UserAuthenticationService := services.NewUserAuthenticationService(UserManagementRepository, jsonLogger)
	UserManagementService := services.NewUserManagementService(UserManagementRepository, jsonLogger)
	MediaBrowsingService := services.NewMediaBrowsingService(MediaBrowsingRepository, jsonLogger)
	MediaRetrievalService := services.NewMediaRetrievalService(MediaBrowsingRepository, jsonLogger)
	MediaScanningService := services.NewMediaScanningService(MediaBrowsingRepository, config, jsonLogger)

	// Middleware
	UserAuthenticationMiddleware := handlers.NewUserManagementMiddleware(UserAuthenticationService, jsonLogger)

	// Handlers
	UserManagementHandler := handlers.NewUserManagementHandler(UserManagementService, jsonLogger)
	MediaBrowsingHandler := handlers.NewMediaBrowsingHandler(MediaBrowsingService, jsonLogger)
	MediaRetrievalHandler := handlers.NewMediaRetrievalHandler(MediaRetrievalService, jsonLogger)
	MediaScanningHandler := handlers.NewMediaScanningHandler(MediaScanningService, jsonLogger)
	SystemHandler := handlers.NewSystemHandler(jsonLogger)

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
		ErrorLog:       logger,
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
		jsonLogger.Error("Server error", slog.String("error", err.Error()))
	case sig := <-stop:
		jsonLogger.Info("Shutting down server", slog.String("signal", sig.String()))
	}
}
