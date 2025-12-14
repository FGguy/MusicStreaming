package util

import (
	"encoding/json"
	"os/exec"
)

type Format struct {
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
	Format *Format `json:"format"`
}

//ffprobe -v error -hide_banner -print_format json -show_format '.\Ashes In Your Mouth.mp3'

func FFProbeProcessFile(filepath string) (*FFProbeInfo, error) {
	//build command
	cmd := exec.Command("ffprobe", "-v", "error", "-hide_banner", "-print_format", "json", "-show_format", filepath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info FFProbeInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
