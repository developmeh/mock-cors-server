package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/developmeh/mock-cors-server/internal/config"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Port: 8081,
		Routes: []config.Route{
			{
				Path: "/test",
				Type: "dummy",
			},
		},
	}

	server := New(cfg)
	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.config != cfg {
		t.Error("Expected server config to match provided config")
	}

	if server.mux == nil {
		t.Error("Expected server mux to be initialized")
	}
}

func TestGetContentTypeFromFile(t *testing.T) {
	server := &Server{}

	tests := []struct {
		filePath   string
		expectedCT string
	}{
		{"test.html", "text/html"},
		{"test.htm", "text/html"},
		{"test.css", "text/css"},
		{"test.js", "application/javascript"},
		{"test.json", "application/json"},
		{"test.xml", "application/xml"},
		{"test.txt", "text/plain"},
		{"test.png", "image/png"},
		{"test.jpg", "image/jpeg"},
		{"test.jpeg", "image/jpeg"},
		{"test.gif", "image/gif"},
		{"test.svg", "image/svg+xml"},
		{"test.pdf", "application/pdf"},
		{"test.unknown", "application/octet-stream"},
		{"test", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			ct := server.getContentTypeFromFile(tt.filePath)
			if ct != tt.expectedCT {
				t.Errorf("Expected content type %s for %s, got %s", tt.expectedCT, tt.filePath, ct)
			}
		})
	}
}

func TestHandleDummyResponse(t *testing.T) {
	server := &Server{}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectJSON     bool
	}{
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			expectJSON:     true,
		},
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
			expectJSON:     false,
		},
		{
			name:           "PUT request",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectJSON:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			server.handleDummyResponse(w, req, "application/json")

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectJSON {
				var response ResponseData
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Expected valid JSON response, got error: %v", err)
				}

				if response.Status != "success" {
					t.Errorf("Expected status 'success', got %s", response.Status)
				}

				if response.Challenge == "" {
					t.Error("Expected challenge to be set")
				}

				if response.SessionID == "" {
					t.Error("Expected session ID to be set")
				}

				if response.ExpiresIn != 300 {
					t.Errorf("Expected expires in 300, got %d", response.ExpiresIn)
				}
			}
		})
	}
}

