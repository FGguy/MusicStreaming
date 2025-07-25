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

	consts "music-streaming/consts"
	types "music-streaming/types"
	auth "music-streaming/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type requiredParams struct {
	U string `form:"u" binding:"required"`
	T string `form:"t" binding:"required"`
	S string `form:"s" binding:"required"`
	V string `form:"v" binding:"required"`
	C string `form:"c" binding:"required"`
	F string `form:"f"`
	P string `form:"p"`
}

func (s *Server) subValidateQParamsMiddleware(c *gin.Context) {
	var params requiredParams
	if err := c.ShouldBindQuery(&params); err != nil {
		buildAndSendError(c, "10")
		return
	}
	c.Set("requiredParams", params)

	//validate Subsonic API compatibility with client
	clientVersion := strings.Split(params.V, ".")
	clientMajorVersion, err := strconv.Atoi(clientVersion[0])
	if err != nil {
		debugLogError("Failed converting subsonic client major version into int", err)
		buildAndSendError(c, "0")
		return
	}
	clientMinorVersion, err := strconv.Atoi(clientVersion[1])
	if err != nil {
		debugLogError("Failed converting subsonic client minor version into int", err)
		buildAndSendError(c, "0")
		return
	}

	if clientMajorVersion > consts.SubsonicMajorVersion {
		buildAndSendError(c, "30")
		return
	} else if clientMajorVersion < consts.SubsonicMajorVersion {
		buildAndSendError(c, "20")
		return
	}

	if clientMinorVersion > consts.SubsonicMinorVersion {
		buildAndSendError(c, "30")
		return
	}
}

func (s *Server) subWithAuth(c *gin.Context) {
	var (
		requiredParams  = c.MustGet("requiredParams").(requiredParams)
		qUser           = requiredParams.U
		qHashedPassword = requiredParams.T
		qSalt           = requiredParams.S
		password        string
		cachedUser      types.SubsonicUser
		ctx             = context.Background()
	)

	if err := s.cache.Get(ctx, qUser).Scan(&cachedUser); err != nil {
		conn, err := s.pg_pool.Acquire(ctx)
		if err != nil {
			debugLogError("Failed acquiring connection from postgres connection pool", err)
			buildAndSendError(c, "0")
			return
		}
		defer conn.Release()
		query := sqlc.New(conn)

		user, err := query.GetUserByUsername(ctx, pgtype.Text{String: qUser, Valid: true})
		if err != nil {
			debugLogError("User does not exist", err)
			buildAndSendError(c, "40")
			return
		}

		encodedUser, err := json.Marshal(types.MapSqlUserToSubsonicUser(&user, user.Password))
		if err != nil {
			debugLogError("Failed encoding user credentials", err)
			buildAndSendError(c, "0")
			return
		}
		if err = s.cache.Set(ctx, user.Username.String, encodedUser, time.Minute*10).Err(); err != nil {
			debugLogError("Failed creating cache entry for user credentials", err)
			buildAndSendError(c, "0")
			return
		}
		password = user.Password
	} else {
		password = cachedUser.Password
	}

	if !auth.ValidatePassword(qHashedPassword, qSalt, password) {
		debugLog("Wrong Password.")
		buildAndSendError(c, "40")
		return
	}
}

// Util
func buildAndSendError(c *gin.Context, errorCode string) {
	c.Abort()

	subsonicRes := types.SubsonicResponse{
		Xmlns:   consts.Xmlns,
		Status:  "failed",
		Version: consts.SubsonicVersion,
	}

	subsonicRes.Error = &types.SubsonicError{
		Code:    errorCode,
		Message: consts.SubsonicErrorMessages[errorCode],
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

func debugLog(message string) {
	if gin.Mode() == gin.DebugMode {
		log.Printf(message)
	}
}

func debugLogError(message string, err error) {
	if gin.Mode() == gin.DebugMode {
		log.Printf("%s. Error: %s", message, err)
	}
}
