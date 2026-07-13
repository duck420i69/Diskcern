package providers

import (
	"fmt"
	"path/filepath"
	"strings"
)

type NodeProvider struct{}

func (p *NodeProvider) ID() string {
	return "node_modules"
}

func (p *NodeProvider) Name() string {
	return "Node.js Modules"
}

func (p *NodeProvider) Detect(path string, isDir bool) (bool, ScanDirective) {
	if isDir && strings.EqualFold(filepath.Base(path), "node_modules") {
		return true, StopTraversal
	}
	return false, ContinueTraversal
}

func (p *NodeProvider) Analyze(path string) (AnalysisResult, error) {
	return AnalysisResult{
		Risk:    "Safe",
		Context: "NPM dependency cache. Can be re-downloaded via npm install.",
	}, nil
}

func (p *NodeProvider) GetCleanupActions(path string) []Action {
	return []Action{
		{
			Name:        "Delete node_modules",
			Description: "Removes the entire node_modules folder.",
			Command:     fmt.Sprintf("Remove-Item -Recurse -Force \"%s\"", path),
		},
	}
}
