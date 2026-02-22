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
	"strings"
	"sync"
	"time"
)

var version = "dev"

func main() {
	var filePath, urlStr, method, bodyTemplate, headersStr string
	var concurrency, delayMs int
	var printResp bool

	flag.StringVar(&filePath, "file", "", "Path to CSV file")
	flag.StringVar(&urlStr, "url", "", "Target URL")
	flag.StringVar(&method, "method", "POST", "HTTP method")
	flag.StringVar(&bodyTemplate, "body", "", "JSON template with ${col}")
	flag.StringVar(&headersStr, "headers", "", "Headers as key:value")
	flag.IntVar(&concurrency, "concurrency", 10, "Concurrent workers")
	flag.IntVar(&delayMs, "delay", 0, "Delay in ms")
	flag.BoolVar(&printResp, "print", false, "Print responses")
	flag.Parse()

	if filePath == "" || urlStr == "" || bodyTemplate == "" {
		fmt.Println("Usage: bulkcaller -file <csv> -url <url> -body '{\"key\":\"${col}\"}'")
		os.Exit(1)
	}

	headers := parseHeaders(headersStr)

	records, err := readCSV(filePath)
	if err != nil {
		log.Fatal(err)
	}
	if len(records) < 2 {
		log.Fatal("no data rows")
	}

	headerRow := records[0]
	client := &http.Client{Timeout: 30 * time.Second}
	var wg sync.WaitGroup
	jobs := make(chan []string, concurrency*2)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(client, method, urlStr, headers, bodyTemplate, headerRow, jobs, &wg, delayMs, printResp)
	}

	for _, row := range records[1:] {
		jobs <- row
	}
	close(jobs)
	wg.Wait()

	log.Println("Done")
}

func parseHeaders(s string) map[string]string {
	result := make(map[string]string)
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
	return csv.NewReader(f).ReadAll()
}

func worker(client *http.Client, method, url string, headers map[string]string, template string, cols []string, jobs <-chan []string, wg *sync.WaitGroup, delay int, print bool) {
	defer wg.Done()
	for row := range jobs {
		body := template
		for i, col := range cols {
			if i < len(row) {
				body = strings.ReplaceAll(body, "${"+col+"}", row[i])
			}
		}

		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(body), &payload); err != nil {
			log.Printf("Invalid JSON: %v", err)
			continue
		}

		req, _ := http.NewRequest(method, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Request failed: %v", err)
			continue
		}
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if print {
			fmt.Printf("Response: %s\n", string(respBody))
		}

		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}
}
