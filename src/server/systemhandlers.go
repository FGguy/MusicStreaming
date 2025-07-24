package server

import (
	"encoding/xml"
	consts "music-streaming/consts"
	types "music-streaming/types"

	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlePing(c *gin.Context) {
	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	xmlBody, err := xml.Marshal(subsonicRes)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}
	c.Data(http.StatusOK, "application/xml", xmlBody)
}
