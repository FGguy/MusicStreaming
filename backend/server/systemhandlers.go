package server

import (
	"encoding/xml"
	subsonic "music-streaming/util/subsonic"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlePing(c *gin.Context) {
	subsonicRes := subsonic.SubsonicXmlResponse{
		Xmlns:   subsonic.Xmlns,
		Status:  "ok",
		Version: subsonic.SubsonicVersion,
	}

	xmlBody, err := xml.Marshal(subsonicRes)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}
	c.Data(http.StatusOK, "application/xml", xmlBody)
}
