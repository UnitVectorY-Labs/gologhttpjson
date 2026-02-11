package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

// Version is the application version, injected at build time via ldflags
var Version = "dev"

type LogEntry struct {
	Body     json.RawMessage   `json:"body"`
	Headers  map[string]string `json:"headers,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Path     string            `json:"path"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusBadRequest)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}

	// Validate that the body is valid JSON
	var jsonBody json.RawMessage
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		http.Error(w, "Request body must be valid JSON", http.StatusBadRequest)
		return
	}

	// Create log entry
	logEntry := LogEntry{
		Body: jsonBody,
		Path: r.URL.Path,
	}

	// Check if header logging is enabled
	if os.Getenv("LOG_HEADERS") != "" {
		headers := make(map[string]string)
		for key, values := range r.Header {
			headers[key] = values[0] // Take the first value for each header
		}
		logEntry.Headers = headers
	}

	// Collect metadata from environment variables
	metadata := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 && strings.HasPrefix(pair[0], "METADATA_") {
			key := strings.TrimPrefix(pair[0], "METADATA_")
			metadata[key] = pair[1]
		}
	}
	if len(metadata) > 0 {
		logEntry.Metadata = metadata
	}

	// Serialize log entry to JSON
	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		http.Error(w, "Unable to create log entry", http.StatusInternalServerError)
		return
	}

	// Log the JSON without timestamp
	os.Stdout.Write(logJSON)
	os.Stdout.WriteString("\n")

	// Add application version to X-App-Version header
	w.Header().Set("X-App-Version", Version)

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK\n"))
}

func main() {
	// Set the build version from the build info if not set by the build system
	if Version == "dev" || Version == "" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
				Version = bi.Main.Version
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}
