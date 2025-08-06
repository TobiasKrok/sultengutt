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
	createTask() []string
}

func NewScheduler(options config.InstallOptions) Scheduler {
	execPath, err := utils.ResolveExecutablePath("sultengutt")
	if err != nil {
		panic(fmt.Errorf("Failed to resolve executable path for Sultengutt: %w", err))
	}
	switch runtime.GOOS {
	case "windows":
		schTask, err := utils.ResolveExecutablePath("schtasks")
		if err != nil {
			panic(fmt.Errorf("Failed to find schtasks (Windows): %w", err))
		}
		return &WindowsScheduler{execPath: execPath, installOptions: options, schedulerExecPath: schTask}
	default:
		panic("Unsupported OS")
	}
}
