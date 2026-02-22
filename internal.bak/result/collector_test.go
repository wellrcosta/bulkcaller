package result

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector("/tmp/output", true)
	if c == nil {
		t.Fatal("NewCollector returned nil")
	}
	if c.outputDir != "/tmp/output" {
		t.Errorf("outputDir = %v", c.outputDir)
	}
	if !c.printToStdout {
		t.Error("printToStdout should be true")
	}
}

func TestCollector_Init(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewCollector(tmpDir, false)
	
	err := c.Init()
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	
	// Check dir was created
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}

func TestCollector_Init_NoOutput(t *testing.T) {
	c := NewCollector("", false)
	
	err := c.Init()
	if err != nil {
		t.Fatalf("Init() with empty dir should not error: %v", err)
	}
}

func TestCollector_Collect_Success(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewCollector(tmpDir, false)
	c.Init()
	
	c.Collect(5, 200, []byte(`{"success":true}`), nil)
	
	success, fail := c.GetCounts()
	if success != 1 {
		t.Errorf("success = %d, want 1", success)
	}
	if fail != 0 {
		t.Errorf("fail = %d, want 0", fail)
	}
	
	// Check file was created
	expectedFile := filepath.Join(tmpDir, "response_5.json")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Error("Response file was not created")
	}
}

func TestCollector_Collect_Failure(t *testing.T) {
	c := NewCollector("", false)
	
	c.Collect(5, 0, nil, &testError{msg: "connection refused"})
	
	success, fail := c.GetCounts()
	if success != 0 {
		t.Errorf("success = %d, want 0", success)
	}
	if fail != 1 {
		t.Errorf("fail = %d, want 1", fail)
	}
}

func TestCollector_GetCounts(t *testing.T) {
	c := NewCollector("", false)
	c.Collect(1, 200, []byte{}, nil)
	c.Collect(2, 200, []byte{}, nil)
	c.Collect(3, 500, []byte{}, &testError{msg: "error"})
	
	success, fail := c.GetCounts()
	if success != 2 {
		t.Errorf("success = %d, want 2", success)
	}
	if fail != 1 {
		t.Errorf("fail = %d, want 1", fail)
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string { return e.msg }
