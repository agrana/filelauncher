package history

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type History struct {
	dir     string
	enabled bool
	log     *log.Logger
}

type Event struct {
	Path       string
	EventType  string
	Rule       string
	Outputs    []string
	Status     string
	Err        error
	DurationMs int64
}

func New(dir string, logger *log.Logger) (*History, error) {
	_, err := exec.LookPath("dolt")
	if err != nil {
		if logger != nil {
			logger.Printf("dolt not found in PATH; history disabled")
		}
		return &History{enabled: false, log: logger}, nil
	}
	if dir == "" {
		return nil, errors.New("history: dolt dir required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	if _, err := os.Stat(filepath.Join(dir, ".dolt")); os.IsNotExist(err) {
		if err := run(dir, "init"); err != nil {
			return nil, err
		}
	}
	if err := run(dir, "sql", "-q", createTableSQL()); err != nil {
		return nil, err
	}
	if err := run(dir, "sql", "-q", createStateTableSQL()); err != nil {
		return nil, err
	}
	return &History{dir: dir, enabled: true, log: logger}, nil
}

func (h *History) Record(event Event) {
	if !h.enabled {
		return
	}
	errText := ""
	if event.Err != nil {
		errText = event.Err.Error()
	}
	query := fmt.Sprintf(
		"insert into file_events (ts, path, event, rule, outputs, status, error, duration_ms) values (now(), '%s', '%s', '%s', '%s', '%s', '%s', %d)",
		escape(event.Path),
		escape(event.EventType),
		escape(event.Rule),
		escape(strings.Join(event.Outputs, ",")),
		escape(event.Status),
		escape(errText),
		event.DurationMs,
	)
	if err := run(h.dir, "sql", "-q", query); err != nil && h.log != nil {
		h.log.Printf("history insert failed: %v", err)
	}
}

func createTableSQL() string {
	return strings.Join([]string{
		"create table if not exists file_events (",
		"id bigint auto_increment primary key,",
		"ts timestamp,",
		"path text,",
		"event text,",
		"rule text,",
		"outputs text,",
		"status text,",
		"error text,",
		"duration_ms bigint",
		");",
	}, " ")
}

func createStateTableSQL() string {
	return strings.Join([]string{
		"create table if not exists file_state (",
		"path varchar(768),",
		"rule varchar(255),",
		"last_mod_unix bigint,",
		"primary key (path, rule)",
		");",
	}, " ")
}

func (h *History) GetLastMod(path string, rule string) (int64, bool, error) {
	if !h.enabled {
		return 0, false, nil
	}
	query := fmt.Sprintf(
		"select last_mod_unix from file_state where path = '%s' and rule = '%s' limit 1",
		escape(path),
		escape(rule),
	)
	rows, err := queryCSV(h.dir, query)
	if err != nil {
		return 0, false, err
	}
	if len(rows) < 2 || len(rows[1]) == 0 {
		return 0, false, nil
	}
	value := strings.TrimSpace(rows[1][0])
	if value == "" {
		return 0, false, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false, err
	}
	return parsed, true, nil
}

func (h *History) UpsertState(path string, rule string, modUnix int64) error {
	if !h.enabled {
		return nil
	}
	query := fmt.Sprintf(
		"replace into file_state (path, rule, last_mod_unix) values ('%s', '%s', %d)",
		escape(path),
		escape(rule),
		modUnix,
	)
	return run(h.dir, "sql", "-q", query)
}

func run(dir string, args ...string) error {
	cmd := exec.Command("dolt", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func queryCSV(dir string, query string) ([][]string, error) {
	cmd := exec.Command("dolt", "sql", "-q", query, "-r", "csv")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewReader(output))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func escape(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func ToEventStatus(err error) string {
	if err != nil {
		return "error"
	}
	return "ok"
}

func DurationMs(d time.Duration) int64 {
	return d.Milliseconds()
}
