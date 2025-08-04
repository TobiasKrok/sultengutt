package scheduler

import (
	"fmt"
	"os/exec"
	"runtime"
	"sultengutt/internal/installer"
	"sultengutt/internal/utils"
)

type Scheduler interface {
	RegisterTask() error
	UnregisterTask() error
	TaskExists() (bool, error)
	createTask() []string
}

func NewScheduler(options installer.InstallOptions) Scheduler {
	execPath, err := utils.ResolveExecutablePath()
	if err != nil {
		panic(fmt.Errorf("Failed to resolve executable path for Sultengutt: %w", err))
	}
	switch runtime.GOOS {
	case "windows":
		schTask, err := exec.LookPath("schtasks")
		if err != nil {
			panic(fmt.Errorf("Failed to find schtasks (Windows): %w", err))
		}
		return &WindowsScheduler{execPath: execPath, InstallOptions: options, schedulerExecPath: schTask}
	default:
		panic("Unsupported OS")
	}
}
