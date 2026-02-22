package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile_CSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	content := "name,email\nJohn,john@example.com\nJane,jane@example.com\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	records, err := ReadFile(csvPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(records))
	}

	if len(records[0]) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(records[0]))
	}

	if records[0][0] != "name" || records[0][1] != "email" {
		t.Errorf("Header mismatch: got %v", records[0])
	}
}

func TestReadFile_FileNotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/file.csv")
	if err == nil {
		t.Error("Expected error for missing file")
	}
}

func TestReadFile_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	txtPath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(txtPath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ReadFile(txtPath)
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}

func TestReadFile_EmptyCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "empty.csv")
	if err := os.WriteFile(csvPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	records, err := ReadFile(csvPath)
	if err != nil {
		t.Fatalf("Should not error on empty file: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(records))
	}
}
