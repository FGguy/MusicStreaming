package server

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	consts "music-streaming/consts"
	sqlc "music-streaming/sql/sqlc"
	types "music-streaming/types"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// GET
func (s *Server) hangleGetUser(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	username := c.Query("username")
	if username == "" {
		buildAndSendError(c, "10")
		return
	}

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole && params.U != username {
		buildAndSendError(c, "50")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}
	defer conn.Release()
	q := sqlc.New(conn)

	user, err := q.GetUserByUsername(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		buildAndSendError(c, "0")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		User:    types.MapSqlUserToSubsonicUser(&user, ""),
	}

	//build xml body for answer
	xmlBody, err := xml.Marshal(subsonicRes)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}
	c.Data(http.StatusOK, "application/xml", xmlBody)
}

// GET
func (s *Server) hangleGetUsers(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole {
		buildAndSendError(c, "50")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}
	defer conn.Release()
	q := sqlc.New(conn)

	users, err := q.GetUsers(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}

	xmlUsers := make([]*types.SubsonicUser, 0, len(users))
	for _, user := range users {
		xmlUsers = append(xmlUsers, types.MapSqlUserToSubsonicUser(&user, ""))
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		Users:   xmlUsers,
	}

	xmlBody, err := xml.Marshal(subsonicRes)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}
	c.Data(http.StatusOK, "application/xml", xmlBody)
}

// POST
func (s *Server) handleCreateUser(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	err = json.Unmarshal([]byte(userString), &cachedUser)
	if err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole {
		buildAndSendError(c, "50")
		return
	}

	var (
		username = c.Query("username")
		password = c.Query("password")
		email    = c.Query("email")
	)
	if username == "" || password == "" || email == "" {
		buildAndSendError(c, "10")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	userParams := make(map[string]string)
	userParams["username"] = fmt.Sprintf("'%s'", username)
	userParams["password"] = fmt.Sprintf("'%s'", password)
	userParams["email"] = fmt.Sprintf("'%s'", email)

	for _, userRole := range consts.SubsonicUserRoles {
		enabled := c.Query(userRole)
		if enabled == "true" || enabled == "false" {
			userParams[userRole] = enabled
		}
	}

	musicFolders := c.QueryArray("musicFolderId")
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
		userParams["musicFolders"] = fmt.Sprintf("\"%s\"", strings.Join(ids, ";"))
	}

	maxBitRate := c.Query("maxBitRate")
	if slices.Contains(consts.SubsonicValidBitRates, maxBitRate) {
		userParams["maxBitRate"] = maxBitRate
	}

	qParams := make([]string, 0, len(userParams))
	values := make([]string, 0, len(userParams))
	for param, value := range userParams {
		qParams = append(qParams, param)
		values = append(values, value)
	}
	paramsString := strings.Join(qParams, ", ")
	valuesString := strings.Join(values, ", ")
	createUserQueryString := fmt.Sprintf("INSERT INTO Users (%s) VALUES (%s) ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username RETURNING *;", paramsString, valuesString)

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}
	defer conn.Release()

	if _, err = conn.Exec(ctx, createUserQueryString); err != nil {
		debugLogError("Failed creating user", err)
		buildAndSendError(c, "0")
		return
	}

	xmlBody, err := xml.Marshal(subsonicRes)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}
	c.Data(http.StatusOK, "application/xml", xmlBody)
}

// POST
func (s *Server) handleUpdateUser(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	username := c.Query("username")
	if username == "" {
		debugLog("Failed getting user from params")
		buildAndSendError(c, "10")
		return
	}

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole && (cachedUser.Username != username || !cachedUser.SettingsRole) {
		buildAndSendError(c, "50")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	userUpdates := make(map[string]string)
	for _, role := range consts.SubsonicUserRoles {
		roleUpdate := c.Query(role)
		if roleUpdate == "true" || roleUpdate == "false" {
			userUpdates[role] = roleUpdate
		}
	}
	//check for update to musicFolderID
	musicFolders := c.QueryArray("musicFolderId")
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
		userUpdates["musicFolders"] = fmt.Sprintf("'%s'", strings.Join(ids, ";"))
	}

	maxBitRate := c.Query("maxBitRate")
	if slices.Contains(consts.SubsonicValidBitRates, maxBitRate) {
		userUpdates["maxBitRate"] = maxBitRate
	}

	//if no valid updates abort
	if len(userUpdates) < 1 {
		xmlBody, err := xml.Marshal(subsonicRes)
		if err != nil {
			c.Data(http.StatusInternalServerError, "application/xml", []byte{})
			return
		}
		c.Data(http.StatusOK, "application/xml", xmlBody)
		return
	}

	updates := make([]string, 0, len(userUpdates))
	for role, update := range userUpdates {
		updates = append(updates, fmt.Sprintf("%s = %s", role, update))
	}
	updatesString := strings.Join(updates, ",")
	updateUserQueryString := fmt.Sprintf("UPDATE Users SET %s WHERE username = '%s';", updatesString, username)

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}
	defer conn.Release()

	if _, err = conn.Exec(ctx, updateUserQueryString); err != nil {
		debugLogError("Failed updating user", err)
		buildAndSendError(c, "0")
		return
	}

	xmlBody, err := xml.Marshal(subsonicRes)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}
	c.Data(http.StatusOK, "application/xml", xmlBody)
}

// POST
func (s *Server) handleDeleteUser(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole {
		buildAndSendError(c, "50")
		return
	}

	username := c.Query("username")
	if username == "" {
		debugLog("Failed to get username from url-encoded post form parameters")
		buildAndSendError(c, "10")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}
	defer conn.Release()
	q := sqlc.New(conn)

	if _, err = q.DeleteUser(ctx, pgtype.Text{String: username, Valid: true}); err != nil {
		buildAndSendError(c, "0")
		return
	}

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

// POST
func (s *Server) handleChangePassword(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)

	username := c.Query("username")
	if username == "" {
		buildAndSendError(c, "10")
		return
	}

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole && (cachedUser.Username != username) {
		buildAndSendError(c, "50")
		return
	}

	password := c.Query("password")
	if username == "" || password == "" {
		buildAndSendError(c, "10")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}
	defer conn.Release()
	q := sqlc.New(conn)

	if _, err = q.ChangeUserPassword(ctx, sqlc.ChangeUserPasswordParams{Username: pgtype.Text{String: username, Valid: true}, Password: password}); err != nil {
		buildAndSendError(c, "0")
		return
	}

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
