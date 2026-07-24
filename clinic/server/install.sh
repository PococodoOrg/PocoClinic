#!/usr/bin/env bash
# Install or upgrade PocoClinic from a release tarball on Linux (Raspberry Pi).
#
# Usage (on the Pi, as root or with sudo):
#   tar -xzf pococlinic-1.0.0-linux-arm64.tar.gz
#   cd pococlinic-1.0.0-linux-arm64
#   sudo ./install.sh
#
# Upgrade: run again after extracting a newer tarball over a previous install.
# Data in /var/lib/pococlinic is never removed.

set -euo pipefail

INSTALL_ROOT="${INSTALL_ROOT:-/opt/pococlinic}"
ENV_FILE="${ENV_FILE:-/etc/pococlinic/env}"
POCO_USER="${POCO_USER:-pococlinic}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ "${EUID:-$(id -u)}" -ne 0 ]]; then
  echo "Run as root: sudo ./install.sh"
  exit 1
fi

echo "Installing PocoClinic to ${INSTALL_ROOT}…"

if ! id "${POCO_USER}" &>/dev/null; then
  useradd --system --home-dir "${INSTALL_ROOT}" --shell /usr/sbin/nologin "${POCO_USER}"
fi

mkdir -p "${INSTALL_ROOT}" /var/lib/pococlinic/backups /var/lib/pococlinic/documents /etc/pococlinic /var/log/pococlinic

# Application tree (binaries + static UI) — replaced on upgrade
for item in pococlinic ops-helper ops-helper-static migrate backup restore audit-purge static install.sh env.template MANIFEST.json; do
  if [[ -e "${SCRIPT_DIR}/${item}" ]]; then
    rm -rf "${INSTALL_ROOT}/${item}"
    cp -a "${SCRIPT_DIR}/${item}" "${INSTALL_ROOT}/${item}"
  fi
done

if [[ -d "${SCRIPT_DIR}/bin" ]]; then
  rm -rf "${INSTALL_ROOT}/bin"
  cp -a "${SCRIPT_DIR}/bin" "${INSTALL_ROOT}/bin"
  chmod +x "${INSTALL_ROOT}/bin/"* 2>/dev/null || true
fi

if [[ -d "${SCRIPT_DIR}/docs" ]]; then
  rm -rf "${INSTALL_ROOT}/docs"
  cp -a "${SCRIPT_DIR}/docs" "${INSTALL_ROOT}/docs"
fi

if [[ -d "${SCRIPT_DIR}/cron" ]]; then
  rm -rf "${INSTALL_ROOT}/cron"
  cp -a "${SCRIPT_DIR}/cron" "${INSTALL_ROOT}/cron"
fi

chmod +x "${INSTALL_ROOT}/pococlinic" "${INSTALL_ROOT}/ops-helper" "${INSTALL_ROOT}/migrate" \
  "${INSTALL_ROOT}/backup" "${INSTALL_ROOT}/restore" "${INSTALL_ROOT}/audit-purge"

if [[ ! -f "${ENV_FILE}" ]]; then
  cp "${SCRIPT_DIR}/env.template" "${ENV_FILE}"
  chmod 600 "${ENV_FILE}"
  echo "Created ${ENV_FILE} — edit JWT secrets and ALLOWED_ORIGIN before starting."
else
  echo "Keeping existing ${ENV_FILE}"
fi

chown -R "${POCO_USER}:${POCO_USER}" "${INSTALL_ROOT}" /var/lib/pococlinic /var/log/pococlinic

# systemd
if [[ -d "${SCRIPT_DIR}/systemd" ]]; then
  cp "${SCRIPT_DIR}/systemd/"*.service /etc/systemd/system/
  systemctl daemon-reload
fi

echo ""
echo "Next steps:"
echo "  1. Edit ${ENV_FILE} (secrets, ALLOWED_ORIGIN, paths)"
echo "  2. Ensure DATABASE_URL points at the SQLite file"
echo "  3. sudo ${INSTALL_ROOT}/bin/migrate"
echo "  4. sudo systemctl enable pococlinic pococlinic-ops-helper"
echo "  5. sudo systemctl start pococlinic pococlinic-ops-helper"
echo "  6. Optional cron: sudo cp ${INSTALL_ROOT}/cron/pococlinic-backup /etc/cron.d/"
echo ""
echo "Install complete."
