package report

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type TextReporter struct{}

func (TextReporter) Write(r Report, w io.Writer) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("GhostOps Report - source: %s\n", r.Source))
	b.WriteString(strings.Repeat("=", 44) + "\n")
	b.WriteString(fmt.Sprintf("Generated at:    %s\n", r.GeneratedAt.UTC().Format("2006-01-02T15:04:05Z")))
	b.WriteString(fmt.Sprintf("Valid entries:   %d\n", r.Stats.TotalValid))
	b.WriteString(fmt.Sprintf("Malformed lines: %d\n", r.Malformed))
	b.WriteString("\nBy level:\n")
	if r.Stats.TotalValid == 0 {
		b.WriteString("  (no valid entries)\n")
	} else {
		levels := make([]string, 0, len(r.Stats.ByLevel))
		for l := range r.Stats.ByLevel {
			levels = append(levels, l)
		}
		sort.Strings(levels)
		for _, l := range levels {
			b.WriteString(fmt.Sprintf("  %-5s %d\n", l, r.Stats.ByLevel[l]))
		}
	}
	b.WriteString("\nBy service:\n")
	if len(r.Stats.ByService) == 0 {
		b.WriteString("  (no valid entries)\n")
	} else {
		services := make([]string, 0, len(r.Stats.ByService))
		for s := range r.Stats.ByService {
			services = append(services, s)
		}
		sort.Strings(services)
		for _, s := range services {
			b.WriteString(fmt.Sprintf("  %-22s %d\n", s, r.Stats.ByService[s]))
		}
	}
	b.WriteString(fmt.Sprintf("\nIncidents (%d groups):\n", len(r.Incidents)))
	if len(r.Incidents) == 0 {
		b.WriteString("  No incidents detected.\n")
	} else {
		for _, in := range r.Incidents {
			b.WriteString(fmt.Sprintf("  [%s] %s @ %s x%d\n", in.Severity, in.RuleID, in.Service, in.Count))
			b.WriteString(fmt.Sprintf("    %s\n", in.Title))
			b.WriteString(fmt.Sprintf("    %s -> %s\n",
				in.FirstSeen.UTC().Format("2006-01-02T15:04:05Z"),
				in.LastSeen.UTC().Format("2006-01-02T15:04:05Z")))
			for _, s := range in.Samples {
				b.WriteString(fmt.Sprintf("    e.g. %q\n", s))
			}
		}
	}
	if r.Malformed > 0 {
		b.WriteString(fmt.Sprintf("\nMalformed (first %d of %d):\n", len(r.MalformedSample), r.Malformed))
		for _, m := range r.MalformedSample {
			b.WriteString(fmt.Sprintf("  line %d: %s - %q\n", m.Line, m.Reason, m.Content))
		}
	}
	_, err := io.WriteString(w, b.String())
	return err
}
