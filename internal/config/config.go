package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version int        `yaml:"version"`
	Dolt    DoltConfig `yaml:"dolt"`
	Rules   []Rule     `yaml:"rules"`
}

type DoltConfig struct {
	Dir string `yaml:"dir"`
}

type Rule struct {
	Name    string            `yaml:"name"`
	Paths   []string          `yaml:"paths"`
	Include string            `yaml:"include"`
	Exclude []string          `yaml:"exclude"`
	Outputs []string          `yaml:"outputs"`
	Action  ActionConfig      `yaml:"action"`
	Env     map[string]string `yaml:"env"`
}

type ActionConfig struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	applyDefaults(&cfg)
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Dolt.Dir == "" {
		cfg.Dolt.Dir = "data"
	}
}

func validate(cfg *Config) error {
	if len(cfg.Rules) == 0 {
		return errors.New("config: at least one rule is required")
	}
	for i, rule := range cfg.Rules {
		if rule.Name == "" {
			return errors.New("config: rule name is required")
		}
		if len(rule.Paths) == 0 {
			return errors.New("config: rule paths are required")
		}
		if rule.Action.Command == "" {
			return errors.New("config: rule action command is required")
		}
		_ = i
	}
	return nil
}
