//go:build darwin
// +build darwin

package scheduler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sultengutt/internal/config"
)

type MacScheduler struct {
	installOptions    config.InstallOptions
	execPath          string
	schedulerExecPath string
}

func (m *MacScheduler) RegisterTask() error {
	plistPath := m.getPlistPath()
	plistContent := m.createPlist()

	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	cmd := exec.Command("launchctl", "load", plistPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to register task: %w\n%s", err, out)
	}
	return nil
}

func (m *MacScheduler) UnregisterTask() error {
	exists, err := m.TaskExists()
	if err != nil {
		return err
	}
	if !exists {

		return nil
	}

	plistPath := m.getPlistPath()

	cmd := exec.Command("launchctl", "unload", plistPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unregister task: %w\n%s", err, out)
	}

	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}

	return nil
}

func (m *MacScheduler) TaskExists() (bool, error) {
	plistPath := m.getPlistPath()
	if _, err := os.Stat(plistPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Check if it's loaded in launchctl
	cmd := exec.Command("launchctl", "list", "no.tobias.sultengutt")
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 0 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (m *MacScheduler) Snooze() error {
	// TODO: Implement snooze functionality
	return nil
}

func (m *MacScheduler) getPlistPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, "Library", "LaunchAgents", "no.tobias.sultengutt.plist")
}

func (m *MacScheduler) createPlist() string {
	parts := strings.Split(m.installOptions.Hour, ":")
	hour := parts[0]
	minute := "0"
	if len(parts) > 1 {
		minute = parts[1]
	}

	hourInt, _ := strconv.Atoi(hour)
	minuteInt, _ := strconv.Atoi(minute)

	var calendarIntervals string
	dayMap := map[string]int{
		"Sunday": 0, "Monday": 1, "Tuesday": 2, "Wednesday": 3,
		"Thursday": 4, "Friday": 5, "Saturday": 6,
	}

	for _, day := range m.installOptions.Days {
		if weekday, ok := dayMap[day]; ok {
			calendarIntervals += fmt.Sprintf(`
		<dict>
			<key>Weekday</key>
			<integer>%d</integer>
			<key>Hour</key>
			<integer>%d</integer>
			<key>Minute</key>
			<integer>%d</integer>
		</dict>`, weekday, hourInt, minuteInt)
		}
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>no.tobias.sultengutt</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
		<string>execute</string>
	</array>
	<key>StartCalendarInterval</key>
	<array>%s
	</array>
	<key>RunAtLoad</key>
	<false/>
</dict>
</plist>`, m.execPath, calendarIntervals)
}
