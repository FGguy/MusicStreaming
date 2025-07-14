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

	//TODO: load config file
	//TODO: Create tables if they dont exist
	//TODO: Create admin user if he doesnt exist

	server.NewServer().Run(fmt.Sprintf(":%d", PORT))
}
