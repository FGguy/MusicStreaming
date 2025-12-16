package domain

import "encoding/xml"

type ScanStatus struct {
	XMLName  xml.Name `xml:"scanStatus" json:"-"`
	Scanning bool     `xml:"scanning,attr" json:"scanning"`
	Count    int      `xml:"count,attr" json:"count"`
}
