package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"music-streaming/scripts"
	server "music-streaming/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const (
	PORT = 8080
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	pg_pool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect to postgres Database in main, Err: %s", err)
	}
	defer pg_pool.Close()

	opt, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect to redis Database in main, Err: %s", err)
	}
	cache := redis.NewClient(opt)

	scripts.SqlSetup(pg_pool)

	serverError := make(chan error, 1)
	go func() {
		if err := server.NewServer(pg_pool, cache).Router.Run(fmt.Sprintf(":%d", PORT)); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	case sig := <-stop:
		log.Printf("Received shutdown signal: %v", sig)
	}

	//TODO: add graceful shutdown
}
