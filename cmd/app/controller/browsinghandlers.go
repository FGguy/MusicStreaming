package controller

import (
	consts "music-streaming/internal/consts"
	types "music-streaming/internal/types"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Application) handleGetArtist(c *gin.Context) {
	ctx := c.Request.Context()
	paramId := c.Query("id")
	if paramId == "" {
		buildAndSendError(c, "10")
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Warn().Err(err).Msgf("Invalid artist id")
		buildAndSendError(c, "0")
		return
	}

	artist, err := s.dataLayer.GetArtist(ctx, int32(id))
	if err != nil {
		log.Warn().Err(err).Msgf("Failed fetching artist with id: %d", id)
		buildAndSendError(c, "0")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		Artist:  artist,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (s *Application) handleGetAlbum(c *gin.Context) {
	ctx := c.Request.Context()
	paramId := c.Query("id")
	if paramId == "" {
		buildAndSendError(c, "10")
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Warn().Err(err).Msgf("Invalid album id")
		buildAndSendError(c, "0")
		return
	}

	album, err := s.dataLayer.GetAlbum(ctx, int32(id))
	if err != nil {
		log.Warn().Err(err).Msgf("Failed fetching album with id: %d", id)
		buildAndSendError(c, "0")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		Album:   album,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (s *Application) handleGetSong(c *gin.Context) {
	ctx := c.Request.Context()
	paramId := c.Query("id")
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

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		Song:    song,
	}

	SerializeAndSendBody(c, subsonicRes)
}
