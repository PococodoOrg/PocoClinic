package pagination

import "testing"

func TestClamp(t *testing.T) {
	if got := Clamp(0, 20); got != 20 {
		t.Fatalf("zero should use default, got %d", got)
	}
	if got := Clamp(-5, 20); got != 20 {
		t.Fatalf("negative should use default, got %d", got)
	}
	if got := Clamp(50, 20); got != 50 {
		t.Fatalf("in-range should pass through, got %d", got)
	}
	if got := Clamp(500, 20); got != MaxPageSize {
		t.Fatalf("oversized should clamp to max, got %d", got)
	}
}
