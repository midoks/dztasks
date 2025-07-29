package tools

import (
	"reflect"
	"testing"
)

// TestToSlice tests the ToSlice function
func TestToSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int64
		hasError bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []int64{},
			hasError: false,
		},
		{
			name:     "single number",
			input:    "123",
			expected: []int64{123},
			hasError: false,
		},
		{
			name:     "multiple numbers",
			input:    "1,2,3,4,5",
			expected: []int64{1, 2, 3, 4, 5},
			hasError: false,
		},
		{
			name:     "numbers with spaces",
			input:    "1, 2, 3, 4, 5",
			expected: []int64{1, 2, 3, 4, 5},
			hasError: false,
		},
		{
			name:     "negative numbers",
			input:    "-1,-2,3",
			expected: []int64{-1, -2, 3},
			hasError: false,
		},
		{
			name:     "empty elements",
			input:    "1,,3",
			expected: []int64{1, 3},
			hasError: false,
		},
		{
			name:     "invalid number",
			input:    "1,abc,3",
			expected: nil,
			hasError: true,
		},
		{
			name:     "large numbers",
			input:    "9223372036854775807,-9223372036854775808",
			expected: []int64{9223372036854775807, -9223372036854775808},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToSlice(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToSlice(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckStringIsExist tests the CheckStringIsExist function
func TestCheckStringIsExist(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		check    []string
		expected bool
	}{
		{
			name:     "string exists",
			source:   "apple",
			check:    []string{"apple", "banana", "cherry"},
			expected: true,
		},
		{
			name:     "string does not exist",
			source:   "grape",
			check:    []string{"apple", "banana", "cherry"},
			expected: false,
		},
		{
			name:     "empty slice",
			source:   "apple",
			check:    []string{},
			expected: false,
		},
		{
			name:     "empty string in slice",
			source:   "",
			check:    []string{"", "apple"},
			expected: true,
		},
		{
			name:     "case sensitive",
			source:   "Apple",
			check:    []string{"apple", "banana"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckStringIsExist(tt.source, tt.check)
			if result != tt.expected {
				t.Errorf("CheckStringIsExist(%q, %v) = %v, want %v", tt.source, tt.check, result, tt.expected)
			}
		})
	}
}

// TestStringInSlice tests the StringInSlice function (alias)
func TestStringInSlice(t *testing.T) {
	// Test that StringInSlice works the same as CheckStringIsExist
	str := "test"
	slice := []string{"test", "example"}

	result1 := CheckStringIsExist(str, slice)
	result2 := StringInSlice(str, slice)

	if result1 != result2 {
		t.Errorf("StringInSlice and CheckStringIsExist should return the same result")
	}
}

// TestMd5 tests the Md5 function
func TestMd5(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "hello world",
			input:    "hello world",
			expected: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
		{
			name:     "test string",
			input:    "test",
			expected: "098f6bcd4621d373cade4e832627b4f6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Md5(tt.input)
			if result != tt.expected {
				t.Errorf("Md5(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRandString tests the RandString function
func TestRandString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 0", 0},
		{"length 1", 1},
		{"length 8", 8},
		{"length 16", 16},
		{"length 32", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandString(tt.length)
			if len(result) != tt.length {
				t.Errorf("RandString(%d) returned string of length %d, want %d", tt.length, len(result), tt.length)
			}

			// Test that multiple calls return different strings (for length > 0)
			if tt.length > 0 {
				result2 := RandString(tt.length)
				if result == result2 {
					t.Logf("Warning: RandString(%d) returned same string twice: %q", tt.length, result)
				}
			}
		})
	}
}

// TestRemoveDuplicatesAndEmpty tests the RemoveDuplicatesAndEmpty function
// Note: This function only removes consecutive duplicates and empty strings
func TestRemoveDuplicatesAndEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates or empty",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with consecutive duplicates",
			input:    []string{"a", "a", "b", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with non-consecutive duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "a", "c", "b"},
		},
		{
			name:     "with empty strings",
			input:    []string{"a", "", "b", "", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with consecutive duplicates and empty",
			input:    []string{"a", "a", "", "b", "b", "", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string(nil),
		},
		{
			name:     "only empty strings",
			input:    []string{"", "", ""},
			expected: []string(nil),
		},
		{
			name:     "consecutive empty strings",
			input:    []string{"a", "", "", "b"},
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveDuplicatesAndEmpty(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveDuplicatesAndEmpty(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
