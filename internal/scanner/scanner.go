package scanner

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/diskcern/diskcern/internal/engine"
	"github.com/diskcern/diskcern/internal/models"
	"github.com/diskcern/diskcern/internal/providers"
)

type Scanner struct {
	rulesEngine *engine.RulesEngine
	registry    *providers.Registry
}

func NewScanner(re *engine.RulesEngine, reg *providers.Registry) *Scanner {
	return &Scanner{rulesEngine: re, registry: reg}
}

func (s *Scanner) Scan(root string, progressCb func(path string)) ([]models.FileRecord, error) {
	var records []models.FileRecord
	var mu sync.Mutex
	var wg sync.WaitGroup

	recordCh := make(chan models.FileRecord, 10000)

	go func() {
		for r := range recordCh {
			mu.Lock()
			records = append(records, r)
			mu.Unlock()
		}
	}()

	sem := make(chan struct{}, 50)
	var walk func(path string, info os.FileInfo, forceProvider string, forceStop bool)
	
	walk = func(path string, info os.FileInfo, forceProvider string, forceStop bool) {
		defer wg.Done()
		sem <- struct{}{}
		defer func() { <-sem }()

		if progressCb != nil && info.IsDir() {
			progressCb(path)
		}

		if !info.IsDir() {
			recordCh <- models.FileRecord{
				Path:  path,
				Size:  info.Size(),
				IsDir: false,
			}
			return
		}

		if forceStop {
			size := calculateDirSizeSync(path)
			recordCh <- models.FileRecord{
				Path:       path,
				Size:       size,
				IsDir:      true,
				ProviderID: forceProvider,
			}
			return
		}

		matchedRule := ""
		if s.rulesEngine != nil {
			matchedRule = s.rulesEngine.Match(path, true)
		}

		providerID := ""
		directive := providers.ContinueTraversal

		if s.registry != nil && path != root {
			if p, d := s.registry.Detect(path, true); p != nil {
				providerID = p.ID()
				directive = d
			}
		}

		if directive == providers.StopTraversal || (matchedRule != "" && path != root) {
			size := calculateDirSizeSync(path)
			recordCh <- models.FileRecord{
				Path:        path,
				Size:        size,
				IsDir:       true,
				MatchedRule: matchedRule,
				ProviderID:  providerID,
			}
			return
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return
		}

		var dirSize int64
		for _, e := range entries {
			childPath := filepath.Join(path, e.Name())
			info, err := e.Info()
			if err != nil {
				continue
			}
			if e.IsDir() {
				wg.Add(1)
				childForceProvider := ""
				childForceStop := false
				if directive == providers.LabelChildren {
					childForceProvider = providerID
					childForceStop = true
				}
				go walk(childPath, info, childForceProvider, childForceStop)
			} else {
				recordCh <- models.FileRecord{
					Path:  childPath,
					Size:  info.Size(),
					IsDir: false,
				}
				dirSize += info.Size()
			}
		}
		
		recordCh <- models.FileRecord{
			Path:  path,
			Size:  dirSize,
			IsDir: true,
		}
	}

	info, err := os.Stat(root)
	if err == nil {
		wg.Add(1)
		walk(root, info, "", false)
	}
	
	wg.Wait()
	close(recordCh)

	return records, nil
}

func calculateDirSizeSync(path string) int64 {
	var total int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			atomic.AddInt64(&total, info.Size())
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}
