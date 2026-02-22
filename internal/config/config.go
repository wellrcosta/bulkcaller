package config

import (
	"errors"
	"fmt"
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
	Delay         int
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

// GetHTTPURL returns URL with query params
func (c *Config) GetHTTPURL() string {
	if len(c.QueryParams) == 0 {
		return c.URL
	}

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
