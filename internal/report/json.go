package report

import (
	"encoding/json"
	"io"
)

type jsonIncident struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Severity  string   `json:"severity"`
	Service   string   `json:"service"`
	Count     int      `json:"count"`
	FirstSeen string   `json:"first_seen"`
	LastSeen  string   `json:"last_seen"`
	Samples   []string `json:"samples"`
}

type jsonMalformed struct {
	Line    int    `json:"line"`
	Content string `json:"content"`
	Reason  string `json:"reason"`
}

type jsonReport struct {
	Meta struct {
		Source       string `json:"source"`
		GeneratedAt  string `json:"generated_at"`
		ValidEntries int    `json:"valid_entries"`
		Malformed    int    `json:"malformed_lines"`
	} `json:"meta"`
	Stats struct {
		TotalValid int            `json:"total_valid"`
		Malformed  int            `json:"malformed"`
		ByLevel    map[string]int `json:"by_level"`
		ByService  map[string]int `json:"by_service"`
	} `json:"stats"`
	Incidents        []jsonIncident  `json:"incidents"`
	MalformedSamples []jsonMalformed `json:"malformed_samples"`
}

type JSONReporter struct{}

func (JSONReporter) Write(r Report, w io.Writer) error {
	var out jsonReport
	out.Meta.Source = r.Source
	out.Meta.GeneratedAt = r.GeneratedAt.UTC().Format("2006-01-02T15:04:05Z")
	out.Meta.ValidEntries = r.Stats.TotalValid
	out.Meta.Malformed = r.Malformed
	out.Stats.TotalValid = r.Stats.TotalValid
	out.Stats.Malformed = r.Malformed
	if r.Stats.ByLevel == nil {
		out.Stats.ByLevel = map[string]int{}
	} else {
		out.Stats.ByLevel = r.Stats.ByLevel
	}
	if r.Stats.ByService == nil {
		out.Stats.ByService = map[string]int{}
	} else {
		out.Stats.ByService = r.Stats.ByService
	}
	out.Incidents = make([]jsonIncident, 0, len(r.Incidents))
	for _, in := range r.Incidents {
		samples := in.Samples
		if samples == nil {
			samples = []string{}
		}
		out.Incidents = append(out.Incidents, jsonIncident{
			ID:        in.RuleID,
			Title:     in.Title,
			Severity:  in.Severity,
			Service:   in.Service,
			Count:     in.Count,
			FirstSeen: in.FirstSeen.UTC().Format("2006-01-02T15:04:05Z"),
			LastSeen:  in.LastSeen.UTC().Format("2006-01-02T15:04:05Z"),
			Samples:   samples,
		})
	}
	out.MalformedSamples = make([]jsonMalformed, 0, len(r.MalformedSample))
	for _, m := range r.MalformedSample {
		out.MalformedSamples = append(out.MalformedSamples, jsonMalformed{
			Line:    m.Line,
			Content: m.Content,
			Reason:  m.Reason,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
