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

func TestOpenBackupFileRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	_, err := pathsafe.OpenBackupFile(root, "../pococlinic-backup-evil.tar.gz")
	if err == nil {
		t.Fatal("expected traversal filename to be rejected")
	}
}

func TestOpenBackupFileOpensValidBundle(t *testing.T) {
	root := t.TempDir()
	filename := "pococlinic-backup-20260101-120000.tar.gz"
	path, err := pathsafe.JoinRoot(root, filename)
	if err != nil {
		t.Fatalf("JoinRoot: %v", err)
	}
	if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	file, err := pathsafe.OpenBackupFile(root, filename)
	if err != nil {
		t.Fatalf("OpenBackupFile: %v", err)
	}
	_ = file.Close()
}

func TestStatBackupFile(t *testing.T) {
	root := t.TempDir()
	filename := "pococlinic-backup-20260101-120000.tar.gz"
	path, err := pathsafe.JoinRoot(root, filename)
	if err != nil {
		t.Fatalf("JoinRoot: %v", err)
	}
	if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	info, err := pathsafe.StatBackupFile(root, filename)
	if err != nil {
		t.Fatalf("StatBackupFile: %v", err)
	}
	if info.Size() != 4 {
		t.Fatalf("unexpected size: %d", info.Size())
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

func TestResolveBasenameUnderRootRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	secretDir := filepath.Join(filepath.Dir(root), "secret-area")
	if err := os.MkdirAll(secretDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	secretFile := filepath.Join(secretDir, "leaked.txt")
	if err := os.WriteFile(secretFile, []byte("secret"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	attempts := []string{
		"../secret-area/leaked.txt",
		"..\\secret-area\\leaked.txt",
		".../...//",
	}
	for _, name := range attempts {
		_, err := pathsafe.OpenBackupFile(root, name)
		if err == nil {
			t.Fatalf("expected rejection for %q", name)
		}
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
