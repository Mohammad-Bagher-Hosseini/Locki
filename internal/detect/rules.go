package detect

func DefaultRules() []Rule {
	return []Rule{
		{
			ID:       "db-pool-exhausted",
			Title:    "Database connection pool exhausted",
			Severity: "critical",
			Levels:   []string{"FATAL"},
			Keywords: []string{"connection pool exhausted"},
		},
		{
			ID:       "db-connect-timeout",
			Title:    "Database connection timeout",
			Severity: "high",
			Levels:   []string{"ERROR"},
			Keywords: []string{"failed to connect to database"},
		},
		{
			ID:       "sql-syntax-error",
			Title:    "SQL syntax error",
			Severity: "high",
			Levels:   []string{"ERROR"},
			Keywords: []string{"sql syntax error", "syntax error"},
		},
		{
			ID:       "null-pointer",
			Title:    "NullPointerException in controller",
			Severity: "high",
			Levels:   []string{"ERROR"},
			Keywords: []string{"nullpointerexception"},
		},
		{
			ID:       "payment-gateway",
			Title:    "Payment gateway failure",
			Severity: "high",
			Levels:   []string{"ERROR"},
			Keywords: []string{"payment gateway"},
		},
		{
			ID:       "db-pool-threshold",
			Title:    "Database connection pool threshold warning",
			Severity: "medium",
			Levels:   []string{"WARN"},
			Keywords: []string{"pool reaching threshold"},
		},
		{
			ID:       "high-memory",
			Title:    "High memory usage",
			Severity: "medium",
			Levels:   []string{"WARN"},
			Keywords: []string{"high memory usage"},
		},
		{
			ID:       "rate-limit",
			Title:    "Rate limit approaching",
			Severity: "low",
			Levels:   []string{"WARN"},
			Keywords: []string{"rate limit approaching"},
		},
		{
			ID:       "deprecated-api",
			Title:    "Deprecated API usage",
			Severity: "low",
			Levels:   []string{"WARN"},
			Keywords: []string{"deprecated api"},
		},
		{
			ID:       "fatal-other",
			Title:    "Unclassified fatal error",
			Severity: "critical",
			Levels:   []string{"FATAL"},
			Keywords: nil,
		},
		{
			ID:       "error-other",
			Title:    "Unclassified error",
			Severity: "medium",
			Levels:   []string{"ERROR"},
			Keywords: nil,
		},
	}
}
