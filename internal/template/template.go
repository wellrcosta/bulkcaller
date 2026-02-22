package template

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var placeholderRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

// ExtractPlaceholders finds all ${placeholder} patterns in a template
func ExtractPlaceholders(template string) []string {
	matches := placeholderRegex.FindAllStringSubmatch(template, -1)
	seen := make(map[string]bool)
	var result []string
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}
	return result
}

// Substitute replaces ${placeholder} with values from the map
func Substitute(template string, values map[string]string) string {
	result := template
	for key, value := range values {
		placeholder := "${" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// SubstituteWithHeaders replaces placeholders using column headers and row values
func SubstituteWithHeaders(template string, headers []string, row []string) string {
	result := template
	for i, h := range headers {
		if i < len(row) {
			placeholder := "${" + h + "}"
			result = strings.ReplaceAll(result, placeholder, row[i])
		}
	}
	return result
}

// ValidateJSON checks if the string is valid JSON
func ValidateJSON(s string) error {
	var v interface{}
	return json.Unmarshal([]byte(s), &v)
}

// BuildRowMap creates a map from headers and row values
func BuildRowMap(headers []string, row []string, index int) map[string]string {
	result := make(map[string]string)
	result["__index__"] = fmt.Sprintf("%d", index)
	for i, h := range headers {
		if i < len(row) {
			result[h] = row[i]
		}
	}
	return result
}
