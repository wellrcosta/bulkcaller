package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Reader defines the interface for reading data files
type Reader interface {
	Read(path string) ([][]string, error)
}

// ReadFile reads a CSV, XLS, or XLSX file and returns rows as string slices
func ReadFile(path string) ([][]string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".csv":
		return readCSV(path)
	case ".xls", ".xlsx":
		return readExcel(path)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

func readCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable fields per record

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
