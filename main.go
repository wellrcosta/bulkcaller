package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/wellrcosta/bulkcaller/internal/reader"
)

var version = "dev"

type Config struct {
	FilePath      string
	URL           string
	Method        string
	BodyTemplate  string
	Headers       map[string]string
	QueryParams   map[string]string
	Concurrency   int
	DelayMs       int
	OutputDir     string
	PrintResponse bool
	MaxRetries    int
}

func main() {
	var config Config
	var headersStr, queryStr string
	var showVersion bool

	flag.StringVar(&config.FilePath, "file", "", "Path to CSV/XLS/XLSX file")
	flag.StringVar(&config.URL, "url", "", "Target URL")
	flag.StringVar(&config.Method, "method", "POST", "HTTP method (GET, POST, PUT, PATCH, DELETE)")
	flag.StringVar(&config.BodyTemplate, "body", "", "JSON body template with placeholders like {{columnName}}")
	flag.StringVar(&headersStr, "headers", "", "Headers as key:value pairs, comma-separated")
	flag.StringVar(&queryStr, "query", "", "Query params as key=value pairs, comma-separated")
	flag.IntVar(&config.Concurrency, "concurrency", 10, "Number of concurrent workers")
	flag.IntVar(&config.DelayMs, "delay", 0, "Delay between requests in milliseconds")
	flag.StringVar(&config.OutputDir, "output", "", "Output directory for responses (optional)")
	flag.IntVar(&config.MaxRetries, "retries", 3, "Max retries on failure")
	flag.BoolVar(&config.PrintResponse, "print", false, "Print responses to stdout")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("bulkcaller %s\n", version)
		os.Exit(0)
	}

	if config.FilePath == "" || config.URL == "" || config.BodyTemplate == "" {
		fmt.Fprintf(os.Stderr, "usage: bulkcaller -file <data> -url <endpoint> -body <template>\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  bulkcaller -file data.csv -url https://api.example.com -body '{\"name\":\"{{name}}\"}'\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	config.Headers = parseKeyValue(headersStr, ":")
	config.QueryParams = parseKeyValue(queryStr, "=")

	log.Printf("🚀 Starting bulk requests to %s", config.URL)
	log.Printf("📁 Reading from: %s", config.FilePath)
	log.Printf("🔧 Workers: %d | Delay: %dms | Retries: %d\n", config.Concurrency, config.DelayMs, config.MaxRetries)

	if err := run(config); err != nil {
		log.Fatalf("❌ Error: %v", err)
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

func extractPlaceholders(template string) []string {
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(template, -1)
	seen := make(map[string]bool)
	var result []string
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}
	return result
}

func run(config Config) error {
	placeholders := extractPlaceholders(config.BodyTemplate)
	log.Printf("🔍 Found placeholders: %v", placeholders)

	records, err := reader.ReadFile(config.FilePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("no data found in file")
	}

	headers := records[0]
	data := records[1:]
	log.Printf("📊 Total rows to process: %d\n", len(data))

	if config.OutputDir != "" {
		if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
			return fmt.Errorf("creating output dir: %w", err)
		}
	}

	targetURL := config.URL
	if len(config.QueryParams) > 0 {
		u, _ := url.Parse(config.URL)
		q := u.Query()
		for k, v := range config.QueryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		targetURL = u.String()
	}

	var wg sync.WaitGroup
	jobs := make(chan map[string]string, config.Concurrency*2)
	results := make(chan RequestResult, config.Concurrency*2)

	go collectResults(results, config, len(data))

	client := &http.Client{Timeout: 30 * time.Second}
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go worker(i, client, config, targetURL, placeholders, headers, jobs, results, &wg)
	}

	startTime := time.Now()
	processed := 0
	total := len(data)

	for idx, row := range data {
		rowMap := make(map[string]string)
		rowMap["__index__"] = fmt.Sprintf("%d", idx)
		for i, h := range headers {
			if i < len(row) {
				rowMap[h] = row[i]
			}
		}
		jobs <- rowMap
		processed++
		if processed%1000 == 0 || processed == total {
			log.Printf("⏳ Queued %d/%d rows...", processed, total)
		}
	}
	close(jobs)

	wg.Wait()
	close(results)

	time.Sleep(200 * time.Millisecond)

	elapsed := time.Since(startTime)
	log.Printf("\n✅ Completed %d requests in %v (%.2f req/s)", len(data), elapsed, float64(len(data))/elapsed.Seconds())

	return nil
}

type RequestResult struct {
	Index      int
	StatusCode int
	Error      error
	Body       []byte
	Duration   time.Duration
}

func worker(id int, client *http.Client, config Config, targetURL string, placeholders []string, headers []string, jobs <-chan map[string]string, results chan<- RequestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for row := range jobs {
		body := config.BodyTemplate
		for _, col := range headers {
			if val, ok := row[col]; ok {
				placeholder := "{{" + col + "}}"
				body = strings.ReplaceAll(body, placeholder, val)
			}
		}
		body = strings.ReplaceAll(body, "{{__index__}}", row["__index__"])

		var jsonBody map[string]interface{}
		if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
			log.Printf("Worker %d: Invalid JSON after substitution (row %s): %v", id, row["__index__"], err)
			results <- RequestResult{Index: parseIndex(row["__index__"]), Error: err}
			continue
		}

		duration, status, respBody, err := doRequest(client, config.Method, targetURL, config.Headers, body, config.MaxRetries)

		if config.DelayMs > 0 {
			time.Sleep(time.Duration(config.DelayMs) * time.Millisecond)
		}

		results <- RequestResult{Index: parseIndex(row["__index__"]), StatusCode: status, Error: err, Body: respBody, Duration: duration}
	}
}

func parseIndex(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func doRequest(client *http.Client, method, url string, headers map[string]string, body string, maxRetries int) (time.Duration, int, []byte, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		start := time.Now()
		req, err := http.NewRequest(method, url, strings.NewReader(body))
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		duration := time.Since(start)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return duration, resp.StatusCode, respBody, nil
		}

		lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
		