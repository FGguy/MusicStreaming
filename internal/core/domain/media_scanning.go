package domain

import "encoding/xml"

type ScanStatus struct {
	XMLName  xml.Name `xml:"scanStatus" json:"-"`
	Scanning bool     `xml:"scanning,attr" json:"scanning"`
	Count    int      `xml:"count,attr" json:"count"`
}

type FFProbeFormat struct {
	Filename       string `json:"filename"`
	NbStreams      int    `json:"nb_streams"`
	NbPrograms     int    `json:"nb_programs"`
	NbStreamGroups int    `json:"nb_stream_groups"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime      string `json:"start_time"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
}

type FFProbeInfo struct {
	Format *FFProbeFormat `json:"format"`
}
