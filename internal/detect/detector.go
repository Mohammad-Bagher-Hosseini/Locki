package detect

import (
	"sort"
	"time"

	"ghostops/internal/logparse"
)

type Incident struct {
	RuleID    string
	Title     string
	Severity  string
	Service   string
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
	Samples   []string
}

type Detector interface {
	Observe(e logparse.Entry)
	Incidents() []Incident
}

type GroupedDetector struct {
	rules  []Rule
	groups map[string]*Incident
}

func NewGroupedDetector(rules []Rule) *GroupedDetector {
	return &GroupedDetector{
		rules:  rules,
		groups: make(map[string]*Incident),
	}
}

func severityRank(s string) int {
	switch s {
	case "critical":
		return 0
	case "high":
		return 1
	case "medium":
		return 2
	default:
		return 3
	}
}

// first matching rule wins
func (d *GroupedDetector) Observe(e logparse.Entry) {
	for _, r := range d.rules {
		if !r.Matches(e) {
			continue
		}
		key := r.ID + "\x00" + e.Service
		g, ok := d.groups[key]
		if !ok {
			g = &Incident{
				RuleID:    r.ID,
				Title:     r.Title,
				Severity:  r.Severity,
				Service:   e.Service,
				FirstSeen: e.Time,
				LastSeen:  e.Time,
			}
			d.groups[key] = g
		}
		g.Count++
		if e.Time.Before(g.FirstSeen) {
			g.FirstSeen = e.Time
		}
		if e.Time.After(g.LastSeen) {
			g.LastSeen = e.Time
		}
		if len(g.Samples) < 3 {
			dup := false
			for _, s := range g.Samples {
				if s == e.Message {
					dup = true
					break
				}
			}
			if !dup {
				g.Samples = append(g.Samples, e.Message)
			}
		}
		return
	}
}

// fixed order so repeated runs print the same report
func (d *GroupedDetector) Incidents() []Incident {
	out := make([]Incident, 0, len(d.groups))
	for _, g := range d.groups {
		cp := *g
		cp.Samples = append([]string(nil), g.Samples...)
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool {
		ri := severityRank(out[i].Severity)
		rj := severityRank(out[j].Severity)
		if ri != rj {
			return ri < rj
		}
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		if out[i].RuleID != out[j].RuleID {
			return out[i].RuleID < out[j].RuleID
		}
		return out[i].Service < out[j].Service
	})
	return out
}
