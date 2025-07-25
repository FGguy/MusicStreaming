package server

import (
	consts "music-streaming/consts"
	types "music-streaming/types"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlePing(c *gin.Context) {
	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}
