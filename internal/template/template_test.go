package template

import (
	"testing"
)

func TestExtractPlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		template string
		want     []string
	}{
		{
			name:     "simple placeholder",
			template: `{"name": "${name}"}`,
			want:     []string{"name"},
		},
		{
			name:     "multiple placeholders",
			template: `{"name": "${name}", "email": "${email}"}`,
			want:     []string{"name", "email"},
		},
		{
			name:     "duplicate placeholders",
			template: `{"name": "${name}", "display": "${name}"}`,
			want:     []string{"name"},
		},
		{
			name:     "no placeholders",
			template: `{"name": "John"}`,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractPlaceholders(tt.template)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractPlaceholders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstitute(t *testing.T) {
	values := map[string]string{
		"name":  "John",
		"email": "john@example.com",
	}

	template := `{"name": "${name}", "email": "${email}"}`
	got := Substitute(template, values)
	want := `{"name": "John", "email": "john@example.com"}`

	if got != want {
		t.Errorf("Substitute() = %s, want %s", got, want)
	}
}

func TestSubstituteWithHeaders(t *testing.T) {
	headers := []string{"name", "email"}
	row := []string{"Jane", "jane@example.com"}

	template := `{"name": "${name}", "email": "${email}"}`
	got := SubstituteWithHeaders(template, headers, row)
	want := `{"name": "Jane", "email": "jane@example.com"}`

	if got != want {
		t.Errorf("SubstituteWithHeaders() = %s, want %s", got, want)
	}
}

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid JSON",
			input:   `{"name": "John"}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"name": "John"`,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildRowMap(t *testing.T) {
	headers := []string{"id", "name", "email"}
	row := []string{"1", "John", "john@example.com"}

	got := BuildRowMap(headers, row, 5)

	if got["__index__"] != "5" {
		t.Errorf("BuildRowMap() __index__ = %v, want 5", got["__index__"])
	}
	if got["id"] != "1" {
		t.Errorf("BuildRowMap() id = %v, want 1", got["id"])
	}
	if got["name"] != "John" {
		t.Errorf("BuildRowMap() name = %v, want John", got["name"])
	}
}
