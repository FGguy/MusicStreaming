package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Handler interface {
	RegisterRoutes(group *gin.RouterGroup)
}

type Config struct {
	MusicDirectories []string `mapstructure:"music-directories"`
}

type Application struct {
	Router     *gin.Engine
	config     *Config
	middleware []gin.HandlerFunc
	handlers   []Handler
}

// Global middleware is run in order that they are added
func (a *Application) WithMiddleware(middleware ...gin.HandlerFunc) *Application {
	a.middleware = middleware
	return a
}

func (a *Application) WithHandlers(handlers ...Handler) *Application {
	a.handlers = handlers
	return a
}

func NewApplication(config *Config) *Application {
	router := gin.Default()

	app := &Application{
		Router: router,
		config: config,
	}

	return app
}

func (a *Application) RegisterHandlers() *Application {
	api := a.Router.Group("/rest", a.middleware...)
	for _, handler := range a.handlers {
		handler.RegisterRoutes(api)
	}
	return a
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
