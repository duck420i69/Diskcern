package providers

type AnalysisResult struct {
	Risk        string
	Recoverable int64
	Context     string
}

type Action struct {
	Name        string
	Description string
	Command     string
}

type ScanDirective string

const (
	ContinueTraversal ScanDirective = "CONTINUE"
	StopTraversal     ScanDirective = "STOP"
	LabelChildren     ScanDirective = "CHILDREN"
)

type Provider interface {
	ID() string
	Name() string
	Detect(path string, isDir bool) (bool, ScanDirective)
	Analyze(path string) (AnalysisResult, error)
	GetCleanupActions(path string) []Action
}
