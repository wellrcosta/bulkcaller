package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile_CSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	content := "name,email\nJohn,john@example.com\n"
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	records, err := ReadFile(csvPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(records))
	}
}

func TestReadFile_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	txtPath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(txtPath, []byte("test"), 0644)

	_, err := ReadFile(txtPath)
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}
