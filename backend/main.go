package main

import (
	"context"
	"fmt"
	"log"
	server "music-streaming/server"
	"os"

	sqlc "music-streaming/sql/sqlc"

	"github.com/jackc/pgx/v5"
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

	ctx := context.Background()

	adminName, adminNameDefined := os.LookupEnv("ADMIN_NAME")
	adminPassword, adminPasswordDefined := os.LookupEnv("ADMIN_PASSWORD")
	if !adminNameDefined || !adminPasswordDefined {
		log.Fatal("Failed to get admin credentials from ENV. Make sure the variables ADMIN_NAME and ADMIN_PASSWORD are defined in your .env or in the docker-compose file.")
	}

	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect to Database in main, Err: %s", err)
	}
	defer conn.Close(ctx)

	createTablesScript, err := os.ReadFile("./sql/tables.sql")
	if err != nil {
		log.Fatalf("Failed to open script for creating tables, Err: %s", err)
	}

	_, err = conn.Exec(ctx, string(createTablesScript)) //create all tables
	if err != nil {
		log.Fatalf("Failed to create tables in main, Err: %s", err)
	}

	q := sqlc.New(conn)
	_, err = q.InsertUser(ctx, sqlc.InsertUserParams{Name: adminName, Password: adminPassword})
	if err != nil {
		log.Fatalf("Failed to create admin user, Err: %s", err)
	}

	server.NewServer().Run(fmt.Sprintf(":%d", PORT))
}
