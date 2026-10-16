package stats

import (
	"ghostops/internal/logparse"
)

type Stats struct {
	TotalValid int
	ByLevel    map[string]int
	ByService  map[string]int
}

type Collector interface {
	Observe(e logparse.Entry)
	Snapshot() Stats
}

type Counter struct {
	total     int
	byLevel   map[string]int
	byService map[string]int
}

func NewCounter() *Counter {
	return &Counter{
		byLevel:   make(map[string]int),
		byService: make(map[string]int),
	}
}

func (c *Counter) Observe(e logparse.Entry) {
	c.total++
	c.byLevel[e.Level]++
	c.byService[e.Service]++
}

func (c *Counter) Snapshot() Stats {
	byLevel := make(map[string]int, len(c.byLevel))
	for k, v := range c.byLevel {
		byLevel[k] = v
	}
	byService := make(map[string]int, len(c.byService))
	for k, v := range c.byService {
		byService[k] = v
	}
	return Stats{
		TotalValid: c.total,
		ByLevel:    byLevel,
		ByService:  byService,
	}
}
