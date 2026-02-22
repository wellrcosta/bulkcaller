package config

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() returned nil")
	}
	if c.Method != "POST" {
		t.Errorf("Method = %v, want POST", c.Method)
	}
	if c.Concurrency != 10 {
		t.Errorf("Concurrency = %v, want 10", c.Concurrency)
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", c.Timeout)
	}
	if c.MaxRetries != 3 {
		t.Errorf("MaxRetries = %v, want 3", c.MaxRetries)
	}
}

func TestParseKeyValue(t *testing.T) {
	tests := []struct {
		input string
		sep   string
		want  map[string]string
	}{
		{
			input: "Content-Type:application/json,Accept:*/*",
			sep:   ":",
			want:  map[string]string{"Content-Type": "application/json", "Accept": "*/*"},
		},
		{
			input: "key1=value1,key2=value2",
			sep:   "=",
			want:  map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			input: "",
			sep:   ":",
			want:  map[string]string{},
		},
	}

	for _, tt := range tests {
		got := parseKeyValue(tt.input, tt.sep)
		if len(got) != len(tt.want) {
			t.Errorf("parseKeyValue(%q, %q) = %v, want %v", tt.input, tt.sep, got, tt.want)
		}
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				FilePath:     "data.csv",
				URL:          "https://api.example.com",
				BodyTemplate: `{"name": "${name}"}`,
			},
			wantErr: false,
		},
		{
			name: "missing file",
			config: &Config{
				URL:          "https://api.example.com",
				BodyTemplate: `{"name": "${name}"}`,
			},
			wantErr: true,
		},
		{
			name: "missing URL",
			config: &Config{
				FilePath:     "data.csv",
				BodyTemplate: `{"name": "${name}"}`,
			},
			wantErr: true,
		},
		{
			name: "missing body",
			config: &Config{
				FilePath: "data.csv",
				URL:      "https://api.example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"30s", 30 * time.Second},
		{"5m", 5 * time.Minute},
		{"1h", time.Hour},
		{"", 30 * time.Second},  // default
		{"invalid", 30 * time.Second},  // default
	}

	for _, tt := range tests {
		got := parseDuration(tt.input)
		if got != tt.want {
			t.Errorf("parseDuration(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseFlags(t *testing.T) {
	c := New()
	args := []string{"-file", "data.csv", "-url", "https://api.com", "-body", `{"test": true}`}
	
	err := c.ParseFlags(args)
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	
	if c.FilePath != "data.csv" {
		t.Errorf("FilePath = %v, want data.csv", c.FilePath)
	}
	if c.URL != "https://api.com" {
		t.Errorf("URL = %v, want https://api.com", c.URL)
	}
}

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	if v == "" {
		t.Error("GetVersion() returned empty string")
	}
}
