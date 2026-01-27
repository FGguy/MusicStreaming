package handlers

import (
	"context"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type MediaScanningHandler struct {
	mediaScanningService ports.MediaScanningPort
	logger               *slog.Logger
}

func NewMediaScanningHandler(mediaScanningService ports.MediaScanningPort, logger *slog.Logger) *MediaScanningHandler {
	return &MediaScanningHandler{
		mediaScanningService: mediaScanningService,
		logger:               logger,
	}
}

func (h *MediaScanningHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/getScanStatus", h.handleGetScanStatus)
	group.POST("/startScan", h.handleStartScan)
}

func (h *MediaScanningHandler) handleGetScanStatus(c *gin.Context) {
	h.logger.Info("Get scan status handler called")
	scanStatus, err := h.mediaScanningService.GetScanStatus(c.Request.Context())
	if err != nil {
		h.logger.Warn("Get scan status handler error", slog.String("error", err.Error()))
		buildAndSendError(c, "0")
		return
	}

	h.logger.Info("Get scan status handler success", slog.Bool("scanning", scanStatus.Scanning), slog.Int("count", scanStatus.Count))

	// Convert to DTO
	scanStatusDTO := ScanStatusToDTO(scanStatus)

	subsonicRes := SubsonicResponse{
		Xmlns:      Xmlns,
		Status:     "ok",
		Version:    SubsonicVersion,
		ScanStatus: &scanStatusDTO,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *MediaScanningHandler) handleStartScan(c *gin.Context) {
	var (
		rUser = c.MustGet(RequestingUserKey).(*domain.User)
		ctx   = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	h.logger.Info("Start scan handler called", slog.String("username", rUser.Username))
	scanStatus, err := h.mediaScanningService.StartScan(ctx)
	if err != nil {
		h.logger.Warn("Start scan handler error", slog.String("username", rUser.Username), slog.String("error", err.Error()))
		switch err.(type) {
		case *ports.NotAuthorizedError:
			buildAndSendError(c, "50")
		default:
			buildAndSendError(c, "0")
		}
		return
	}

	h.logger.Info("Start scan handler success", slog.String("username", rUser.Username), slog.Bool("scanning", scanStatus.Scanning))

	// Convert to DTO
	scanStatusDTO := ScanStatusToDTO(scanStatus)

	subsonicRes := SubsonicResponse{
		Xmlns:      Xmlns,
		Status:     "ok",
		Version:    SubsonicVersion,
		ScanStatus: &scanStatusDTO,
	}

	SerializeAndSendBody(c, subsonicRes)
}
