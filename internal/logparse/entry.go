package logparse

import (
	"errors"
	"strings"
	"time"
)

type Entry struct {
	Time    time.Time
	Level   string
	Service string
	Message string
	LineNo  int
}

type Parser interface {
	ParseLine(line string, lineNo int) (Entry, error)
}

type LineParser struct{}

var validLevels = map[string]bool{
	"INFO":  true,
	"WARN":  true,
	"ERROR": true,
	"FATAL": true,
}

func (LineParser) ParseLine(line string, lineNo int) (Entry, error) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return Entry{}, errors.New("empty line")
	}
	parts := strings.SplitN(trimmed, " ", 4)
	if len(parts) < 4 {
		return Entry{}, errors.New("invalid format: expected '<timestamp> <level> <service> <message>'")
	}
	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return Entry{}, errors.New("bad timestamp: " + parts[0])
	}
	if !validLevels[parts[1]] {
		return Entry{}, errors.New("unknown level: " + parts[1])
	}
	if parts[2] == "" {
		return Entry{}, errors.New("missing service")
	}
	msg := strings.TrimSpace(parts[3])
	if msg == "" {
		return Entry{}, errors.New("missing message")
	}
	return Entry{
		Time:    ts,
		Level:   parts[1],
		Service: parts[2],
		Message: msg,
		LineNo:  lineNo,
	}, nil
}
