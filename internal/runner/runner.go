package runner

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"filelauncher/internal/config"
)

type Result struct {
	Duration time.Duration
}

func Run(rule config.Rule, path string, event string) (Result, error) {
	start := time.Now()
	cmd := exec.Command(rule.Action.Command, expandArgs(rule.Action.Args, rule, path, event)...)
	cmd.Env = append(os.Environ(), buildEnv(rule, path, event)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return Result{Duration: time.Since(start)}, err
}

func buildEnv(rule config.Rule, path string, event string) []string {
	outputs := strings.Join(rule.Outputs, ",")
	pairs := []string{
		"FILE_PATH=" + path,
		"RULE_NAME=" + rule.Name,
		"EVENT_TYPE=" + event,
		"OUTPUT_SUFFIXES=" + outputs,
	}
	for k, v := range rule.Env {
		pairs = append(pairs, k+"="+expand(v, rule, path, event))
	}
	return pairs
}

func expandArgs(args []string, rule config.Rule, path string, event string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		out = append(out, expand(arg, rule, path, event))
	}
	return out
}

func expand(value string, rule config.Rule, path string, event string) string {
	replacements := map[string]string{
		"{path}":    path,
		"{rule}":    rule.Name,
		"{event}":   event,
		"{outputs}": strings.Join(rule.Outputs, ","),
	}
	for key, val := range replacements {
		value = strings.ReplaceAll(value, key, val)
	}
	return value
}
