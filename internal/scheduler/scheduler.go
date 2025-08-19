package scheduler

import (
	"fmt"
	"runtime"
	"sultengutt/internal/config"
	"sultengutt/internal/utils"
)

type Scheduler interface {
	RegisterTask() error
	UnregisterTask() error
	TaskExists() (bool, error)
	Snooze() error
}

func NewScheduler(options config.InstallOptions, configDir string) Scheduler {
	execPath, err := utils.ResolveExecutablePath("sultengutt")
	if err != nil {
		panic(fmt.Errorf("failed to resolve executable path for Sultengutt: %w", err))
	}
	switch runtime.GOOS {
	case "windows":
		schTask, err := utils.ResolveExecutablePath("schtasks")
		if err != nil {
			panic(fmt.Errorf("failed to find schtasks (Windows): %w", err))
		}
		return &WindowsScheduler{execPath: execPath, installOptions: options, schedulerExecPath: schTask, configDir: configDir}
	case "darwin":
		return &MacScheduler{execPath: execPath, installOptions: options}
	default:
		panic("Unsupported OS")
	}
}
