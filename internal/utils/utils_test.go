package utils

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected time.Duration
		hasError bool
	}{
		// String inputs
		{"30 minutes", "30 minutes", 30 * time.Minute, false},
		{"30m", "30m", 30 * time.Minute, false},
		{"2 hours", "2 hours", 2 * time.Hour, false},
		{"2h", "2h", 2 * time.Hour, false},
		{"1 day", "1 day", 24 * time.Hour, false},
		{"1d", "1d", 24 * time.Hour, false},
		{"2 weeks", "2 weeks", 14 * 24 * time.Hour, false},
		{"1 month", "1 month", 30 * 24 * time.Hour, false},

		// Plural forms
		{"5 mins", "5 mins", 5 * time.Minute, false},
		{"3 hrs", "3 hrs", 3 * time.Hour, false},
		{"2 days", "2 days", 2 * 24 * time.Hour, false},

		// String slice inputs
		{"slice input", []string{"30", "minutes"}, 30 * time.Minute, false},
		{"slice input 2", []string{"2", "hours"}, 2 * time.Hour, false},

		// Error cases
		{"invalid format", "invalid", 0, true},
		{"empty string", "", 0, true},
		{"negative number", "-5 hours", 0, true},
		{"no number", "hours", 0, true},
		{"unsupported unit", "5 years", 0, true},
		{"empty slice", []string{}, 0, true},
		{"invalid type", 123, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDuration(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input %v, but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %v: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("For input %v, expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}

func TestCalculatePauseUntil(t *testing.T) {
	tests := []struct {
		name          string
		duration      time.Duration
		scheduledHour string
		expectError   bool
	}{
		{"valid time", 2 * time.Hour, "14:30", false},
		{"midnight", time.Hour, "00:00", false},
		{"noon", 30 * time.Minute, "12:00", false},
		{"invalid time format", time.Hour, "25:00", false}, // parseTimeFromString doesn't validate ranges, CalculatePauseUntil will work
		{"invalid hour", time.Hour, "abc:30", true},
		{"invalid minute", time.Hour, "14:abc", true},
		{"missing colon", time.Hour, "1430", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculatePauseUntil(tt.duration, tt.scheduledHour)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for scheduledHour %s, but got none", tt.scheduledHour)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for scheduledHour %s: %v", tt.scheduledHour, err)
				return
			}

			// Verify result is a valid timestamp in the future
			if result <= 0 {
				t.Errorf("Expected positive timestamp, got %d", result)
			}

			// Convert back to time and verify it's in the future
			unpauseTime := time.Unix(result, 0)
			if unpauseTime.Before(time.Now()) {
				t.Errorf("Expected future time, got %v", unpauseTime)
			}
		})
	}
}

func TestParseTimeFromString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedHour int
		expectedMin  int
		expectError  bool
	}{
		{"valid time", "14:30", 14, 30, false},
		{"midnight", "00:00", 0, 0, false},
		{"noon", "12:00", 12, 0, false},
		{"single digit hour", "9:15", 9, 15, false},
		{"single digit minute", "14:5", 14, 5, false},
		{"missing colon", "1430", 0, 0, true},
		{"invalid hour", "25:30", 25, 30, false},   // parseTimeFromString doesn't validate ranges
		{"invalid minute", "14:60", 14, 60, false}, // parseTimeFromString doesn't validate ranges
		{"non-numeric hour", "abc:30", 0, 0, true},
		{"non-numeric minute", "14:abc", 0, 0, true},
		{"empty string", "", 0, 0, true},
		{"too many parts", "14:30:45", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, minute, err := parseTimeFromString(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", tt.input, err)
				return
			}

			if hour != tt.expectedHour {
				t.Errorf("Expected hour %d, got %d", tt.expectedHour, hour)
			}

			if minute != tt.expectedMin {
				t.Errorf("Expected minute %d, got %d", tt.expectedMin, minute)
			}
		})
	}
}

func TestResolveExecutablePath(t *testing.T) {
	// Test with a common executable that should exist on all systems
	tests := []struct {
		name        string
		executable  string
		expectError bool
	}{
		{"valid executable - go", "go", false},
		{"invalid executable", "definitely-not-an-executable-12345", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolveExecutablePath(tt.executable)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for executable %s, but got none", tt.executable)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for executable %s: %v", tt.executable, err)
				return
			}

			if result == "" {
				t.Errorf("Expected non-empty path for executable %s", tt.executable)
			}
		})
	}
}

// Benchmark tests
func BenchmarkParseDuration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseDuration("30 minutes")
	}
}

func BenchmarkCalculatePauseUntil(b *testing.B) {
	duration := 2 * time.Hour
	for i := 0; i < b.N; i++ {
		CalculatePauseUntil(duration, "14:30")
	}
}
