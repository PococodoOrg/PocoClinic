#!/usr/bin/env bash
# Thin wrapper — all logic lives in run.py
set -euo pipefail
cd "$(dirname "$0")"
if [[ ! -x .venv/bin/python ]]; then
  python3 -m venv .venv
fi
exec .venv/bin/python run.py "$@"
