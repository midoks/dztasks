package context

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/macaron.v1"
)

// TestJSONData tests the JSONData structure
func TestJSONData(t *testing.T) {
	tests := []struct {
		name     string
		data     JSONData
		expected string
	}{
		{
			name: "success response",
			data: JSONData{
				Code: 0,
				Msg:  "success",
				Data: "test data",
			},
			expected: `{"code":0,"msg":"success","data":"test data"}`,
		},
		{
			name: "error response",
			data: JSONData{
				Code: -1,
				Msg:  "error occurred",
				Data: nil,
			},
			expected: `{"code":-1,"msg":"error occurred"}`,
		},
		{
			name: "response with complex data",
			data: JSONData{
				Code: 1,
				Msg:  "ok",
				Data: map[string]interface{}{
					"id":   123,
					"name": "test",
				},
			},
			expected: `{"code":1,"msg":"ok","data":{"id":123,"name":"test"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.data)
			if err != nil {
				t.Fatalf("Failed to marshal JSONData: %v", err)
			}

			if string(jsonBytes) != tt.expected {
				t.Errorf("JSONData marshal = %s, want %s", string(jsonBytes), tt.expected)
			}
		})
	}
}

// TestLayuiData tests the LayuiData structure
func TestLayuiData(t *testing.T) {
	tests := []struct {
		name     string
		data     LayuiData
		expected string
	}{
		{
			name: "table data response",
			data: LayuiData{
				Code:  0,
				Count: 100,
				Msg:   "success",
				Data:  []string{"item1", "item2"},
			},
			expected: `{"code":0,"count":100,"msg":"success","data":["item1","item2"]}`,
		},
		{
			name: "empty table response",
			data: LayuiData{
				Code:  0,
				Count: 0,
				Msg:   "no data",
				Data:  []interface{}{},
			},
			expected: `{"code":0,"count":0,"msg":"no data","data":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.data)
			if err != nil {
				t.Fatalf("Failed to marshal LayuiData: %v", err)
			}

			if string(jsonBytes) != tt.expected {
				t.Errorf("LayuiData marshal = %s, want %s", string(jsonBytes), tt.expected)
			}
		})
	}
}

// TestContextJSONMethods tests the JSON response methods
func TestContextJSONMethods(t *testing.T) {
	// Create a test macaron instance with proper middleware
	m := macaron.New()
	m.Use(macaron.Renderer())

	// Test ReturnJSON method
	m.Get("/test-json", func(ctx *macaron.Context) {
		c := &Context{Context: ctx}
		c.ReturnJSON(0, "success", "test data")
	})

	// Test Ok method
	m.Get("/test-ok", func(ctx *macaron.Context) {
		c := &Context{Context: ctx}
		c.Ok("operation successful")
	})

	// Test Fail method
	m.Get("/test-fail", func(ctx *macaron.Context) {
		c := &Context{Context: ctx}
		c.Fail("operation failed")
	})

	// Test ReturnLayuiJSON method
	m.Get("/test-layui", func(ctx *macaron.Context) {
		c := &Context{Context: ctx}
		c.ReturnLayuiJSON(0, "success", 10, []string{"item1", "item2"})
	})

	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "ReturnJSON method",
			path:         "/test-json",
			expectedCode: http.StatusOK,
			expectedBody: `{"code":0,"msg":"success","data":"test data"}`,
		},
		{
			name:         "Ok method",
			path:         "/test-ok",
			expectedCode: http.StatusOK,
			expectedBody: `{"code":1,"msg":"operation successful","data":""}`,
		},
		{
			name:         "Fail method",
			path:         "/test-fail",
			expectedCode: http.StatusOK,
			expectedBody: `{"code":-1,"msg":"operation failed","data":""}`,
		},
		{
			name:         "ReturnLayuiJSON method",
			path:         "/test-layui",
			expectedCode: http.StatusOK,
			expectedBody: `{"code":0,"count":10,"msg":"success","data":["item1","item2"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, http.NoBody)
			w := httptest.NewRecorder()

			m.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			body := w.Body.String()
			if body != tt.expectedBody {
				t.Errorf("Expected body %s, got %s", tt.expectedBody, body)
			}

			// Verify content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json; charset=UTF-8" {
				t.Errorf("Expected content type 'application/json; charset=UTF-8', got '%s'", contentType)
			}
		})
	}
}

// TestContextHelperMethods tests various helper methods
func TestContextHelperMethods(t *testing.T) {
	m := macaron.New()
	m.Use(macaron.Renderer())

	// Test PageIs method
	m.Get("/test-page-is", func(ctx *macaron.Context) {
		c := &Context{Context: ctx}
		c.Data = make(map[string]interface{})
		c.PageIs("Home")
		if !c.Data["PageIsHome"].(bool) {
			t.Error("PageIs method did not set the correct data")
		}
		c.ReturnJSON(0, "ok", nil)
	})

	// Test Require method
	m.Get("/test-require", func(ctx *macaron.Context) {
		c := &Context{Context: ctx}
		c.Data = make(map[string]interface{})
		c.Require("jQuery")
		if !c.Data["RequirejQuery"].(bool) {
			t.Error("Require method did not set the correct data")
		}
		c.ReturnJSON(0, "ok", nil)
	})

	// Test FormErr method
	m.Get("/test-form-err", func(ctx *macaron.Context) {
		c := &Context{Context: ctx}
		c.Data = make(map[string]interface{})
		c.FormErr("username", "password")
		if !c.Data["Err_username"].(bool) || !c.Data["Err_password"].(bool) {
			t.Error("FormErr method did not set the correct data")
		}
		c.ReturnJSON(0, "ok", nil)
	})

	tests := []string{"/test-page-is", "/test-require", "/test-form-err"}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, http.NoBody)
			w := httptest.NewRecorder()

			m.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})
	}
}
