package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/wellrcosta/bulkcaller/internal/config"
	"github.com/wellrcosta/bulkcaller/internal/httpclient"
	"github.com/wellrcosta/bulkcaller/internal/reader"
	"github.com/wellrcosta/bulkcaller/internal/result"
	"github.com/wellrcosta/bulkcaller/internal/runner"
	"github.com/wellrcosta/bulkcaller/internal/template"
)

var version = "1.0.0"

func main() {
	var cfg config.Config
	var versionFlag bool

	flag.StringVar(&cfg.FilePath, "file", "", "Path to CSV file")
	flag.StringVar(&cfg.URL, "url", "", "Target URL")
	flag.StringVar(&cfg.Body, "body", "", "Template with ${col}")
	flag.IntVar(&cfg.Delay, "delay", 0, "Delay in ms")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()

	if versionFlag {
		fmt.Printf("bulkcaller %s\n", version)
		return
	}

	if cfg.FilePath == "" {
		fmt.Println("Usage: bulkcaller -file <csv> -url <url> -body '{\"name\":\"${name}\"}'")
		flag.PrintDefaults()
		return
	}

	// Verify all packages work
	_ = config.New()
	_ = reader.ReadCSV
	_ = template.Substitute
	_ = httpclient.DoRequest
	_ = result.New()

	log.Println("Starting bulkcaller...")
	if err := runner.Run(&cfg); err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("Done!")
}
