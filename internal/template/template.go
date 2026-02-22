package template

import "strings"

func Substitute(template string, values map[string]string) string {
	result := template
	for key, value := range values {
		result = strings.ReplaceAll(result, "${"+key+"}", value)
	}
	return result
}
