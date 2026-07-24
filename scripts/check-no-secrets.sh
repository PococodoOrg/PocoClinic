#!/usr/bin/env bash
# Fail if git would stage secrets, PHI, or local env files. Run before commit.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

BLOCKED_PATTERNS=(
  '\.env$'
  '/setup\.env$'
  '/personas\.env$'
  '\.db$'
  '\.db-wal$'
  '\.db-shm$'
  'bootstrap-admin-once\.txt'
  '/backend/data/'
  '^data/documents/'
  'pococlinic-backup-.*\.tar\.gz'
  'BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY'
)

mapfile -t CANDIDATES < <(
  {
    git diff --cached --name-only --diff-filter=ACMRT
    git ls-files --others --exclude-standard
  } | sort -u
)

if ((${#CANDIDATES[@]} == 0)); then
  echo "check-no-secrets: nothing to scan."
  exit 0
fi

FAILED=0
for path in "${CANDIDATES[@]}"; do
  for pattern in "${BLOCKED_PATTERNS[@]}"; do
    if [[ "$path" =~ $pattern ]]; then
      echo "BLOCKED: $path (matches /$pattern/)"
      FAILED=1
    fi
  done
done

if ((${#CANDIDATES[@]} > 0)); then
  while IFS= read -r path; do
    [[ -f "$path" ]] || continue
    if grep -qE 'BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY' "$path" 2>/dev/null; then
      echo "BLOCKED: $path (contains a private key block)"
      FAILED=1
    fi
  done < <(printf '%s\n' "${CANDIDATES[@]}")
fi

if ((FAILED)); then
  echo
  echo "check-no-secrets: remove blocked files from the commit or update .gitignore."
  exit 1
fi

echo "check-no-secrets: OK (${#CANDIDATES[@]} paths scanned)"
