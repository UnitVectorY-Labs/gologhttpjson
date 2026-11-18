package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandlerJSONValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
		want int
	}{
		{"object", `{"nested":[true,null,{"value":1}]}`, http.StatusOK},
		{"large number", `123456789012345678901234567890`, http.StatusOK},
		{"null", `null`, http.StatusOK},
		{"duplicate names", `{"outer":{"key":1,"key":2}}`, http.StatusBadRequest},
		{"escaped duplicate names", `{"a":1,"\u0061":2}`, http.StatusBadRequest},
		{"invalid UTF-8", "\"\xff\"", http.StatusBadRequest},
		{"invalid Unicode escape", `"\ud800"`, http.StatusBadRequest},
		{"multiple values", `{} {}`, http.StatusBadRequest},
		{"malformed", `{"a":`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.body))
				w := httptest.NewRecorder()
				handler(w, r)
				if w.Code != tt.want {
					t.Errorf("status = %d, want %d; response: %s", w.Code, tt.want, w.Body.String())
				}
			})
			if tt.want != http.StatusOK {
				if output != "" {
					t.Errorf("invalid request logged: %q", output)
				}
				return
			}
			var entry struct {
				Body jsontext.Value `json:"body"`
				Path string         `json:"path"`
			}
			if err := json.Unmarshal([]byte(output), &entry); err != nil {
				t.Fatalf("invalid log entry %q: %v", output, err)
			}
			if string(entry.Body) != tt.body || entry.Path != "/test" {
				t.Errorf("logged body = %s, path = %q", entry.Body, entry.Path)
			}
		})
	}
}

func TestHandlerOptionalFields(t *testing.T) {
	t.Setenv("LOG_HEADERS", "1")
	t.Setenv("METADATA_TEST", "example")
	output := captureStdout(t, func() {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
		r.Header.Set("X-Test", "value")
		w := httptest.NewRecorder()
		handler(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})
	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatal(err)
	}
	if entry.Headers["X-Test"] != "value" || entry.Metadata["TEST"] != "example" {
		t.Errorf("headers = %v, metadata = %v", entry.Headers, entry.Metadata)
	}
}

func TestLogEntryOmitsEmptyOptionalFields(t *testing.T) {
	for _, entry := range []LogEntry{
		{Body: jsontext.Value(`{}`), Path: "/"},
		{Body: jsontext.Value(`{}`), Path: "/", Headers: map[string]string{}, Metadata: map[string]string{}},
	} {
		output, err := json.Marshal(entry)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]jsontext.Value
		if err := json.Unmarshal(output, &fields); err != nil {
			t.Fatal(err)
		}
		if _, ok := fields["headers"]; ok {
			t.Errorf("empty headers included in %s", output)
		}
		if _, ok := fields["metadata"]; ok {
			t.Errorf("empty metadata included in %s", output)
		}
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = write
	defer func() { os.Stdout = original; read.Close(); write.Close() }()
	f()
	if err := write.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	return string(output)
}
