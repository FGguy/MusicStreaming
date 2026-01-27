package handlers

import (
	"encoding/xml"
	"music-streaming/internal/core/domain"
)

// UserDTO represents the HTTP layer representation of a User
type UserDTO struct {
	XMLName             xml.Name `xml:"user" json:"-"`
	Username            string   `xml:"username,attr" json:"username" form:"username" binding:"required"`
	Email               string   `xml:"email,attr" json:"email" form:"email" binding:"required"`
	Password            string   `xml:"-" json:"password,omitempty" form:"password" binding:"required"`
	ScrobblingEnabled   bool     `xml:"scrobblingEnabled,attr" json:"scrobblingEnabled" form:"scrobblingEnabled"`
	LdapAuthenticated   bool     `xml:"ldapAuthenticated,attr" json:"ldapAuthenticated" form:"ldapAuthenticated"`
	AdminRole           bool     `xml:"adminRole,attr" json:"adminRole" form:"adminRole"`
	SettingsRole        bool     `xml:"settingsRole,attr" json:"settingsRole" form:"settingsRole"`
	StreamRole          bool     `xml:"streamRole,attr" json:"streamRole" form:"streamRole"`
	JukeboxRole         bool     `xml:"jukeboxRole,attr" json:"jukeboxRole" form:"jukeboxRole"`
	DownloadRole        bool     `xml:"downloadRole,attr" json:"downloadRole" form:"downloadRole"`
	UploadRole          bool     `xml:"uploadRole,attr" json:"uploadRole" form:"uploadRole"`
	PlaylistRole        bool     `xml:"playlistRole,attr" json:"playlistRole" form:"playlistRole"`
	CoverArtRole        bool     `xml:"coverArtRole,attr" json:"coverArtRole" form:"coverArtRole"`
	CommentRole         bool     `xml:"commentRole,attr" json:"commentRole" form:"commentRole"`
	PodcastRole         bool     `xml:"podcastRole,attr" json:"podcastRole" form:"podcastRole"`
	ShareRole           bool     `xml:"shareRole,attr" json:"shareRole" form:"shareRole"`
	VideoConversionRole bool     `xml:"videoConversionRole,attr" json:"videoConversionRole" form:"videoConversionRole"`
	MusicfolderId       []string `xml:"folder,omitempty" json:"folder,omitempty" form:"folder"`
	MaxBitRate          int32    `xml:"maxBitRate,attr" json:"maxBitRate" form:"maxBitRate"`
}

// ArtistDTO represents the HTTP layer representation of an Artist
type ArtistDTO struct {
	Id         int    `json:"id" xml:"id,attr"`
	Name       string `json:"name" xml:"name,attr"`
	CoverArt   string `json:"coverArt" xml:"coverArt,attr"`
	AlbumCount int    `json:"albumCount" xml:"albumCount,attr"`
}

// AlbumDTO represents the HTTP layer representation of an Album
type AlbumDTO struct {
	Id        int    `json:"id" xml:"id,attr"`
	ArtistId  int    `json:"artistId" xml:"artistId,attr"`
	Name      string `json:"name" xml:"name,attr"`
	CoverArt  string `json:"coverArt" xml:"coverArt,attr"`
	SongCount int    `json:"songCount" xml:"songCount,attr"`
	Created   string `json:"created" xml:"created,attr"`
	Duration  int    `json:"duration" xml:"duration,attr"`
	Artist    string `json:"artist" xml:"artist,attr"`
}

// SongDTO represents the HTTP layer representation of a Song
type SongDTO struct {
	Id          int    `json:"id" xml:"id,attr"`
	AlbumId     int    `json:"albumId" xml:"albumId,attr"`
	Title       string `json:"title" xml:"title,attr"`
	Album       string `json:"album" xml:"album,attr"`
	Artist      string `json:"artist" xml:"artist,attr"`
	IsDir       bool   `json:"isDir" xml:"isDir,attr"`
	CoverArt    string `json:"coverArt" xml:"coverArt,attr"`
	Created     string `json:"created" xml:"created,attr"`
	Duration    int    `json:"duration" xml:"duration,attr"`
	BitRate     int    `json:"bitRate" xml:"bitRate,attr"`
	Size        int64  `json:"size" xml:"size,attr"`
	Suffix      string `json:"suffix" xml:"suffix,attr"`
	ContentType string `json:"contentType" xml:"contentType,attr"`
	IsVideo     bool   `json:"isVideo" xml:"isVideo,attr"`
	Path        string `json:"path" xml:"path,attr"`
}

