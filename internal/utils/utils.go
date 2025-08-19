package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ResolveExecutablePath finds and resolves the full path to an executable
func ResolveExecutablePath(execPath string) (string, error) {
	path, err := exec.LookPath(execPath)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}

// ParseDuration parses duration strings in various formats
func ParseDuration(input interface{}) (time.Duration, error) {
	var durationStr string
	
	// Handle both string and []string inputs
	switch v := input.(type) {
	case string:
		durationStr = v
	case []string:
		if len(v) < 1 {
			return 0, fmt.Errorf("invalid duration format")
		}
		durationStr = strings.Join(v, " ")
	default:
		return 0, fmt.Errorf("invalid input type")
	}

	durationStr = strings.TrimSpace(durationStr)
	
	// Single regex to handle all formats
	re := regexp.MustCompile(`^\s*(\d+)\s*([mhd]|minutes?|mins?|hours?|hrs?|days?|weeks?|months?)\s*$`)
	matches := re.FindStringSubmatch(durationStr)
	
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format")
	}

	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", matches[1])
	}

	unit := strings.ToLower(strings.TrimSuffix(matches[2], "s"))
	
	// Map all unit variations to duration
	switch unit {
	case "m", "min", "minute":
		return time.Duration(amount) * time.Minute, nil
	case "h", "hr", "hour":
		return time.Duration(amount) * time.Hour, nil
	case "d", "day":
		return time.Duration(amount) * 24 * time.Hour, nil
	case "week":
		return time.Duration(amount) * 7 * 24 * time.Hour, nil
	case "month":
		return time.Duration(amount) * 30 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}
}

// parseTimeFromString parses "HH:MM" format
func parseTimeFromString(hourStr string) (hour, minute int, err error) {
	parts := strings.Split(hourStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format: %s", hourStr)
	}

	hour, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid hour: %s", parts[0])
	}

	minute, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minute: %s", parts[1])
	}

	return hour, minute, nil
}

// CalculatePauseUntil calculates when to unpause based on duration and scheduled time
func CalculatePauseUntil(duration time.Duration, scheduledHour string) (int64, error) {
	hour, minute, err := parseTimeFromString(scheduledHour)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	
	// Calculate next scheduled time
	scheduledTime := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if scheduledTime.Before(now) || scheduledTime.Equal(now) {
		scheduledTime = scheduledTime.Add(24 * time.Hour)
	}

	// Add the pause duration and align to the scheduled time of day
	unpauseTime := scheduledTime.Add(duration)
	finalUnpauseTime := time.Date(unpauseTime.Year(), unpauseTime.Month(), unpauseTime.Day(), hour, minute, 0, 0, unpauseTime.Location())

	return finalUnpauseTime.Unix(), nil
}
