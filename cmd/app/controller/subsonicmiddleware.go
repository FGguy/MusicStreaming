package controller

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	consts "music-streaming/internal/consts"
	types "music-streaming/internal/types"
	auth "music-streaming/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

func (s *Application) subValidateQParamsMiddleware(c *gin.Context) {
	c.Set("contentType", "application/xml")

	var params requiredParams
	if err := c.ShouldBindQuery(&params); err != nil {
		buildAndSendError(c, "10")
		return
	}
	c.Set("requiredParams", params)

	if params.F == "json" {
		c.Set("contentType", "application/json")
	}

	//validate Subsonic API compatibility with client
	clientVersion := strings.Split(params.V, ".")
	clientMajorVersion, err := strconv.Atoi(clientVersion[0])
	if err != nil {
		log.Warn().Err(err).Msgf("Failed converting subsonic client major version into int")
		buildAndSendError(c, "0")
		return
	}
	clientMinorVersion, err := strconv.Atoi(clientVersion[1])
	if err != nil {
		log.Warn().Err(err).Msgf("Failed converting subsonic client minor version into int")
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

func (s *Application) subWithAuth(c *gin.Context) {
	var (
		requiredParams  = c.MustGet("requiredParams").(requiredParams)
		qUser           = requiredParams.U
		qHashedPassword = requiredParams.T
		qSalt           = requiredParams.S
		ctx             = context.Background()
	)

	user, err := s.dataLayer.GetUser(ctx, qUser)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed fetching user for authorization")
		buildAndSendError(c, "0")
		return
	}

	if !auth.ValidatePassword(qHashedPassword, qSalt, user.Password) {
		log.Trace().Msg("Login attempt with wrong password")
		buildAndSendError(c, "40")
		return
	}
	c.Set("requestingUser", user) //so further routes can check permission
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

	SerializeAndSendBody(c, subsonicRes)
}

func SerializeAndSendBody(c *gin.Context, body any) {
	var (
		serializedBody []byte
		err            error
		contentType    = c.MustGet("contentType").(string)
	)

	if contentType == "application/json" {
		serializedBody, err = json.Marshal(body)
	} else {
		serializedBody, err = xml.Marshal(body)
	}

	if err != nil {
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}
	c.Data(http.StatusOK, contentType, serializedBody)
}
