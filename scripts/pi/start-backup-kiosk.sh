#!/usr/bin/env bash
# Launch Chromium in kiosk mode for the PocoClinic backup helper on a Pi touchscreen.
# Install: sudo apt install chromium-browser unclutter
# Usage: ./start-backup-kiosk.sh

set -euo pipefail

HELPER_URL="${HELPER_URL:-http://127.0.0.1:9090/?pi=1}"
DISPLAY="${DISPLAY:-:0}"

# Hide mouse cursor after idle (optional)
if command -v unclutter >/dev/null 2>&1; then
  unclutter -idle 0.5 -root &
fi

# Disable screen blanking while kiosk runs
xset s off
xset -dpms
xset s noblank

exec chromium-browser \
  --kiosk \
  --noerrdialogs \
  --disable-infobars \
  --check-for-update-interval=31536000 \
  --app="${HELPER_URL}"
