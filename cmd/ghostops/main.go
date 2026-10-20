package main

import (
	"flag"
	"fmt"
	"os"

	"ghostops/internal/processor"
	"ghostops/internal/report"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: ghostops <log-file> [-o output.json]\n")
}

// allow -o both before and after the file argument
func splitArgs(args []string) (positional []string, outPath string, err error) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-o" {
			if i+1 >= len(args) {
				return nil, "", flag.ErrHelp
			}
			outPath = args[i+1]
			i++
			continue
		}
		if len(a) > 3 && a[:3] == "-o=" {
			outPath = a[3:]
			continue
		}
		if a == "-h" || a == "-help" || a == "--help" {
			return nil, "", flag.ErrHelp
		}
		if len(a) > 0 && a[0] == '-' {
			return nil, "", fmt.Errorf("unknown flag %s", a)
		}
		positional = append(positional, a)
	}
	return positional, outPath, nil
}

func run(args []string) int {
	rest, outPath, err := splitArgs(args)
	if err != nil {
		usage()
		return 1
	}
	if len(rest) != 1 {
		usage()
		return 1
	}
	inPath := rest[0]
	f, err := os.Open(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ghostops: cannot open %s: %v\n", inPath, err)
		return 1
	}
	rep, err := processor.Run(inPath, f)
	cerr := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ghostops: analysis failed: %v\n", err)
		return 1
	}
	if cerr != nil {
		fmt.Fprintf(os.Stderr, "ghostops: cannot close %s: %v\n", inPath, cerr)
		return 1
	}
	if err := (report.TextReporter{}).Write(rep, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "ghostops: cannot write report: %v\n", err)
		return 1
	}
	if outPath != "" {
		out, err := os.Create(outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ghostops: cannot create %s: %v\n", outPath, err)
			return 1
		}
		werr := (report.JSONReporter{}).Write(rep, out)
		cerr := out.Close()
		if werr != nil {
			fmt.Fprintf(os.Stderr, "ghostops: cannot write %s: %v\n", outPath, werr)
			return 1
		}
		if cerr != nil {
			fmt.Fprintf(os.Stderr, "ghostops: cannot close %s: %v\n", outPath, cerr)
			return 1
		}
		fmt.Fprintf(os.Stderr, "JSON report written to %s\n", outPath)
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
