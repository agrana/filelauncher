package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	input := flag.String("input", "", "Path to the input markdown file")
	outputs := flag.String("outputs", "", "Comma-separated output suffixes (e.g. medium,linkedin,x)")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "missing --input")
		os.Exit(2)
	}

	if _, err := exec.LookPath("langchain"); err != nil {
		fmt.Fprintln(os.Stderr, "langchain CLI not found in PATH")
		os.Exit(1)
	}

	cmd := exec.Command("langchain", "run", "--input", *input, "--outputs", *outputs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
