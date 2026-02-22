package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client wraps an HTTP client with configuration
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// ResponseResult holds the result of an HTTP request
type ResponseResult struct {
	Duration   time.Duration
	StatusCode int
	Body       []byte
	Error      error
}

// NewClient creates a new HTTP client with timeout
func NewClient(timeout time.Duration, maxRetries int) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		maxRetries: maxRetries,
	}
}

// DoRequest performs an HTTP request with retries
func (c *Client) DoRequest(method, url string, headers map[string]string, body string) ResponseResult {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
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

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		duration := time.Since(start)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return ResponseResult{
				Duration:   duration,
				StatusCode: resp.StatusCode,
				Body:       respBody,
				Error:      nil,
			}
		}

		lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return ResponseResult{Error: lastErr}
}
