package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	input := flag.String("input", "", "Path to the input markdown file")
	outputs := flag.String("outputs", "", "Comma-separated output suffixes (e.g. medium,linkedin,x)")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "missing --input")
		os.Exit(2)
	}

	data, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	suffixes := splitSuffixes(*outputs)
	if len(suffixes) == 0 {
		fmt.Fprintln(os.Stderr, "missing --outputs")
		os.Exit(2)
	}

	for _, suffix := range suffixes {
		outPath := withSuffix(*input, suffix)
		if err := os.WriteFile(outPath, data, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("wrote", outPath)
	}
}

func splitSuffixes(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		out = append(out, part)
	}
	return out
}

func withSuffix(path string, suffix string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if ext == "" {
		return filepath.Join(dir, name+"."+suffix)
	}
	return filepath.Join(dir, name+"."+suffix+ext)
}
