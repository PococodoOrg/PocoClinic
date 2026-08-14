#!/usr/bin/env bash
# Interactive PocoClinic setup wizard (repo root).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"
if ! command -v node >/dev/null 2>&1; then
  echo "Node.js is required. Install LTS from https://nodejs.org/"
  exit 1
fi
exec node "$ROOT/scripts/setup.mjs" "$@"
