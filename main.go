package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wellrcosta/bulkcaller/internal/config"
	"github.com/wellrcosta/bulkcaller/internal/runner"
)

var version = "1.0.0"

func main() {
	// Parse command line flags using standard flag package
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
	flag.StringVar(&cfg.OutputDir, "output", "", "Output directory for responses (optional)")
	flag.IntVar(&cfg.MaxRetries, "retries", 3, "Max retries on failure")
	flag.BoolVar(&cfg.PrintResponse, "print", false, "Print responses to stdout")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()

	// Show version
	if versionFlag {
		fmt.Printf("bulkcaller %s\n", version)
		os.Exit(0)
	}

	// Parse custom flags
	if headersStr != "" {
		cfg.Headers = parseKeyValue(headersStr, ":")
	}
	if queryStr != "" {
		cfg.QueryParams = parseKeyValue(queryStr, "=")
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n\n", err)
		printUsage()
		os.Exit(1)
	}

	// Build full URL with query params
	cfg.URL = cfg.GetHTTPURL()

	// Run
	r := runner.New(cfg)
	if err := r.Run(); err != nil {
		log.Fatalf("❌ Error: %v", err)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: bulkcaller -file <data> -url <endpoint> -body <template>\n\n")
	fmt.Fprintf(os.Stderr, "Example:\n")
	fmt.Fprintf(os.Stderr, "  bulkcaller -file data.csv -url https://api.example.com -body '{\"name\":\"${name}\"}'\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

// parseKeyValue parses key:value or key=value pairs
func parseKeyValue(s, sep string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}
	for _, pair := range rangeSplit(s, ",") {
		parts := splitN(pair, sep, 2)
		if len(parts) == 2 {
			result[trim(parts[0])] = trim(parts[1])
		}
	}
	return result
}

// Helper functions to avoid importing strings for simple operations
func rangeSplit(s, sep string) []string {
	// Simple split implementation
	if s == "" {
		return nil
	}
	result := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if i < len(s)-len(sep)+1 && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}
	result = append(result, s[start:])
	return result
}

func splitN(s, sep string, n int) []string {
	if n <= 0 {
		return []string{s}
	}
	result := []string{}
	start := 0
	for i := 0; i < len(s) && len(result) < n-1; i++ {
		if i < len(s)-len(sep)+1 && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}
	result = append(result, s[start:])
	return result
}

func trim(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
