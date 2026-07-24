# Shared helper for clinic/server/bin wrappers (installed to /opt/pococlinic/bin).
POCOCLINIC_ROOT="${POCOCLINIC_ROOT:-/opt/pococlinic}"
POCOCLINIC_ENV="${POCOCLINIC_ENV:-/etc/pococlinic/env}"

if [[ -f "${POCOCLINIC_ENV}" ]]; then
  set -a
  # shellcheck source=/dev/null
  source "${POCOCLINIC_ENV}"
  set +a
fi
