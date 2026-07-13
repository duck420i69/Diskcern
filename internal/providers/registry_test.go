package providers

import (
	"testing"
)

func TestRegistryDetect(t *testing.T) {
	r := NewRegistry()
	r.Register(&NodeProvider{})
	r.Register(&SteamProvider{})
	r.Register(&DriverProvider{})
	
	p, directive := r.Detect("C:/projects/myapp/node_modules", true)
	if p == nil || p.ID() != "node_modules" {
		t.Errorf("Expected node_modules provider, got %v", p)
	}
	if directive != StopTraversal {
		t.Errorf("Expected StopTraversal, got %v", directive)
	}
	
	p, directive = r.Detect("C:/SteamLibrary/steamapps/common", true)
	if p == nil || p.ID() != "steam_game" {
		t.Errorf("Expected steam_game provider, got %v", p)
	}
	if directive != LabelChildren {
		t.Errorf("Expected LabelChildren, got %v", directive)
	}
	
	p, directive = r.Detect("C:/AMD", true)
	if p == nil || p.ID() != "driver_cache" {
		t.Errorf("Expected driver_cache provider, got %v", p)
	}
	if directive != StopTraversal {
		t.Errorf("Expected StopTraversal, got %v", directive)
	}

	p, directive = r.Detect("C:/projects/myapp/src", true)
	if p != nil {
		t.Errorf("Expected no provider, got %v", p)
	}
	if directive != ContinueTraversal {
		t.Errorf("Expected ContinueTraversal, got %v", directive)
	}
}
