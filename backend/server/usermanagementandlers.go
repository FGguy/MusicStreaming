package server

import (
	"context"
	"encoding/xml"
	"fmt"
	sqlc "music-streaming/sql/sqlc"
	subsonic "music-streaming/util/subsonic"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// GET
func (s *Server) hangleGetUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		buildAndSendXMLError(c, "10")
		return
	}

	ctx := context.Background()
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

	subsonicRes := subsonic.SubsonicResponse{
		Xmlns:   subsonic.Xmlns,
		Status:  "ok",
		Version: subsonic.SubsonicVersion,
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
	ctx := context.Background()
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

	xmlUsers := make([]*subsonic.SubsonicUser, 0, len(users))
	for _, user := range users {
		xmlUsers = append(xmlUsers, mapSqlUserToXmlUser(&user))
	}

	subsonicRes := subsonic.SubsonicResponse{
		Xmlns:   subsonic.Xmlns,
		Status:  "ok",
		Version: subsonic.SubsonicVersion,
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
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	if username == "" || password == "" || email == "" {
		buildAndSendXMLError(c, "10")
		return
	}

	subsonicRes := subsonic.SubsonicResponse{
		Xmlns:   subsonic.Xmlns,
		Status:  "ok",
		Version: subsonic.SubsonicVersion,
	}

	userParams := make(map[string]string)

	userParams["username"] = fmt.Sprintf("\"%s\"", username)
	userParams["password"] = fmt.Sprintf("\"%s\"", password)
	userParams["email"] = fmt.Sprintf("\"%s\"", email)

	for _, userRole := range subsonic.SubsonicUserRoles {
		enabled := c.PostForm(userRole)
		if enabled == "true" || enabled == "false" {
			userParams[userRole] = enabled
		}
	}

	musicFolders := c.PostFormArray("musicFolderId")
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

	maxBitRate := c.PostForm("maxBitRate")
	if slices.Contains(subsonic.SubsonicValidBitRates, maxBitRate) {
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

	ctx := context.Background()
	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, createUserQueryString) //create all tables
	if err != nil {
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
func (s *Server) hangleUpdateUser(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		buildAndSendXMLError(c, "10")
		return
	}

	subsonicRes := subsonic.SubsonicResponse{
		Xmlns:   subsonic.Xmlns,
		Status:  "ok",
		Version: subsonic.SubsonicVersion,
	}

	userUpdates := make(map[string]string)
	for _, role := range subsonic.SubsonicUserRoles {
		roleUpdate := c.PostForm(role)
		if roleUpdate == "true" || roleUpdate == "false" {
			userUpdates[role] = roleUpdate
		}
	}
	//check for update to musicFolderID
	musicFolders := c.PostFormArray("musicFolderId")
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
		userUpdates["musicFolders"] = fmt.Sprintf("\"%s\"", strings.Join(ids, ";"))
	}

	maxBitRate := c.PostForm("maxBitRate")
	if slices.Contains(subsonic.SubsonicValidBitRates, maxBitRate) {
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
	updateUserQueryString := fmt.Sprintf("UPDATE Users SET %s WHERE username = %s;", updatesString, username)

	ctx := context.Background()
	conn, err := s.pg_pool.Acquire(ctx)
	if err != nil {
		buildAndSendXMLError(c, "0")
		return
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, updateUserQueryString) //create all tables
	if err != nil {
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
func (s *Server) hangleDeleteUser(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		buildAndSendXMLError(c, "10")
		return
	}
	ctx := context.Background()
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

	subsonicRes := subsonic.SubsonicResponse{
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

// POST
func (s *Server) handleChangePassword(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		buildAndSendXMLError(c, "10")
		return
	}
	ctx := context.Background()
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

	subsonicRes := subsonic.SubsonicResponse{
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

func mapSqlUserToXmlUser(user *sqlc.User) *subsonic.SubsonicUser {
	return &subsonic.SubsonicUser{
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
