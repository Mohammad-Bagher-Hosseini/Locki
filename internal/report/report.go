package report

import (
	"io"
	"time"

	"ghostops/internal/detect"
	"ghostops/internal/stats"
)

type MalformedSample struct {
	Line    int
	Content string
	Reason  string
}

type Report struct {
	Source          string
	GeneratedAt     time.Time
	Stats           stats.Stats
	Malformed       int
	MalformedSample []MalformedSample
	Incidents       []detect.Incident
}

type Reporter interface {
	Write(r Report, w io.Writer) error
}
