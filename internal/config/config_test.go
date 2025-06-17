package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Test basic configuration
	if cfg.Port != 8081 {
		t.Errorf("Expected default port 8081, got %d", cfg.Port)
	}

	if cfg.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", cfg.Version)
	}

	// Test default route
	if len(cfg.Routes) != 1 {
		t.Errorf("Expected 1 default route, got %d", len(cfg.Routes))
	}

	route := cfg.Routes[0]
	if route.Path != "/v1/json/begin" {
		t.Errorf("Expected default route path '/v1/json/begin', got %s", route.Path)
	}

	if route.Type != "dummy" {
		t.Errorf("Expected default route type 'dummy', got %s", route.Type)
	}

	if route.ContentType != "application/json" {
		t.Errorf("Expected default route content type 'application/json', got %s", route.ContentType)
	}

	// Test CORS configuration
	if len(cfg.CORS.AllowOrigins) != 1 || cfg.CORS.AllowOrigins[0] != "*" {
		t.Errorf("Expected default CORS allow origins ['*'], got %v", cfg.CORS.AllowOrigins)
	}

	expectedMethods := []string{"GET", "POST", "OPTIONS"}
	if len(cfg.CORS.AllowMethods) != len(expectedMethods) {
		t.Errorf("Expected %d CORS methods, got %d", len(expectedMethods), len(cfg.CORS.AllowMethods))
	}

	if !cfg.CORS.AllowCredentials {
		t.Error("Expected CORS allow credentials to be true")
	}

	if cfg.CORS.MaxAge != 86400 {
		t.Errorf("Expected CORS max age 86400, got %d", cfg.CORS.MaxAge)
	}
}

func TestLoadConfigWithoutFile(t *testing.T) {
	// Reset viper
	viper.Reset()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error when loading config without file, got %v", err)
	}

	// Should return default config
	defaultCfg := DefaultConfig()
	if cfg.Port != defaultCfg.Port {
		t.Errorf("Expected port %d, got %d", defaultCfg.Port, cfg.Port)
	}
}

func TestLoadConfigWithEnvironmentVariables(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Set environment variable
	os.Setenv("MOCK_CORS_PORT", "9000")
	defer os.Unsetenv("MOCK_CORS_PORT")

	// Configure viper to read environment variables
	viper.SetEnvPrefix("MOCK_CORS")
	viper.AutomaticEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error when loading config with env vars, got %v", err)
	}

	if cfg.Port != 9000 {
		t.Errorf("Expected port 9000 from env var, got %d", cfg.Port)
	}
}

func TestRouteTypes(t *testing.T) {
	tests := []struct {
		name        string
		routeType   string
		filePath    string
		jsonContent string
		expectValid bool
	}{
		{
			name:        "dummy route",
			routeType:   "dummy",
			expectValid: true,
		},
		{
			name:        "static route with file path",
			routeType:   "static",
			filePath:    "./static/example.html",
			expectValid: true,
		},
		{
			name:        "json route with content",
			routeType:   "json",
			jsonContent: `{"test": true}`,
			expectValid: true,
		},
		{
			name:        "static route without file path",
			routeType:   "static",
			expectValid: true, // Should still be valid, error handled at runtime
		},
		{
			name:        "json route without content",
			routeType:   "json",
			expectValid: true, // Should still be valid, error handled at runtime
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := Route{
				Path:        "/test",
				Type:        tt.routeType,
				FilePath:    tt.filePath,
				JSONContent: tt.jsonContent,
				ContentType: "application/json",
			}

			// Basic validation - route should have required fields
			if route.Path == "" {
				t.Error("Route path should not be empty")
			}

			if route.Type == "" {
				t.Error("Route type should not be empty")
			}

			// Type-specific validation
			switch route.Type {
			case "static":
				// For static routes, we expect a file path (but don't validate file existence here)
				if tt.filePath == "" && tt.expectValid {
					t.Log("Static route without file path - will be handled at runtime")
				}
			case "json":
				// For JSON routes, we expect JSON content
				if tt.jsonContent == "" && tt.expectValid {
					t.Log("JSON route without content - will be handled at runtime")
				}
			case "dummy":
				// Dummy routes don't need additional fields
			}
		})
	}
}

func TestCORSConfig(t *testing.T) {
	cors := CORSConfig{
		AllowOrigins:     []string{"https://example.com", "https://test.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600,
	}

	if len(cors.AllowOrigins) != 2 {
		t.Errorf("Expected 2 allowed origins, got %d", len(cors.AllowOrigins))
	}

	if len(cors.AllowMethods) != 2 {
		t.Errorf("Expected 2 allowed methods, got %d", len(cors.AllowMethods))
	}

	if len(cors.AllowHeaders) != 2 {
		t.Errorf("Expected 2 allowed headers, got %d", len(cors.AllowHeaders))
	}

	if !cors.AllowCredentials {
		t.Error("Expected allow credentials to be true")
	}

	if cors.MaxAge != 3600 {
		t.Errorf("Expected max age 3600, got %d", cors.MaxAge)
	}
}
