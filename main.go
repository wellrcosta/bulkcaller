package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var version = "dev"

type Config struct {
	FilePath      string
	URL           string
	Method        string
	BodyTemplate  string
	Headers       map[string]string
	Concurrency   int
	DelayMs       int
	OutputDir     string
	PrintResponse bool
}

func main() {
	var config Config
	var headersStr string
	var showVersion bool

	flag.StringVar(&config.FilePath, "file", "", "Path to CSV file")
	flag.StringVar(&config.URL, "url", "", "Target URL")
	flag.StringVar(&config.Method, "method", "POST", "HTTP method")
	flag.StringVar(&config.BodyTemplate, "body", "", "JSON template with ${col}")
	flag.StringVar(&headersStr, "headers", "", "Headers as key:value")
	flag.IntVar(&config.Concurrency, "concurrency", 10, "Concurrent workers")
	flag.IntVar(&config.DelayMs, "delay", 0, "Delay in ms")
	flag.StringVar(&config.OutputDir, "output", "", "Output dir")
	flag.BoolVar(&config.PrintResponse, "print", false, "Print responses")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("bulkcaller %s\n", version)
		return
	}

	if config.FilePath == "" || config.URL == "" || config.BodyTemplate == "" {
		fmt.Println("Usage: bulkcaller -file <csv> -url <url> -body '{\"key\":\"${col}\"}'")
		flag.PrintDefaults()
		os.Exit(1)
	}

	config.Headers = parseHeaders(headersStr)

	if err := run(config); err != nil {
		log.Fatal(err)
	}
}

func parseHeaders(s string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}
	for _, pair := range strings.Split(s, ",") {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	return reader.ReadAll()
}

func run(config Config) error {
	records, err := readCSV(config.FilePath)
	if err != nil {
		return fmt.Errorf("reading CSV: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("no data rows")
	}

	headers := records[0]

	client := &http.Client{Timeout: 30 * time.Second}
	var wg sync.WaitGroup
	jobs := make(chan map[string]string, config.Concurrency*2)

	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				body := config.BodyTemplate
				for i, h := range headers {
					if i < len(row) {
						body = strings.ReplaceAll(body, "${"+h+"}", row[i])
					}
				}

				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(body), &payload); err != nil {
					log.Printf("Invalid JSON: %v", err)
					continue
				}

				req, _ := http.NewRequest(config.Method, config.URL, strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				for k, v := range config.Headers {
					req.Header.Set(k, v)
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Request failed: %v", err)
					continue
				}
				respBody, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				if config.PrintResponse {
					fmt.Printf("Response: %s\n", string(respBody))
				}
				_= respBody

				if config.DelayMs > 0 {
					time.Sleep(time.Duration(config.DelayMs) * time.Millisecond)
				}
			}
		}()
	}

	for _, row := range records[1:] {
		rowMap := make(map[string]string)
		for i, h := range headers {
			if i < len(row) {
				rowMap[h] = row[i]
			}
		}
		jobs <- rowMap
	}
	close(jobs)
	wg.Wait()

	return nil
}
