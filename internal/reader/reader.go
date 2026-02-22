package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func ReadCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening CSV: %w", err)
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
