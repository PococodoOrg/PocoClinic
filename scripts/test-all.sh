#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "=== PocoClinic: run all tests ==="
echo

bash "$ROOT/scripts/check-no-secrets.sh"
echo

echo "[1/3] Backend (go test ./...)"
(
  cd backend
  go test ./...
)
echo "Backend tests passed."
echo

echo "[2/3] Frontend (vitest + typecheck + build)"
(
  cd frontend
  npx tsc --noEmit
  npm test
  npm run build
)
echo "Frontend tests passed."
echo

echo "[3/3] Loadtest helpers (pytest)"
if [[ -x "$ROOT/loadtest/.venv/bin/python" ]]; then
  "$ROOT/loadtest/.venv/bin/python" -m pytest -q "$ROOT/loadtest"
elif [[ -x "$ROOT/loadtest/.venv/Scripts/python.exe" ]]; then
  "$ROOT/loadtest/.venv/Scripts/python.exe" -m pytest -q "$ROOT/loadtest"
else
  python -m pytest -q "$ROOT/loadtest"
fi
echo "Loadtest tests passed."
echo

echo "=== All tests passed ==="
