package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient(30*time.Second, 3)
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.httpClient == nil {
		t.Error("httpClient is nil")
	}
}

func TestDoRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 0)
	result := client.DoRequest("POST", server.URL, nil, `{"test": true}`)

	if result.Error != nil {
		t.Fatalf("Expected no error, got: %v", result.Error)
	}
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", result.StatusCode)
	}
}

func TestDoRequest_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 0)
	result := client.DoRequest("POST", server.URL, nil, `{}`)

	if result.Error == nil {
		t.Error("Expected error for 500 status")
	}
	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", result.StatusCode)
	}
}

func TestDoRequest_WithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer token123" {
			t.Errorf("Expected Authorization header, got: %s", auth)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 0)
	headers := map[string]string{"Authorization": "Bearer token123"}
	result := client.DoRequest("GET", server.URL, headers, "")

	if result.Error != nil {
		t.Fatalf("Expected no error: %v", result.Error)
	}
}

func TestDoRequest_InvalidURL(t *testing.T) {
	client := NewClient(5*time.Second, 0)
	result := client.DoRequest("GET", "://invalid-url", nil, "")

	if result.Error == nil {
		t.Error("Expected error for invalid URL")
	}
}
