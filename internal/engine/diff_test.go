package engine

import (
	"testing"

	"github.com/diskcern/diskcern/internal/models"
)

func TestCompareSnapshots(t *testing.T) {
	oldRecords := []models.FileRecord{
		{Path: "/a.txt", Size: 100},
		{Path: "/b.txt", Size: 200},
		{Path: "/c.txt", Size: 300},
	}

	newRecords := []models.FileRecord{
		{Path: "/a.txt", Size: 100}, // Unchanged
		{Path: "/b.txt", Size: 250}, // Modified
		{Path: "/d.txt", Size: 400}, // Added
		// c.txt is removed
	}

	diff := CompareSnapshots(oldRecords, newRecords)

	if len(diff.Added) != 1 || diff.Added[0].Path != "/d.txt" {
		t.Errorf("Expected 1 added file (/d.txt), got %v", diff.Added)
	}

	if len(diff.Modified) != 1 || diff.Modified[0].Path != "/b.txt" {
		t.Errorf("Expected 1 modified file (/b.txt), got %v", diff.Modified)
	}

	if len(diff.Removed) != 1 || diff.Removed[0].Path != "/c.txt" {
		t.Errorf("Expected 1 removed file (/c.txt), got %v", diff.Removed)
	}
}
