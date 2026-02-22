package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/wellrcosta/bulkcaller/internal/reader"
)

var version = "1.0.0"

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

	records, err := reader.ReadCSV(filePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Read %d records", len(records))
}
