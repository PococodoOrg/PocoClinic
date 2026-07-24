# Tools & scripts catalog

Every helper, CLI entrypoint, batch wrapper, and ops script in the PocoClinic repository. Use this as the single index — when you add a new helper, add a row here.

**Deployment context:** the end goal is a **microSD card** inserted in a **Raspberry Pi** that boots into a clinic server. Pi autostart, OS image flashing, and kiosk display setup are owned by a **separate image-build process** — see [Deployment boundary](../deploy/DEPLOYMENT-BOUNDARY.md).

**Audience split:** [ADR-0016](../../adr/0016-clinic-runtime-packaging.md) — **`clinic/`** holds post-deploy ops (shipped on the server); **`scripts/`** holds build/CI tooling; repo-root `.bat` dev launchers forward to `clinic/dev-windows/` for local testing.

---

## Quick reference

| I want to… | Use |
|------------|-----|
| **Build Pi release tarball** | `build-release.bat` or `node scripts/build-release.mjs` |
| Run everything locally (Windows dev) | `run-all.bat` |
| Set DATABASE_URL to a SQLite file path | `set DATABASE_URL=./data/pococlinic.db` |
| Apply database migrations (dev) | `migrate.bat` → [`clinic/dev-windows/migrate.bat`](../../clinic/dev-windows/migrate.bat) |
| Create a backup tarball (dev) | `backup.bat` → [`clinic/dev-windows/backup.bat`](../../clinic/dev-windows/backup.bat) |
| Restore from backup (dev) | `restore.bat` → [`clinic/dev-windows/restore.bat`](../../clinic/dev-windows/restore.bat) |
| Purge old audit logs (dev) | `go run ./cmd/audit-purge --days 365 --dry-run` |
| Purge old audit logs (Pi) | `sudo /opt/pococlinic/bin/audit-purge --days 365 --dry-run` |
| **Migrate on deployed Pi** | `sudo /opt/pococlinic/bin/migrate` |
| **Backup on deployed Pi** | `sudo /opt/pococlinic/bin/backup` |
| Build backup helper UI (desktop) | `build-ops-helper.bat` |
| Build backup helper UI (Pi touchscreen) | `build-ops-helper-pi.bat` |
| Run backup helper locally | `run-ops-helper.bat` |
| Full-screen Pi kiosk browser | `scripts/pi/start-backup-kiosk.sh` |
| **Print ops / device binders** | `print-binder.bat` → [`clinic/workstation/`](../../clinic/workstation/README.md) → [binder-printer](../../binder-printer/README.md) |
| Generate static evaluator help site | `node scripts/build-help-site.mjs` |
| Install production systemd units | [scripts/deploy/README.md](../../scripts/deploy/README.md) |
| **Server install tree (tarball)** | [`clinic/server/`](../../clinic/server/README.md) |

---

## Windows batch wrappers (repo root)

Convenience **forwarders** for local development on Windows. Ops scripts live under [`clinic/dev-windows/`](../../clinic/dev-windows/README.md); root `.bat` files delegate there so old paths keep working.

| Script | Forwards to | Prerequisites |
|--------|-------------|---------------|
| `run-all.bat` | _(repo root)_ Opens backend + frontend in separate terminals | Go, Node.js |
| `run-backend.bat` | _(repo root)_ Starts API server | Go, `DATABASE_URL` optional |
| `run-frontend.bat` | _(repo root)_ Vite dev server | Node.js |
| `migrate.bat` | `clinic/dev-windows/migrate.bat` | `DATABASE_URL` required |
| `backup.bat` | `clinic/dev-windows/backup.bat` | `DATABASE_URL`, `BACKUP_DIR` |
| `restore.bat` | `clinic/dev-windows/restore.bat` | `DATABASE_URL`, `BACKUP_DIR` |
| `build-ops-helper.bat` | _(repo root)_ Builds ops-helper SPA | Node.js |
| `build-ops-helper-pi.bat` | _(repo root)_ Pi touch layout | Node.js |
| `run-ops-helper.bat` | `clinic/dev-windows/run-ops-helper.bat` | Go, built static UI |
| `run-ops-helper-pi.bat` | `clinic/dev-windows/run-ops-helper-pi.bat` | Go, Pi-built static UI |
| `print-binder.bat` | `clinic/workstation/print-binder.bat` | Node.js (admin PC) |

