package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

const (
	Time24hRegex = `^(?:[01]?\d|2[0-3]):[0-5]\d$`
)

var validDays = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

type InstallOptions struct {
	Days []string `json:"days"`
	Hour string   `json:"hour"`
}

type Config struct {
	InstallOptions InstallOptions `json:"install_options"`
	PausedUntil    int64          `json:"paused_until"` // -1: not paused, 0: paused indefinitely, >0: unix timestamp

	configPath     string
	isFreshInstall bool
}

type ConfigManager struct {
	configDir  string
	configFile string
}

func NewConfigManager() (*ConfigManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	return &ConfigManager{
		configDir:  filepath.Join(homeDir, ".sultengutt"),
		configFile: "sultengutt.json",
	}, nil
}

func (cm *ConfigManager) Load() (*Config, error) {
	if err := os.MkdirAll(cm.configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(cm.configDir, cm.configFile)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := &Config{
				PausedUntil:    -1, // not paused by default
				configPath:     configPath,
				isFreshInstall: true,
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg.configPath = configPath
	cfg.isFreshInstall = false

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return &cfg, nil
}

func (cm *ConfigManager) Save(cfg *Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	configPath := filepath.Join(cm.configDir, cm.configFile)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	// atomic saving
	tempPath := configPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp config file: %w", err)
	}

	if err := os.Rename(tempPath, configPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

func (c *Config) IsFreshInstall() bool {
	return c.isFreshInstall

}

func (c *Config) IsPaused() bool {
	return c.PausedUntil != -1
}

func (c *Config) SetPausedUntil(timestamp int64) {
	c.PausedUntil = timestamp
}

func (c *Config) Resume() {
	c.PausedUntil = -1
}

func (c *Config) validate() error {

	if len(c.InstallOptions.Days) == 0 {
		return errors.New("no days specified")
	}
	for _, day := range c.InstallOptions.Days {
		if !slices.Contains(validDays, day) {
			return errors.New("invalid day specified")
		}
	}
	if c.InstallOptions.Hour == "" {
		return errors.New("no hour specified")
	}
	pattern := regexp.MustCompile(Time24hRegex)
	if !pattern.MatchString(c.InstallOptions.Hour) {
		return errors.New("invalid hour specified")
	}
	return nil
}