// ScanStatusDTO represents the HTTP layer representation of ScanStatus
type ScanStatusDTO struct {
	XMLName  xml.Name `xml:"scanStatus" json:"-"`
	Scanning bool     `xml:"scanning,attr" json:"scanning"`
	Count    int      `xml:"count,attr" json:"count"`
}

// Mapper functions from Domain to DTO

// UserToDTO converts a domain User to a UserDTO
func UserToDTO(user domain.User) UserDTO {
	return UserDTO{
		Username:            user.Username,
		Email:               user.Email,
		Password:            user.Password,
		ScrobblingEnabled:   user.ScrobblingEnabled,
		LdapAuthenticated:   user.LdapAuthenticated,
		AdminRole:           user.AdminRole,
		SettingsRole:        user.SettingsRole,
		StreamRole:          user.StreamRole,
		JukeboxRole:         user.JukeboxRole,
		DownloadRole:        user.DownloadRole,
		UploadRole:          user.UploadRole,
		PlaylistRole:        user.PlaylistRole,
		CoverArtRole:        user.CoverArtRole,
		CommentRole:         user.CommentRole,
		PodcastRole:         user.PodcastRole,
		ShareRole:           user.ShareRole,
		VideoConversionRole: user.VideoConversionRole,
		MusicfolderId:       user.MusicfolderId,
		MaxBitRate:          user.MaxBitRate,
	}
}

// UsersToDTO converts a slice of domain Users to a slice of UserDTOs
func UsersToDTO(users []domain.User) []UserDTO {
	dtos := make([]UserDTO, len(users))
	for i, user := range users {
		dtos[i] = UserToDTO(user)
	}
	return dtos
}

// ArtistToDTO converts a domain Artist to an ArtistDTO
func ArtistToDTO(artist domain.Artist) ArtistDTO {
	return ArtistDTO{
		Id:         artist.Id,
		Name:       artist.Name,
		CoverArt:   artist.CoverArt,
		AlbumCount: artist.AlbumCount,
	}
}

// AlbumToDTO converts a domain Album to an AlbumDTO
func AlbumToDTO(album domain.Album) AlbumDTO {
	return AlbumDTO{
		Id:        album.Id,
		ArtistId:  album.ArtistId,
		Name:      album.Name,
		CoverArt:  album.CoverArt,
		SongCount: album.SongCount,
		Created:   album.Created,
		Duration:  album.Duration,
		Artist:    album.Artist,
	}
}

// SongToDTO converts a domain Song to a SongDTO
func SongToDTO(song domain.Song) SongDTO {
	return SongDTO{
		Id:          song.Id,
		AlbumId:     song.AlbumId,
		Title:       song.Title,
		Album:       song.Album,
		Artist:      song.Artist,
		IsDir:       song.IsDir,
		CoverArt:    song.CoverArt,
		Created:     song.Created,
		Duration:    song.Duration,
		BitRate:     song.BitRate,
		Size:        song.Size,
		Suffix:      song.Suffix,
		ContentType: song.ContentType,
		IsVideo:     song.IsVideo,
		Path:        song.Path,
	}
}

// ScanStatusToDTO converts a domain ScanStatus to a ScanStatusDTO
func ScanStatusToDTO(status domain.ScanStatus) ScanStatusDTO {
	return ScanStatusDTO{
		Scanning: status.Scanning,
		Count:    status.Count,
	}
}

// Mapper functions from DTO to Domain

// DTOToUser converts a UserDTO to a domain User
func DTOToUser(dto UserDTO) domain.User {
	return domain.User{
		Username:            dto.Username,
		Email:               dto.Email,
		Password:            dto.Password,
		ScrobblingEnabled:   dto.ScrobblingEnabled,
		LdapAuthenticated:   dto.LdapAuthenticated,
		AdminRole:           dto.AdminRole,
		SettingsRole:        dto.SettingsRole,
		StreamRole:          dto.StreamRole,
		JukeboxRole:         dto.JukeboxRole,
		DownloadRole:        dto.DownloadRole,
		UploadRole:          dto.UploadRole,
		PlaylistRole:        dto.PlaylistRole,
		CoverArtRole:        dto.CoverArtRole,
		CommentRole:         dto.CommentRole,
		PodcastRole:         dto.PodcastRole,
		ShareRole:           dto.ShareRole,
		VideoConversionRole: dto.VideoConversionRole,
		MusicfolderId:       dto.MusicfolderId,
		MaxBitRate:          dto.MaxBitRate,
	}
}
