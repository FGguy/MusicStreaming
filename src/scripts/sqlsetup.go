package scripts

import (
	"context"
	"fmt"
	sqlc "music-streaming/sql/sqlc"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SqlSetup(pg_pool *pgxpool.Pool) error {
	var (
		ctx                                 = context.Background()
		cleanStart                          = os.Getenv("CLEAN_START")
		adminName, adminNameDefined         = os.LookupEnv("ADMIN_NAME")
		adminPassword, adminPasswordDefined = os.LookupEnv("ADMIN_PASSWORD")
		adminEmail, adminEmailDefined       = os.LookupEnv("ADMIN_EMAIL")
	)

	if !adminNameDefined || !adminPasswordDefined || !adminEmailDefined {
		return fmt.Errorf("SqlSetup: Failed to get admin credentials from ENV. Make sure the variables ADMIN_NAME and ADMIN_PASSWORD are defined in your .env or in the docker-compose file.")
	}

	conn, err := pg_pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("SqlSetup: Failed to acquire connection from postgres connection pool in main, Err: %s", err)
	}
	defer conn.Release()

	dropTables := "./sql/droptables.sql"
	tables := "./sql/tables.sql"

	if cleanStart == "true" {
		dropTablesScript, err := os.ReadFile(dropTables)
		if err != nil {
			return fmt.Errorf("SqlSetup: Failed to open script for dropping tables, Err: %s", err)
		}

		if _, err = conn.Exec(ctx, string(dropTablesScript)); err != nil {
			return fmt.Errorf("SqlSetup: Failed to drop tables in main, Err: %s", err)
		}
	}

	createTablesScript, err := os.ReadFile(tables)
	if err != nil {
		return fmt.Errorf("SqlSetup: Failed to open script for creating tables, Err: %s", err)
	}

	if _, err = conn.Exec(ctx, string(createTablesScript)); err != nil {
		return fmt.Errorf("SqlSetup: Failed to create tables in main, Err: %s", err)
	}

	q := sqlc.New(conn)
	if _, err = q.CreateAdminUser(ctx, sqlc.CreateAdminUserParams{Username: pgtype.Text{String: adminName, Valid: true}, Password: adminPassword, Email: adminEmail}); err != nil {
		return fmt.Errorf("SqlSetup: Failed to create admin user, Err: %s", err)
	}

	return nil
}
