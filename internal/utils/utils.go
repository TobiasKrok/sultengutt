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

func ResolveExecutablePath(execPath string) (string, error) {

	path, err := exec.LookPath(execPath)
	if err != nil {
		return "", err
	}
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return resolvedPath, nil
}

func ParseDuration(args []string) (time.Duration, error) {
	if len(args) < 2 {
		return 0, fmt.Errorf("invalid duration format. Use: sultengutt pause <number> <unit>")
	}

	durationStr := strings.Join(args, " ")

	re := regexp.MustCompile(`^\s*(\d+)\s+(day|days|week|weeks|month|months)\s*$`)
	matches := re.FindStringSubmatch(durationStr)

	if matches == nil {
		return 0, fmt.Errorf("invalid duration format. Use: sultengutt pause <number> <unit>")
	}

	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", matches[1])
	}

	unit := matches[2]

	unit = strings.TrimSuffix(unit, "s")

	var duration time.Duration
	switch unit {
	case "day":
		duration = time.Duration(amount) * 24 * time.Hour
	case "week":
		duration = time.Duration(amount) * 7 * 24 * time.Hour
	case "month":
		// Using 30 days as an approximation for a month
		duration = time.Duration(amount) * 30 * 24 * time.Hour
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return duration, nil
}

func CalculateNextScheduledTime(hourStr string, fromTime time.Time) (time.Time, error) {
	parts := strings.Split(hourStr, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid hour format: %s", hourStr)
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour: %s", parts[0])
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute: %s", parts[1])
	}

	year, month, day := fromTime.Date()
	location := fromTime.Location()
	scheduledTime := time.Date(year, month, day, hour, minute, 0, 0, location)

	if scheduledTime.Before(fromTime) || scheduledTime.Equal(fromTime) {
		scheduledTime = scheduledTime.Add(24 * time.Hour)
	}

	return scheduledTime, nil
}

func CalculatePauseUntil(duration time.Duration, scheduledHour string) (int64, error) {
	now := time.Now()

	nextScheduled, err := CalculateNextScheduledTime(scheduledHour, now)
	if err != nil {
		return 0, err
	}

	unpauseTime := nextScheduled.Add(duration)

	year, month, day := unpauseTime.Date()
	location := unpauseTime.Location()

	parts := strings.Split(scheduledHour, ":")
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])

	finalUnpauseTime := time.Date(year, month, day, hour, minute, 0, 0, location)

	return finalUnpauseTime.Unix(), nil
}
