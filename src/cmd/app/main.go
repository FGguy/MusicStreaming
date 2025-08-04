package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	controller "music-streaming/controller"
	"music-streaming/data"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

const (
	PORT = 8080
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dataLayer, err := data.New(context.Background())
	if err != nil {
		log.Fatalf("Failed initializing data layer. Error: %s", err)
	}
	defer dataLayer.Pg_pool.Close()

	server := controller.NewServer(dataLayer)
	if err = server.LoadConfig(); err != nil {
		log.Fatalf("Failed loading server configuration file. Error: %s", err)
	}

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
		log.Printf("Server error: %v", err)
	case sig := <-stop:
		log.Printf("Received shutdown signal: %v", sig)
	}

	//TODO: add graceful shutdown
}
