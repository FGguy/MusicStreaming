package server

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Config struct {
}

type Server struct {
	router  *gin.Engine
	pg_pool *pgxpool.Pool
	cache   *redis.Client
	config  *Config
}

func NewServer() *Server {
	handleErr := func(err error) {
		if err != nil {
			log.Fatalf("Failed to create server\nError: %s", err)
		}
	}

	router := gin.Default()

	//do i need to defer closing connections somewhere?
	pg_pool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_CONNECTION_STRING"))
	handleErr(err)

	opt, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION_STRING"))
	handleErr(err)
	cache := redis.NewClient(opt)

	config := &Config{}

	server := &Server{
		router:  router,
		pg_pool: pg_pool,
		cache:   cache,
		config:  config,
	}
	server.mountHandlers()

	return server
}

func (s *Server) mountHandlers() {
	api := s.router.Group("/rest")
	{
		api.GET("/ping", s.WithAuth, s.handlePing)
	}
}

func (s *Server) Run(port string) {
	err := s.router.Run(port)
	if err != nil {
		log.Fatal("Failed to run gin router")
	}
}
