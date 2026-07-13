package providers

import (
	"fmt"
	"path/filepath"
	"strings"
)

type DriverProvider struct{}

func (p *DriverProvider) ID() string {
	return "driver_cache"
}

func (p *DriverProvider) Name() string {
	return "Hardware Driver Cache"
}

func (p *DriverProvider) Detect(path string, isDir bool) (bool, ScanDirective) {
	if !isDir {
		return false, ContinueTraversal
	}
	
	normalized := strings.ToLower(filepath.ToSlash(path))
	
	if strings.HasSuffix(normalized, "nvidia/displaydriver") || 
	   strings.HasSuffix(normalized, "amd") ||
	   strings.HasSuffix(normalized, "programdata/nvidia corporation/downloader") {
		return true, StopTraversal
	}
	return false, ContinueTraversal
}

func (p *DriverProvider) Analyze(path string) (AnalysisResult, error) {
	return AnalysisResult{
		Risk:    "Safe",
		Context: "Extracted driver installers. Can be safely deleted after installation.",
	}, nil
}

func (p *DriverProvider) GetCleanupActions(path string) []Action {
	return []Action{
		{
			Name:        "Clear Installer Cache",
			Description: "Deletes the extracted installer files.",
			Command:     fmt.Sprintf("Remove-Item -Recurse -Force \"%s\"", path),
		},
	}
}
