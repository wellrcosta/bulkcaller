package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var version = "1.0.0"

type Reader struct{}

func (r *Reader) ReadCSV(path string) ([][]string, error) {
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

func main() {
	var filePath string
	var versionFlag bool

	flag.StringVar(&filePath, "file", "", "Path to CSV file")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()

	if versionFlag {
		fmt.Printf("bulkcaller %s\n", version)
		return
	}

	if filePath == "" {
		log.Println("Usage: bulkcaller -file <path>")
		return
	}

	r := &Reader{}
	records, err := r.ReadCSV(filePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Read %d records", len(records))
}
