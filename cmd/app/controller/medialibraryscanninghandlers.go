package controller

import (
	consts "music-streaming/internal/consts"
	types "music-streaming/internal/types"

	"github.com/gin-gonic/gin"
)

// GET
func (s *Application) handleGetScanStatus(c *gin.Context) {
	s.mu.Lock()
	scanStatus := &types.SubsonicScanStatus{
		Scanning: s.state.scanning,
		Count:    s.state.count,
	}
	s.mu.Unlock()

	subsonicRes := types.SubsonicResponse{
		Xmlns:      consts.Xmlns,
		Status:     "ok",
		Version:    consts.SubsonicVersion,
		ScanStatus: scanStatus,
	}

	SerializeAndSendBody(c, subsonicRes)
}

// GET
func (s *Application) handleStartScan(c *gin.Context) {
	rUser := c.MustGet("requestingUser").(*types.SubsonicUser)
	if !rUser.AdminRole {
		buildAndSendError(c, "50")
		return
	}

	var scanStatus *types.SubsonicScanStatus
	if !s.state.scanning {
		go s.MediaScan()

		scanStatus = &types.SubsonicScanStatus{
			Scanning: true,
			Count:    0,
		}
	} else {
		s.mu.Lock()
		scanStatus = &types.SubsonicScanStatus{
			Scanning: s.state.scanning,
			Count:    s.state.count,
		}
		s.mu.Unlock()
	}

	subsonicRes := types.SubsonicResponse{
		Xmlns:      consts.Xmlns,
		Status:     "ok",
		Version:    consts.SubsonicVersion,
		ScanStatus: scanStatus,
	}

	SerializeAndSendBody(c, subsonicRes)
}
