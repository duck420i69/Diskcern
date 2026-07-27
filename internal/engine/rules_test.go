package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRules(t *testing.T) {
	// Test loading valid JSON
	jsonContent := `{
		"rules": [
			{"name": "Node Modules", "pattern": "node_modules", "type": "directory"},
			{"name": "Log Files", "pattern": "*.log", "type": "file"}
		]
	}`

	tmpfile, err := os.CreateTemp("", "rules_valid.json")
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

	if len(engine.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(engine.Rules))
	}

	// Test file not found
	_, err = LoadRules(filepath.Join(os.TempDir(), "non_existent_file_xyz.json"))
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test invalid JSON
	invalidJson := `{ "rules": [ {"name": "Test", `
	tmpfileInvalid, err := os.CreateTemp("", "rules_invalid.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfileInvalid.Name())

	if _, err := tmpfileInvalid.Write([]byte(invalidJson)); err != nil {
		t.Fatal(err)
	}
	tmpfileInvalid.Close()

	_, err = LoadRules(tmpfileInvalid.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestRulesEngineMatch(t *testing.T) {
	engine := &RulesEngine{
		Rules: []Rule{
			{"Node Modules", "node_modules", "directory"},
			{"Git", ".git", "directory"},
			{"Log File", "app.log", "file"},
			{"Config", "config.yml", "file"},
			{"Generic Build", "build", ""}, // No specific type
		},
	}

	tests := []struct {
		name     string
		path     string
		isDir    bool
		expected string
	}{
		{
			name:     "Match directory exact",
			path:     "path/to/node_modules",
			isDir:    true,
			expected: "Node Modules",
		},
		{
			name:     "Match directory case-insensitive",
			path:     "path/to/NODE_MODULES",
			isDir:    true,
			expected: "Node Modules",
		},
		{
			name:     "Mismatch directory as file",
			path:     "path/to/node_modules",
			isDir:    false,
			expected: "",
		},
		{
			name:     "Match file exact",
			path:     "var/log/app.log",
			isDir:    false,
			expected: "Log File",
		},
		{
			name:     "Match file case-insensitive",
			path:     "var/log/APP.LOG",
			isDir:    false,
			expected: "Log File",
		},
		{
			name:     "Mismatch file as directory",
			path:     "var/log/app.log",
			isDir:    true,
			expected: "",
		},
		{
			name:     "Match generic rule as directory",
			path:     "proj/build",
			isDir:    true,
			expected: "Generic Build",
		},
		{
			name:     "Match generic rule as file",
			path:     "proj/build",
			isDir:    false,
			expected: "Generic Build",
		},
		{
			name:     "No match",
			path:     "path/to/other",
			isDir:    true,
			expected: "",
		},
		{
			name:     "Match only on base name",
			path:     "node_modules/other",
			isDir:    true,
			expected: "", // Base is "other", not "node_modules"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := engine.Match(tt.path, tt.isDir)
			if match != tt.expected {
				t.Errorf("Match(%q, %v) = %q; want %q", tt.path, tt.isDir, match, tt.expected)
			}
		})
	}
}
