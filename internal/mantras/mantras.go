package mantras

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type MantraLoader struct {
	mantras []string
}

func New() (*MantraLoader, error) {
	ml := &MantraLoader{}
	if err := ml.loadMantras(); err != nil {
		return nil, err
	}
	return ml, nil
}

func (ml *MantraLoader) loadMantras() error {
	// Try to load from multiple locations in order of preference
	locations := ml.getMantraLocations()
	
	for _, location := range locations {
		if mantras, err := ml.loadFromFile(location); err == nil {
			ml.mantras = mantras
			return nil
		}
	}
	
	return fmt.Errorf("no mantras.json file found in any of the expected locations")
}

func (ml *MantraLoader) getMantraLocations() []string {
	var locations []string
	
	// 1. User config directory (~/.sultengutt/mantras.json)
	if homeDir, err := os.UserHomeDir(); err == nil {
		userConfig := filepath.Join(homeDir, ".sultengutt", "mantras.json")
		locations = append(locations, userConfig)
	}
	
	// 2. System-wide config locations (for package managers)
	switch runtime.GOOS {
	case "darwin":
		// Homebrew installs config files here
		locations = append(locations, 
			"/usr/local/etc/sultengutt/mantras.json",
			"/opt/homebrew/etc/sultengutt/mantras.json",
		)
	case "windows":
		// Scoop installs config files in persist directory
		if scoopDir := os.Getenv("SCOOP"); scoopDir != "" {
			locations = append(locations, filepath.Join(scoopDir, "persist", "sultengutt", "mantras.json"))
		}
		// Also check ProgramData for system-wide config
		locations = append(locations, filepath.Join(os.Getenv("ProgramData"), "sultengutt", "mantras.json"))
	case "linux":
		locations = append(locations, 
			"/etc/sultengutt/mantras.json",
			"/usr/local/etc/sultengutt/mantras.json",
		)
	}
	
	// 3. Relative to executable (for development)
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		locations = append(locations, 
			filepath.Join(execDir, "config", "mantras.json"),
			filepath.Join(execDir, "..", "config", "mantras.json"),
		)
	}
	
	return locations
}

func (ml *MantraLoader) loadFromFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var mantras []string
	if err := json.Unmarshal(data, &mantras); err != nil {
		return nil, fmt.Errorf("failed to parse mantras from %s: %w", path, err)
	}
	
	if len(mantras) == 0 {
		return nil, fmt.Errorf("no mantras found in %s", path)
	}
	
	return mantras, nil
}

func (ml *MantraLoader) GetRandom() string {
	if len(ml.mantras) == 0 {
		return ""
	}
	
	rand.Seed(time.Now().UnixNano())
	return ml.mantras[rand.Intn(len(ml.mantras))]
}

func (ml *MantraLoader) GetAll() []string {
	return ml.mantras
}

// GetConfigLocations returns all the locations where mantras.json is searched for
func (ml *MantraLoader) GetConfigLocations() []string {
	return ml.getMantraLocations()
}

// CopyMantrasToUserConfig copies the mantras.json from source to user config directory
func CopyMantrasToUserConfig(sourcePath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	configDir := filepath.Join(homeDir, ".sultengutt")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	destPath := filepath.Join(configDir, "mantras.json")
	
	// Check if file already exists
	if _, err := os.Stat(destPath); err == nil {
		return nil // File already exists
	}
	
	// Copy from source
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source mantras file: %w", err)
	}
	
	return os.WriteFile(destPath, data, 0644)
}