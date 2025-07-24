package server

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
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
	u := c.MustGet("u").(string)

	username := c.Query("username")
	if username == "" {
		buildAndSendXMLError(c, "10")
		return
	}

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, u).Result() //bug
	if err != nil {                                 //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed fetching user credentials from cache, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicRedisUser
	err = json.Unmarshal([]byte(userString), &cachedUser)
	if err != nil { //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed unmarshalling user, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole && u != username {
		buildAndSendXMLError(c, "50")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}
	defer conn.Release()
	q := sqlc.New(conn)

	user, err := q.GetUserByUsername(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}

	subsonicRes := types.SubsonicXmlResponse{
		Xmlns:   consts.Xmlns,
		Status:  "ok",
		Version: consts.SubsonicVersion,
		User:    mapSqlUserToXmlUser(&user),
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
	u := c.MustGet("u").(string)

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, u).Result() //bug
	if err != nil {                                 //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed fetching user credentials from cache, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicRedisUser
	err = json.Unmarshal([]byte(userString), &cachedUser)
	if err != nil { //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed unmarshalling user, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole {
		buildAndSendXMLError(c, "50")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}
	defer conn.Release()
	q := sqlc.New(conn)

	users, err := q.GetUsers(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}

	xmlUsers := make([]*types.SubsonicXmlUser, 0, len(users))
	for _, user := range users {
		xmlUsers = append(xmlUsers, mapSqlUserToXmlUser(&user))
	}

	subsonicRes := types.SubsonicXmlResponse{
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
	u := c.MustGet("u").(string)

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, u).Result() //bug
	if err != nil {                                 //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed fetching user credentials from cache, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicRedisUser
	err = json.Unmarshal([]byte(userString), &cachedUser)
	if err != nil { //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed unmarshalling user, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole {
		buildAndSendXMLError(c, "50")
		return
	}

	username := c.Query("username")
	password := c.Query("password")
	email := c.Query("email")
	if username == "" || password == "" || email == "" {
		buildAndSendXMLError(c, "10")
		return
	}

	subsonicRes := types.SubsonicXmlResponse{
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
				buildAndSendXMLError(c, "0")
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

	params := make([]string, 0, len(userParams))
	values := make([]string, 0, len(userParams))
	for param, value := range userParams {
		params = append(params, param)
		values = append(values, value)
	}
	paramsString := strings.Join(params, ", ")
	valuesString := strings.Join(values, ", ")
	createUserQueryString := fmt.Sprintf("INSERT INTO Users (%s) VALUES (%s) ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username RETURNING *;", paramsString, valuesString)
	if gin.Mode() == gin.DebugMode {
		log.Printf("Query String: %s", createUserQueryString)
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, createUserQueryString) //create all tables
	if err != nil {
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed creating user Err: %s", err)
		}
		buildAndSendXMLError(c, "0")
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
	u := c.MustGet("u").(string)

	username := c.Query("username")
	if username == "" {
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed getting user from params")
		}
		buildAndSendXMLError(c, "10")
		return
	}

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, u).Result() //bug
	if err != nil {                                 //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed fetching user credentials from cache, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicRedisUser
	err = json.Unmarshal([]byte(userString), &cachedUser)
	if err != nil { //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed unmarshalling user, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole && (cachedUser.Username != username || !cachedUser.SettingsRole) {
		buildAndSendXMLError(c, "50")
		return
	}

	subsonicRes := types.SubsonicXmlResponse{
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
				buildAndSendXMLError(c, "0")
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
	if gin.Mode() == gin.DebugMode {
		log.Printf("Query String: %s", updateUserQueryString)
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, updateUserQueryString) //create all tables
	if err != nil {
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed updating user. Error: %s", err)
		}
		buildAndSendXMLError(c, "0")
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
	u := c.MustGet("u").(string)

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, u).Result() //bug
	if err != nil {                                 //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed fetching user credentials from cache, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicRedisUser
	err = json.Unmarshal([]byte(userString), &cachedUser)
	if err != nil { //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed unmarshalling user, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole {
		buildAndSendXMLError(c, "50")
		return
	}

	username := c.Query("username")
	if username == "" {
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed to get username from url-encoded post form parameters")
		}
		buildAndSendXMLError(c, "10")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}
	defer conn.Release()
	q := sqlc.New(conn)

	_, err = q.DeleteUser(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}

	subsonicRes := types.SubsonicXmlResponse{
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
	u := c.MustGet("u").(string)

	username := c.Query("username")
	if username == "" {
		buildAndSendXMLError(c, "10")
		return
	}

	ctx := context.Background()

	userString, err := s.cache.Get(ctx, u).Result() //bug
	if err != nil {                                 //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed fetching user credentials from cache, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	var cachedUser types.SubsonicRedisUser
	err = json.Unmarshal([]byte(userString), &cachedUser)
	if err != nil { //if user is authenticated their info should be cached
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed unmarshalling user, Err: %s", err)
		}
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}

	if !cachedUser.AdminRole && (cachedUser.Username != username) {
		buildAndSendXMLError(c, "50")
		return
	}

	password := c.Query("password")
	if username == "" || password == "" {
		buildAndSendXMLError(c, "10")
		return
	}

	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}
	defer conn.Release()
	q := sqlc.New(conn)

	_, err = q.ChangeUserPassword(ctx, sqlc.ChangeUserPasswordParams{Username: pgtype.Text{String: username, Valid: true}, Password: password})
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}

	subsonicRes := types.SubsonicXmlResponse{
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

func mapSqlUserToXmlUser(user *sqlc.User) *types.SubsonicXmlUser {
	return &types.SubsonicXmlUser{
		Username:            user.Username.String,
		Email:               user.Email,
		ScrobblingEnabled:   user.Scrobblingenabled,
		LdapAuthenticated:   user.Ldapauthenticated,
		AdminRole:           user.Adminrole,
		SettingsRole:        user.Settingsrole,
		StreamRole:          user.Streamrole,
		JukeboxRole:         user.Jukeboxrole,
		DownloadRole:        user.Downloadrole,
		UploadRole:          user.Uploadrole,
		PlaylistRole:        user.Playlistrole,
		CoverArtRole:        user.Coverartrole,
		CommentRole:         user.Commentrole,
		PodcastRole:         user.Podcastrole,
		ShareRole:           user.Sharerole,
		VideoConversionRole: user.Videoconversionrole,
		MusicfolderId:       strings.Split(user.Musicfolderid.String, ";"),
		MaxBitRate:          user.Maxbitrate,
	}
}
