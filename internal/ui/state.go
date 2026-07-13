package ui

import (
	"encoding/json"
	"os"
)

type CustomState struct {
	Scores       map[string]int    `json:"scores"`
	Descriptions map[string]string `json:"descriptions"`
}

func LoadState() CustomState {
	state := CustomState{
		Scores:       make(map[string]int),
		Descriptions: make(map[string]string),
	}
	data, err := os.ReadFile("config/custom_state.json")
	if err == nil {
		json.Unmarshal(data, &state)
	}
	if state.Scores == nil { state.Scores = make(map[string]int) }
	if state.Descriptions == nil { state.Descriptions = make(map[string]string) }
	return state
}

func SaveState(state CustomState) {
	data, err := json.MarshalIndent(state, "", "  ")
	if err == nil {
		os.MkdirAll("config", 0755)
		os.WriteFile("config/custom_state.json", data, 0644)
	}
}

func AddJSONRule(path string) {
	type Rule struct {
		Name      string `json:"name"`
		Pattern   string `json:"pattern"`
		Type      string `json:"type"`
		Risk      string `json:"risk"`
		ActionCmd string `json:"action_cmd"`
	}
	
	type RuleConfig struct {
		Rules []Rule `json:"rules"`
	}
	
	var cfg RuleConfig
	data, err := os.ReadFile("config/default_rules.json")
	if err == nil {
		json.Unmarshal(data, &cfg)
	}
	
	cfg.Rules = append(cfg.Rules, Rule{
		Name:      "Custom User Rule",
		Pattern:   path,
		Type:      "directory",
		Risk:      "Warning",
		ActionCmd: "Remove-Item -Recurse -Force '" + path + "'",
	})
	
	out, _ := json.MarshalIndent(cfg, "", "  ")
	os.MkdirAll("config", 0755)
	os.WriteFile("config/default_rules.json", out, 0644)
}
