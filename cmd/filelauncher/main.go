package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"filelauncher/internal/config"
	"filelauncher/internal/history"
	"filelauncher/internal/runner"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	logger := log.New(os.Stdout, "filelauncher: ", log.LstdFlags)

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	hist, err := history.New(cfg.Dolt.Dir, logger)
	if err != nil {
		logger.Fatalf("failed to init history: %v", err)
	}

	if err := processAll(cfg, hist, logger); err != nil {
		logger.Fatalf("processing failed: %v", err)
	}
}

func processAll(cfg *config.Config, hist *history.History, logger *log.Logger) error {
	for _, rule := range cfg.Rules {
		for _, root := range rule.Paths {
			root = filepath.Clean(root)
			err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
				if err != nil {
					logger.Printf("walk error (%s): %v", root, err)
					return nil
				}
				if entry.IsDir() {
					return nil
				}
				path = filepath.Clean(path)
				if !ruleMatches(rule, path) {
					return nil
				}
				info, err := entry.Info()
				if err != nil {
					logger.Printf("stat error (%s): %v", path, err)
					return nil
				}
				modUnix := info.ModTime().Unix()
				lastMod, ok, err := hist.GetLastMod(path, rule.Name)
				if err != nil {
					logger.Printf("history lookup failed (%s): %v", path, err)
					return nil
				}
				if ok && modUnix <= lastMod {
					return nil
				}
				logger.Printf("rule %s matched %s", rule.Name, path)
				eventType := "change"
				res, err := runner.Run(rule, path, eventType)
				hist.Record(history.Event{
					Path:       path,
					EventType:  eventType,
					Rule:       rule.Name,
					Outputs:    rule.Outputs,
					Status:     history.ToEventStatus(err),
					Err:        err,
					DurationMs: history.DurationMs(res.Duration),
				})
				if err != nil {
					logger.Printf("action failed (%s): %v", rule.Name, err)
					return nil
				}
				if err := hist.UpsertState(path, rule.Name, modUnix); err != nil {
					logger.Printf("history state update failed (%s): %v", path, err)
				}
				return nil
			})
			if err != nil {
				logger.Printf("walk failed (%s): %v", root, err)
			}
		}
	}
	return nil
}

func ruleMatches(rule config.Rule, path string) bool {
	for _, root := range rule.Paths {
		root = filepath.Clean(root)
		rel, err := filepath.Rel(root, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		if !includeMatch(rule.Include, rel) {
			continue
		}
		if excludeMatch(rule.Exclude, rel) {
			continue
		}
		return true
	}
	return false
}

func includeMatch(pattern string, rel string) bool {
	if pattern == "" {
		return true
	}
	return globMatch(pattern, rel)
}

func excludeMatch(patterns []string, rel string) bool {
	for _, pattern := range patterns {
		if globMatch(pattern, rel) {
			return true
		}
	}
	return false
}

func globMatch(pattern string, rel string) bool {
	if strings.Contains(pattern, string(os.PathSeparator)) {
		matched, _ := filepath.Match(pattern, rel)
		return matched
	}
	matched, _ := filepath.Match(pattern, filepath.Base(rel))
	return matched
}
