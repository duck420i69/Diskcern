package providers

import (
	"os"
	"path/filepath"
	"strings"
)

type EpicGamesProvider struct{}

func (p *EpicGamesProvider) ID() string { return "epic_game" }
func (p *EpicGamesProvider) Name() string { return "Epic Games" }

func (p *EpicGamesProvider) Detect(path string, isDir bool) (bool, ScanDirective) {
	if !isDir {
		return false, ContinueTraversal
	}
	
	// Fast check: look for .egstore directory (standard for Epic Games installs)
	egstorePath := filepath.Join(path, ".egstore")
	if info, err := os.Stat(egstorePath); err == nil && info.IsDir() {
		return true, StopTraversal
	}
	
	// Alternate check: inside an "Epic Games" folder and contains an executable
	normalized := strings.ToLower(filepath.ToSlash(path))
	if strings.Contains(normalized, "/epic games/") {
		// Just assuming it's a game folder if it's directly inside Epic Games
		parent := filepath.Base(filepath.Dir(normalized))
		if parent == "epic games" {
			return true, StopTraversal
		}
	}
	return false, ContinueTraversal
}

func (p *EpicGamesProvider) Analyze(path string) (AnalysisResult, error) {
	return AnalysisResult{
		Risk:    "High",
		Context: "Epic Games installation. Uninstall via Epic Launcher to ensure cloud saves and registry entries are cleaned up properly.",
	}, nil
}

func (p *EpicGamesProvider) GetCleanupActions(path string) []Action {
	return []Action{}
}
