# Adding a New Provider

Diskcern uses a pluggable "Provider" system to detect and analyze specific types of files or directories (e.g., game libraries, developer artifacts). This guide explains how to create and register a new Provider.

## The Provider Interface

All providers must implement the `Provider` interface defined in `internal/providers/provider.go`:

```go
type Provider interface {
	ID() string
	Name() string
	Detect(path string, isDir bool) (bool, ScanDirective)
	Analyze(path string) (AnalysisResult, error)
	GetCleanupActions(path string) []Action
}
```

### Implementing the Methods

1. **`ID() string`**: Returns a unique identifier for the provider (e.g., `"steam"`, `"node_modules"`).
2. **`Name() string`**: Returns a human-readable name for the provider.
3. **`Detect(path string, isDir bool) (bool, ScanDirective)`**: Evaluates a file or directory path during a scan.
   * Return `true` if the provider recognizes this path.
   * Return a `ScanDirective` to control traversal:
     * `ContinueTraversal`: Keep scanning children.
     * `StopTraversal`: Stop scanning this directory (useful if the provider handles the whole directory as a single unit).
     * `LabelChildren`: Label all children as belonging to this provider.
4. **`Analyze(path string) (AnalysisResult, error)`**: If `Detect` returns true, this method is called to extract insights, such as the `Recoverable` size or a `Risk` level.
5. **`GetCleanupActions(path string) []Action`**: Returns a list of potential actions the user can take (e.g., "Delete Cache", "Uninstall Game").

## Registering the Provider

Once you have implemented your provider, you must add it to the global Registry.

1. Open `internal/providers/registry.go` (or wherever the application initializes the registry).
2. Instantiate your provider and call `Registry.Register()`.

```go
registry := providers.NewRegistry()
registry.Register(&MyNewProvider{})
```

## Global Providers

If your provider does not rely on a specific path detection (for example, it queries a system API like Docker), `Detect("", false)` should return `true`. The registry will recognize it as a global provider via the `GlobalProviders()` method.
