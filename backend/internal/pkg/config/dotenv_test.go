package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDotEnvFileDoesNotOverrideExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(path, []byte("DOTENV_TEST_KEY=from-file\n"), 0o600))

	t.Setenv("DOTENV_TEST_KEY", "from-os")
	require.NoError(t, loadDotEnvFile(path))
	assert.Equal(t, "from-os", os.Getenv("DOTENV_TEST_KEY"))
}

func TestLoadDotEnvFileSetsMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(path, []byte("# comment\nDOTENV_TEST_MISSING=hello # trailing\n"), 0o600))

	os.Unsetenv("DOTENV_TEST_MISSING")
	require.NoError(t, loadDotEnvFile(path))
	assert.Equal(t, "hello", os.Getenv("DOTENV_TEST_MISSING"))
	os.Unsetenv("DOTENV_TEST_MISSING")
}
