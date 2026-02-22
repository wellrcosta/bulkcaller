package result

import (
	"fmt"
	"log"
	"os"
)

// Collector handles result aggregation and output
type Collector struct {
	outputDir    string
	printToStdout bool
	successCount int
	failCount    int
}

// NewCollector creates a new result collector
func NewCollector(outputDir string, printToStdout bool) *Collector {
	return &Collector{
		outputDir:     outputDir,
		printToStdout: printToStdout,
	}
}

// Init prepares the output directory
func (c *Collector) Init() error {
	if c.outputDir == "" {
		return nil
	}
	return os.MkdirAll(c.outputDir, 0755)
}

// Collect processes a result and outputs it
func (c *Collector) Collect(index int, statusCode int, body []byte, err error) {
	if err != nil {
		c.failCount++
		log.Printf("❌ Row %d failed: %v", index, err)
		return
	}
	
	c.successCount++
	if c.printToStdout {
		log.Printf("✅ Row %d: HTTP %d", index, statusCode)
	}
	
	if c.outputDir != "" {
		filename := fmt.Sprintf("%s/response_%d.json", c.outputDir, index)
		if err := os.WriteFile(filename, body, 0644); err != nil {
			log.Printf("⚠️ Failed to write response %d: %v", index, err)
		}
	}
}

// PrintSummary prints the final summary
func (c *Collector) PrintSummary() {
	log.Printf("📈 Results: %d success, %d failed", c.successCount, c.failCount)
}

// GetCounts returns the success and failure counts
func (c *Collector) GetCounts() (int, int) {
	return c.successCount, c.failCount
}
