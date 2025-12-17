package handlers

import "github.com/gin-gonic/gin"

type SystemHandler struct {
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/ping", h.handlePing)
}

func (h *SystemHandler) handlePing(c *gin.Context) {
	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}
