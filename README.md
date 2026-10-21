# GhostOps - Streaming Log Analysis CLI

## Build

Requires Go 1.21+ and only the standard library.

```sh
go build -o ghostops ./cmd/ghostops
```

## Usage

```sh
./ghostops <log-file> [-o output.json]
```

The `-o` flag may appear before or after the log file path.

- Human-readable report is always printed to stdout.
- With `-o`, a JSON report is also written to the given file.
- Exit code `0` on success, `1` on invalid arguments, unreadable input, or output errors.

Examples:

```sh
./ghostops sample.log
./ghostops sample.log -o report.json
./ghostops -o report.json sample.log
```

Example output: [example-output.txt](example-output.txt), [example-report.json](example-report.json)

## Log format

Each valid line:

```
<RFC3339 timestamp> <LEVEL> <service> <message>
```

Example:

```
2026-10-01T00:26:06Z INFO payment-service User successfully logged in
```

Supported levels: `INFO`, `WARN`, `ERROR`, `FATAL`.

Lines that do not match this format are counted as malformed and skipped without aborting the analysis. At most the first 10 malformed lines are kept as samples.

## Incident detection

Entries with level `INFO` never raise incidents. All others are matched against the rule table below, in order. The first matching rule wins, which makes overlapping rules deterministic. Matching entries are grouped by `(rule, service)`, so each incident shows type, severity, frequency, affected service, and time range.

| # | Rule ID | Level | Message keyword | Severity |
|---|---------|-------|-----------------|----------|
| 1 | `db-pool-exhausted` | FATAL | `connection pool exhausted` | critical |
| 2 | `db-connect-timeout` | ERROR | `failed to connect to database` | high |
| 3 | `sql-syntax-error` | ERROR | `sql syntax error` / `syntax error` | high |
| 4 | `null-pointer` | ERROR | `nullpointerexception` | high |
| 5 | `payment-gateway` | ERROR | `payment gateway` | high |
| 6 | `db-pool-threshold` | WARN | `pool reaching threshold` | medium |
| 7 | `high-memory` | WARN | `high memory usage` | medium |
| 8 | `rate-limit` | WARN | `rate limit approaching` | low |
| 9 | `deprecated-api` | WARN | `deprecated api` | low |
| 10 | `fatal-other` | FATAL | any unmatched FATAL | critical |
| 11 | `error-other` | ERROR | any unmatched ERROR | medium |

Matching is case-insensitive substring search. The fallback rules guarantee no `ERROR` or `FATAL` entry is silently dropped.

Incidents are sorted by severity (`critical` → `high` → `medium` → `low`), then count (descending), then rule ID and service, so repeated runs produce identical reports.

## Architecture

```
cmd/ghostops      CLI: argument parsing, file open/close, exit codes
internal/logparse Entry type, Parser interface, ParseLine validation
internal/stats    Collector interface, count aggregation
internal/detect   Rule table, Detector interface, (rule, service) grouping
internal/report   Report model, TextReporter, JSONReporter
internal/processor Pipeline: bufio.Scanner loop wiring the pieces
```

Single-pass streaming: the input file is read line by line with `bufio.Scanner` and never fully loaded into memory. Memory use is proportional to the number of services, rules, and incident groups, not the number of log lines.

Interfaces (`Parser`, `Collector`, `Detector`, `Reporter`) mark component boundaries: the pipeline depends on behavior, not concrete types, so parsing, detection, or output can be tested and replaced independently. `main` only wires the components together.

## Error handling

- Malformed log lines are skipped, counted, and sampled; the exit code stays `0`.
- Invalid arguments, unreadable input files, and output write failures print to stderr and exit with `1`. No success message is printed when the analysis has failed.
