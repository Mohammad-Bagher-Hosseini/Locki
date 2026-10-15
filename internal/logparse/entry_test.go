package logparse

import (
	"testing"
	"time"
)

func TestParseLine(t *testing.T) {
	parser := LineParser{}

	t.Run("valid", func(t *testing.T) {
		e, err := parser.ParseLine("2026-10-01T00:26:06Z INFO payment-service User successfully logged in", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if e.Level != "INFO" || e.Service != "payment-service" {
			t.Fatalf("wrong fields: %+v", e)
		}
		want := time.Date(2026, 10, 1, 0, 26, 6, 0, time.UTC)
		if !e.Time.Equal(want) {
			t.Fatalf("wrong time: %v", e.Time)
		}
	})

	cases := []struct {
		name string
		line string
	}{
		{"empty", ""},
		{"too few fields", "2026-10-01T00:26:06Z INFO"},
		{"bad timestamp", "2026-10-16 10:01:12 ERROR [malformed timestamp format]"},
		{"unknown level", "2026-10-16T10:00:00Z INVALID_LEVEL auth-service something broke"},
		{"missing message", "2026-10-01T00:26:06Z INFO payment-service"},
		{"garbage", "THIS IS A COMPLETELY INVALID LINE WITHOUT STRUCTURE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parser.ParseLine(tc.line, 7); err == nil {
				t.Fatalf("expected error for %q", tc.line)
			}
		})
	}
}
