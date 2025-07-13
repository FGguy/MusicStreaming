package main

import (
	"fmt"
	"log"
	server "music-streaming/server"

	"github.com/joho/godotenv"
)

const (
	PORT = 8080
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server.NewServer().Run(fmt.Sprintf(":%d", PORT))
}
