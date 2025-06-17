package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/developmeh/mock-cors-server/internal/config"
)

// ResponseData represents the structure of our JSON response
type ResponseData struct {
	Status    string `json:"status"`
	Challenge string `json:"challenge"`
	Timestamp string `json:"timestamp"`
	ExpiresIn int    `json:"expiresIn"`
	SessionID string `json:"sessionId"`
}

// Server represents the HTTP server
type Server struct {
	config *config.Config
	mux    *http.ServeMux
}

// New creates a new server with the given configuration
func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		mux:    http.NewServeMux(),
	}
}

// loggingMiddleware logs HTTP requests to stdout
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the request details
		fmt.Printf("[%s] %s %s %s\n",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
		)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the response time
		fmt.Printf("[%s] Completed in %v\n",
			time.Now().Format(time.RFC3339),
			time.Since(start),
		)
	})
}

// setCORSHeaders sets CORS headers based on configuration
func (s *Server) setCORSHeaders(w http.ResponseWriter, r *http.Request, routeCORS *config.CORSConfig) {
	// Use route-specific CORS if provided, otherwise use global CORS
	cors := s.config.CORS
	if routeCORS != nil {
		cors = *routeCORS
	}

	// Set CORS headers
	origin := r.Header.Get("Origin")
	if origin != "" && (contains(cors.AllowOrigins, origin) || contains(cors.AllowOrigins, "*")) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", joinStrings(cors.AllowMethods))
	w.Header().Set("Access-Control-Allow-Headers", joinStrings(cors.AllowHeaders))

	if cors.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cors.MaxAge))
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// joinStrings joins a slice of strings with commas
func joinStrings(slice []string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

// setupRoutes sets up the routes based on configuration
func (s *Server) setupRoutes() {
	// Set up routes from configuration
	for _, route := range s.config.Routes {
		routePath := route.Path
		routeCORS := route.CORS
		routeType := route.Type
		filePath := route.FilePath
		jsonContent := route.JSONContent
		contentType := route.ContentType

		// Default content type based on route type
		if contentType == "" {
			switch routeType {
			case "static":
				contentType = s.getContentTypeFromFile(filePath)
			case "json", "dummy":
				contentType = "application/json"
			default:
				contentType = "application/json"
			}
		}

		// Create a handler for this route
		s.mux.HandleFunc(routePath, func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			s.setCORSHeaders(w, r, routeCORS)

			// Handle OPTIONS method (CORS preflight)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Handle different route types
			switch routeType {
			case "static":
				s.handleStaticFile(w, r, filePath, contentType)
			case "json":
				s.handleJSONBlob(w, r, jsonContent, contentType)
			case "dummy":
				s.handleDummyResponse(w, r, contentType)
			default:
				// Default to dummy response for backward compatibility
				s.handleDummyResponse(w, r, contentType)
			}
		})
	}
}

// getContentTypeFromFile determines content type based on file extension
func (s *Server) getContentTypeFromFile(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".txt":
		return "text/plain"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

// handleStaticFile serves a static file
func (s *Server) handleStaticFile(w http.ResponseWriter, r *http.Request, filePath, contentType string) {
	// Only allow GET and HEAD methods for static files
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set content type header
	w.Header().Set("Content-Type", contentType)

	// Set status code
	w.WriteHeader(http.StatusOK)

	// Copy file content to response (skip for HEAD requests)
	if r.Method != http.MethodHead {
		io.Copy(w, file)
	}
}

// handleJSONBlob serves a JSON blob from configuration
func (s *Server) handleJSONBlob(w http.ResponseWriter, r *http.Request, jsonContent, contentType string) {
	// Only allow POST method for JSON blob endpoints
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate JSON content
	if jsonContent == "" {
		http.Error(w, "No JSON content configured", http.StatusInternalServerError)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", contentType)

	// Set status code
	w.WriteHeader(http.StatusOK)

	// Write JSON content
	w.Write([]byte(jsonContent))
}

// handleDummyResponse serves the hardcoded dummy response
func (s *Server) handleDummyResponse(w http.ResponseWriter, r *http.Request, contentType string) {
	// Only allow POST method for dummy response
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create dummy response data
	responseData := ResponseData{
		Status:    "success",
		Challenge: "abc123xyz789",
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		ExpiresIn: 300, // 5 minutes in seconds
		SessionID: "12345-67890-abcde-fghij",
	}

	// Set content type header
	w.Header().Set("Content-Type", contentType)

	// Set status code
	w.WriteHeader(http.StatusOK)

	// Encode and send the response
	json.NewEncoder(w).Encode(responseData)
}

// Start starts the server
func (s *Server) Start() error {
	// Set up routes
	s.setupRoutes()

	// Wrap the mux with the logging middleware
	loggingHandler := s.loggingMiddleware(s.mux)

	// Start the server
	addr := fmt.Sprintf(":%d", s.config.Port)
	fmt.Printf("Server running on http://localhost:%d\n", s.config.Port)
	return http.ListenAndServe(addr, loggingHandler)
}
