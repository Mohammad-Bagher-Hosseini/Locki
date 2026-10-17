package detect

import (
	"strings"

	"ghostops/internal/logparse"
)

type Rule struct {
	ID       string
	Title    string
	Severity string
	Levels   []string
	Keywords []string
}

func (r Rule) Matches(e logparse.Entry) bool {
	if len(r.Levels) > 0 {
		found := false
		for _, l := range r.Levels {
			if e.Level == l {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(r.Keywords) == 0 {
		return true
	}
	lower := strings.ToLower(e.Message)
	for _, k := range r.Keywords {
		if strings.Contains(lower, k) {
			return true
		}
	}
	return false
}
