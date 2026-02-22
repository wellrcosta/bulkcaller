package runner

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wellrcosta/bulkcaller/internal/config"
	"github.com/wellrcosta/bulkcaller/internal/httpclient"
	"github.com/wellrcosta/bulkcaller/internal/reader"
	"github.com/wellrcosta/bulkcaller/internal/result"
	"github.com/wellrcosta/bulkcaller/internal/template"
)

// Runner orchestrates the bulk HTTP requests
type Runner struct {
	config   *config.Config
	client   *httpclient.Client
	result   *result.Collector
}

// New creates a new runner
func New(cfg *config.Config) *Runner {
	return &Runner{
		config: cfg,
		client: httpclient.NewClient(cfg.Timeout, cfg.MaxRetries),
		result: result.NewCollector(cfg.OutputDir, cfg.PrintResponse),
	}
}

// Run executes the bulk requests
func (r *Runner) Run() error {
	// Initialize output directory
	if err := r.result.Init(); err != nil {
		return fmt.Errorf("initializing output: %w", err)
	}

	// Read data file
	records, err := reader.ReadFile(r.config.FilePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("no data rows found (need at least header + 1 data row)")
	}

	headers := records[0]
	dataRows := records[1:]

	log.Printf("🚀 Starting bulk requests to %s", r.config.URL)
	log.Printf("📁 Reading from: %s", r.config.FilePath)
	log.Printf("📊 Total rows: %d | Workers: %d | Delay: %v | Retries: %d\n",
		len(dataRows), r.config.Concurrency, r.config.Delay, r.config.MaxRetries)

	// Extract placeholders for logging
	placeholders := template.ExtractPlaceholders(r.config.BodyTemplate)
	if len(placeholders) > 0 {
		log.Printf("🔍 Found placeholders: %v", placeholders)
	}

	// Setup workers
	var wg sync.WaitGroup
	jobs := make(chan []string, r.config.Concurrency*2)
	results := make(chan rowResult, r.config.Concurrency*2)

	// Start result collector goroutine
	go func() {
		for res := range results {
			r.result.Collect(res.index, res.statusCode, res.body, res.err)
		}
	}()

	// Start workers
	for i := 0; i < r.config.Concurrency; i++ {
		wg.Add(1)
		go r.worker(i, jobs, results, headers, &wg)
	}

	// Queue jobs
	startTime := time.Now()
	queued := 0
	for idx, row := range dataRows {
		jobs <- row
		queued++
		if queued%100 == 0 || queued == len(dataRows) {
			log.Printf("⏳ Queued %d/%d rows...", queued, len(dataRows))
		}
		_ = idx // used in logging above
	}

	// Wait for completion
	close(jobs)
	wg.Wait()
	close(results)

	// Wait for result collector to finish
	time.Sleep(100 * time.Millisecond)

	// Print summary
	elapsed := time.Since(startTime)
	r.result.PrintSummary()
	log.Printf("\n✅ Completed %d requests in %v (%.2f req/s)",
		len(dataRows), elapsed, float64(len(dataRows))/elapsed.Seconds())

	return nil
}

type rowResult struct {
	index      int
	statusCode int
	body       []byte
	err        error
}

func (r *Runner) worker(id int, jobs <-chan []string, results chan<- rowResult, headers []string, wg *sync.WaitGroup) {
	defer wg.Done()

	index := 0
	for row := range jobs {
		// Substitute placeholders
		body := template.SubstituteWithHeaders(r.config.BodyTemplate, headers, row)

		// Validate JSON
		if err := template.ValidateJSON(body); err != nil {
			results <- rowResult{index: index, err: fmt.Errorf("invalid JSON: %w", err)}
			index++
			continue
		}

		// Make request
		result := r.client.DoRequest(r.config.Method, r.config.URL, r.config.Headers, body)

		if r.config.Delay > 0 {
			time.Sleep(r.config.Delay)
		}

		results <- rowResult{
			index:      index,
			statusCode: result.StatusCode,
			body:       result.Body,
			err:        result.Error,
		}
		index++
	}
}
