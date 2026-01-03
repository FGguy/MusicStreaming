package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	logger *slog.Logger
}

func NewSystemHandler(logger *slog.Logger) *SystemHandler {
	return &SystemHandler{
		logger: logger,
	}
}

func (h *SystemHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/ping", h.handlePing)
}

func (h *SystemHandler) handlePing(c *gin.Context) {
	h.logger.Debug("Ping handler called")
	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}
