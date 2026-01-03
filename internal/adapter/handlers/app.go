package handlers

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(group *gin.RouterGroup)
}

type Application struct {
	Router     *gin.Engine
	middleware []gin.HandlerFunc
	handlers   []Handler
}

func NewApplication() *Application {
	router := gin.Default()

	app := &Application{
		Router: router,
	}

	return app
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

func (a *Application) RegisterHandlers() *Application {
	api := a.Router.Group("/rest", a.middleware...)
	for _, handler := range a.handlers {
		handler.RegisterRoutes(api)
	}
	return a
}
