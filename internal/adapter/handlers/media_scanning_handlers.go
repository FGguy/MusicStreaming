package handlers

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type MediaScanningHandler struct {
	mediaScanningService ports.MediaScanningPort
}

func NewMediaScanningHandler(mediaScanningService ports.MediaScanningPort) *MediaScanningHandler {
	return &MediaScanningHandler{
		mediaScanningService: mediaScanningService,
	}
}

func (h *MediaScanningHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/getScanStatus", h.handleGetScanStatus)
	group.POST("/startScan", h.handleStartScan)
}

func (h *MediaScanningHandler) handleGetScanStatus(c *gin.Context) {
	scanStatus, err := h.mediaScanningService.GetScanStatus(c.Request.Context())
	if err != nil {
		buildAndSendError(c, "0")
	}

	subsonicRes := SubsonicResponse{
		Xmlns:      Xmlns,
		Status:     "ok",
		Version:    SubsonicVersion,
		ScanStatus: &scanStatus,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *MediaScanningHandler) handleStartScan(c *gin.Context) {
	var (
		rUser = c.MustGet(RequestingUserKey).(*domain.User)
		ctx   = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	scanStatus, err := h.mediaScanningService.StartScan(ctx)
	if err != nil {
		switch err.(type) {
		case *ports.NotAuthorizedError:
			buildAndSendError(c, "50")
		default:
			buildAndSendError(c, "0")
		}
		return
	}

	subsonicRes := SubsonicResponse{
		Xmlns:      Xmlns,
		Status:     "ok",
		Version:    SubsonicVersion,
		ScanStatus: &scanStatus,
	}

	SerializeAndSendBody(c, subsonicRes)
}
