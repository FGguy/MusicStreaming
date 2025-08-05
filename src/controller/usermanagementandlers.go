package controller

import (
	"context"
	"fmt"
	consts "music-streaming/consts"
	types "music-streaming/types"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GET
func (s *Server) hangleGetUser(c *gin.Context) {
	var (
		rUser    = c.MustGet("requestingUser").(*types.SubsonicUser)
		username = c.Query("username")
		ctx      = context.Background()
	)

	if username == "" {
		buildAndSendError(c, "10")
		return
	}

	if !rUser.AdminRole && rUser.Username != username {
		buildAndSendError(c, "50")
		return
	}

	user, err := s.dataLayer.GetUser(ctx, username)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}

	user.Password = "" //do not send user passwords

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		User:    user,
	}

	SerializeAndSendBody(c, subsonicRes)
}

// GET
func (s *Server) hangleGetUsers(c *gin.Context) {
	rUser := c.MustGet("requestingUser").(*types.SubsonicUser)
	ctx := context.Background()

	if !rUser.AdminRole {
		buildAndSendError(c, "50")
		return
	}

	users, err := s.dataLayer.GetUsers(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		Users:   users,
	}

	SerializeAndSendBody(c, subsonicRes)
}

// POST
func (s *Server) handleCreateUser(c *gin.Context) {
	rUser := c.MustGet("requestingUser").(*types.SubsonicUser)
	ctx := context.Background()

	if !rUser.AdminRole {
		buildAndSendError(c, "50")
		return
	}

	userParams := make(map[string]any)
	userParams["username"] = c.PostForm("username")
	userParams["password"] = c.PostForm("password")
	userParams["email"] = c.PostForm("email")

	if userParams["username"] == "" || userParams["password"] == "" || userParams["email"] == "" {
		buildAndSendError(c, "10")
		return
	}

	for _, userRole := range consts.SubsonicUserRoles {
		enabled := c.PostForm(userRole)
		if enabled == "true" || enabled == "false" {
			userParams[strings.ToLower(userRole)] = enabled
		}
	}

	musicFolders := c.PostFormArray("musicFolderId")
	if len(musicFolders) > 0 {
		ids := make([]string, 0, len(musicFolders))
		for _, folderId := range musicFolders {
			id, err := strconv.Atoi(folderId)
			if err != nil || id < 1 { //invalid folder id passed as param
				buildAndSendError(c, "0")
				return
			}
			ids = append(ids, folderId)
		}
		//none of the ids are invalid
		userParams["musicfolders"] = fmt.Sprintf("\"%s\"", strings.Join(ids, ";"))
	}

	maxBitRate := c.PostForm("maxbitrate")
	if slices.Contains(consts.SubsonicValidBitRates, maxBitRate) {
		userParams["maxbitrate"] = maxBitRate
	}

	if err := s.dataLayer.CreateUser(ctx, userParams); err != nil {
		log.Error().Err(err).Msgf("Failed to create user")
		buildAndSendError(c, "0")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}

// POST
func (s *Server) handleUpdateUser(c *gin.Context) {
	var (
		rUser    = c.MustGet("requestingUser").(*types.SubsonicUser)
		username = c.PostForm("username")
		ctx      = context.Background()
	)

	if username == "" {
		log.Error().Msg("Failed getting user from params")
		buildAndSendError(c, "10")
		return
	}

	if !rUser.AdminRole && !(rUser.Username == username && rUser.SettingsRole) {
		buildAndSendError(c, "50")
		return
	}

	userUpdates := make(map[string]string)
	for _, role := range consts.SubsonicUserRoles {
		roleUpdate := c.PostForm(role)
		if roleUpdate == "true" || roleUpdate == "false" {
			userUpdates[strings.ToLower(role)] = roleUpdate
		}
	}
	//check for update to musicFolderID
	musicFolders := c.PostFormArray("musicFolderId")
	if len(musicFolders) > 0 {
		ids := make([]string, 0, len(musicFolders))
		for _, folderId := range musicFolders {
			id, err := strconv.Atoi(folderId)
			if err != nil || id < 1 { //invalid folder id passed as param
				buildAndSendError(c, "0")
				return
			}
			ids = append(ids, folderId)
		}
		//none of the ids are invalid
		userUpdates["musicfolders"] = fmt.Sprintf("'%s'", strings.Join(ids, ";"))
	}

	maxBitRate := c.PostForm("maxBitRate")
	if slices.Contains(consts.SubsonicValidBitRates, maxBitRate) {
		userUpdates["maxbitrate"] = maxBitRate
	}

	//if no valid updates abort
	if len(userUpdates) >= 1 {
		if err := s.dataLayer.UpdateUser(ctx, username, userUpdates); err != nil {
			log.Error().Err(err).Msgf("Failed to update user")
			buildAndSendError(c, "0")
			return
		}
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}

// POST
func (s *Server) handleDeleteUser(c *gin.Context) {
	var (
		rUser    = c.MustGet("requestingUser").(*types.SubsonicUser)
		username = c.PostForm("username")
		ctx      = context.Background()
	)

	if !rUser.AdminRole {
		buildAndSendError(c, "50")
		return
	}

	if username == "" {
		log.Warn().Msgf("Failed to get username from url-encoded post form parameters")
		buildAndSendError(c, "10")
		return
	}

	if err := s.dataLayer.DeleteUser(ctx, username); err != nil {
		buildAndSendError(c, "0")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}

// POST
func (s *Server) handleChangePassword(c *gin.Context) {
	var (
		rUser    = c.MustGet("requestingUser").(*types.SubsonicUser)
		username = c.PostForm("username")
		ctx      = context.Background()
	)

	if username == "" {
		buildAndSendError(c, "10")
		return
	}

	if !rUser.AdminRole && (rUser.Username != username) {
		buildAndSendError(c, "50")
		return
	}

	password := c.PostForm("password")
	if username == "" || password == "" {
		buildAndSendError(c, "10")
		return
	}

	if err := s.dataLayer.ChangeUserPassword(ctx, username, password); err != nil {
		buildAndSendError(c, "0")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}
