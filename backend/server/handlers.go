package server

import (
	"encoding/xml"
	subxml "music-streaming/util/subxml"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlePing(c *gin.Context) {
	subsonicRes := subxml.SubsonicResponse{
		Xmlns:   subxml.Xmlns,
		Status:  "ok",
		Version: subxml.SubsonicVersion,
	}

	xmlBody, err := xml.Marshal(subsonicRes)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}
	c.Data(http.StatusOK, "application/xml", xmlBody)
}
