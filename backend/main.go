package main

import (
	"context"
	"fmt"
	"log"
	server "music-streaming/server"
	"os"

	sqlc "music-streaming/sql/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const (
	PORT = 8080
)

func main() {
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

	sqlSetup(pg_pool)

	server.NewServer(pg_pool, cache).Run(fmt.Sprintf(":%d", PORT))
}

func sqlSetup(pg_pool *pgxpool.Pool) {
	ctx := context.Background()

	adminName, adminNameDefined := os.LookupEnv("ADMIN_NAME")
	adminPassword, adminPasswordDefined := os.LookupEnv("ADMIN_PASSWORD")
	adminEmail, adminEmailDefined := os.LookupEnv("ADMIN_EMAIL")
	if !adminNameDefined || !adminPasswordDefined || !adminEmailDefined {
		log.Fatal("Failed to get admin credentials from ENV. Make sure the variables ADMIN_NAME and ADMIN_PASSWORD are defined in your .env or in the docker-compose file.")
	}

	createTablesScript, err := os.ReadFile("./sql/tables.sql")
	if err != nil {
		log.Fatalf("Failed to open script for creating tables, Err: %s", err)
	}

	conn, err := pg_pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Failed to acquire connection from postgres connection pool in main, Err: %s", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, string(createTablesScript)) //create all tables
	if err != nil {
		log.Fatalf("Failed to create tables in main, Err: %s", err)
	}

	q := sqlc.New(conn)
	_, err = q.CreateAdminUser(ctx, sqlc.CreateAdminUserParams{Username: pgtype.Text{String: adminName, Valid: true}, Password: adminPassword, Email: adminEmail})
	if err != nil {
		log.Fatalf("Failed to create admin user, Err: %s", err)
	}
}
