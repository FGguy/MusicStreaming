package server

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"time"

	sqlc "music-streaming/sql/sqlc"

	auth "music-streaming/util/auth"
	subxml "music-streaming/util/subxml"

	"github.com/gin-gonic/gin"
)

/*
Authentication middlware for:
	- Authentication
	- Verifying user got permission for the current endpoint
*/

func (s *Server) WithAuth(c *gin.Context) {
	subsonicRes := subxml.SubsonicResponse{
		Xmlns:   subxml.Xmlns,
		Status:  "failed",
		Version: subxml.SubsonicVersion,
	}

	buildAndSendXMLError := func(errorCode string) {
		subsonicRes.Error = &subxml.SubsonicError{
			Code:    errorCode,
			Message: subxml.SubsonicErrorMessages[errorCode],
		}
		xmlBody, err := xml.Marshal(subsonicRes)
		if err != nil {
			c.Data(http.StatusInternalServerError, "application/xml", []byte{})
			return
		}
		c.Data(http.StatusOK, "application/xml", xmlBody)
		c.Abort()
	}

	qUser := c.Query("u")
	qHashedPassword := c.Query("p")
	qSalt := c.Query("s")
	if qUser == "" || qHashedPassword == "" || qSalt == "" {
		buildAndSendXMLError("10")
		return
	}

	var password string
	ctx := context.Background()
	password, err := s.cache.Get(ctx, qUser).Result()
	if err != nil {
		conn, err := s.pg_pool.Acquire(ctx)
		if err != nil {
			if gin.Mode() == gin.DebugMode {
				log.Printf("Failed acquiring connection from postgres connection pool, Err: %s", err)
			}
			buildAndSendXMLError("0")
			return
		}
		defer conn.Release()
		query := sqlc.New(conn)

		user, err := query.GetUserByName(ctx, qUser)
		if err != nil {
			buildAndSendXMLError("40") //user doesnt exist
			return
		}

		err = s.cache.Set(ctx, user.Name, user.Password, time.Minute*10).Err()
		if err != nil {
			if gin.Mode() == gin.DebugMode {
				log.Printf("Failed creating cache entry for user credentials, Err: %s", err)
			}
			buildAndSendXMLError("0")
			return
		}
		password = user.Password
	}

	//TODO: Change to support permissions
	if !auth.ValidatePassword(qHashedPassword, qSalt, password) {
		buildAndSendXMLError("40") //password incorrect
		return
	}
}
