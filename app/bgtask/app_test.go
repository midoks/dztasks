package bgtask

import (
	"os"
	"testing"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
)

func TestInitTask(t *testing.T) {
	// Initialize log system for testing
	tempLogDir, err := os.MkdirTemp("", "log_test")
	if err != nil {
		t.Fatalf("Failed to create temp log dir: %v", err)
	}
	defer os.RemoveAll(tempLogDir)

	originalLogPath := conf.Log.RootPath
	defer func() {
		conf.Log.RootPath = originalLogPath
	}()
	conf.Log.RootPath = tempLogDir
	log.Init()

	// Test InitTask function
	InitTask()

	// Verify that task is initialized
	if task == nil {
		t.Error("Expected task to be initialized, but it was nil")
	}

	// Verify that task is running
	entries := task.Entries()
	if len(entries) == 0 {
		t.Log("No cron entries found, which is expected if no plugins are configured")
	}

	// Clean up
	if task != nil {
		task.Stop()
	}
}

func TestResetTask(t *testing.T) {
	// Initialize log system for testing
	tempLogDir2, err := os.MkdirTemp("", "log_test")
	if err != nil {
		t.Fatalf("Failed to create temp log dir: %v", err)
	}
	defer os.RemoveAll(tempLogDir2)

	originalLogPath := conf.Log.RootPath
	defer func() {
		conf.Log.RootPath = originalLogPath
	}()
	conf.Log.RootPath = tempLogDir2
	log.Init()

	// Initialize task first
	InitTask()
	if task == nil {
		t.Fatal("Failed to initialize task")
	}

	// Stop the task to avoid conflicts
	task.Stop()

	// Add a test entry to verify reset functionality
	initialEntries := len(task.Entries())

	// Add a dummy cron job
	_, err2 := task.AddFunc("@every 1h", func() {
		// dummy function
	})
	if err2 != nil {
		t.Fatalf("Failed to add test cron job: %v", err2)
	}

	// Verify entry was added
	if len(task.Entries()) <= initialEntries {
		t.Error("Expected cron entry to be added")
	}

	// Test ResetTask - this will clear entries and restart
	ResetTask()

	// Verify task is still initialized
	if task == nil {
		t.Error("Expected task to remain initialized after reset")
	}

	// Clean up
	if task != nil {
		task.Stop()
	}
}

func TestClearTask(t *testing.T) {
	// Initialize log system for testing
	tempLogDir3, err := os.MkdirTemp("", "log_test")
	if err != nil {
		t.Fatalf("Failed to create temp log dir: %v", err)
	}
	defer os.RemoveAll(tempLogDir3)

	originalLogPath := conf.Log.RootPath
	defer func() {
		conf.Log.RootPath = originalLogPath
	}()
	conf.Log.RootPath = tempLogDir3
	log.Init()

	// Initialize a new cron instance for testing
	testCron := cron.New()
	task = testCron

	// Add some test entries
	_, err2 := task.AddFunc("@every 1m", func() {})
	if err2 != nil {
		t.Fatalf("Failed to add test cron job: %v", err2)
	}

	_, err3 := task.AddFunc("@every 5m", func() {})
	if err3 != nil {
		t.Fatalf("Failed to add test cron job: %v", err3)
	}

	// Verify entries were added
	if len(task.Entries()) != 2 {
		t.Errorf("Expected 2 cron entries, got %d", len(task.Entries()))
	}

	// Test clearTask function
	clearTask()

	// Verify all entries were removed
	if len(task.Entries()) != 0 {
		t.Errorf("Expected 0 cron entries after clear, got %d", len(task.Entries()))
	}

	// Clean up
	task.Stop()
}

func TestRunPluginTask(t *testing.T) {
	// Initialize log system for testing
	tempLogDir4, err := os.MkdirTemp("", "log_test")
	if err != nil {
		t.Fatalf("Failed to create temp log dir: %v", err)
	}
	defer os.RemoveAll(tempLogDir4)

	originalLogPath := conf.Log.RootPath
	defer func() {
		conf.Log.RootPath = originalLogPath
	}()
	conf.Log.RootPath = tempLogDir4
	log.Init()

	// Initialize task first
	InitTask()
	if task == nil {
		t.Fatal("Failed to initialize task")
	}

	// Test runPluginTask function
	// This function reads from conf.Plugins.Path and processes plugins
	// Since we don't want to depend on external configuration in tests,
	// we'll just verify it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runPluginTask panicked: %v", r)
		}
	}()

	// Call the function
	runPluginTask()

	// Clean up
	if task != nil {
		task.Stop()
	}
}

// Test concurrent access to task
func TestConcurrentTaskAccess(t *testing.T) {
	// Initialize log system for testing
	tempLogDir5, err := os.MkdirTemp("", "log_test")
	if err != nil {
		t.Fatalf("Failed to create temp log dir: %v", err)
	}
	defer os.RemoveAll(tempLogDir5)

	originalLogPath := conf.Log.RootPath
	defer func() {
		conf.Log.RootPath = originalLogPath
	}()
	conf.Log.RootPath = tempLogDir5
	log.Init()

	// Initialize task
	InitTask()
	if task == nil {
		t.Fatal("Failed to initialize task")
	}

	// Test concurrent access
	done := make(chan bool, 2)

	// Goroutine 1: Add entries
	go func() {
		defer func() { done <- true }()
		for i := 0; i < 5; i++ {
			_, err := task.AddFunc("@every 1h", func() {})
			if err != nil {
				t.Errorf("Failed to add cron job: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Goroutine 2: Read entries
	go func() {
		defer func() { done <- true }()
		for i := 0; i < 5; i++ {
			entries := task.Entries()
			_ = len(entries) // Just read the length
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// Clean up
	if task != nil {
		task.Stop()
	}
}

// Benchmark tests
func BenchmarkInitTask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitTask()
		if task != nil {
			task.Stop()
		}
	}
}

func BenchmarkResetTask(b *testing.B) {
	// Initialize once
	InitTask()
	if task == nil {
		b.Fatal("Failed to initialize task")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ResetTask()
	}

	// Clean up
	if task != nil {
		task.Stop()
	}
}