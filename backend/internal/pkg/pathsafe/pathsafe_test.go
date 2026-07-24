package pathsafe_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/pathsafe"
)

func TestJoinRootBlocksTraversal(t *testing.T) {
	root := t.TempDir()

	_, err := pathsafe.JoinRoot(root, "..", "secret.txt")
	if err == nil {
		t.Fatal("expected traversal to be rejected")
	}

	safe, err := pathsafe.JoinRoot(root, "nested", "file.txt")
	if err != nil {
		t.Fatalf("JoinRoot: %v", err)
	}
	if !strings.HasPrefix(safe, filepath.Clean(root)) {
		t.Fatalf("safe path left root: %q", safe)
	}
}

func TestValidateBackupFilename(t *testing.T) {
	valid := []string{
		"pococlinic-backup-20260101-120000.tar.gz",
	}
	for _, name := range valid {
		if err := pathsafe.ValidateBackupFilename(name); err != nil {
			t.Fatalf("valid name rejected: %q (%v)", name, err)
		}
	}

	invalid := []string{
		"",
		"../pococlinic-backup-20260101-120000.tar.gz",
		"backup.tar.gz",
		`etc\passwd`,
	}
	for _, name := range invalid {
		if err := pathsafe.ValidateBackupFilename(name); err == nil {
			t.Fatalf("invalid name accepted: %q", name)
		}
	}
}

func TestValidateArchiveEntry(t *testing.T) {
	valid := []string{
		"manifest.json",
		"database/pococlinic.sqlite",
		"documents/patient/file.bin",
	}
	for _, name := range valid {
		if err := pathsafe.ValidateArchiveEntry(name); err != nil {
			t.Fatalf("valid entry rejected: %q (%v)", name, err)
		}
	}

	invalid := []string{
		"../outside.txt",
		"documents/../../etc/passwd",
		"/etc/passwd",
	}
	for _, name := range invalid {
		if err := pathsafe.ValidateArchiveEntry(name); err == nil {
			t.Fatalf("invalid entry accepted: %q", name)
		}
	}
}

func TestValidateRelativeKey(t *testing.T) {
	if err := pathsafe.ValidateRelativeKey("patient-id/doc-id"); err != nil {
		t.Fatalf("expected valid key: %v", err)
	}
	if err := pathsafe.ValidateRelativeKey("../secret"); err == nil {
		t.Fatal("expected traversal key to fail")
	}
}

func TestJoinRootWindowsDrive(t *testing.T) {
	if os.PathSeparator != '\\' {
		t.Skip("windows-specific")
	}
	root := t.TempDir()
	_, err := pathsafe.JoinRoot(root, `C:\Windows\System32`)
	if err == nil {
		t.Fatal("expected escape via drive path to fail")
	}
}
