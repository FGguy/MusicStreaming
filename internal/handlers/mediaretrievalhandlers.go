package controller

import (
	types "music-streaming/internal/types"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Application) handleDownload(c *gin.Context) {
	var (
		rUser   = c.MustGet("requestingUser").(*types.SubsonicUser)
		ctx     = c.Request.Context()
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

// struct for fields
type StreamParameters struct {
	Id                     string `form:"id" binding:"required"`
	MaxBitRate             int    `form:"maxBitRate"`
	Format                 string `form:"format"`
	EstimatedContentLength bool   `form:"estimateContentLength"`
}

func (s *Application) handleStream(c *gin.Context) {
	var (
		rUser = c.MustGet("requestingUser").(*types.SubsonicUser)
		ctx   = c.Request.Context()
	)

	var params StreamParameters
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Debug().Err(err)
		buildAndSendError(c, "10")
		return
	}

	if !rUser.AdminRole && !rUser.StreamRole {
		buildAndSendError(c, "50")
		return
	}

	if params.Id == "" {
		buildAndSendError(c, "10")
		return
	}

	id, err := strconv.Atoi(params.Id)
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

	log.Info().Msgf("Streaming song: %s, to user: %s", song.Title, rUser.Username)
	c.File(song.Path)
}

func (s *Application) handleGetCoverArt(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		id  = c.Query("id")
	)

	cover, err := s.dataLayer.GetCover(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("Failed fetching cover with id: %s", id)
		buildAndSendError(c, "0")
		return
	}

	c.File(cover.Path)
}
