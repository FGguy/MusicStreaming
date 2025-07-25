package server

import (
	"context"
	"encoding/json"
	"fmt"
	consts "music-streaming/consts"
	sqlc "music-streaming/sql/sqlc"
	types "music-streaming/types"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// GET
func (s *Server) hangleGetUser(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	contentType := c.MustGet("contentType").(string)
	username := c.Query("username")
	if username == "" {
		buildAndSendError(c, "10")
		return
	}

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result()
	if err != nil {
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil {
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
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

	SerializeAndSendBody(c, subsonicRes)
}

// GET
func (s *Server) hangleGetUsers(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	contentType := c.MustGet("contentType").(string)
	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
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

	sqlUsers, err := q.GetUsers(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}

	Users := make([]*types.SubsonicUser, 0, len(sqlUsers))
	for _, user := range sqlUsers {
		Users = append(Users, types.MapSqlUserToSubsonicUser(&user, ""))
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		Users:   Users,
	}

	SerializeAndSendBody(c, subsonicRes)
}

// POST
func (s *Server) handleCreateUser(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	contentType := c.MustGet("contentType").(string)
	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	err = json.Unmarshal([]byte(userString), &cachedUser)
	if err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}

	if !cachedUser.AdminRole {
		buildAndSendError(c, "50")
		return
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
	}

	userParams := make(map[string]any)
	userParams["username"] = c.Query("username")
	userParams["password"] = c.Query("password")
	userParams["email"] = c.Query("email")

	if userParams["username"] == "" || userParams["password"] == "" || userParams["email"] == "" {
		buildAndSendError(c, "10")
		return
	}

	for _, userRole := range consts.SubsonicUserRoles {
		enabled := c.Query(userRole)
		if enabled == "true" || enabled == "false" {
			userParams[strings.ToLower(userRole)] = enabled
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
		userParams["musicfolders"] = fmt.Sprintf("\"%s\"", strings.Join(ids, ";"))
	}

	maxBitRate := c.Query("maxbitrate")
	if slices.Contains(consts.SubsonicValidBitRates, maxBitRate) {
		userParams["maxbitrate"] = maxBitRate
	}

	insertSQL := goqu.Insert("users").Rows(userParams)
	queryString, _, err := insertSQL.ToSQL()
	if err != nil {
		debugLog("Failed to generate query string.")
		buildAndSendError(c, "0")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}
	defer conn.Release()

	if _, err = conn.Exec(ctx, queryString); err != nil {
		debugLogError("Failed creating user", err)
		buildAndSendError(c, "0")
		return
	}

	SerializeAndSendBody(c, subsonicRes)
}

// POST
func (s *Server) handleUpdateUser(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	contentType := c.MustGet("contentType").(string)
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
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
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
			userUpdates[strings.ToLower(role)] = roleUpdate
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
		userUpdates["musicfolders"] = fmt.Sprintf("'%s'", strings.Join(ids, ";"))
	}

	maxBitRate := c.Query("maxBitRate")
	if slices.Contains(consts.SubsonicValidBitRates, maxBitRate) {
		userUpdates["maxbitrate"] = maxBitRate
	}

	//if no valid updates abort
	if len(userUpdates) < 1 {
		SerializeAndSendBody(c, subsonicRes)
		return
	}

	sqlUpdate := goqu.Update("users").
		Set(userUpdates).
		Where(goqu.Ex{"username": username})

	queryString, _, err := sqlUpdate.ToSQL()
	if err != nil {
		debugLog("Failed to generate query string.")
		buildAndSendError(c, "0")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendError(c, "0")
		return
	}
	defer conn.Release()

	if _, err = conn.Exec(ctx, queryString); err != nil {
		debugLogError("Failed updating user", err)
		buildAndSendError(c, "0")
		return
	}

	SerializeAndSendBody(c, subsonicRes)
}

// POST
func (s *Server) handleDeleteUser(c *gin.Context) {
	params := c.MustGet("requiredParams").(requiredParams)
	contentType := c.MustGet("contentType").(string)
	ctx := context.Background()

	userString, err := s.cache.Get(ctx, params.U).Result() //bug
	if err != nil {                                        //if user is authenticated their info should be cached
		debugLogError("Failed fetching user credentials from cache", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
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

	SerializeAndSendBody(c, subsonicRes)
}

// POST
func (s *Server) handleChangePassword(c *gin.Context) {
	contentType := c.MustGet("contentType").(string)
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
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}

	var cachedUser types.SubsonicUser
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil { //if user is authenticated their info should be cached
		debugLogError("Failed unmarshalling user", err)
		c.Data(http.StatusInternalServerError, contentType, []byte{})
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

	SerializeAndSendBody(c, subsonicRes)
}
