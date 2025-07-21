package cmd

import (
	"context"
	"fmt"
	"log"
	server "music-streaming/server"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const (
	PORT = 8080
)

func Execute() {
	err := godotenv.Load()
	if err != nil {
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

	server.SqlSetup(pg_pool)

	//Blocks here
	server.NewServer(pg_pool, cache).Run(fmt.Sprintf(":%d", PORT))
}
