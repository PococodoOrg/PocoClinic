# Production deployment (Linux / Raspberry Pi)

Example systemd units for running PocoClinic on a clinic server. Adjust paths and user for your environment.

**First-time Pi install?** Prefer the beginner walkthrough: [devices/raspberry-pi/setup.md](../../devices/raspberry-pi/setup.md).

## Prerequisites

- SQLite file path configured
- Release tarball built with `node scripts/build-release.mjs`, or manual binary build
- Environment file: `/etc/pococlinic/env` with real `JWT_*`, `DOCUMENT_ENCRYPTION_KEY`, and matching `ALLOWED_ORIGIN`

## Quick install from release tarball

```bash
# On the Pi (after copying tarball from dev machine or CI)
tar -xzf pococlinic-1.0.0-linux-arm64.tar.gz
cd pococlinic-1.0.0-linux-arm64
sudo ./install.sh
sudo nano /etc/pococlinic/env    # JWT_*, DOCUMENT_ENCRYPTION_KEY, ALLOWED_ORIGIN
sudo /opt/pococlinic/bin/migrate
sudo systemctl enable --now pococlinic pococlinic-ops-helper
sudo cp /opt/pococlinic/cron/pococlinic-backup /etc/cron.d/
```

## Manual build (without release script)

- Built backend binary: `go build -o /opt/pococlinic/pococlinic ./cmd/main.go`
- Frontend static files in `/opt/pococlinic/static/` (served by main binary when `ENV=production`)

Example `/etc/pococlinic/env`:

```bash
ENV=production
APP_VERSION=1.0.0
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DATABASE_URL=/var/lib/pococlinic/pococlinic.db
BACKUP_DIR=/var/lib/pococlinic/backups
DOCUMENTS_DIR=/var/lib/pococlinic/documents
STATIC_DIR=/opt/pococlinic/static
JWT_ACCESS_SECRET=use-openssl-rand-base64-48
JWT_REFRESH_SECRET=use-different-openssl-rand-base64-48
DOCUMENT_ENCRYPTION_KEY=use-openssl-rand-base64-32
ALLOWED_ORIGIN=http://192.168.1.50:8080
COOKIE_SECURE=false
AUDIT_RETENTION_DAYS=365
```

(Generate secrets with `openssl rand`; set `ALLOWED_ORIGIN` to the exact browser URL. Keep `COOKIE_SECURE=false` only until LAN TLS — see [Pi setup](../../devices/raspberry-pi/setup.md).)
## Install services

```bash
sudo cp scripts/deploy/pococlinic.service /etc/systemd/system/
sudo cp scripts/deploy/pococlinic-ops-helper.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable pococlinic pococlinic-ops-helper
sudo systemctl start pococlinic
```

## Migrations

Run before first start or after upgrades (prefer `bin/` wrapper):

```bash
sudo /opt/pococlinic/bin/migrate
```

Or manually with env loaded (fallback if `bin/` wrapper unavailable):

```bash
set -a && source /etc/pococlinic/env && set +a
/opt/pococlinic/migrate   # bare binary — prefer bin/migrate wrapper
```

## Backup schedule

Daily backup via cron (example 6 PM) — use tarball `cron/` drop-in or:

```cron
0 18 * * * root /opt/pococlinic/bin/backup >> /var/log/pococlinic/backup.log 2>&1
```

## Audit log rotation (optional)

Monthly purge (example first Sunday 2 AM):

```cron
0 2 1 * * root /opt/pococlinic/bin/audit-purge >> /var/log/pococlinic/audit-purge.log 2>&1
```

## Related

- [First-time setup](../../docs/guide/for-administrators/first-time-setup.md)
- [Clinic network requirements](../../docs/guide/clinic-network.md)
- [clinic/server/](../../clinic/server/README.md) — install script, env template, cron, bin wrappers
- [Deployment boundary](../../docs/deploy/DEPLOYMENT-BOUNDARY.md) — microSD / Pi release artifact vs image-build process