func TestHandleJSONBlob(t *testing.T) {
	server := &Server{}

	tests := []struct {
		name           string
		method         string
		jsonContent    string
		expectedStatus int
		expectContent  bool
	}{
		{
			name:           "POST with valid JSON",
			method:         http.MethodPost,
			jsonContent:    `{"test": true, "message": "hello"}`,
			expectedStatus: http.StatusOK,
			expectContent:  true,
		},
		{
			name:           "POST with empty JSON content",
			method:         http.MethodPost,
			jsonContent:    "",
			expectedStatus: http.StatusInternalServerError,
			expectContent:  false,
		},
		{
			name:           "GET request",
			method:         http.MethodGet,
			jsonContent:    `{"test": true}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectContent:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			server.handleJSONBlob(w, req, tt.jsonContent, "application/json")

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectContent {
				body := w.Body.String()
				if body != tt.jsonContent {
					t.Errorf("Expected body %s, got %s", tt.jsonContent, body)
				}

				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected content type application/json, got %s", contentType)
				}
			}
		})
	}
}

func TestHandleStaticFile(t *testing.T) {
	server := &Server{}

	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	testContent := "Hello, this is a test file!"
	if _, err := tempFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name           string
		method         string
		filePath       string
		expectedStatus int
		expectContent  bool
	}{
		{
			name:           "GET existing file",
			method:         http.MethodGet,
			filePath:       tempFile.Name(),
			expectedStatus: http.StatusOK,
			expectContent:  true,
		},
		{
			name:           "HEAD existing file",
			method:         http.MethodHead,
			filePath:       tempFile.Name(),
			expectedStatus: http.StatusOK,
			expectContent:  false,
		},
		{
			name:           "GET non-existing file",
			method:         http.MethodGet,
			filePath:       "/non/existing/file.txt",
			expectedStatus: http.StatusNotFound,
			expectContent:  false,
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			filePath:       tempFile.Name(),
			expectedStatus: http.StatusMethodNotAllowed,
			expectContent:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			server.handleStaticFile(w, req, tt.filePath, "text/plain")

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectContent {
				body := w.Body.String()
				if body != testContent {
					t.Errorf("Expected body %s, got %s", testContent, body)
				}

				contentType := w.Header().Get("Content-Type")
				if contentType != "text/plain" {
					t.Errorf("Expected content type text/plain, got %s", contentType)
				}
			} else if tt.method == http.MethodHead && tt.expectedStatus == http.StatusOK {
				// HEAD requests should have no body
				if w.Body.Len() != 0 {
					t.Error("Expected empty body for HEAD request")
				}
			}
		})
	}
}

func TestSetCORSHeaders(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowOrigins:     []string{"https://example.com", "*"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           3600,
		},
	}

	server := &Server{config: cfg}

	tests := []struct {
		name         string
		origin       string
		routeCORS    *config.CORSConfig
		expectOrigin string
	}{
		{
			name:         "allowed origin",
			origin:       "https://example.com",
			expectOrigin: "https://example.com",
		},
		{
			name:         "wildcard origin",
			origin:       "https://test.com",
			expectOrigin: "https://test.com",
		},
		{
			name:         "no origin header",
			origin:       "",
			expectOrigin: "",
		},
		{
			name:   "route-specific CORS",
			origin: "https://route-specific.com",
			routeCORS: &config.CORSConfig{
				AllowOrigins:     []string{"https://route-specific.com"},
				AllowMethods:     []string{"POST"},
				AllowHeaders:     []string{"Content-Type"},
				AllowCredentials: false,
				MaxAge:           1800,
			},
			expectOrigin: "https://route-specific.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodOptions, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			w := httptest.NewRecorder()

			server.setCORSHeaders(w, req, tt.routeCORS)

			// Check Access-Control-Allow-Origin
			allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if allowOrigin != tt.expectOrigin {
				t.Errorf("Expected Access-Control-Allow-Origin %s, got %s", tt.expectOrigin, allowOrigin)
			}

			// Check other CORS headers
			allowMethods := w.Header().Get("Access-Control-Allow-Methods")
			if allowMethods == "" {
				t.Error("Expected Access-Control-Allow-Methods to be set")
			}

			allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
			if allowHeaders == "" {
				t.Error("Expected Access-Control-Allow-Headers to be set")
			}

			maxAge := w.Header().Get("Access-Control-Max-Age")
			if maxAge == "" {
				t.Error("Expected Access-Control-Max-Age to be set")
			}

			// Check route-specific CORS
			if tt.routeCORS != nil {
				if tt.routeCORS.AllowCredentials {
					allowCredentials := w.Header().Get("Access-Control-Allow-Credentials")
					if allowCredentials != "true" {
						t.Error("Expected Access-Control-Allow-Credentials to be true")
					}
				}
			} else {
				// Global CORS allows credentials
				allowCredentials := w.Header().Get("Access-Control-Allow-Credentials")
				if allowCredentials != "true" {
					t.Error("Expected Access-Control-Allow-Credentials to be true")
				}
			}
		})
	}
}

func TestUtilityFunctions(t *testing.T) {
	// Test contains function
	t.Run("contains function", func(t *testing.T) {
		slice := []string{"a", "b", "c"}

		if !contains(slice, "b") {
			t.Error("Expected contains to return true for existing item")
		}

		if contains(slice, "d") {
			t.Error("Expected contains to return false for non-existing item")
		}

		if contains([]string{}, "a") {
			t.Error("Expected contains to return false for empty slice")
		}
	})

	// Test joinStrings function
	t.Run("joinStrings function", func(t *testing.T) {
		tests := []struct {
			input    []string
			expected string
		}{
			{[]string{"a", "b", "c"}, "a, b, c"},
			{[]string{"single"}, "single"},
			{[]string{}, ""},
			{[]string{"one", "two"}, "one, two"},
		}

		for _, tt := range tests {
			result := joinStrings(tt.input)
			if result != tt.expected {
				t.Errorf("Expected joinStrings(%v) = %s, got %s", tt.input, tt.expected, result)
			}
		}
	})
}

func TestSetupRoutes(t *testing.T) {
	// Create test static file
	tempFile, err := os.CreateTemp("", "test*.html")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	testHTML := "<html><body>Test</body></html>"
	if _, err := tempFile.WriteString(testHTML); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	cfg := &config.Config{
		Port: 8081,
		Routes: []config.Route{
			{
				Path:        "/dummy",
				Type:        "dummy",
				ContentType: "application/json",
			},
			{
				Path:     "/static",
				Type:     "static",
				FilePath: tempFile.Name(),
			},
			{
				Path:        "/json",
				Type:        "json",
				JSONContent: `{"test": true}`,
				ContentType: "application/json",
			},
		},
		CORS: config.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "OPTIONS"},
			AllowHeaders: []string{"Content-Type"},
		},
	}

	server := New(cfg)
	server.setupRoutes()

	// Test dummy route
	t.Run("dummy route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/dummy", nil)
		w := httptest.NewRecorder()

		server.mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response ResponseData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Expected valid JSON response, got error: %v", err)
		}
	})

	// Test static route
	t.Run("static route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/static", nil)
		w := httptest.NewRecorder()

		server.mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if body != testHTML {
			t.Errorf("Expected body %s, got %s", testHTML, body)
		}
	})

	// Test JSON blob route
	t.Run("json blob route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json", nil)
		w := httptest.NewRecorder()

		server.mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		expected := `{"test": true}`
		if body != expected {
			t.Errorf("Expected body %s, got %s", expected, body)
		}
	})

	// Test OPTIONS requests (CORS preflight)
	t.Run("CORS preflight", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/dummy", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		server.mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
		if allowOrigin != "https://example.com" {
			t.Errorf("Expected Access-Control-Allow-Origin https://example.com, got %s", allowOrigin)
		}
	})
}
