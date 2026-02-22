package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	FilePath      string
	URL           string
	Method        string
	BodyTemplate  string
	Headers       map[string]string
	QueryParams   map[string]string
	Concurrency   int
	Delay         time.Duration
	Timeout       time.Duration
	OutputDir     string
	PrintResponse bool
	MaxRetries    int
}

// New creates a Config with defaults
func New() *Config {
	return &Config{
		Method:      "POST",
		Concurrency: 10,
		Timeout:     30 * time.Second,
		Delay:       0,
		MaxRetries:  3,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
}

// ParseFlags parses command line flags
func (c *Config) ParseFlags(args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-file":
			if i+1 < len(args) {
				c.FilePath = args[i+1]
				i++
			}
		case "-url":
			if i+1 < len(args) {
				c.URL = args[i+1]
				i++
			}
		case "-method":
			if i+1 < len(args) {
				c.Method = args[i+1]
				i++
			}
		case "-body":
			if i+1 < len(args) {
				c.BodyTemplate = args[i+1]
				i++
			}
		case "-headers":
			if i+1 < len(args) {
				c.Headers = parseKeyValue(args[i+1], ":")
				i++
			}
		case "-query":
			if i+1 < len(args) {
				c.QueryParams = parseKeyValue(args[i+1], "=")
				i++
			}
		case "-concurrency":
			if i+1 < len(args) {
				if n, err := strconv.Atoi(args[i+1]); err == nil {
					c.Concurrency = n
				}
				i++
			}
		case "-delay":
			if i+1 < len(args) {
				if n, err := strconv.Atoi(args[i+1]); err == nil {
					c.Delay = time.Duration(n) * time.Millisecond
				}
				i++
			}
		case "-timeout":
			if i+1 < len(args) {
				c.Timeout = parseDuration(args[i+1])
				i++
			}
		case "-output":
			if i+1 < len(args) {
				c.OutputDir = args[i+1]
				i++
			}
		case "-max-retries":
			if i+1 < len(args) {
				if n, err := strconv.Atoi(args[i+1]); err == nil {
					c.MaxRetries = n
				}
				i++
			}
		case "-print":
			c.PrintResponse = true
		}
	}
	return nil
}

// Validate checks if required fields are set
func (c *Config) Validate() error {
	if c.FilePath == "" {
		return errors.New("file path is required")
	}
	if c.URL == "" {
		return errors.New("URL is required")
	}
	if c.BodyTemplate == "" {
		return errors.New("body template is required")
	}
	return nil
}

// parseKeyValue parses key:value or key=value pairs
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

// parseDuration parses duration string
func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	if d == 0 {
		d = 30 * time.Second
	}
	return d
}

// GetVersion returns the app version
func GetVersion() string {
	return "1.0.0"
}

// GetHTTPURL returns URL with query params
func (c *Config) GetHTTPURL() string {
	if len(c.QueryParams) == 0 {
		return c.URL
	}
	
	// Add query params to URL
	sep := "?"
	if strings.Contains(c.URL, "?") {
		sep = "&"
	}
	
	var params []string
	for k, v := range c.QueryParams {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	
	return c.URL + sep + strings.Join(params, "&")
}
