//go:build linux
// +build linux

package scheduler

import (
	"fmt"
	"sultengutt/internal/config"
)

type LinuxScheduler struct {
	installOptions config.InstallOptions
	execPath       string
}

func (l *LinuxScheduler) RegisterTask() error {
	return fmt.Errorf("Linux scheduler not yet implemented")
}

func (l *LinuxScheduler) UnregisterTask() error {
	return fmt.Errorf("Linux scheduler not yet implemented")
}

func (l *LinuxScheduler) TaskExists() (bool, error) {
	return false, fmt.Errorf("Linux scheduler not yet implemented")
}
