package controller

import (
	"context"
	consts "music-streaming/consts"
	types "music-streaming/types"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Server) handleGetArtist(c *gin.Context) {
	ctx := context.Background()
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

func (s *Server) handleGetAlbum(c *gin.Context) {
	ctx := context.Background()
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

func (s *Server) handleGetSong(c *gin.Context) {
	ctx := context.Background()
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
