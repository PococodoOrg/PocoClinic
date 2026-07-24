#!/usr/bin/env bash
# Baseline hardening for a PocoClinic Raspberry Pi clinic server.
# Run once after OS install and before go-live. Review output — adjust LAN_CIDR for your site.
#
# Usage:
#   sudo LAN_CIDR=192.168.1.0/24 ADMIN_SSH_CIDR=192.168.1.0/24 ./harden-pi.sh
#
# Does NOT replace clinic firewall/router policy. Does NOT enable automatic internet updates.

set -euo pipefail

LAN_CIDR="${LAN_CIDR:-192.168.0.0/16}"
ADMIN_SSH_CIDR="${ADMIN_SSH_CIDR:-${LAN_CIDR}}"
ENV_FILE="${ENV_FILE:-/etc/pococlinic/env}"
SSHD_CONFIG="/etc/ssh/sshd_config"

if [[ "${EUID:-$(id -u)}" -ne 0 ]]; then
  echo "Run as root: sudo $0"
  exit 1
fi

echo "== PocoClinic Pi hardening =="
echo "LAN_CIDR=${LAN_CIDR}"
echo "ADMIN_SSH_CIDR=${ADMIN_SSH_CIDR}"
echo ""

# --- Kernel / network basics ---
SYSCTL_DROP="/etc/sysctl.d/99-pococlinic-hardening.conf"
cat > "${SYSCTL_DROP}" <<'EOF'
# PocoClinic clinic server — private LAN, not a router
net.ipv4.ip_forward = 0
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
net.ipv4.tcp_syncookies = 1
net.ipv6.conf.all.forwarding = 0
EOF
sysctl --system >/dev/null 2>&1 || true

# --- SSH (back up first) ---
if [[ -f "${SSHD_CONFIG}" ]]; then
  cp -a "${SSHD_CONFIG}" "${SSHD_CONFIG}.bak-pococlinic-$(date +%Y%m%d)"
  apply_sshd() {
    local key="$1" value="$2"
    if grep -qE "^[#[:space:]]*${key}[[:space:]]" "${SSHD_CONFIG}"; then
      sed -i "s/^[#[:space:]]*${key}[[:space:]].*/${key} ${value}/" "${SSHD_CONFIG}"
    else
      echo "${key} ${value}" >> "${SSHD_CONFIG}"
    fi
  }
  apply_sshd "PermitRootLogin" "no"
  apply_sshd "PasswordAuthentication" "no"
  apply_sshd "PubkeyAuthentication" "yes"
  apply_sshd "MaxAuthTries" "3"
  apply_sshd "X11Forwarding" "no"
  apply_sshd "AllowTcpForwarding" "no"
  if sshd -t 2>/dev/null; then
    systemctl reload ssh 2>/dev/null || systemctl reload sshd 2>/dev/null || true
    echo "SSH: hardened (backup at ${SSHD_CONFIG}.bak-pococlinic-*)"
  else
    echo "WARNING: sshd config test failed — restore from backup before reloading SSH"
  fi
fi

# --- Firewall (ufw) ---
if command -v ufw >/dev/null 2>&1; then
  ufw --force default deny incoming
  ufw default allow outgoing
  ufw allow from "${ADMIN_SSH_CIDR}" to any port 22 proto tcp comment 'PocoClinic admin SSH'
  ufw allow from "${LAN_CIDR}" to any port 443 proto tcp comment 'PocoClinic HTTPS'
  # Plain HTTP only from LAN during migration — remove after TLS cutover
  ufw allow from "${LAN_CIDR}" to any port 8080 proto tcp comment 'PocoClinic HTTP (migrate off after TLS)'
  ufw --force enable
  echo "UFW: enabled (443 + 8080 from LAN; SSH from ${ADMIN_SSH_CIDR})"
else
  echo "NOTE: ufw not installed — configure nftables/iptables manually"
fi

# --- PocoClinic file permissions ---
POCO_USER="${POCO_USER:-pococlinic}"
if [[ -f "${ENV_FILE}" ]]; then
  if getent group "${POCO_USER}" >/dev/null 2>&1; then
    chown root:"${POCO_USER}" "${ENV_FILE}"
  fi
  chmod 640 "${ENV_FILE}"
  echo "Secured ${ENV_FILE} (640)"
fi
if id "${POCO_USER}" &>/dev/null; then
  for dir in /var/lib/pococlinic /var/lib/pococlinic/backups /var/lib/pococlinic/documents; do
    if [[ -d "${dir}" ]]; then
      chown -R "${POCO_USER}:${POCO_USER}" "${dir}"
      chmod 750 "${dir}"
    fi
  done
fi
if [[ -d /etc/pococlinic/tls ]]; then
  TLS_GROUP="${POCO_USER}"
  if ! getent group "${TLS_GROUP}" >/dev/null 2>&1; then
    TLS_GROUP="root"
  fi
  chown -R root:"${TLS_GROUP}" /etc/pococlinic/tls
  chmod 750 /etc/pococlinic/tls
  chmod 640 /etc/pococlinic/tls/*.crt 2>/dev/null || true
  chmod 600 /etc/pococlinic/tls/*.key 2>/dev/null || true
fi

# --- Fail2ban (optional) ---
if command -v fail2ban-client >/dev/null 2>&1; then
  systemctl enable --now fail2ban 2>/dev/null || true
  echo "fail2ban: enabled"
else
  echo "NOTE: install fail2ban for SSH brute-force protection (optional)"
fi

echo ""
echo "Hardening pass complete. Manual checklist:"
echo "  • Confirm DNS or hosts entry for pococlinic.local on staff devices"
echo "  • Run clinic/server/scripts/generate-lan-tls.sh and trust the CA"
echo "  • Set SERVER_HOST=127.0.0.1 when using Caddy reverse proxy (see tls-lan.md)"
echo "  • Block port forwarding on the clinic router — no WAN access to this Pi"
echo "  • Record LAN IP and hostname in the safe vault binder"
