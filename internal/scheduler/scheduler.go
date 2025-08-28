package scheduler

import (
	"sultengutt/internal/config"
)

type Scheduler interface {
	RegisterTask() error
	UnregisterTask() error
	TaskExists() (bool, error)
}

// NewScheduler creates a platform-specific scheduler
// The actual implementation is in scheduler_darwin.go, scheduler_windows.go, etc.
func NewScheduler(options config.InstallOptions, configDir string) Scheduler {
	return newScheduler(options, configDir)
}
