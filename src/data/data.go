package data

import (
	"context"
	"music-streaming/scripts"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

/*
	For queries implemented directly by sqlc just get a SqlcConn
	and run the query. For more complicated queries, use queries implemented on the DataLayer Type
*/

type DataLayer struct {
	Pg_pool *pgxpool.Pool
	cache   *redis.Client
}

// defer pg_pool.Close() has to be called in parent scope
func New(ctx context.Context) (*DataLayer, error) {
	pg_pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}
	cache := redis.NewClient(opt)

	if err = scripts.SqlSetup(pg_pool); err != nil {
		return nil, err
	}

	return &DataLayer{
		Pg_pool: pg_pool,
		cache:   cache,
	}, nil
}

func NewTest(ctx context.Context) (*DataLayer, error) {
	pg_pool, err := pgxpool.New(ctx, os.Getenv("TEST_POSTGRES_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(os.Getenv("TEST_REDIS_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}
	cache := redis.NewClient(opt)

	if err = scripts.SqlSetup(pg_pool); err != nil {
		return nil, err
	}

	return &DataLayer{
		Pg_pool: pg_pool,
		cache:   cache,
	}, nil
}
