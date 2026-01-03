package handlers

import (
	"log/slog"
	"music-streaming/internal/core/ports"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MediaBrowsingHandler struct {
	MediaBrowsingService ports.MediaBrowsingPort
	logger               *slog.Logger
}

func NewMediaBrowsingHandler(mediaBrowsingServ ports.MediaBrowsingPort, logger *slog.Logger) *MediaBrowsingHandler {
	return &MediaBrowsingHandler{
		MediaBrowsingService: mediaBrowsingServ,
		logger:               logger,
	}
}

func (h *MediaBrowsingHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/getArtist", h.handleGetArtist)
	group.GET("/getAlbum", h.handleGetAlbum)
	group.GET("/getSong", h.handleGetSong)
}

func (h *MediaBrowsingHandler) handleGetArtist(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		paramId = c.Query("id")
	)

	id, err := strconv.Atoi(paramId)
	if paramId == "" || err != nil {
		h.logger.Warn("Get artist handler - invalid id parameter", slog.String("id", paramId))
		buildAndSendError(c, "10")
		return
	}

	h.logger.Info("Get artist handler called", slog.Int("id", id))
	artist, err := h.MediaBrowsingService.GetArtist(ctx, id)
	if err != nil {
		h.logger.Warn("Get artist handler error", slog.Int("id", id), slog.String("error", err.Error()))
		switch err.(type) {
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

	h.logger.Info("Get artist handler success", slog.Int("id", id), slog.String("name", artist.Name))
	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
		Artist:  &artist,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *MediaBrowsingHandler) handleGetAlbum(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		paramId = c.Query("id")
	)

	id, err := strconv.Atoi(paramId)
	if paramId == "" || err != nil {
		h.logger.Warn("Get album handler - invalid id parameter", slog.String("id", paramId))
		buildAndSendError(c, "10")
		return
	}

	h.logger.Info("Get album handler called", slog.Int("id", id))
	album, err := h.MediaBrowsingService.GetAlbum(ctx, id)
	if err != nil {
		h.logger.Warn("Get album handler error", slog.Int("id", id), slog.String("error", err.Error()))
		switch err.(type) {
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

	h.logger.Info("Get album handler success", slog.Int("id", id), slog.String("name", album.Name))
	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
		Album:   &album,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *MediaBrowsingHandler) handleGetSong(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		paramId = c.Query("id")
	)

	id, err := strconv.Atoi(paramId)
	if paramId == "" || err != nil {
		h.logger.Warn("Get song handler - invalid id parameter", slog.String("id", paramId))
		buildAndSendError(c, "10")
		return
	}

	h.logger.Info("Get song handler called", slog.Int("id", id))
	song, err := h.MediaBrowsingService.GetSong(ctx, id)
	if err != nil {
		h.logger.Warn("Get song handler error", slog.Int("id", id), slog.String("error", err.Error()))
		switch err.(type) {
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

	h.logger.Info("Get song handler success", slog.Int("id", id), slog.String("title", song.Title))
	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
		Song:    &song,
	}

	SerializeAndSendBody(c, subsonicRes)
}
