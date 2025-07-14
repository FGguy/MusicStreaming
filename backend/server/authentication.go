package server

import "github.com/gin-gonic/gin"

/*
Authentication middlware for:
	- Authentication
	- Verifying user got permission for the current endpoint
*/

func (s *Server) WithAuth(c *gin.Context) {
	//check if user credentials is in cache
	//if not check if user credential is in db
	//if in db push to cache
	//authenticate
}
