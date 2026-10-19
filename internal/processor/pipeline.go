package processor

import (
	"bufio"
	"io"
	"strings"
	"time"

	"ghostops/internal/detect"
	"ghostops/internal/logparse"
	"ghostops/internal/report"
	"ghostops/internal/stats"
)

const maxMalformedSamples = 10

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func Run(source string, r io.Reader) (report.Report, error) {
	parser := logparse.LineParser{}
	collector := stats.NewCounter()
	detector := detect.NewGroupedDetector(detect.DefaultRules())
	// line by line so big files never sit in memory
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	malformed := 0
	var samples []report.MalformedSample
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		entry, err := parser.ParseLine(line, lineNo)
		if err != nil {
			malformed++
			if len(samples) < maxMalformedSamples {
				samples = append(samples, report.MalformedSample{
					Line:    lineNo,
					Content: truncate(strings.TrimSpace(line), 160),
					Reason:  err.Error(),
				})
			}
			continue
		}
		collector.Observe(entry)
		detector.Observe(entry)
	}
	if err := scanner.Err(); err != nil {
		return report.Report{}, err
	}
	if samples == nil {
		samples = []report.MalformedSample{}
	}
	st := collector.Snapshot()
	return report.Report{
		Source:          source,
		GeneratedAt:     time.Now().UTC(),
		Stats:           st,
		Malformed:       malformed,
		MalformedSample: samples,
		Incidents:       detector.Incidents(),
	}, nil
}
