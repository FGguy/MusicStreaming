package subxml

import (
	"encoding/xml"
)

const (
	Xmlns           = "http://subsonic.org/restapi"
	SubsonicVersion = "1.16.1"
)

type SubsonicResponse struct {
	XMLName xml.Name `xml:"subsonic-response"`
	Xmlns   string   `xml:"xmlns,attr"`
	Status  string   `xml:"status,attr"`
	Version string   `xml:"version,attr"`
	//add all possible types of children
}
