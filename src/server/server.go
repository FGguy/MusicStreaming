package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Config struct {
}

type Server struct {
	Router  *gin.Engine
	pg_pool *pgxpool.Pool
	cache   *redis.Client
	config  *Config
}

func NewServer(pg_pool *pgxpool.Pool, cache *redis.Client) *Server {
	router := gin.Default()
	config := &Config{}

	server := &Server{
		Router:  router,
		pg_pool: pg_pool,
		cache:   cache,
		config:  config,
	}
	server.mountHandlers()

	return server
}

func (s *Server) mountHandlers() {
	api := s.Router.Group("/rest", s.subValidateQParamsMiddleware, s.subWithAuth)
	{
		api.GET("/ping", s.handlePing)

		//User management routes
		api.GET("/getUser", s.hangleGetUser)
		api.GET("/getUsers", s.hangleGetUsers)
		api.POST("/createUser", s.handleCreateUser)
		api.POST("/updateUser", s.handleUpdateUser)
		api.POST("/deleteUser", s.handleDeleteUser)
		api.POST("/changePassword", s.handleChangePassword)
	}
}
