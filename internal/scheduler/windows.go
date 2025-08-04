package scheduler

import (
	"fmt"
	"os/exec"
	"strings"
	"sultengutt/internal/config"
)

type WindowsScheduler struct {
	installOptions    config.InstallOptions
	execPath          string
	schedulerExecPath string
}

func (w *WindowsScheduler) RegisterTask() error {
	args := w.createTask()
	cmd := exec.Command(w.schedulerExecPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to register task: %w\n%s", err, out)
	}
	return nil
}

func (w *WindowsScheduler) UnregisterTask() error {
	exists, err := w.TaskExists()
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	cmd := exec.Command(w.schedulerExecPath, "/delete", "/tn", "Sultengutt", "/f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to unregister task: %w\n%s", err, out)
	}
	return nil
}

func (w *WindowsScheduler) TaskExists() (bool, error) {
	// we assume there is only one Sultengutt task :)
	cmd := exec.Command(w.schedulerExecPath, "/query", "/tn", "Sultengutt")
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (w *WindowsScheduler) createTask() []string {
	var days []string
	for _, day := range w.installOptions.Days {
		days = append(days, strings.ToUpper(day)[0:3]) // schtask accepts strings like MON,TUE,THU
	}
	return []string{"/create",
		"/tn", "Sultengutt",
		"/tr", fmt.Sprintf("%s execute", w.execPath),
		"/sc", "weekly",
		"/d", strings.Join(days, ","),
		"/st", w.installOptions.Hour,
		"/f"}
}
