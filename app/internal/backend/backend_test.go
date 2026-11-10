package backend

import "testing"

func TestTargetsPreservesURLs(t *testing.T) {
	targetURLs := []string{"http://backend-1:8081", "http://backend-2:8082"}

	targets := Targets(targetURLs)

	if len(targets) != len(targetURLs) {
		t.Fatalf("got %d targets, want %d", len(targets), len(targetURLs))
	}

	for i, target := range targets {
		if target.URL != targetURLs[i] {
			t.Errorf("target %d URL = %q, want %q", i, target.URL, targetURLs[i])
		}
	}
}

func TestTargetsEmptyInput(t *testing.T) {
	if targets := Targets(nil); len(targets) != 0 {
		t.Fatalf("got %d targets, want 0", len(targets))
	}
}
