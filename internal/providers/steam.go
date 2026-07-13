package providers

import (
	"path/filepath"
	"strings"
)

type SteamProvider struct{}

func (p *SteamProvider) ID() string {
	return "steam_game"
}

func (p *SteamProvider) Name() string {
	return "Steam Game"
}

func (p *SteamProvider) Detect(path string, isDir bool) (bool, ScanDirective) {
	if !isDir {
		return false, ContinueTraversal
	}
	
	normalized := strings.ToLower(filepath.ToSlash(path))
	if strings.Contains(normalized, "steamapps/common/") {
		parent := filepath.Base(filepath.Dir(normalized))
		if parent == "common" {
			return true, StopTraversal
		}
	}
	return false, ContinueTraversal
}

func (p *SteamProvider) Analyze(path string) (AnalysisResult, error) {
	return AnalysisResult{
		Risk:    "High",
		Context: "Installed Steam Game.",
	}, nil
}

func (p *SteamProvider) GetCleanupActions(path string) []Action {
	return []Action{
		{
			Name:        "Uninstall via Steam",
			Description: "It is recommended to uninstall Steam games directly through the Steam client.",
			Command:     "",
		},
	}
}
