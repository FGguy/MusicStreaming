package subxml

import (
	"encoding/xml"
)

const (
	Xmlns           = "http://subsonic.org/restapi"
	SubsonicVersion = "1.16.1"
)

var SubsonicErrorMessages = map[string]string{
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

type SubsonicResponse struct {
	XMLName xml.Name       `xml:"subsonic-response"`
	Xmlns   string         `xml:"xmlns,attr"`
	Status  string         `xml:"status,attr"`
	Version string         `xml:"version,attr"`
	Error   *SubsonicError `xml:"error,omitempty"`
}

type SubsonicError struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:"message,attr"`
}
