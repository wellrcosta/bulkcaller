package result

import "log"

type Collector struct {
	Success int
	Failed  int
}

func New() *Collector {
	return &Collector{}
}

func (c *Collector) Print() {
	log.Printf("Results: %d success, %d failed", c.Success, c.Failed)
}
