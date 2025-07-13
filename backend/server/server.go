package server

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Router  *gin.Engine
	pg_pool *pgxpool.Pool
	cache   *redis.Client
}

func NewServer() *Server {
	handleErr := func(err error) {
		if err != nil {
			log.Fatalf("Failed to create server\nError: %s", err)
		}
	}

	//do i need to defer closing connections somewhere?
	pg_pool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_CONNECTION_STRING"))
	handleErr(err)

	opt, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION_STRING"))
	handleErr(err)
	cache := redis.NewClient(opt)

	router := gin.Default()

	server := &Server{
		Router:  router,
		pg_pool: pg_pool,
		cache:   cache,
	}
	server.mountHandlers()

	return server
}

func (s *Server) mountHandlers() {
	api := s.Router.Group("/rest")
	{
		api.GET("/ping", s.handlePing)
	}
}

func (s *Server) Run(port string) {
	err := s.Router.Run(port)
	if err != nil {
		log.Fatal("Failed to run gin router")
	}
}
