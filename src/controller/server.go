package controller

import (
	consts "music-streaming/consts"
	"music-streaming/data"
	types "music-streaming/types"
	"sync"

	"github.com/gin-gonic/gin"
)

type State struct {
	scanning bool
	count    int
}

type Config struct {
}

type Server struct {
	Router    *gin.Engine
	config    *Config
	dataLayer data.DataLayer

	mu    sync.Mutex
	state *State
}

func NewServer(dataLayer *data.DataLayerPg) *Server {
	router := gin.Default()
	config := &Config{}
	state := &State{scanning: false, count: 0}

	server := &Server{
		Router:    router,
		config:    config,
		state:     state,
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

		//Media library scanning
		api.GET("/getScanStatus", s.handleGetScanStatus)
		api.GET("/startScan", s.handleStartScan)
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

func (s *Server) MediaScan() {
	s.mu.Lock()
	s.state.scanning = true
	s.state.count = 0
	s.mu.Unlock()

	topLevelDirs := []string{}
	mediaCount := make(chan int)
	done := make(chan struct{})

	go s.dataLayer.MediaScan(topLevelDirs, mediaCount, done)

	for {
		select {
		case count := <-mediaCount:
			s.mu.Lock()
			s.state.count += count
			s.mu.Unlock()
		case <-done:
			s.mu.Lock()
			s.state.scanning = false
			s.mu.Unlock()
			return
		}
	}
}
