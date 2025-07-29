package template

import (
	"html/template"
	"reflect"
	"testing"
	"time"
)

var DateFmtMailData = []struct {
	in1 string
	in2 string
	out string
	ok  bool
}{
	{"2021-10-10 13:49:01", "2021-10-10 13:49:01", "13:49", true},
	{"2021-10-09 13:49:01", "2021-10-10 13:49:01", "昨天", true},
	{"2021-10-08 13:49:01", "2021-10-10 13:49:01", "2021-10-08", true},
	{"2021-10-07 13:49:01", "2021-10-10 13:49:01", "2021-10-07", true},
}

func ForDateFmtMail(t time.Time, n time.Time) string {
	in := t.Format("2006-01-02")
	now := n.Format("2006-01-02")

	if in == now {
		return t.Format("15:04")
	}
	in2, _ := time.Parse("2006-01-02 15:04:05", in+" 00:00:00")
	now2, _ := time.Parse("2006-01-02 15:04:05", now+" 00:00:00")
	if in2.Unix()+86400 == now2.Unix() {
		return "昨天"
	} else {
		return t.Format("2006-01-02")
	}
}

// TestDateFmtMail tests the date formatting function
func TestDateFmtMail(t *testing.T) {
	for _, test := range DateFmtMailData {
		in1, _ := time.Parse("2006-01-02 15:04:05", test.in1)
		in2, _ := time.Parse("2006-01-02 15:04:05", test.in2)
		out := ForDateFmtMail(in1, in2)
		if out != test.out {
			t.Errorf("ForDateFmtMail(%+q,%+q) expected %+q; got %+q", test.in1, test.in2, test.out, out)
		}
	}
}

// TestFuncMap tests the template function map
func TestFuncMap(t *testing.T) {
	funcMaps := FuncMap()

	if len(funcMaps) == 0 {
		t.Error("FuncMap should return at least one function map")
	}

	funcMap := funcMaps[0]

	// Test that expected functions exist
	expectedFunctions := []string{
		"BuildCommit",
		"Year",
		"AppSubURL",
		"AppName",
		"AppVer",
		"Safe",
		"Str2HTML",
		"Sanitize",
	}

	for _, funcName := range expectedFunctions {
		if _, exists := funcMap[funcName]; !exists {
			t.Errorf("Expected function %s not found in FuncMap", funcName)
		}
	}

	// Test Year function
	if yearFunc, ok := funcMap["Year"]; ok {
		if yearFn, ok := yearFunc.(func() int); ok {
			year := yearFn()
			currentYear := time.Now().Year()
			if year != currentYear {
				t.Errorf("Year function returned %d, expected %d", year, currentYear)
			}
		} else {
			t.Error("Year function has wrong type")
		}
	}
}

// TestSafe tests the Safe function
func TestSafe(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected template.HTML
	}{
		{
			name:     "plain text",
			input:    "hello world",
			expected: template.HTML("hello world"),
		},
		{
			name:     "html content",
			input:    "<p>Hello <strong>World</strong></p>",
			expected: template.HTML("<p>Hello <strong>World</strong></p>"),
		},
		{
			name:     "empty string",
			input:    "",
			expected: template.HTML(""),
		},
		{
			name:     "special characters",
			input:    "<script>alert('test')</script>",
			expected: template.HTML("<script>alert('test')</script>"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Safe(tt.input)
			if result != tt.expected {
				t.Errorf("Safe(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestStr2HTML tests the Str2HTML function
// Note: Str2HTML uses bluemonday.UGCPolicy().Sanitize() which may modify content
func TestStr2HTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected template.HTML
	}{
		{
			name:     "plain text",
			input:    "hello world",
			expected: template.HTML("hello world"),
		},
		{
			name:     "text with newlines",
			input:    "line1\nline2\nline3",
			expected: template.HTML("line1\nline2\nline3"), // bluemonday preserves newlines
		},
		{
			name:     "empty string",
			input:    "",
			expected: template.HTML(""),
		},
		{
			name:     "html content",
			input:    "<p>Hello <strong>World</strong></p>",
			expected: template.HTML("<p>Hello <strong>World</strong></p>"), // bluemonday allows safe HTML
		},
		{
			name:     "unsafe html",
			input:    "<script>alert('test')</script>",
			expected: template.HTML(""), // bluemonday removes unsafe content
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Str2HTML(tt.input)
			if result != tt.expected {
				t.Errorf("Str2HTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNewLine2br tests the NewLine2br function
func TestNewLine2br(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single newline",
			input:    "line1\nline2",
			expected: "line1<br>line2",
		},
		{
			name:     "multiple newlines",
			input:    "line1\nline2\nline3",
			expected: "line1<br>line2<br>line3",
		},
		{
			name:     "no newlines",
			input:    "single line",
			expected: "single line",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewLine2br(tt.input)
			if result != tt.expected {
				t.Errorf("NewLine2br(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestEscapePound tests the EscapePound function
// Note: EscapePound escapes %, #, space, and ? characters
func TestEscapePound(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "string with pound",
			input:    "test#value",
			expected: "test%23value",
		},
		{
			name:     "multiple pounds",
			input:    "#test#value#",
			expected: "%23test%23value%23",
		},
		{
			name:     "string with space",
			input:    "test value",
			expected: "test%20value",
		},
		{
			name:     "string with question mark",
			input:    "test?value",
			expected: "test%3Fvalue",
		},
		{
			name:     "string with percent",
			input:    "test%value",
			expected: "test%25value",
		},
		{
			name:     "no special characters",
			input:    "testvalue",
			expected: "testvalue",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "all special characters",
			input:    "test%#?value test",
			expected: "test%25%23%3Fvalue%20test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapePound(tt.input)
			if result != tt.expected {
				t.Errorf("EscapePound(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// BenchmarkFuncMap benchmarks the FuncMap function
func BenchmarkFuncMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FuncMap()
	}
}

// BenchmarkSafe benchmarks the Safe function
func BenchmarkSafe(b *testing.B) {
	testString := "<p>This is a test string with <strong>HTML</strong> content</p>"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Safe(testString)
	}
}

// BenchmarkStr2HTML benchmarks the Str2HTML function
func BenchmarkStr2HTML(b *testing.B) {
	testString := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Str2HTML(testString)
	}
}

// TestFuncMapConsistency tests that FuncMap returns consistent results
func TestFuncMapConsistency(t *testing.T) {
	funcMap1 := FuncMap()
	funcMap2 := FuncMap()

	if !reflect.DeepEqual(funcMap1, funcMap2) {
		t.Error("FuncMap should return consistent results across multiple calls")
	}
}
