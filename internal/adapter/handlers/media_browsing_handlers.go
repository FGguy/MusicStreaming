package handlers

import (
	"music-streaming/internal/core/ports"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MediaBrowsingHandler struct {
	MediaBrowsingService ports.MediaBrowsingPort
}

func (h *MediaBrowsingHandler) handleGetArtist(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		paramId = c.Query("id")
	)

	id, err := strconv.Atoi(paramId)
	if paramId == "" || err != nil {
		buildAndSendError(c, "10")
		return
	}

	artist, err := h.MediaBrowsingService.GetArtist(ctx, id)
	if err != nil {
		switch err.(type) {
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

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
		buildAndSendError(c, "10")
		return
	}

	album, err := h.MediaBrowsingService.GetAlbum(ctx, id)
	if err != nil {
		switch err.(type) {
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

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
		buildAndSendError(c, "10")
		return
	}

	song, err := h.MediaBrowsingService.GetSong(ctx, id)
	if err != nil {
		switch err.(type) {
		case *ports.NotFoundError:
			buildAndSendError(c, "70")
		case *ports.FailedOperationError:
			buildAndSendError(c, "0")
		}
		return
	}

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
		Song:    &song,
	}

	SerializeAndSendBody(c, subsonicRes)
}
