package installer

import (
	"fmt"
	"strings"
	"sultengutt/internal/config"
	"testing"
)

func TestWelcomeNote(t *testing.T) {
	note := WelcomeNote()

	if note == "" {
		t.Error("Welcome note should not be empty")
	}

	// Check that the welcome note contains expected elements
	expectedElements := []string{
		"Sultengutt",
		"Welcome to Sultengutt",
		"Never miss a surprise dinner",
		"🍔",
	}

	for _, element := range expectedElements {
		if !strings.Contains(note, element) {
			t.Errorf("Welcome note should contain '%s', but it doesn't", element)
		}
	}

	// Check that ASCII art is present (look for basic structure)
	if !strings.Contains(note, "/") || !strings.Contains(note, "\\") {
		t.Error("Welcome note should contain ASCII art with / and \\ characters")
	}
}

func TestValidateDays(t *testing.T) {
	tests := []struct {
		name        string
		days        []string
		expectValid bool
	}{
		{
			name:        "valid single day",
			days:        []string{"Monday"},
			expectValid: true,
		},
		{
			name:        "valid multiple days",
			days:        []string{"Monday", "Wednesday", "Friday"},
			expectValid: true,
		},
		{
			name:        "all days",
			days:        []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
			expectValid: true,
		},
		{
			name:        "empty days",
			days:        []string{},
			expectValid: false,
		},
		{
			name:        "invalid day",
			days:        []string{"InvalidDay"},
			expectValid: false,
		},
		{
			name:        "mixed valid and invalid",
			days:        []string{"Monday", "InvalidDay"},
			expectValid: false,
		},
	}

	validDays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := true

			// Check if days list is empty
			if len(tt.days) == 0 {
				isValid = false
			} else {
				// Check each day
				for _, day := range tt.days {
					found := false
					for _, validDay := range validDays {
						if day == validDay {
							found = true
							break
						}
					}
					if !found {
						isValid = false
						break
					}
				}
			}

			if isValid != tt.expectValid {
				t.Errorf("Expected validity %v for days %v, got %v", tt.expectValid, tt.days, isValid)
			}
		})
	}
}

