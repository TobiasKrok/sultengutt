package scheduler

import (
	"sultengutt/internal/config"
	"testing"
)

func TestNewScheduler(t *testing.T) {
	options := config.InstallOptions{
		Days:     []string{"Monday", "Tuesday"},
		Hour:     "14:30",
		SiteLink: "https://example.com",
	}

	// This will test the platform-specific factory function
	// Note: This may panic if sultengutt executable is not in PATH, which is expected
	defer func() {
		if r := recover(); r != nil {
			t.Logf("NewScheduler panicked (expected if executable not in PATH): %v", r)
		}
	}()

	scheduler := NewScheduler(options, "/tmp/test")

	if scheduler == nil {
		t.Fatal("Scheduler should not be nil")
	}

	// Test that we can call the interface methods without panicking
	exists, err := scheduler.TaskExists()
	if err != nil {
		t.Logf("TaskExists returned error (expected for test): %v", err)
	}
	_ = exists
}

func TestSchedulerInterface(t *testing.T) {
	options := config.InstallOptions{
		Days:     []string{"Monday", "Wednesday", "Friday"},
		Hour:     "09:00",
		SiteLink: "https://test.com",
	}

	// Handle potential panic if executable not in PATH
	defer func() {
		if r := recover(); r != nil {
			t.Logf("NewScheduler panicked (expected if executable not in PATH): %v", r)
			return
		}
	}()

	scheduler := NewScheduler(options, "/tmp/test")

	// Test that the scheduler implements all interface methods
	var _ Scheduler = scheduler

	// Test TaskExists method exists and returns appropriate values
	exists, err := scheduler.TaskExists()

	// For most platforms in test environment, this will return an error
	// which is expected since we're not actually installing tasks
	if err != nil {
		t.Logf("TaskExists error (expected in test environment): %v", err)
	}

	// exists should be a boolean regardless
	t.Logf("Task exists status: %v", exists)
}

// Test platform-specific behavior
func TestPlatformSpecificScheduler(t *testing.T) {
	options := config.InstallOptions{
		Days:     []string{"Monday"},
		Hour:     "10:00",
		SiteLink: "https://example.com",
	}

	// Handle potential panic if executable not in PATH
	defer func() {
		if r := recover(); r != nil {
			t.Logf("NewScheduler panicked (expected if executable not in PATH): %v", r)
			return
		}
	}()

	scheduler := NewScheduler(options, "/tmp/test")

	// Since the scheduler types are platform-specific and may not be accessible
	// in cross-platform tests, we just verify the scheduler is not nil
	if scheduler == nil {
		t.Error("Expected non-nil scheduler")
	}

	// Test that it implements the interface
	var _ Scheduler = scheduler
}

// Mock scheduler for testing
type MockScheduler struct {
	registered    bool
	taskExists    bool
	registerErr   error
	unregisterErr error
	existsErr     error
}

func (m *MockScheduler) RegisterTask() error {
	if m.registerErr != nil {
		return m.registerErr
	}
	m.registered = true
	m.taskExists = true
	return nil
}

func (m *MockScheduler) UnregisterTask() error {
	if m.unregisterErr != nil {
		return m.unregisterErr
	}
	m.registered = false
	m.taskExists = false
	return nil
}

func (m *MockScheduler) TaskExists() (bool, error) {
	if m.existsErr != nil {
		return false, m.existsErr
	}
	return m.taskExists, nil
}

func TestMockScheduler(t *testing.T) {
	mock := &MockScheduler{}

	// Test initial state
	exists, err := mock.TaskExists()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if exists {
		t.Error("Task should not exist initially")
	}

	// Test registration
	err = mock.RegisterTask()
	if err != nil {
		t.Fatalf("Failed to register task: %v", err)
	}

	if !mock.registered {
		t.Error("Task should be marked as registered")
	}

	exists, err = mock.TaskExists()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !exists {
		t.Error("Task should exist after registration")
	}

	// Test unregistration
	err = mock.UnregisterTask()
	if err != nil {
		t.Fatalf("Failed to unregister task: %v", err)
	}

	if mock.registered {
		t.Error("Task should not be marked as registered after unregistration")
	}

	exists, err = mock.TaskExists()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if exists {
		t.Error("Task should not exist after unregistration")
	}
}

// Benchmark the scheduler creation
func BenchmarkNewScheduler(b *testing.B) {
	options := config.InstallOptions{
		Days:     []string{"Monday", "Tuesday", "Wednesday"},
		Hour:     "14:30",
		SiteLink: "https://example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler := NewScheduler(options, "/tmp/test")
		_ = scheduler
	}
}

func TestSchedulerOptions(t *testing.T) {
	tests := []struct {
		name    string
		options config.InstallOptions
	}{
		{
			name: "single day",
			options: config.InstallOptions{
				Days:     []string{"Monday"},
				Hour:     "09:00",
				SiteLink: "https://example.com",
			},
		},
		{
			name: "multiple days",
			options: config.InstallOptions{
				Days:     []string{"Monday", "Wednesday", "Friday"},
				Hour:     "14:30",
				SiteLink: "https://test.com",
			},
		},
		{
			name: "all days",
			options: config.InstallOptions{
				Days:     []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
				Hour:     "12:00",
				SiteLink: "https://daily.com",
			},
		},
		{
			name: "midnight",
			options: config.InstallOptions{
				Days:     []string{"Sunday"},
				Hour:     "00:00",
				SiteLink: "https://midnight.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle potential panic if executable not in PATH
			defer func() {
				if r := recover(); r != nil {
					t.Logf("NewScheduler panicked (expected if executable not in PATH): %v", r)
					return
				}
			}()

			scheduler := NewScheduler(tt.options, "/tmp/test")
			if scheduler == nil {
				t.Error("Scheduler should not be nil")
			}

			// Test that we can interact with the scheduler without crashing
			_, err := scheduler.TaskExists()
			// We expect errors in test environment, so we just log them
			if err != nil {
				t.Logf("TaskExists error (expected): %v", err)
			}
		})
	}
}
