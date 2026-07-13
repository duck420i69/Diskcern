package engine

import "github.com/diskcern/diskcern/internal/models"

func CompareSnapshots(oldRecords, newRecords []models.FileRecord) models.DiffResult {
	var result models.DiffResult
	
	oldMap := make(map[string]models.FileRecord)
	for _, r := range oldRecords {
		oldMap[r.Path] = r
	}

	for _, r := range newRecords {
		if old, exists := oldMap[r.Path]; exists {
			if old.Size != r.Size {
				result.Modified = append(result.Modified, r)
			}
			delete(oldMap, r.Path)
		} else {
			result.Added = append(result.Added, r)
		}
	}

	for _, r := range oldMap {
		result.Removed = append(result.Removed, r)
	}

	return result
}
