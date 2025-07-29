package router

import (
	"reflect"
	"testing"
)

// TestReverseStrings tests the reverseStrings function
func TestReverseStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []string{"a"},
			expected: []string{"a"},
		},
		{
			name:     "two elements",
			input:    []string{"a", "b"},
			expected: []string{"b", "a"},
		},
		{
			name:     "multiple elements",
			input:    []string{"first", "second", "third", "fourth"},
			expected: []string{"fourth", "third", "second", "first"},
		},
		{
			name:     "odd number of elements",
			input:    []string{"1", "2", "3", "4", "5"},
			expected: []string{"5", "4", "3", "2", "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of input to avoid modifying the test data
			inputCopy := make([]string, len(tt.input))
			copy(inputCopy, tt.input)

			reverseStrings(inputCopy)

			if !reflect.DeepEqual(inputCopy, tt.expected) {
				t.Errorf("reverseStrings(%v) = %v, want %v", tt.input, inputCopy, tt.expected)
			}
		})
	}
}

// TestGetLinesText tests the getLinesText function
func TestGetLinesText(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: "\n",
		},
		{
			name:     "single line",
			input:    []string{"hello"},
			expected: "hello\n",
		},
		{
			name:     "multiple lines",
			input:    []string{"first", "second", "third"},
			expected: "third\nsecond\nfirst\n",
		},
		{
			name:     "log lines example",
			input:    []string{"[INFO] Starting server", "[DEBUG] Loading config", "[ERROR] Connection failed"},
			expected: "[ERROR] Connection failed\n[DEBUG] Loading config\n[INFO] Starting server\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of input since getLinesText modifies the slice
			inputCopy := make([]string, len(tt.input))
			copy(inputCopy, tt.input)

			result := getLinesText(inputCopy)

			if result != tt.expected {
				t.Errorf("getLinesText(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// BenchmarkReverseStrings benchmarks the reverseStrings function
func BenchmarkReverseStrings(b *testing.B) {
	// Create test data
	testData := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		testData[i] = "line" + string(rune(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Make a copy for each iteration
		data := make([]string, len(testData))
		copy(data, testData)
		reverseStrings(data)
	}
}

// BenchmarkGetLinesText benchmarks the getLinesText function
func BenchmarkGetLinesText(b *testing.B) {
	// Create test data
	testData := make([]string, 100)
	for i := 0; i < 100; i++ {
		testData[i] = "This is log line number " + string(rune(i)) + " with some content"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Make a copy for each iteration
		data := make([]string, len(testData))
		copy(data, testData)
		getLinesText(data)
	}
}
