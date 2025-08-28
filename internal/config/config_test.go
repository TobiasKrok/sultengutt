package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewConfigManager(t *testing.T) {
	cm, err := NewConfigManager()
	if err != nil {
		t.Fatalf("Failed to create ConfigManager: %v", err)
	}

	if cm == nil {
		t.Fatal("ConfigManager is nil")
	}

	if cm.configDir == "" {
		t.Error("ConfigDir is empty")
	}

	if cm.configFile == "" {
		t.Error("ConfigFile is empty")
	}
}

func TestConfigLoad(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	cm := &ConfigManager{
		configDir:  tempDir,
		configFile: "test_config.json",
	}

	t.Run("fresh install", func(t *testing.T) {
		cfg, err := cm.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if !cfg.IsFreshInstall() {
			t.Error("Expected fresh install to be true")
		}

		if cfg.PausedUntil != -1 {
			t.Errorf("Expected PausedUntil to be -1, got %d", cfg.PausedUntil)
		}
	})

	t.Run("existing valid config", func(t *testing.T) {
		// Create a valid config file
		testConfig := Config{
			InstallOptions: InstallOptions{
				Days:     []string{"Monday", "Tuesday"},
				Hour:     "14:30",
				SiteLink: "https://example.com",
			},
			PausedUntil: -1,
		}

		data, err := json.MarshalIndent(testConfig, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal test config: %v", err)
		}

		configPath := filepath.Join(tempDir, "test_config.json")
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		cfg, err := cm.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if cfg.IsFreshInstall() {
			t.Error("Expected fresh install to be false")
		}

		if len(cfg.InstallOptions.Days) != 2 {
			t.Errorf("Expected 2 days, got %d", len(cfg.InstallOptions.Days))
		}

		if cfg.InstallOptions.Hour != "14:30" {
			t.Errorf("Expected hour '14:30', got '%s'", cfg.InstallOptions.Hour)
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		// Create invalid config file
		configPath := filepath.Join(tempDir, "invalid_config.json")
		if err := os.WriteFile(configPath, []byte("invalid json"), 0644); err != nil {
			t.Fatalf("Failed to write invalid config: %v", err)
		}

		cm := &ConfigManager{
			configDir:  tempDir,
			configFile: "invalid_config.json",
		}

		_, err := cm.Load()
		if err == nil {
			t.Error("Expected error for invalid config, but got none")
		}
	})
}

func TestConfigSave(t *testing.T) {
	tempDir := t.TempDir()

	cm := &ConfigManager{
		configDir:  tempDir,
		configFile: "test_config.json",
	}

	cfg := &Config{
		InstallOptions: InstallOptions{
			Days:     []string{"Monday", "Wednesday", "Friday"},
			Hour:     "12:00",
			SiteLink: "https://test.com",
		},
		PausedUntil: -1,
		configPath:  filepath.Join(tempDir, "test_config.json"),
	}

	err := cm.Save(cfg)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists and is valid JSON
	configPath := filepath.Join(tempDir, "test_config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	var savedConfig Config
	if err := json.Unmarshal(data, &savedConfig); err != nil {
		t.Fatalf("Saved config is not valid JSON: %v", err)
	}

	if len(savedConfig.InstallOptions.Days) != 3 {
		t.Errorf("Expected 3 days, got %d", len(savedConfig.InstallOptions.Days))
	}

	if savedConfig.InstallOptions.Hour != "12:00" {
		t.Errorf("Expected hour '12:00', got '%s'", savedConfig.InstallOptions.Hour)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				InstallOptions: InstallOptions{
					Days:     []string{"Monday", "Tuesday"},
					Hour:     "14:30",
					SiteLink: "https://example.com",
				},
			},
			expectError: false,
		},
		{
			name: "empty days",
			config: Config{
				InstallOptions: InstallOptions{
					Days:     []string{},
					Hour:     "14:30",
					SiteLink: "https://example.com",
				},
			},
			expectError: true,
		},
		{
			name: "invalid day",
			config: Config{
				InstallOptions: InstallOptions{
					Days:     []string{"InvalidDay"},
					Hour:     "14:30",
					SiteLink: "https://example.com",
				},
			},
			expectError: true,
		},
		{
			name: "empty hour",
			config: Config{
				InstallOptions: InstallOptions{
					Days:     []string{"Monday"},
					Hour:     "",
					SiteLink: "https://example.com",
				},
			},
			expectError: true,
		},
		{
			name: "invalid hour format",
			config: Config{
				InstallOptions: InstallOptions{
					Days:     []string{"Monday"},
					Hour:     "25:30",
					SiteLink: "https://example.com",
				},
			},
			expectError: true,
		},
		{
			name: "empty site link",
			config: Config{
				InstallOptions: InstallOptions{
					Days:     []string{"Monday"},
					Hour:     "14:30",
					SiteLink: "",
				},
			},
			expectError: true,
		},
		{
			name: "invalid URL",
			config: Config{
				InstallOptions: InstallOptions{
					Days:     []string{"Monday"},
					Hour:     "14:30",
					SiteLink: "://invalid-url",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()

			if tt.expectError && err == nil {
				t.Error("Expected validation error, but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestConfigPauseResume(t *testing.T) {
	cfg := &Config{
		PausedUntil: -1,
	}

	// Test initial state
	if cfg.IsPaused() {
		t.Error("Expected config to not be paused initially")
	}

	// Test indefinite pause
	cfg.SetPausedUntil(0)
	if !cfg.IsPaused() {
		t.Error("Expected config to be paused after setting PausedUntil to 0")
	}

	// Test timed pause
	futureTime := time.Now().Add(2 * time.Hour).Unix()
	cfg.SetPausedUntil(futureTime)
	if !cfg.IsPaused() {
		t.Error("Expected config to be paused after setting future timestamp")
	}

	// Test resume
	cfg.Resume()
	if cfg.IsPaused() {
		t.Error("Expected config to not be paused after resume")
	}

	if cfg.PausedUntil != -1 {
		t.Errorf("Expected PausedUntil to be -1 after resume, got %d", cfg.PausedUntil)
	}
}

func TestConfigClean(t *testing.T) {
	tempDir := t.TempDir()

	// Create some files in the config directory
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cm := &ConfigManager{
		configDir:  tempDir,
		configFile: "config.json",
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Test file should exist before cleaning")
	}

	// Clean the directory
	if err := cm.Clean(); err != nil {
		t.Fatalf("Failed to clean config: %v", err)
	}

	// Verify directory is gone
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Config directory should be removed after cleaning")
	}
}

// Benchmark tests
func BenchmarkConfigLoad(b *testing.B) {
	tempDir := b.TempDir()

	// Create a test config file
	testConfig := Config{
		InstallOptions: InstallOptions{
			Days:     []string{"Monday", "Tuesday", "Wednesday"},
			Hour:     "14:30",
			SiteLink: "https://example.com",
		},
		PausedUntil: -1,
	}

	data, _ := json.MarshalIndent(testConfig, "", "  ")
	configPath := filepath.Join(tempDir, "bench_config.json")
	os.WriteFile(configPath, data, 0644)

	cm := &ConfigManager{
		configDir:  tempDir,
		configFile: "bench_config.json",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Load()
	}
}

func BenchmarkConfigSave(b *testing.B) {
	tempDir := b.TempDir()

	cm := &ConfigManager{
		configDir:  tempDir,
		configFile: "bench_config.json",
	}

	cfg := &Config{
		InstallOptions: InstallOptions{
			Days:     []string{"Monday", "Tuesday", "Wednesday"},
			Hour:     "14:30",
			SiteLink: "https://example.com",
		},
		PausedUntil: -1,
		configPath:  filepath.Join(tempDir, "bench_config.json"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Save(cfg)
	}
}
