package domain

// ScanStatus represents the current status of a media library scan
type ScanStatus struct {
	Scanning bool
	Count    int
}

// FFProbeFormat represents the format information from ffprobe
// JSON tags are kept here as they are used for parsing external tool output
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

// FFProbeInfo represents the full information from ffprobe
// JSON tags are kept here as they are used for parsing external tool output
type FFProbeInfo struct {
	Format *FFProbeFormat `json:"format"`
}