---

## Clinic server runtime (`clinic/server/`)

Copied into the release tarball and installed under `/opt/pococlinic`. See [`clinic/server/README.md`](../../clinic/server/README.md).

| Path | Purpose |
|------|---------|
| `clinic/server/install.sh` | First install / upgrade on Pi |
| `clinic/server/env.template` | Seeds `/etc/pococlinic/env` |
| `clinic/server/cron/` | Example `/etc/cron.d` drop-ins |
| `clinic/server/bin/*` | Wrappers that source env, then exec binary |

Production cron and manual ops should use **`/opt/pococlinic/bin/backup`**, not bare binaries without env.

---

## Go CLI tools (`backend/cmd/`)

Standalone binaries intended for production cron, systemd `ExecStart`, and manual admin use. All read configuration from environment variables via `internal/pkg/config` (see [Environment variables](#environment-variables)).

| Command | Purpose | Typical production use |
|---------|---------|------------------------|
| `cmd/main.go` | EMR API + production static frontend (`ENV=production`, `STATIC_DIR`) | `pococlinic.service` → `/opt/pococlinic/pococlinic` |
| `cmd/migrate` | Apply embedded SQL migrations | Before first boot; after every upgrade |
| `cmd/backup` | Create `.tar.gz` backup (SQLite file via `VACUUM INTO` + optional legacy docs + manifest) | Daily cron; admin dashboard |
| `cmd/restore` | Restore DB + document files from tarball | Disaster recovery; ops helper wizard |
| `cmd/audit-purge` | Purge audit rows older than retention | Monthly cron |
| `cmd/ops-helper` | Localhost-only backup/restore web UI | `pococlinic-ops-helper.service` → port 9090 |

### Examples

**Production (Pi — prefer `bin/` wrappers):**

```bash
sudo /opt/pococlinic/bin/migrate
sudo /opt/pococlinic/bin/backup
sudo /opt/pococlinic/bin/restore --confirm
sudo /opt/pococlinic/bin/audit-purge --days 365 --dry-run
```

**Development (local checkout):**

```bash
cd backend
export DATABASE_URL=./data/pococlinic.db

go run ./cmd/migrate
go run ./cmd/backup
go run ./cmd/restore --confirm
go run ./cmd/audit-purge --days 365 --dry-run
```

Windows: `migrate.bat`, `backup.bat`, `restore.bat` at repo root (forward to `clinic/dev-windows/`).

---

## Shell scripts (`scripts/`)

| Path | Purpose | Notes |
|------|---------|-------|
| `scripts/build-release.mjs` | **Production release tarball** for Linux ARM64 (Pi) | Output: `dist/pococlinic-<ver>-linux-arm64.tar.gz` |
| `scripts/build-help-site.mjs` | Converts `docs/guide/**/*.md` → static HTML in `docs/guide-site/` | No npm deps; run with Node 18+ |
| `scripts/pi/start-backup-kiosk.sh` | Launches Chromium kiosk → `http://127.0.0.1:9090/?pi=1` | Requires X11, Chromium, optional `unclutter` |
| `scripts/deploy/pococlinic.service` | systemd unit for main EMR | [Install guide](../../scripts/deploy/README.md) |
| `scripts/deploy/pococlinic-ops-helper.service` | systemd unit for ops helper | Binds `127.0.0.1` only |

---

## Ops helper (touchscreen backup kiosk)

Separate small React app for **server-local** backup and restore — not the clinical EMR.

| Piece | Location |
|-------|----------|
| Source | `ops-helper/` (Vite + React + Mantine) |
| Pi touch build | `npm run build:pi` (sets touch layout) |
| Embedded static output | `backend/cmd/ops-helper/static/` |
| Server entrypoint | `backend/cmd/ops-helper/main.go` |
| Admin docs | [devices/raspberry-pi/touchscreen-backup.md](../../devices/raspberry-pi/touchscreen-backup.md) |

The helper listens on **`127.0.0.1:9090` only** — never expose it on the clinic LAN.

---

## Docker (not used)

There is **no** `docker-compose` database service. SQLite is a single file — see `docker-compose.yml` (empty stub) and [ADR-0015](../../adr/0015-clinic-database-engine.md).

---

## In-app admin tools (not CLI, but ops-facing)

These ship in the EMR frontend and call admin API endpoints — document them here so admins know they exist.

| UI location | Backend endpoint | Purpose |
|-------------|------------------|---------|
| Admin → Overview → System health check | `GET /admin/health-check` | DB, storage, migrations, document integrity, disk usage |
| Admin → Overview → HIPAA compliance checklist | `GET /admin/compliance-check` | Automated operational controls |
| Admin → Reports → Clinic activity | `GET /admin/reports/activity-summary` | Sign-ins, patient work, backups by period |
| Admin → Audit log → Download CSV | `GET /admin/audit-logs/export.csv` | Admin-only audit log download (30 days default); archive offline before purge |
| Admin → Backup & maintenance | `POST /admin/backups`, verify, restore | Same operations as CLI, with audit logging |

---

## Documentation generators

| Tool | Input | Output |
|------|-------|--------|
| `scripts/build-help-site.mjs` | `docs/guide/` markdown | `docs/guide-site/` static HTML (gitignored) |

In-app help is **not** generated — it lives in `frontend/src/help/content.ts` and must be updated manually when procedures change (see [IN-APP-HELP.md](../IN-APP-HELP.md)).

---

## Environment variables

Shared by all Go binaries via `backend/internal/pkg/config`. Production values live in `/etc/pococlinic/env` (see [scripts/deploy/README.md](../../scripts/deploy/README.md)).

| Variable | Default (dev) | Production notes |
|----------|---------------|------------------|
| `ENV` | `development` | Set `production` on clinic server |
| `APP_VERSION` | `dev` | Match release tarball version |
| `SERVER_HOST` | `localhost` | `0.0.0.0` to listen on LAN |
| `SERVER_PORT` | `8080` | Clinic workstations reach this port |
| `DATABASE_URL` | _(empty = in-memory)_ | Required in production |
| `RUN_MIGRATIONS` | auto: off in production | Prefer explicit `cmd/migrate` on upgrade |
| `BACKUP_DIR` | `./backups` | e.g. `/var/lib/pococlinic/backups` |
| `DOCUMENT_ENCRYPTION_KEY` | (dev default) | `openssl rand -base64 32` — required in production |
| `DOCUMENTS_DIR` | `./data/documents` | Legacy files only; new uploads store encrypted blobs in DB |
| `JWT_ACCESS_SECRET` | dev placeholder | **Required** in production (≥32 chars, non-default; app refuses to start otherwise) |
| `JWT_REFRESH_SECRET` | dev placeholder | **Required** in production (≥32 chars, non-default, must differ from access secret) |
| `ACCESS_TOKEN_COOKIE` | `poco_access_token` | HttpOnly access session cookie name |
| `REFRESH_TOKEN_COOKIE` | `poco_refresh_token` | HttpOnly refresh cookie name |
| `ALLOWED_ORIGIN` | `http://localhost:3000` | Clinic LAN URL (e.g. `http://pococlinic.local`) |
| `TRUSTED_PROXIES` | _(unset)_ | Comma-separated reverse-proxy IPs/CIDRs when the app sits behind nginx/Caddy on the LAN (e.g. `127.0.0.1,192.168.1.10`). Enables accurate `ClientIP` for rate limits and audit logs. Leave unset for direct workstation access. |
| `AUDIT_RETENTION_DAYS` | _(unset)_ | Used by `cmd/audit-purge` |
| `OPS_HELPER_HOST` | `127.0.0.1` | Keep localhost-only |
| `OPS_HELPER_PORT` | `9090` | Pi touchscreen kiosk target |
| `OPS_HELPER_MAIN_APP_URL` | `http://localhost:3000` | Link back to EMR in helper UI |

---

## Adding a new helper

1. Implement the tool (`backend/cmd/…`, `scripts/…`, or `clinic/server/bin/` wrapper).
2. If it runs on the clinic server after deploy, add files under **`clinic/server/`** and wire `build-release.mjs`.
3. Add a row to the appropriate section **in this file**.
4. If it affects clinic procedures, update `docs/ops/administrator-runbook.md` and the matching in-app help article.
5. If it affects the release artifact layout, update [Deployment boundary](../deploy/DEPLOYMENT-BOUNDARY.md) and [ADR-0016](../../adr/0016-clinic-runtime-packaging.md).
