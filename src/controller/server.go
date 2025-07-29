package controller

import (
	consts "music-streaming/consts"
	"music-streaming/data"
	types "music-streaming/types"

	"github.com/gin-gonic/gin"
)

type Config struct {
}

type Server struct {
	Router    *gin.Engine
	config    *Config
	dataLayer *data.DataLayer
}

func NewServer(dataLayer *data.DataLayer) *Server {
	router := gin.Default()
	config := &Config{}

	server := &Server{
		Router:    router,
		config:    config,
		dataLayer: dataLayer,
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

func (s *Server) handlePing(c *gin.Context) {
	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}
