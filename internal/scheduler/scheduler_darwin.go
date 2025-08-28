//go:build darwin
// +build darwin

package scheduler

import (
	"fmt"
	"sultengutt/internal/config"
	"sultengutt/internal/utils"
)

func newScheduler(options config.InstallOptions, configDir string) Scheduler {
	execPath, err := utils.ResolveExecutablePath("sultengutt")
	if err != nil {
		panic(fmt.Errorf("failed to resolve executable path for Sultengutt: %w", err))
	}
	return &MacScheduler{execPath: execPath, installOptions: options}
}
