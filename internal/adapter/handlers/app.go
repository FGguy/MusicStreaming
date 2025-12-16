package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	MusicDirectories []string `mapstructure:"music-directories"`
}

type Application struct {
	Router *gin.Engine
	config *Config

	userManagementHandler *UserManagementHandler
	userAuthMiddleware    *UserManagementMiddleware

	mediaBrowsingHandler  *MediaBrowsingHandler
	mediaRetrievalHandler *MediaRetrievalHandler
	mediaScanningHandler  *MediaScanningHandler
}

func NewApplication(config *Config, userManagementHandler *UserManagementHandler, userAuthMiddleware *UserManagementMiddleware, mediaBrowsingHandler *MediaBrowsingHandler, mediaRetrievalHandler *MediaRetrievalHandler, mediaScanningHandler *MediaScanningHandler) *Application {
	router := gin.Default()

	app := &Application{
		Router:                router,
		config:                config,
		userManagementHandler: userManagementHandler,
		userAuthMiddleware:    userAuthMiddleware,
		mediaBrowsingHandler:  mediaBrowsingHandler,
		mediaRetrievalHandler: mediaRetrievalHandler,
		mediaScanningHandler:  mediaScanningHandler,
	}
	app.mountHandlers()

	return app
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("musicstreaming")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	log.Debug().Msgf("Using config file: %s\n", viper.ConfigFileUsed())
	log.Debug().Msgf("Loaded config: %+v\n", config)
	return &config, nil
}

func (s *Application) mountHandlers() {
	api := s.Router.Group("/rest", ValidateSubsonicQueryParameters, s.userAuthMiddleware.WithAuth)
	{
		api.GET("/ping", s.handlePing)

		//User management routes
		api.GET("/getUser", s.userManagementHandler.hangleGetUser)
		api.GET("/getUsers", s.userManagementHandler.hangleGetUsers)
		api.POST("/createUser", s.userManagementHandler.handleCreateUser)
		api.POST("/updateUser", s.userManagementHandler.handleUpdateUser)
		api.POST("/deleteUser", s.userManagementHandler.handleDeleteUser)
		api.POST("/changePassword", s.userManagementHandler.handleChangePassword)

		//Media browsing routes
		api.GET("/getArtist", s.mediaBrowsingHandler.handleGetArtist)
		api.GET("/getAlbum", s.mediaBrowsingHandler.handleGetAlbum)
		api.GET("/getSong", s.mediaBrowsingHandler.handleGetSong)

		//Media retrieval routes
		api.GET("/stream", s.mediaRetrievalHandler.handleStream)
		api.GET("/download", s.mediaRetrievalHandler.handleDownload)
		api.GET("/getCoverArt", s.mediaRetrievalHandler.handleGetCoverArt)

		//Media scanning routes
		api.GET("/getScanStatus", s.mediaScanningHandler.handleGetScanStatus)
		api.POST("/startScan", s.mediaScanningHandler.handleStartScan)
	}
}

func (s *Application) handlePing(c *gin.Context) {
	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}
