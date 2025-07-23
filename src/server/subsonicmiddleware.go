package server

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	sqlc "music-streaming/sql/sqlc"

	auth "music-streaming/util/auth"
	subsonic "music-streaming/util/subsonic"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) subValidateQParamsMiddleware(c *gin.Context) {
	/*
		Do not care about 'p' or 'f' query params
		password needs to be hashed
		only supported format will be xml
	*/

	//FIXME: use custom struct to bind with parameters
	params := []string{"u", "t", "s", "v", "c"}
	for _, param := range params {
		if c.Query(param) == "" {
			if gin.Mode() == gin.DebugMode {
				log.Printf("Missing required parameter.")
			}
			buildAndSendXMLError(c, "10")
			return
		} else {
			c.Set(param, c.Query(param))
		}
	}

	//validate Subsonic API compatibility with client
	clientVersion := strings.Split(c.Query("v"), ".")
	clientMajorVersion, err := strconv.Atoi(clientVersion[0])
	if err != nil {
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed converting subsonic client major version into int, Err: %s", err)
		}
		buildAndSendXMLError(c, "0")
		return
	}
	clientMinorVersion, err := strconv.Atoi(clientVersion[1])
	if err != nil {
		if gin.Mode() == gin.DebugMode {
			log.Printf("Failed converting subsonic client minor version into int, Err: %s", err)
		}
		buildAndSendXMLError(c, "0")
		return
	}

	if clientMajorVersion > subsonic.SubsonicMajorVersion {
		buildAndSendXMLError(c, "30")
		return
	} else if clientMajorVersion < subsonic.SubsonicMajorVersion {
		buildAndSendXMLError(c, "20")
		return
	}

	if clientMinorVersion > subsonic.SubsonicMinorVersion {
		buildAndSendXMLError(c, "30")
		return
	}
}

func (s *Server) subWithAuth(c *gin.Context) {

	//Shouldn't panic, if it does its kinda cooked
	qUser := c.MustGet("u").(string)
	qHashedPassword := c.MustGet("t").(string)
	qSalt := c.MustGet("s").(string)

	var password string
	var cachedUser subsonic.SubsonicRedisUser
	ctx := context.Background()
	err := s.cache.Get(ctx, qUser).Scan(&cachedUser)
	if err != nil {
		conn, err := s.pg_pool.Acquire(ctx)
		if err != nil {
			if gin.Mode() == gin.DebugMode {
				log.Printf("Failed acquiring connection from postgres connection pool, Err: %s", err)
			}
			buildAndSendXMLError(c, "0")
			return
		}
		defer conn.Release()
		query := sqlc.New(conn)

		user, err := query.GetUserByUsername(ctx, pgtype.Text{String: qUser, Valid: true})
		if err != nil {
			if gin.Mode() == gin.DebugMode {
				log.Printf("User does not exist. Err: %s", err)
			}
			buildAndSendXMLError(c, "40") //user doesnt exist
			return
		}

		encodedUser, err := json.Marshal(mapSqlUserToRedisUser(&user))
		if err != nil {
			if gin.Mode() == gin.DebugMode {
				log.Printf("Failed encoding user credentials, Err: %s", err)
			}
			buildAndSendXMLError(c, "0")
			return
		}
		err = s.cache.Set(ctx, user.Username.String, encodedUser, time.Minute*10).Err()
		if err != nil {
			if gin.Mode() == gin.DebugMode {
				log.Printf("Failed creating cache entry for user credentials, Err: %s", err)
			}
			buildAndSendXMLError(c, "0")
			return
		}
		password = user.Password
	} else {
		password = cachedUser.Password
	}

	//TODO: Change to support permissions
	if !auth.ValidatePassword(qHashedPassword, qSalt, password) {
		if gin.Mode() == gin.DebugMode {
			log.Printf("Wrong password, Err: %s", err)
		}
		buildAndSendXMLError(c, "40") //password incorrect
		return
	}
}

// Util
func buildAndSendXMLError(c *gin.Context, errorCode string) {
	c.Abort()
	subsonicRes := subsonic.SubsonicXmlResponse{
		Xmlns:   subsonic.Xmlns,
		Status:  "failed",
		Version: subsonic.SubsonicVersion,
	}

	subsonicRes.Error = &subsonic.SubsonicXmlError{
		Code:    errorCode,
		Message: subsonic.SubsonicErrorMessages[errorCode],
	}

	xmlBody, err := xml.Marshal(subsonicRes)

	//temp fix, disgusting
	temp := strings.Replace(string(xmlBody), "></error>", "/>", 1)

	if err != nil {
		c.Data(http.StatusInternalServerError, "application/xml", []byte{})
		return
	}
	c.Data(http.StatusOK, "application/xml", []byte(temp))
}

func mapSqlUserToRedisUser(user *sqlc.User) *subsonic.SubsonicRedisUser {
	return &subsonic.SubsonicRedisUser{
		Username:            user.Username.String,
		Email:               user.Email,
		Password:            user.Password,
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
