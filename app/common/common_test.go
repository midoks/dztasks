package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetPluginInstallLock(t *testing.T) {
	tests := []struct {
		name       string
		pluginName string
		want       string
	}{
		{
			name:       "basic plugin name",
			pluginName: "test-plugin",
			want:       "test-plugin/install.lock",
		},
		{
			name:       "plugin with path",
			pluginName: "/path/to/plugin",
			want:       "/path/to/plugin/install.lock",
		},
		{
			name:       "empty plugin name",
			pluginName: "",
			want:       "/install.lock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPluginInstallLock(tt.pluginName)
			if got != tt.want {
				t.Errorf("GetPluginInstallLock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginList(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "plugin_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test plugin directory with info.json
	pluginDir := filepath.Join(tempDir, "test-plugin")
	err = os.MkdirAll(pluginDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugin dir: %v", err)
	}

	// Create a test info.json file
	testPlugin := Plugin{
		Name:   "Test Plugin",
		Author: "Test Author",
		Ps:     "Test Description",
		Icon:   "test-icon",
	}

	infoJSON, err := json.Marshal(testPlugin)
	if err != nil {
		t.Fatalf("Failed to marshal test plugin: %v", err)
	}

	infoPath := filepath.Join(pluginDir, "info.json")
	err = os.WriteFile(infoPath, infoJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write info.json: %v", err)
	}

	// Test PluginList function
	plugins := PluginList(tempDir)

	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(plugins))
	}

	if len(plugins) > 0 {
		plugin := plugins[0]
		if plugin.Name != "Test Plugin" {
			t.Errorf("Expected plugin name 'Test Plugin', got '%s'", plugin.Name)
		}
		if plugin.Path != "test-plugin" {
			t.Errorf("Expected plugin path 'test-plugin', got '%s'", plugin.Path)
		}
		if plugin.Icon != "test-icon" {
			t.Errorf("Expected plugin icon 'test-icon', got '%s'", plugin.Icon)
		}
		if plugin.Installed {
			t.Errorf("Expected plugin to not be installed")
		}
	}
}

func TestPluginListWithInstallLock(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "plugin_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test plugin directory with info.json
	pluginDir := filepath.Join(tempDir, "test-plugin")
	err = os.MkdirAll(pluginDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugin dir: %v", err)
	}

	// Create a test info.json file
	testPlugin := Plugin{
		Name:   "Test Plugin",
		Author: "Test Author",
		Ps:     "Test Description",
	}

	infoJSON, err := json.Marshal(testPlugin)
	if err != nil {
		t.Fatalf("Failed to marshal test plugin: %v", err)
	}

	infoPath := filepath.Join(pluginDir, "info.json")
	err = os.WriteFile(infoPath, infoJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write info.json: %v", err)
	}

	// Create install.lock file
	lockPath := filepath.Join(pluginDir, "install.lock")
	err = os.WriteFile(lockPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write install.lock: %v", err)
	}

	// Test PluginList function
	plugins := PluginList(tempDir)

	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(plugins))
	}

	if len(plugins) > 0 {
		plugin := plugins[0]
		if !plugin.Installed {
			t.Errorf("Expected plugin to be installed")
		}
	}
}

func TestPluginListEmptyDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "plugin_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test PluginList function with empty directory
	plugins := PluginList(tempDir)

	if len(plugins) != 0 {
		t.Errorf("Expected 0 plugins, got %d", len(plugins))
	}
}

func TestPluginListWithDefaultIcon(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "plugin_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test plugin directory with info.json
	pluginDir := filepath.Join(tempDir, "test-plugin")
	err = os.MkdirAll(pluginDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugin dir: %v", err)
	}

	// Create a test info.json file without icon
	testPlugin := Plugin{
		Name:   "Test Plugin",
		Author: "Test Author",
		Ps:     "Test Description",
		// Icon is empty, should default to "layui-icon-tree"
	}

	infoJSON, err := json.Marshal(testPlugin)
	if err != nil {
		t.Fatalf("Failed to marshal test plugin: %v", err)
	}

	infoPath := filepath.Join(pluginDir, "info.json")
	err = os.WriteFile(infoPath, infoJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write info.json: %v", err)
	}

	// Test PluginList function
	plugins := PluginList(tempDir)

	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(plugins))
	}

	if len(plugins) > 0 {
		plugin := plugins[0]
		if plugin.Icon != "layui-icon-tree" {
			t.Errorf("Expected default icon 'layui-icon-tree', got '%s'", plugin.Icon)
		}
	}
}

func TestExecInput(t *testing.T) {
	tests := []struct {
		name    string
		bin     string
		args    []string
		wantErr bool
	}{
		{
			name:    "echo command",
			bin:     "echo",
			args:    []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			bin:     "nonexistent-command-12345",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ExecInput(tt.bin, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(output) == 0 {
				t.Errorf("ExecInput() expected output but got empty")
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetPluginInstallLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetPluginInstallLock("test-plugin")
	}
}

func BenchmarkPluginList(b *testing.B) {
	// Create a temporary directory for benchmarking
	tempDir, err := os.MkdirTemp("", "plugin_bench")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test plugin directory with info.json
	pluginDir := filepath.Join(tempDir, "test-plugin")
	err = os.MkdirAll(pluginDir, 0755)
	if err != nil {
		b.Fatalf("Failed to create plugin dir: %v", err)
	}

	testPlugin := Plugin{
		Name:   "Test Plugin",
		Author: "Test Author",
	}

	infoJSON, _ := json.Marshal(testPlugin)
	infoPath := filepath.Join(pluginDir, "info.json")
	os.WriteFile(infoPath, infoJSON, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PluginList(tempDir)
	}
}