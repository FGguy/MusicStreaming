package controller

import (
	"context"
	types "music-streaming/internal/types"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Application) handleDownload(c *gin.Context) {
	var (
		rUser   = c.MustGet("requestingUser").(*types.SubsonicUser)
		ctx     = context.Background()
		paramId = c.Query("id")
	)

	if !rUser.AdminRole && !rUser.DownloadRole {
		buildAndSendError(c, "50")
		return
	}

	if paramId == "" {
		buildAndSendError(c, "10")
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Warn().Err(err).Msgf("Invalid song id")
		buildAndSendError(c, "0")
		return
	}

	song, err := s.dataLayer.GetSong(ctx, int32(id))
	if err != nil {
		log.Warn().Err(err).Msgf("Failed fetching song with id: %d", id)
		buildAndSendError(c, "0")
		return
	}

	c.FileAttachment(song.Path, song.Title)
}

func (s *Application) handleStream(c *gin.Context) {
	var (
		rUser   = c.MustGet("requestingUser").(*types.SubsonicUser)
		ctx     = context.Background()
		paramId = c.Query("id")
	)

	if !rUser.AdminRole && !rUser.StreamRole {
		buildAndSendError(c, "50")
		return
	}

	if paramId == "" {
		buildAndSendError(c, "10")
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Warn().Err(err).Msgf("Invalid song id")
		buildAndSendError(c, "0")
		return
	}

	song, err := s.dataLayer.GetSong(ctx, int32(id))
	if err != nil {
		log.Warn().Err(err).Msgf("Failed fetching song with id: %d", id)
		buildAndSendError(c, "0")
		return
	}

	c.File(song.Path)
}
