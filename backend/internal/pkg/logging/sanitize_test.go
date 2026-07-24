package logging

import "testing"

func TestSanitizeStringRemovesLineBreaks(t *testing.T) {
	input := "alice\r\nforged log entry"
	got := SanitizeString(input)
	if got != "aliceforged log entry" {
		t.Fatalf("SanitizeString() = %q, want %q", got, "aliceforged log entry")
	}
}

func TestSanitizeLogArgs(t *testing.T) {
	args := sanitizeLogArgs("path", "/x\n/y", "status", 404)
	if args[1] != "/x/y" {
		t.Fatalf("sanitized path = %q", args[1])
	}
	if args[3] != 404 {
		t.Fatalf("non-string value changed: %v", args[3])
	}
}
