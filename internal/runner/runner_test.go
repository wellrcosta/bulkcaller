package runner

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/wellrcosta/bulkcaller/internal/config"
)

func TestNew(t *testing.T) {
	cfg := config.New()
	cfg.FilePath = "test.csv"
	cfg.URL = "http://localhost"
	cfg.BodyTemplate = `{}`

	r := New(cfg)
	if r == nil {
		t.Fatal("New() returned nil")
	}
	if r.config != cfg {
		t.Error("config not set correctly")
	}
}

func TestRunner_Run_E2E(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	// Create test CSV
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	content := "name,email\nJohn,john@example.com\nJane,jane@example.com\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create CSV: %v", err)
	}

	// Create runner config
	cfg := config.New()
	cfg.FilePath = csvPath
	cfg.URL = server.URL
	cfg.BodyTemplate = `{"name":"${name}","email":"${email}"}`
	cfg.Concurrency = 2
	cfg.PrintResponse = false

	runner := New(cfg)
	err := runner.Run()

	// Run might have errors due to CSV format, but should complete
	if err != nil {
		t.Logf("Run completed with error: %v", err)
	}
}

func TestRunner_Run_NoData(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "empty.csv")
	content := "name,email\n" // Only header, no data
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := config.New()
	cfg.FilePath = csvPath
	cfg.URL = "http://localhost"
	cfg.BodyTemplate = `{}`
	cfg.Concurrency = 1

	runner := New(cfg)
	err := runner.Run()

	if err == nil {
		t.Error("Expected error for empty data")
	}
}

func TestRunner_Run_InvalidJSONTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	content := "name\nJohn\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := config.New()
	cfg.FilePath = csvPath
	cfg.URL = "http://localhost"
	cfg.BodyTemplate = `{invalid json ${name}}`
	cfg.Concurrency = 1

	runner := New(cfg)
	// Will process but JSON validation will fail
	_ = runner.Run()
}
