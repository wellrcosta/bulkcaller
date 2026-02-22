package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadFile reads a CSV file and returns rows as string slices
func ReadFile(path string) ([][]string, error) {
	if !strings.HasSuffix(strings.ToLower(path), ".csv") {
		return nil, fmt.Errorf("only CSV files are supported currently")
	}
	return readCSV(path)
}

func readCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading CSV: %w", err)
		}
		records = append(records, record)
	}
	return records, nil
}
