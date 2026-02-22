package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wellrcosta/bulkcaller/internal/config"
	"github.com/wellrcosta/bulkcaller/internal/reader"
)

var version = "1.0.0"

func main() {
	var filePath string
	var versionFlag bool

	flag.StringVar(&filePath, "file", "", "Path to file")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()

	if versionFlag {
		fmt.Printf("bulkcaller %s\n", version)
		os.Exit(0)
	}

	if filePath == "" {
		log.Fatal("file is required")
	}

	// Test config
	cfg := config.New()
	log.Printf("Config created: %+v", cfg)

	// Test reader
	records, err := reader.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	log.Printf("Read %d records", len(records))
}
