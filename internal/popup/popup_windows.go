// +build windows

package popup

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Run displays the Windows popup by executing the PowerShell script
func Run() {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".sultengutt")
	
	// Generate fresh script with random mantra
	scriptPath, err := GenerateWindowsScript(configDir)
	if err != nil {
		// Try fallback if WPF fails
		scriptPath, err = GenerateWindowsFallbackScript(configDir)
		if err != nil {
			return
		}
	}
	
	// Execute the PowerShell script
	cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-WindowStyle", "Hidden", "-File", scriptPath)
	cmd.Run()
}