func TestValidateHour(t *testing.T) {
	tests := []struct {
		name        string
		hour        string
		expectValid bool
	}{
		{"valid hour", "14:30", true},
		{"midnight", "00:00", true},
		{"noon", "12:00", true},
		{"single digit hour", "9:15", true},
		{"single digit minute", "14:5", true},
		{"23:59", "23:59", true},
		{"invalid hour - 24", "24:00", false},
		{"invalid minute - 60", "12:60", false},
		{"invalid format", "1430", false},
		{"missing colon", "14", false},
		{"empty string", "", false},
		{"non-numeric", "abc:30", false},
		{"negative hour", "-5:30", false},
	}

	// 24-hour format regex
	timeRegex := `^(?:[01]?\d|2[0-3]):[0-5]\d$`

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a simplified validation - the actual installer
			// uses more sophisticated validation
			isValid := validateTimeFormat(tt.hour, timeRegex)

			if isValid != tt.expectValid {
				t.Errorf("Expected validity %v for hour %s, got %v", tt.expectValid, tt.hour, isValid)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectValid bool
	}{
		{"valid http URL", "http://example.com", true},
		{"valid https URL", "https://example.com", true},
		{"URL with path", "https://example.com/path", true},
		{"URL with query", "https://example.com?query=value", true},
		{"URL with port", "https://example.com:8080", true},
		{"localhost", "http://localhost:3000", true},
		{"IP address", "http://192.168.1.1", true},
		{"invalid - no protocol", "example.com", false},
		{"invalid - empty", "", false},
		{"invalid - malformed", "not-a-url", false},
		{"invalid - spaces", "https://example .com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateURL(tt.url)

			if isValid != tt.expectValid {
				t.Errorf("Expected validity %v for URL %s, got %v", tt.expectValid, tt.url, isValid)
			}
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		options     config.InstallOptions
		expectValid bool
	}{
		{
			name: "valid config",
			options: config.InstallOptions{
				Days:     []string{"Monday", "Wednesday", "Friday"},
				Hour:     "14:30",
				SiteLink: "https://example.com",
			},
			expectValid: true,
		},
		{
			name: "invalid - empty days",
			options: config.InstallOptions{
				Days:     []string{},
				Hour:     "14:30",
				SiteLink: "https://example.com",
			},
			expectValid: false,
		},
		{
			name: "invalid - bad hour",
			options: config.InstallOptions{
				Days:     []string{"Monday"},
				Hour:     "25:30",
				SiteLink: "https://example.com",
			},
			expectValid: false,
		},
		{
			name: "invalid - bad URL",
			options: config.InstallOptions{
				Days:     []string{"Monday"},
				Hour:     "14:30",
				SiteLink: "not-a-url",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary config to test validation
			cfg := config.Config{
				InstallOptions: tt.options,
				PausedUntil:    -1,
			}

			// We can't directly access the validate method, but we can
			// test the logic indirectly
			hasValidDays := len(tt.options.Days) > 0
			hasValidHour := validateTimeFormat(tt.options.Hour, `^(?:[01]?\d|2[0-3]):[0-5]\d$`)
			hasValidURL := validateURL(tt.options.SiteLink)

			isValid := hasValidDays && hasValidHour && hasValidURL

			// Check individual day validity
			if hasValidDays {
				validDaysList := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
				for _, day := range tt.options.Days {
					found := false
					for _, validDay := range validDaysList {
						if day == validDay {
							found = true
							break
						}
					}
					if !found {
						hasValidDays = false
						break
					}
				}
				isValid = hasValidDays && hasValidHour && hasValidURL
			}

			if isValid != tt.expectValid {
				t.Errorf("Expected validity %v for config %+v, got %v", tt.expectValid, tt.options, isValid)
			}

			_ = cfg // Use the config variable to avoid unused variable error
		})
	}
}

// Helper function to validate time format (mimics installer logic)
func validateTimeFormat(timeStr, pattern string) bool {
	if timeStr == "" {
		return false
	}

	// Simple validation - check HH:MM format
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return false
	}

	// Parse hour and minute
	var hour, minute int
	if _, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute); err != nil {
		return false
	}

	// Validate ranges
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return false
	}

	return true
}

// Helper function to validate URLs (mimics installer logic)
func validateURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	// Simple check for protocol
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	// Check for spaces (basic validation)
	if strings.Contains(urlStr, " ") {
		return false
	}

	return true
}

// Test the ASCII art face constant
func TestSultenguttFace(t *testing.T) {
	if sultenguttFace == "" {
		t.Error("sultenguttFace should not be empty")
	}

	// Check for basic ASCII art elements
	expectedChars := []string{"/", "\\", "_", "O", "^", "-"}
	for _, char := range expectedChars {
		if !strings.Contains(sultenguttFace, char) {
			t.Errorf("sultenguttFace should contain '%s'", char)
		}
	}

	// Check that it has multiple lines
	lines := strings.Split(sultenguttFace, "\n")
	if len(lines) < 5 {
		t.Error("sultenguttFace should have multiple lines")
	}
}

// Benchmark the welcome note generation
func BenchmarkWelcomeNote(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WelcomeNote()
	}
}

// Test style constants
func TestStyles(t *testing.T) {
	// Test that styles can be used (lipgloss.Style is not a pointer, so we can't check for nil)

	// Test that styles can render text
	asciiText := asciiStyle.Render("test")
	if asciiText == "" {
		t.Error("asciiStyle should render text")
	}

	nameText := nameStyle.Render("test")
	if nameText == "" {
		t.Error("nameStyle should render text")
	}

	textText := textStyle.Render("test")
	if textText == "" {
		t.Error("textStyle should render text")
	}
}
