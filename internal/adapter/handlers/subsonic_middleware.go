package handlers

import (
	"encoding/json"
	"encoding/xml"
	"music-streaming/internal/core/domain"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	Xmlns                = "http://subsonic.org/restapi"
	SubsonicVersion      = "1.16.1"
	SubsonicMajorVersion = 1
	SubsonicMinorVersion = 16
)

var (
	SubsonicValidBitRates = []string{
		"0",
		"32",
		"40",
		"48",
		"56",
		"64",
		"80",
		"96",
		"112",
		"128",
		"160",
		"192",
		"224",
		"256",
		"320",
	}

	SubsonicValidFileFormats = []string{
		"mp3",
		"flac",
		"wav",
	}

	SubsonicUserRoles = []string{
		"scrobblingEnabled",
		"ldapAuthenticated",
		"adminRole",
		"settingsRole",
		"streamRole",
		"jukeboxRole",
		"downloadRole",
		"uploadRole",
		"playlistRole",
		"coverArtRole",
		"commentRole",
		"podcastRole",
		"shareRole",
		"videoConversionRole",
	}

	SubsonicErrorMessages = map[string]string{
		"0":  "A generic error.",
		"10": "Required parameter is missing.",
		"20": "Incompatible Subsonic REST protocol version. Client must upgrade.",
		"30": "Incompatible Subsonic REST protocol version. Server must upgrade.",
		"40": "Wrong username or password.",
		"41": "Token authentication not supported for LDAP users.",
		"50": "User is not authorized for the given operation.",
		"60": "The trial period for the Subsonic server is over. Please upgrade to Subsonic Premium. Visit subsonic.org for details.",
		"70": "The requested data was not found.",
	}
)

type SubsonicResponse struct {
	XMLName    xml.Name            `xml:"subsonic-response" json:"-"`
	Xmlns      string              `xml:"xmlns,attr" json:"-"`
	Status     string              `xml:"status,attr" json:"status"`
	Version    string              `xml:"version,attr" json:"version"`
	Error      *SubsonicError      `xml:"error,omitempty" json:"error,omitempty"`
	User       *domain.User        `xml:"user,omitempty" json:"user,omitempty"`
	ScanStatus *SubsonicScanStatus `xml:"scanStatus,omitempty" json:"scanStatus,omitempty"`
	Users      *[]domain.User      `xml:"users,omitempty" json:"users,omitempty"`
	// Artist     *Artist             `xml:"artist,omitempty" json:"artist,omitempty"`
	// Album      *Album              `xml:"album,omitempty" json:"album,omitempty"`
	// Song       *Song               `xml:"song,omitempty" json:"song,omitempty"`
}

type SubsonicError struct {
	XMLName xml.Name `xml:"error,omitempty" json:"-"`
	Code    string   `xml:"code,attr" json:"code"`
	Message string   `xml:"message,attr" json:"message"`
}

type SubsonicScanStatus struct {
	XMLName  xml.Name `xml:"scanStatus" json:"-"`
	Scanning bool     `xml:"scanning,attr" json:"scanning"`
	Count    int      `xml:"count,attr" json:"count"`
}

type requiredParams struct {
	U string `form:"u" binding:"required"`
	T string `form:"t" binding:"required"`
	S string `form:"s" binding:"required"`
	V string `form:"v" binding:"required"`
	C string `form:"c" binding:"required"`
	F string `form:"f"`
	P string `form:"p"`
}

const RequiredParameterKey = "required-parameters"

func ValidateSubsonicQueryParameters(c *gin.Context) {
	c.Set("contentType", "application/xml")

	var params requiredParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Debug().Err(err)
		buildAndSendError(c, "10")
		return
	}
	c.Set(RequiredParameterKey, params)

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

	if clientMajorVersion > SubsonicMajorVersion {
		buildAndSendError(c, "30")
		return
	} else if clientMajorVersion < SubsonicMajorVersion {
		buildAndSendError(c, "20")
		return
	}

	if clientMinorVersion > SubsonicMinorVersion {
		buildAndSendError(c, "30")
		return
	}
}

func buildAndSendError(c *gin.Context, errorCode string) {
	c.Abort()

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "failed",
		Version: SubsonicVersion,
	}

	subsonicRes.Error = &SubsonicError{
		Code:    errorCode,
		Message: SubsonicErrorMessages[errorCode],
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
		log.Error().Err(err).Msg("Failed to serialize response")
		c.Data(http.StatusInternalServerError, contentType, []byte{})
		return
	}
	c.Data(http.StatusOK, contentType, serializedBody)
}
