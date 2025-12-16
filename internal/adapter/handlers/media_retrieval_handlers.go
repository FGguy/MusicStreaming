package handlers

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type MediaRetrievalHandler struct {
	MediaRetrievalService ports.MediaRetrievalPort
}

func NewMediaRetrievalHandler(mediaRetrievalService ports.MediaRetrievalPort) *MediaRetrievalHandler {
	return &MediaRetrievalHandler{
		MediaRetrievalService: mediaRetrievalService,
	}
}

func (h *MediaRetrievalHandler) handleDownload(c *gin.Context) {
	var (
		rUser   = c.MustGet(RequestingUserKey).(*domain.User)
		ctx     = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
		paramId = c.Query("id")
	)

	id, err := strconv.Atoi(paramId)
	if paramId == "" || err != nil {
		buildAndSendError(c, "10")
		return
	}

	song, err := h.MediaRetrievalService.DownloadSong(ctx, id)
	if err != nil {
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
		log.Debug().Err(err)
		buildAndSendError(c, "10")
		return
	}

	id, err := strconv.Atoi(params.Id)
	if params.Id == "" || err != nil {
		buildAndSendError(c, "10")
		return
	}

	song, err := h.MediaRetrievalService.DownloadSong(ctx, id)
	if err != nil {
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
	}
	//Format?
	//Validate string
	//transcode file

	if params.EstimatedContentLength {
		//set header
	}

	c.File(song.Path)
}

func (h *MediaRetrievalHandler) handleGetCoverArt(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		paramId = c.Query("id")
	)

	id, err := strconv.Atoi(paramId)
	if paramId == "" || err != nil {
		buildAndSendError(c, "10")
		return
	}

	cover, err := h.MediaRetrievalService.GetCover(ctx, id)
	if err != nil {
		switch err.(type) {
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

	c.File(cover.Path)
}
