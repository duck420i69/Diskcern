package providers

import (
	"path/filepath"
	"strings"
)

type SDKProvider struct{}

func (p *SDKProvider) ID() string { return "sdk" }
func (p *SDKProvider) Name() string { return "Developer SDKs" }

func (p *SDKProvider) Detect(path string, isDir bool) (bool, ScanDirective) {
	if !isDir {
		return false, ContinueTraversal
	}
	normalized := strings.ToLower(filepath.ToSlash(path))
	base := filepath.Base(normalized)
	
	if strings.Contains(normalized, "appdata/local/android/sdk") && base == "sdk" {
		return true, StopTraversal
	}
	if strings.Contains(normalized, "dotnet/sdk") && base == "sdk" {
		return true, StopTraversal
	}
	return false, ContinueTraversal
}

func (p *SDKProvider) Analyze(path string) (AnalysisResult, error) {
	return AnalysisResult{
		Risk:    "Warning",
		Context: "Developer SDK cache. Can be extremely large. Ensure you are not actively developing for this platform before deleting.",
	}, nil
}

func (p *SDKProvider) GetCleanupActions(path string) []Action {
	return []Action{}
}
