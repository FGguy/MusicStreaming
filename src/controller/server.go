package controller

import (
	consts "music-streaming/consts"
	"music-streaming/data"
	types "music-streaming/types"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type State struct {
	scanning bool
	count    int
}

type Config struct {
	MusicDirectories []string `mapstructure:"music-directories"`
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
	state := &State{scanning: false, count: 0}

	server := &Server{
		Router:    router,
		state:     state,
		dataLayer: dataLayer,
	}
	server.mountHandlers()

	return server
}

func (s *Server) LoadConfig() error {
	viper.SetConfigName("musicstreaming")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	s.config = &config

	log.Debug().Msgf("Using config file: %s\n", viper.ConfigFileUsed())
	log.Debug().Msgf("Loaded config: %+v\n", config)
	return nil
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
	log.Debug().Msgf("Starting Scan")

	mediaCount := make(chan int)
	done := make(chan struct{})

	go s.dataLayer.MediaScan(s.config.MusicDirectories, mediaCount, done)

	for {
		select {
		case count := <-mediaCount:
			s.mu.Lock()
			s.state.count += count
			s.mu.Unlock()
		case <-done:
			log.Debug().Msgf("Finished Scan")
			s.mu.Lock()
			s.state.scanning = false
			s.mu.Unlock()
			return
		}
	}
}
