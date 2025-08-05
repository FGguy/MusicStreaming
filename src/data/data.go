package data

import (
	"context"
	"errors"
	"fmt"
	sqlc "music-streaming/sql/sqlc"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

/*
	For queries implemented directly by sqlc just get a SqlcConn
	and run the query. For more complicated queries, use queries implemented on the DataLayerPg Type
*/

type DataLayer interface {
	SQLUserManagement
	MediaScan(musicFolders []string, count chan<- int, done chan<- struct{})
}

type DataLayerPg struct {
	Pg_pool *pgxpool.Pool
	cache   *redis.Client
}

// defer pg_pool.Close() has to be called in parent scope
func New(ctx context.Context) (*DataLayerPg, error) {
	pg_pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}
	cache := redis.NewClient(opt)

	if err = SqlSetup(pg_pool); err != nil {
		return nil, err
	}

	return &DataLayerPg{
		Pg_pool: pg_pool,
		cache:   cache,
	}, nil
}

func NewTest(ctx context.Context) (*DataLayerPg, error) {
	pg_pool, err := pgxpool.New(ctx, os.Getenv("TEST_POSTGRES_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(os.Getenv("TEST_REDIS_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}
	cache := redis.NewClient(opt)

	if err = SqlSetup(pg_pool); err != nil {
		return nil, err
	}

	return &DataLayerPg{
		Pg_pool: pg_pool,
		cache:   cache,
	}, nil
}

func SqlSetup(pg_pool *pgxpool.Pool) error {
	var (
		ctx                                 = context.Background()
		cleanStart                          = os.Getenv("CLEAN_START")
		adminName, adminNameDefined         = os.LookupEnv("ADMIN_NAME")
		adminPassword, adminPasswordDefined = os.LookupEnv("ADMIN_PASSWORD")
		adminEmail, adminEmailDefined       = os.LookupEnv("ADMIN_EMAIL")
	)

	if !adminNameDefined || !adminPasswordDefined || !adminEmailDefined {
		return errors.New("SqlSetup: Failed to get admin credentials from ENV. Make sure the variables ADMIN_NAME and ADMIN_PASSWORD are defined in your .env or in the docker-compose file")
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
	if _, err = q.CreateAdminUser(ctx, sqlc.CreateAdminUserParams{Username: adminName, Password: adminPassword, Email: adminEmail}); err != nil {
		return fmt.Errorf("SqlSetup: Failed to create admin user, Err: %s", err)
	}

	return nil
}
