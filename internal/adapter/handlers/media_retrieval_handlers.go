package handlers

import (
	"context"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MediaRetrievalHandler struct {
	MediaRetrievalService ports.MediaRetrievalPort
	logger                *slog.Logger
}

func NewMediaRetrievalHandler(mediaRetrievalService ports.MediaRetrievalPort, logger *slog.Logger) *MediaRetrievalHandler {
	return &MediaRetrievalHandler{
		MediaRetrievalService: mediaRetrievalService,
		logger:                logger,
	}
}

func (h *MediaRetrievalHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/stream", h.handleStream)
	group.GET("/download", h.handleDownload)
	group.GET("/getCoverArt", h.handleGetCoverArt)
}

func (h *MediaRetrievalHandler) handleDownload(c *gin.Context) {
	var (
		rUser   = c.MustGet(RequestingUserKey).(*domain.User)
		ctx     = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
		paramId = c.Query("id")
	)

	id, err := strconv.Atoi(paramId)
	if paramId == "" || err != nil {
		h.logger.Warn("Download handler - invalid id parameter", slog.String("id", paramId), slog.String("username", rUser.Username))
		buildAndSendError(c, "10")
		return
	}

	h.logger.Info("Download handler called", slog.Int("id", id), slog.String("username", rUser.Username))
	song, err := h.MediaRetrievalService.DownloadSong(ctx, id)
	if err != nil {
		h.logger.Warn("Download handler error", slog.Int("id", id), slog.String("username", rUser.Username), slog.String("error", err.Error()))
		switch err.(type) {
		case *ports.NotAuthorizedError:
			buildAndSendError(c, "50")
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

	h.logger.Info("Download handler success", slog.Int("id", id), slog.String("title", song.Title), slog.String("username", rUser.Username))
	c.FileAttachment(song.Path, song.Title)
}

type StreamParameters struct {
	Id                     string `form:"id" binding:"required"`
	MaxBitRate             int    `form:"maxBitRate"`
	Format                 string `form:"format"`
	EstimatedContentLength bool   `form:"estimateContentLength"`
}

func (h *MediaRetrievalHandler) handleStream(c *gin.Context) {
	var (
		rUser = c.MustGet(RequestingUserKey).(*domain.User)
		ctx   = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	var params StreamParameters
	if err := c.ShouldBindQuery(&params); err != nil {
		h.logger.Warn("Stream handler - bind error", slog.String("username", rUser.Username), slog.String("error", err.Error()))
		buildAndSendError(c, "10")
		return
	}

	id, err := strconv.Atoi(params.Id)
	if params.Id == "" || err != nil {
		h.logger.Warn("Stream handler - invalid id parameter", slog.String("id", params.Id), slog.String("username", rUser.Username))
		buildAndSendError(c, "10")
		return
	}

	h.logger.Info("Stream handler called", slog.Int("id", id), slog.String("username", rUser.Username), slog.Int("maxBitRate", params.MaxBitRate))
	song, err := h.MediaRetrievalService.StreamSong(ctx, id)
	if err != nil {
		h.logger.Warn("Stream handler error", slog.Int("id", id), slog.String("username", rUser.Username), slog.String("error", err.Error()))
		switch err.(type) {
		case *ports.NotAuthorizedError:
			buildAndSendError(c, "50")
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

	//Move this to service layer and have layer return a stream object instead of song
	//fix bitrate conversion
	if params.MaxBitRate > 0 && song.BitRate > params.MaxBitRate {
		//adjust bitrate
		h.logger.Debug("Bitrate adjustment needed", slog.Int("songBitRate", song.BitRate), slog.Int("maxBitRate", params.MaxBitRate))
	}
	//Format?
	//Validate string
	//transcode file

	if params.EstimatedContentLength {
		//set header
	}

	h.logger.Info("Stream handler success", slog.Int("id", id), slog.String("title", song.Title), slog.String("username", rUser.Username))
	c.File(song.Path)
}

func (h *MediaRetrievalHandler) handleGetCoverArt(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		paramId = c.Query("id")
	)

	id, err := strconv.Atoi(paramId)
	if paramId == "" || err != nil {
		h.logger.Warn("Get cover art handler - invalid id parameter", slog.String("id", paramId))
		buildAndSendError(c, "10")
		return
	}

	h.logger.Info("Get cover art handler called", slog.Int("id", id))
	cover, err := h.MediaRetrievalService.GetCover(ctx, id)
	if err != nil {
		h.logger.Warn("Get cover art handler error", slog.Int("id", id), slog.String("error", err.Error()))
		switch err.(type) {
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

	h.logger.Info("Get cover art handler success", slog.Int("id", id))
	c.File(cover.Path)
}
