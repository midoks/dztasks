package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/midoks/dztasks/app/bgtask"
	"github.com/midoks/dztasks/internal/conf"
	"github.com/midoks/dztasks/internal/log"
)

func TestInitWatch(t *testing.T) {
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
	bgtask.InitTask()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "watch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test JSON file
	testFile := filepath.Join(tempDir, "test.json")
	err = os.WriteFile(testFile, []byte(`{"test": "data"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test InitWatch function
	// Since InitWatch runs indefinitely, we'll test it in a goroutine
	// and verify it doesn't panic during initialization
	done := make(chan bool, 1)
	panicked := false

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			done <- true
		}()

		// Run InitWatch for a short time
		InitWatch(tempDir)
	}()

	// Wait a short time for initialization
	select {
	case <-done:
		if panicked {
			t.Error("InitWatch panicked during initialization")
		}
	case <-time.After(100 * time.Millisecond):
		// InitWatch is running as expected (it should run indefinitely)
		// This is the expected behavior
	}
}

func TestInitWatchWithNonExistentDirectory(t *testing.T) {
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
	bgtask.InitTask()

	// Test InitWatch with a non-existent directory
	nonExistentDir := "/non/existent/directory"

	// Test that InitWatch doesn't panic with non-existent directory
	done := make(chan bool, 1)
	panicked := false

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			done <- true
		}()

		// Run InitWatch for a short time
		InitWatch(nonExistentDir)
	}()

	// Wait a short time for initialization
	select {
	case <-done:
		if panicked {
			t.Error("InitWatch panicked with non-existent directory")
		}
	case <-time.After(100 * time.Millisecond):
		// InitWatch is running as expected
	}
}

func TestInitWatchWithEmptyDirectory(t *testing.T) {
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
	bgtask.InitTask()

	// Create a temporary empty directory for testing
	tempDir, err := os.MkdirTemp("", "watch_test_empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test InitWatch with empty directory
	done := make(chan bool, 1)
	panicked := false

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			done <- true
		}()

		// Run InitWatch for a short time
		InitWatch(tempDir)
	}()

	// Wait a short time for initialization
	select {
	case <-done:
		if panicked {
			t.Error("InitWatch panicked with empty directory")
		}
	case <-time.After(100 * time.Millisecond):
		// InitWatch is running as expected
	}
}

func TestInitWatchWithMultipleJSONFiles(t *testing.T) {
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
	bgtask.InitTask()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "watch_test_multiple")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple test JSON files
	testFiles := []string{"test1.json", "test2.json", "subdir/test3.json"}
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		// Create subdirectory if needed
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		// Create the file
		if err := os.WriteFile(filePath, []byte(`{"test": "data"}`), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}
	}

	// Create some non-JSON files that should be ignored
	nonJSONFiles := []string{"test.txt", "test.go", "README.md"}
	for _, file := range nonJSONFiles {
		filePath := filepath.Join(tempDir, file)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create non-JSON file %s: %v", filePath, err)
		}
	}

	// Test InitWatch function
	done := make(chan bool, 1)
	panicked := false

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			done <- true
		}()

		// Run InitWatch for a short time
		InitWatch(tempDir)
	}()

	// Wait a short time for initialization
	select {
	case <-done:
		if panicked {
			t.Error("InitWatch panicked with multiple JSON files")
		}
	case <-time.After(100 * time.Millisecond):
		// InitWatch is running as expected
	}
}

func TestInitWatchFilePermissions(t *testing.T) {
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
	bgtask.InitTask()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "watch_test_perms")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test JSON file with restricted permissions
	testFile := filepath.Join(tempDir, "restricted.json")
	err = os.WriteFile(testFile, []byte(`{"test": "data"}`), 0000) // No permissions
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test InitWatch function with restricted file
	done := make(chan bool, 1)
	panicked := false

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			done <- true
		}()

		// Run InitWatch for a short time
		InitWatch(tempDir)
	}()

	// Wait a short time for initialization
	select {
	case <-done:
		if panicked {
			t.Error("InitWatch panicked with restricted file permissions")
		}
	case <-time.After(100 * time.Millisecond):
		// InitWatch is running as expected
	}

	// Restore file permissions for cleanup
	os.Chmod(testFile, 0644)
}

// Test concurrent access to InitWatch
func TestInitWatchConcurrent(t *testing.T) {
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
	bgtask.InitTask()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "watch_test_concurrent")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test JSON file
	testFile := filepath.Join(tempDir, "test.json")
	err = os.WriteFile(testFile, []byte(`{"test": "data"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test multiple concurrent InitWatch calls
	const numGoroutines = 3
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Goroutine %d panicked: %v", id, r)
				}
				done <- true
			}()

			// Run InitWatch for a very short time
			InitWatch(tempDir)
		}(i)
	}

	// Wait a short time for all goroutines to start
	time.Sleep(50 * time.Millisecond)

	// We don't wait for all goroutines to complete since InitWatch runs indefinitely
	// The test passes if no panics occur during the initial setup
}

// Benchmark test
func BenchmarkInitWatch(b *testing.B) {
	// Create a temporary directory for benchmarking
	tempDir, err := os.MkdirTemp("", "watch_bench")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple test JSON files
	for i := 0; i < 10; i++ {
		testFile := filepath.Join(tempDir, fmt.Sprintf("test%d.json", i))
		os.WriteFile(testFile, []byte(`{"test": "data"}`), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Since InitWatch runs indefinitely, we'll just test the initialization part
		// by running it in a goroutine and stopping it quickly
		go func() {
			InitWatch(tempDir)
		}()
		time.Sleep(1 * time.Millisecond) // Very short sleep to allow initialization
	}
}
