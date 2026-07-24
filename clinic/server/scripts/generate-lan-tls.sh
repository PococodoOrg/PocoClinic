#!/usr/bin/env bash
# Generate a clinic-local CA and server certificate for HTTPS on the private LAN.
# Works offline — requires openssl only (no Let's Encrypt, no cloud).
#
# Usage (on the Pi, as root):
#   sudo POCO_HOSTNAME=pococlinic.local POCO_LAN_IP=192.168.1.50 ./generate-lan-tls.sh
#
# Output: /etc/pococlinic/tls/{ca.crt,ca.key,server.crt,server.key}
# Trust ca.crt on staff iPads/PCs, then browse https://pococlinic.local

set -euo pipefail

POCO_HOSTNAME="${POCO_HOSTNAME:-pococlinic.local}"
POCO_LAN_IP="${POCO_LAN_IP:-}"
CERT_DIR="${CERT_DIR:-/etc/pococlinic/tls}"
DAYS="${DAYS:-825}" # ~27 months

if [[ "${EUID:-$(id -u)}" -ne 0 ]]; then
  echo "Run as root: sudo $0"
  exit 1
fi

if ! command -v openssl >/dev/null 2>&1; then
  echo "openssl is required."
  exit 1
fi

if [[ ! "${POCO_HOSTNAME}" =~ ^[a-zA-Z0-9]([a-zA-Z0-9.-]*[a-zA-Z0-9])?$ ]]; then
  echo "Invalid POCO_HOSTNAME: ${POCO_HOSTNAME}"
  exit 1
fi

if [[ -n "${POCO_LAN_IP}" && ! "${POCO_LAN_IP}" =~ ^[0-9.]+$ ]]; then
  echo "Invalid POCO_LAN_IP: ${POCO_LAN_IP}"
  exit 1
fi

POCO_USER="${POCO_USER:-pococlinic}"
POCO_GROUP="${POCO_GROUP:-${POCO_USER}}"
if ! getent group "${POCO_GROUP}" >/dev/null 2>&1; then
  POCO_GROUP="root"
fi

install -d -m 750 -o root -g "${POCO_GROUP}" "${CERT_DIR}"
WORK="$(mktemp -d)"
trap 'rm -rf "${WORK}"' EXIT

CA_KEY="${WORK}/ca.key"
CA_CSR="${WORK}/ca.csr"
SRV_KEY="${WORK}/server.key"
SRV_CSR="${WORK}/server.csr"
EXT="${WORK}/server.ext"

openssl genrsa -out "${CA_KEY}" 4096
openssl req -new -x509 -days "${DAYS}" -key "${CA_KEY}" -out "${WORK}/ca.crt" \
  -subj "/CN=PocoClinic Clinic LAN CA/O=PocoClinic/C=US"

openssl genrsa -out "${SRV_KEY}" 2048
openssl req -new -key "${SRV_KEY}" -out "${SRV_CSR}" \
  -subj "/CN=${POCO_HOSTNAME}/O=PocoClinic/C=US"

{
  echo "subjectAltName = @alt_names"
  echo "[alt_names]"
  echo "DNS.1 = ${POCO_HOSTNAME}"
  echo "DNS.2 = localhost"
  if [[ -n "${POCO_LAN_IP}" ]]; then
    echo "IP.1 = ${POCO_LAN_IP}"
  fi
  echo "IP.2 = 127.0.0.1"
} > "${EXT}"

openssl x509 -req -days "${DAYS}" \
  -in "${SRV_CSR}" \
  -CA "${WORK}/ca.crt" -CAkey "${CA_KEY}" -CAcreateserial \
  -out "${WORK}/server.crt" \
  -extfile "${EXT}"

install -m 640 -o root -g "${POCO_GROUP}" "${WORK}/ca.crt" "${CERT_DIR}/ca.crt"
install -m 600 -o root -g root "${CA_KEY}" "${CERT_DIR}/ca.key"
install -m 640 -o root -g "${POCO_GROUP}" "${WORK}/server.crt" "${CERT_DIR}/server.crt"
install -m 600 -o root -g "${POCO_GROUP}" "${SRV_KEY}" "${CERT_DIR}/server.key"

echo ""
echo "TLS files written to ${CERT_DIR}"
echo ""
echo "Next steps:"
echo "  1. Add to /etc/pococlinic/env:"
echo "       SERVER_TLS_CERT=${CERT_DIR}/server.crt"
echo "       SERVER_TLS_KEY=${CERT_DIR}/server.key"
echo "       SERVER_PORT=443"
echo "       ALLOWED_ORIGIN=https://${POCO_HOSTNAME}"
echo "       TRUSTED_PROXIES=127.0.0.1"
echo "  2. Trust ${CERT_DIR}/ca.crt on clinic iPads and PCs (see devices/raspberry-pi/tls-lan.md)"
echo "  3. sudo systemctl restart pococlinic"
echo "  4. Open https://${POCO_HOSTNAME} from a staff browser"
echo ""
echo "Keep ${CERT_DIR}/ca.key offline and locked — it can mint new server certs."
