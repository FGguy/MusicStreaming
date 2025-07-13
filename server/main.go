package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pg_connection_string := os.Getenv("POSTGRES_CONNECTION_STRING")
	redis_connection_string := os.Getenv("REDIS_CONNECTION_STRING")

	conn, err := pgx.Connect(context.Background(), pg_connection_string)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	log.Println("Connected to postgres.")

	opt, err := redis.ParseURL(redis_connection_string)
	if err != nil {
		panic(err)
	}

	_ = redis.NewClient(opt)

	log.Println("Connected to redis.")

	for { //wait

	}
}
