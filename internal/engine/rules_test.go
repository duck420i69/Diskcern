package engine

import (
	"os"
	"testing"
)

func TestRulesEngineMatch(t *testing.T) {
	jsonContent := `{
		"rules": [
			{"name": "Node Modules", "pattern": "node_modules", "type": "directory"},
			{"name": "Git", "pattern": ".git", "type": "directory"}
		]
	}`

	tmpfile, err := os.CreateTemp("", "rules.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(jsonContent)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	engine, err := LoadRules(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if match := engine.Match("path/to/node_modules", true); match != "Node Modules" {
		t.Errorf("Expected 'Node Modules', got '%s'", match)
	}

	if match := engine.Match("path/to/node_modules", false); match != "" {
		t.Errorf("Expected empty (type mismatch), got '%s'", match)
	}

	if match := engine.Match("path/to/other", true); match != "" {
		t.Errorf("Expected empty, got '%s'", match)
	}
}
