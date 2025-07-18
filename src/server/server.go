package server

import (
	"log"

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

func NewServer(pg_pool *pgxpool.Pool, cache *redis.Client) *Server {
	router := gin.Default()
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
	api := s.router.Group("/rest", s.subValidateQParamsMiddleware, s.subWithAuth)
	{
		api.GET("/ping", s.handlePing)

		//User management routes
		api.GET("/getUser", s.hangleGetUser)
		api.GET("/getUsers", s.hangleGetUsers)
		api.GET("/createUser", s.handleCreateUser)
		api.GET("/updateUser", s.handleUpdateUser)
		api.POST("/deleteUser", s.handleDeleteUser)
		api.POST("/changePassword", s.handleChangePassword)
	}
}

func (s *Server) Run(port string) {
	err := s.router.Run(port)
	if err != nil {
		log.Fatal("Failed to run gin router")
	}
}
