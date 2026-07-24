package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/config"
)

func TestLoadConfigTLSFiles(t *testing.T) {
	dir := t.TempDir()
	cert := filepath.Join(dir, "server.crt")
	key := filepath.Join(dir, "server.key")
	if err := os.WriteFile(cert, []byte("test-cert"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(key, []byte("test-key"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("SERVER_TLS_CERT", cert)
	t.Setenv("SERVER_TLS_KEY", key)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if !cfg.Server.TLSEnabled() {
		t.Fatal("expected TLS to be enabled")
	}
}

func TestLoadConfigTLSMissingFile(t *testing.T) {
	t.Setenv("SERVER_TLS_CERT", "/no/such/cert.crt")
	t.Setenv("SERVER_TLS_KEY", "/no/such/key.key")

	if _, err := config.LoadConfig(); err == nil {
		t.Fatal("expected error for missing TLS files")
	}
}
