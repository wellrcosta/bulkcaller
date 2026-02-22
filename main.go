package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wellrcosta/bulkcaller/internal/config"
	"github.com/wellrcosta/bulkcaller/internal/runner"
)

var version = "1.0.0"

func main() {
	cfg := config.New()

	var headersStr, queryStr string
	var versionFlag bool

	flag.StringVar(&cfg.FilePath, "file", "", "Path to CSV/XLS/XLSX file")
	flag.StringVar(&cfg.URL, "url", "", "Target URL")
	flag.StringVar(&cfg.Method, "method", "POST", "HTTP method (GET, POST, PUT, PATCH, DELETE)")
	flag.StringVar(&cfg.BodyTemplate, "body", "", "JSON template with placeholders like ${columnName}")
	flag.StringVar(&headersStr, "headers", "", "Headers as key:value pairs, comma-separated")
	flag.StringVar(&queryStr, "query", "", "Query params as key=value pairs, comma-separated")
	flag.IntVar(&cfg.Concurrency, "concurrency", 10, "Number of concurrent workers")
	flag.IntVar(&cfg.Delay, "delay", 0, "Delay between requests in milliseconds")
	flag.StringVar(&cfg.OutputDir, "output", "", "Output directory for responses")
	flag.IntVar(&cfg.MaxRetries, "retries", 3, "Max retries on failure")
	flag.BoolVar(&cfg.PrintResponse, "print", false, "Print responses to stdout")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()

	if versionFlag {
		fmt.Printf("bulkcaller %s\n", version)
		os.Exit(0)
	}

	if headersStr != "" {
		cfg.Headers = parseKeyValue(headersStr, ":")
	}
	if queryStr != "" {
		cfg.QueryParams = parseKeyValue(queryStr, "=")
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\nUsage:\n", err)
		fmt.Fprintf(os.Stderr, "  bulkcaller -file data.csv -url https://api.example.com -body '{\"name\":\"${name}\"}'\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg.URL = cfg.GetHTTPURL()

	r := runner.New(cfg)
	if err := r.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func parseKeyValue(s, sep string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}
	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), sep, 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}
