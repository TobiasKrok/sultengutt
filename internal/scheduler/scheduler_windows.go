//go:build windows
// +build windows

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
	schTask, err := utils.ResolveExecutablePath("schtasks")
	if err != nil {
		panic(fmt.Errorf("failed to find schtasks (Windows): %w", err))
	}
	return &WindowsScheduler{execPath: execPath, installOptions: options, schedulerExecPath: schTask, configDir: configDir}
}
