package models

import "time"

// Snapshot represents a point-in-time scan of a directory.
type Snapshot struct {
	ID        int64     `json:"id"`
	RootPath  string    `json:"root_path"`
	CreatedAt time.Time `json:"created_at"`
}

// FileRecord represents a file or a grouped directory (from early stopping).
type FileRecord struct {
	ID          int64  `json:"id"`
	SnapshotID  int64  `json:"snapshot_id"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	IsDir       bool   `json:"is_dir"`
	MatchedRule string `json:"matched_rule"`
	ProviderID  string `json:"provider_id"`
}

// DiffResult represents the difference between two snapshots.
type DiffResult struct {
	Added    []FileRecord `json:"added"`
	Removed  []FileRecord `json:"removed"`
	Modified []FileRecord `json:"modified"`
}
