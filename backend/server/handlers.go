package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlePing(c *gin.Context) {
	c.XML(200, gin.H{
		"message": "ping successful",
		"status":  http.StatusOK,
	})
}
