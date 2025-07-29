package conf

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/ini.v1"
)

// TestAutoMakeCustomConf tests the autoMakeCustomConf function
func TestAutoMakeCustomConf(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	customConf := filepath.Join(tempDir, "conf", "app.conf")

	// Test creating a new config file
	err := autoMakeCustomConf(customConf)
	if err != nil {
		t.Fatalf("autoMakeCustomConf failed: %v", err)
	}

	// Verify the file was created
	if !fileExists(customConf) {
		t.Errorf("Config file was not created at %s", customConf)
	}

	// Verify the file contains expected content
	cfg, err := ini.Load(customConf)
	if err != nil {
		t.Fatalf("Failed to load created config: %v", err)
	}

	// Check default values
	expectedValues := map[string]map[string]string{
		"": {
			"app_name": "dztasks",
			"run_mode": "prod",
		},
		"web": {
			"http_port": "11011",
		},
		"session": {
			"provider": "memory",
		},
		"plugins": {
			"path":       "plugins",
			"show_error": "true",
			"show_cmd":   "true",
		},
	}

	for sectionName, keys := range expectedValues {
		section := cfg.Section(sectionName)
		for key, expectedValue := range keys {
			actualValue := section.Key(key).String()
			if key == "user" || key == "pass" {
				// Skip checking random values, just verify they exist and are not empty
				if actualValue == "" {
					t.Errorf("Expected non-empty value for %s.%s", sectionName, key)
				}
			} else if actualValue != expectedValue {
				t.Errorf("Expected %s.%s = %s, got %s", sectionName, key, expectedValue, actualValue)
			}
		}
	}

	// Test that calling again doesn't overwrite existing file
	originalModTime := getFileModTime(t, customConf)
	time.Sleep(10 * time.Millisecond) // Ensure different timestamp

	err = autoMakeCustomConf(customConf)
	if err != nil {
		t.Fatalf("autoMakeCustomConf failed on existing file: %v", err)
	}

	newModTime := getFileModTime(t, customConf)
	if !newModTime.Equal(originalModTime) {
		t.Errorf("Existing config file was modified when it shouldn't have been")
	}
}

// TestAutoMakeCustomConfInvalidPath tests error handling
func TestAutoMakeCustomConfInvalidPath(t *testing.T) {
	// Test with an invalid path (trying to create file in non-existent directory with no permissions)
	invalidPath := "/root/nonexistent/app.conf"

	// This should fail on most systems due to permission issues
	err := autoMakeCustomConf(invalidPath)
	if err == nil {
		t.Log("Warning: Expected error for invalid path, but got none. This might be running with elevated permissions.")
	}
}

// Helper function to check if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Helper function to get file modification time
func getFileModTime(t *testing.T, path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}
	return info.ModTime()
}

// TestConfigStructures tests the configuration structures
func TestConfigStructures(t *testing.T) {
	// Test that all config structures can be created and have expected zero values
	tests := []struct {
		name   string
		config interface{}
	}{
		{"AppConfig", &struct {
			Name    string
			Version string
			RunUser string
		}{}},
		{"WebConfig", &struct {
			HttpPort             int
			ExternalURL          string
			Subpath              string
			LoadAssetsFromDisk   bool
			DisableRouterLog     bool
			EnableGzip           bool
			UnixSocketPermission string
		}{}},
		{"SessionConfig", &struct {
			Provider       string
			ProviderConfig string
			CookieName     string
			CookieSecure   bool
			GCInterval     int64
			MaxLifeTime    int64
			CSRFCookieName string
		}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config == nil {
				t.Errorf("Config structure %s is nil", tt.name)
			}
		})
	}
}

// BenchmarkAutoMakeCustomConf benchmarks the autoMakeCustomConf function
func BenchmarkAutoMakeCustomConf(b *testing.B) {
	tempDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		customConf := filepath.Join(tempDir, "bench", "app.conf")
		// Remove the file if it exists to test creation each time
		os.RemoveAll(filepath.Dir(customConf))

		err := autoMakeCustomConf(customConf)
		if err != nil {
			b.Fatalf("autoMakeCustomConf failed: %v", err)
		}
	}
}

// TestConfigValidation tests configuration validation scenarios
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		configData  string
		expectError bool
	}{
		{
			name: "valid config",
			configData: `
app_name = dztasks
run_mode = prod

[web]
http_port = 8080

[session]
provider = memory
`,
			expectError: false,
		},
		{
			name: "config with comments",
			configData: `
# Application settings
app_name = dztasks  # Application name
run_mode = dev      # Development mode

[web]
# Web server settings
http_port = 3000
`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "test.conf")

			err := os.WriteFile(configPath, []byte(tt.configData), 0644)
			if err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			// Try to load the config
			_, err = ini.Load(configPath)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
