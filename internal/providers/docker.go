package providers

type DockerProvider struct{}

func (p *DockerProvider) ID() string {
	return "docker"
}

func (p *DockerProvider) Name() string {
	return "Docker System"
}

func (p *DockerProvider) Detect(path string, isDir bool) (bool, ScanDirective) {
	if path == "" {
		return true, ContinueTraversal
	}
	return false, ContinueTraversal
}

func (p *DockerProvider) Analyze(path string) (AnalysisResult, error) {
	return AnalysisResult{
		Risk:    "Warning",
		Context: "Global Docker daemon storage (images, containers, volumes).",
	}, nil
}

func (p *DockerProvider) GetCleanupActions(path string) []Action {
	return []Action{
		{
			Name:        "Prune System",
			Description: "Remove all unused containers, networks, and images.",
			Command:     "docker system prune -a -f",
		},
	}
}
