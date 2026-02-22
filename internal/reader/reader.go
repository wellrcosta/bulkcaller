package reader

import "fmt"

func ReadCSV(path string) ([][]string, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
	return [][]string{{"header"}, {"value"}}, nil
}
