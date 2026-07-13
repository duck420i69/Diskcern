package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Rule struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	Type    string `json:"type"` // "directory" or "file"
}

type RuleConfig struct {
	Rules []Rule `json:"rules"`
}

type RulesEngine struct {
	Rules []Rule
}

func LoadRules(path string) (*RulesEngine, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg RuleConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &RulesEngine{Rules: cfg.Rules}, nil
}

func (e *RulesEngine) Match(path string, isDir bool) string {
	base := filepath.Base(path)
	for _, r := range e.Rules {
		if r.Type == "directory" && !isDir {
			continue
		}
		if r.Type == "file" && isDir {
			continue
		}
		if strings.EqualFold(base, r.Pattern) {
			return r.Name
		}
	}
	return ""
}
