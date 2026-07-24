# ADR-0012: Backup and Recovery Operations

## Status
Accepted

## Context

PocoClinic operators are often **non-technical volunteers or office staff**. The system must survive disk failure, bad updates, and human error without requiring a sysadmin. Recovery procedures must fit in a **printed Administrator’s Guide** with checklists and labeled USB drives.

Vision reference: [VISION.md](../docs/VISION.md) — pillar 3. Deployment: [NETWORK-AND-SECURITY.md](../docs/NETWORK-AND-SECURITY.md), [ADR-0014](./0014-private-lan-only-deployment.md).

## Decision

We will design backup and restore as **first-class, opinionated workflows**, not DBA exercises.

### 1. Backup bundle (`pococlinic-backup-YYYYMMDD.tar.gz`)

Each backup is a single file suitable for USB copy, containing:

| Component | Contents |
|-----------|----------|
| `manifest.json` | Version, timestamp, checksums, app version, `databaseFormat: "sqlite"` |
| `database/pococlinic.sqlite` | Consistent SQLite file copy (`VACUUM INTO`) — includes encrypted document blobs and all clinic tables |
| `documents/` | Optional legacy on-disk document files |
| `config/` | Non-secret server config snapshot (when present) |

Secrets (JWT keys, `DOCUMENT_ENCRYPTION_KEY`) are **not** in the bundle by default; keep them on the Safe vault sheet.

### 2. Operator workflow (daily)

1. Insert labeled USB drive (rotation: e.g. Mon / Wed / Fri)
2. Run **one action**: desktop shortcut, `./scripts/backup.sh`, or admin UI “Backup now”
3. Wait for **green success** indicator + manifest written
4. Remove USB; store in locked location
5. Confirm Admin dashboard backup status is green (no daily paper log required)

### 3. Restore workflow (disaster)

Printed runbook sections:

1. **Prepare** — spare Pi/PC or re-flashed SD card with PocoClinic install media
2. **Stop** clinic services
3. **Restore** — single command or admin UI “Restore from USB” selecting bundle
4. **Verify** — health check + login test + patient count spot-check
5. **Sign-off** — date and initials on restore log in binder

Target: completable in **under one hour** following the runbook without SSH expertise.

### 4. Verification

- Backup command validates checksum after write
- **Quarterly restore drill** documented in binder (restore to test device, not production cutover)
- Admin dashboard shows: last backup time, backup age warning (>24h yellow, >48h red)

### 5. USB rotation

- Physical labels: clinic name, day, drive letter/id
- Never rely on a single USB; minimum two-drive rotation
- Runbook warns: encrypt drives if taken off-site (future: optional LUKS guide)

## Consequences

### Positive
- Non-tech staff can protect patient data
- Clear disaster recovery story for boards and donors
- Aligns with “Simple but Secure” philosophy

### Negative
- USB media can fail or be lost—rotation and off-site policy needed
- Full restore downtime during recovery
- Must keep backup format stable across app versions

### Mitigations
- Manifest version field + migration tools for older backups
- Automated backup reminders in admin UI
- Restore drill checklist reduces panic during real incidents

## Implementation status

**SQLite native payload** — `cmd/backup` / `cmd/restore` create versioned `tar.gz` bundles with manifest checksums and a **`database/pococlinic.sqlite`** file (`VACUUM INTO`). See [ADR-0015](./0015-clinic-database-engine.md). Legacy Cockroach-era JSONL archives are not restored.

Admin dashboard can trigger backups and shows backup age warnings. Printable runbooks exist under `docs/ops/` and `devices/`.
