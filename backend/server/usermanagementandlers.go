package server

import (
	"context"
	"encoding/xml"
	"fmt"
	"music-streaming/util/subxml"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) hangleGetUser(c *gin.Context) {

}

func (s *Server) hangleGetUsers(c *gin.Context) {

}

func (s *Server) handleCreateUser(c *gin.Context) {

}

func (s *Server) hangleUpdateUser(c *gin.Context) {
	username := c.MustGet("u").(string)
	queryRoleUpdates := make(map[string]string)
	for _, role := range subxml.SubsonicUserRoles {
		roleUpdate := c.Query(role)
		//needs to be fixed for musicFolderId
		if roleUpdate == "true" || roleUpdate == "false" {
			queryRoleUpdates[role] = roleUpdate
		}
	}

	//if no valid updates abort
	if len(queryRoleUpdates) < 1 {
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
		return
	}

	updates := make([]string, 0, len(queryRoleUpdates))
	for role, update := range queryRoleUpdates {
		updates = append(updates, fmt.Sprintf("%s = %s", role, update))
	}
	updatesString := strings.Join(updates, ",")
	updateUserQueryString := fmt.Sprintf("UPDATE Users SET %s WHERE username = %s;", updatesString, username)

	ctx := context.Background()
	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, updateUserQueryString) //create all tables
	if err != nil {
		buildAndSendXMLError(c, "0")
	}
}

func (s *Server) hangleDeleteUser(c *gin.Context) {

}

func (s *Server) handleChangePassword(c *gin.Context) {

}
