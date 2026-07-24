# Backup and recovery

PocoClinic protects clinic data with **local tar.gz backups** whose database payload is a **SQLite file** (`database/pococlinic.sqlite`, created with `VACUUM INTO`). Encrypted patient documents and exercise logs live inside that file. Optional legacy on-disk document files may also appear under `documents/`.

## What is in a backup

- SQLite database file (patients, staff/users, forms, notes, exercise plans & log entries, encrypted documents, audit log, sessions, migrations)
- Legacy patient document files under `documents/` when still present on disk
- Manifest with checksums and row-count summary

## Daily USB workflow

1. Use **USB rotation** labels (Mon/Wed/Fri)
2. **Backup now** (dashboard) or Backup Helper wizard
3. Copy file to USB; verify before relying on the copy
4. Store USB locked

See [administrator runbook](../../ops/administrator-runbook.md) for a printable checklist.

## Verify a backup

**Admin → Backups → Verify**

Checks:

- Manifest checksums match file contents
- SQLite payload is present and checksum matches
- `PRAGMA quick_check` passes on the embedded database file

Run after USB copy and during quarterly drills.

## Restore (disaster)

**Warning:** Replaces the live SQLite database (including encrypted documents and exercise logs).

1. **Stop the main PocoClinic service** if restoring via Backup Helper / `cmd/restore` (those open the DB in a separate process; replacing the file while the EMR is still running can corrupt data)
2. Sign out **all staff** except one administrator when using in-app Admin restore
3. Restore via Backup Helper (`http://127.0.0.1:9090`) or Admin dashboard
4. Restart the main service if you stopped it; confirm sign-in and patient count
5. Run health check
6. Take a **fresh backup** immediately
7. Log event in ops binder

Keep `DOCUMENT_ENCRYPTION_KEY` from the Safe vault — restoring without the matching key leaves document blobs unreadable.

### Command-line restore (optional)

| Environment | Command |
|-------------|---------|
| **Production (Pi)** | `sudo systemctl stop pococlinic` then `sudo /opt/pococlinic/bin/restore --confirm` |
| **Development** | `restore.bat` from repo root |

## Related

- [clinic/server/](../../../clinic/server/README.md) — production backup/restore wrappers

- [ADR-0012](../../../adr/0012-backup-and-recovery-ops.md)
- [ADR-0015](../../../adr/0015-clinic-database-engine.md)
- [Ops binder](../../ops/binder/README.md)
- [Tools & scripts](../../ops/tools-and-scripts.md)
