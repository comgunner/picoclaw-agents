// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT

//go:build !linux && !darwin && !windows

package service

import "fmt"

func newPlatform(cfg ServiceConfig, exePath string) (ServicePlatform, error) {
	return nil, fmt.Errorf("service management is not supported on this platform")
}
