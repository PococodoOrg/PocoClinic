package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// loadDotEnvFiles loads KEY=VALUE pairs from .env files into the process environment.
// Existing OS environment variables are never overwritten.
// Looks in the current working directory and common monorepo locations.
func loadDotEnvFiles() {
	candidates := []string{
		".env",
		filepath.Join("backend", ".env"),
	}
	// When cwd is backend/, also try parent .env
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(wd), ".env"))
		candidates = append(candidates, filepath.Join(wd, ".env"))
	}

	seen := map[string]struct{}{}
	for _, path := range candidates {
		abs, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		if _, ok := seen[abs]; ok {
			continue
		}
		seen[abs] = struct{}{}
		_ = loadDotEnvFile(abs)
	}
}

func loadDotEnvFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Support optional "export KEY=VALUE"
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		// Inline comments: KEY=value # comment (only if value is unquoted)
		value = strings.TrimSpace(value)
		if len(value) >= 2 {
			if (value[0] == '"' && strings.HasSuffix(value, `"`)) || (value[0] == '\'' && strings.HasSuffix(value, `'`)) {
				value = value[1 : len(value)-1]
			} else if i := strings.Index(value, " #"); i >= 0 {
				value = strings.TrimSpace(value[:i])
			}
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}
	return scanner.Err()
}
