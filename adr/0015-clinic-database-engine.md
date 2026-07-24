# ADR-0015: Clinic database engine (SQLite)

## Status

**Accepted — testing in progress.**

SQLite is the **only** supported database engine for PocoClinic. CockroachDB / Postgres-compatible URLs have been removed from the product path. Clinic validation (Pi-sized hardware, backup/restore drills, encrypted documents) continues under this ADR.

## Context

PocoClinic production shape ([ADR-0014](./0014-private-lan-only-deployment.md), [ADR-0010](./0010-local-airgap-deployment.md)):

- **One** clinic server (often Raspberry Pi 4/5, 4–8 GB RAM)
- **No** multi-region cluster
- Operators are volunteers — they must not run a distributed database
- Backups are **USB tarballs** ([ADR-0012](./0012-backup-and-recovery-ops.md))
- Preference: **database-native** backup (the DB file itself), not app-layer table dumps

An early prototype used CockroachDB via `pgx`. That was overkill for a single-clinic LAN EMR: heavy process, poor mental model for volunteers, and it pushed us toward fragile JSONL exports.

## Decision

1. **Default (and only) engine: SQLite** via pure-Go `modernc.org/sqlite` (no CGO — important for Pi cross-builds and air-gap packaging).
2. **`DATABASE_URL` is a filesystem path** to the DB file (e.g. `./data/pococlinic.db` or `/var/lib/pococlinic/pococlinic.db`). Postgres/Cockroach URLs are rejected.
3. **Backup payload:** USB `tar.gz` contains `database/pococlinic.sqlite` produced with `VACUUM INTO` (consistent copy while the clinic can stay up). Format id: `sqlite`. Legacy JSONL / Cockroach-era bundles are not restored.
4. **Schema:** squashed SQLite migrations under `backend/internal/pkg/database/migrations/` ([ADR-0013](./0013-database-migrations.md)).
5. **No separate DB daemon** in compose, systemd, or runbooks.

## Alternatives considered

| Option | Verdict |
|--------|---------|
| Keep CockroachDB | **Rejected** — overkill for one Pi clinic |
| PostgreSQL as default | **Rejected as default** — still a service to babysit; SQLite is simpler |
| Forever JSONL dumps | **Rejected** — miss tables easily; restore trust is “our importer” |

## Consequences

**Easier**

- Pi images and runbooks shrink (no DB unit)
- Backup = copy of one file inside the USB tar + SHA-256
- Dev setup: set a path, run migrate, start app

**Harder / testing focus**

- Concurrent writers are serialized (fine for clinic scale — verify under light multi-tablet load)
- Restore replaces the live file (controlled downtime; app reopens the handle)
- Contributors must write SQLite-friendly SQL (`?` rebound from `$n` in `database.DB`)

## Production SQLite practices (required)

Implemented in `backend/internal/pkg/database/pool.go` via modernc DSN `_pragma=` (not a one-shot `Exec`, which only configures the first pooled connection):

| Setting | Value | Why |
|---------|-------|-----|
| `journal_mode` | `WAL` | Readers do not block a writer |
| `foreign_keys` | `ON` | SQLite defaults to OFF per connection |
| `busy_timeout` | `5000` | Wait on lock instead of immediate `SQLITE_BUSY` |
| `synchronous` | `NORMAL` | Durable with WAL; fewer fsyncs than `FULL` |
| `temp_store` | `MEMORY` | Avoid temp files on SD cards |
| `wal_autocheckpoint` | `1000` | Bound WAL growth |
| `_txlock` | `immediate` | Take write lock at `BEGIN` (avoids upgrade deadlocks) |
| `MaxOpenConns` | `1` | Single writer connection; safe under `database/sql` pooling |

Also:

- Backup (`VACUUM INTO`) and restore run `PRAGMA quick_check` before trusting the file
- Ops-helper restore: stop the main EMR first (two processes must not hold the same DB file across a file replace)
- Deploy: `cmd/migrate` as a step; app runs with `RUN_MIGRATIONS=false` in production

## Testing in progress (checklist)

- [x] Open schema on SQLite + migrate runner
- [x] Encrypted document blobs in DB (AES-GCM)
- [x] `cmd/backup` / restore using **SQLite file** payload (`VACUUM INTO`)
- [x] Production pragmas via DSN + `quick_check` on backup/restore
- [ ] Clinic / Pi soak: memory/CPU under light concurrent staff use
- [ ] Printed restore drill with SQLite bundle on spare hardware (EMR stopped during restore)
- [x] Confirm systemd/deploy scripts and Pi binder copy no longer mention Cockroach

## Related

- [ADR-0010](./0010-local-airgap-deployment.md)
- [ADR-0012](./0012-backup-and-recovery-ops.md)
- [ADR-0013](./0013-database-migrations.md)
- [ADR-0014](./0014-private-lan-only-deployment.md)
- [Devices · Raspberry Pi](../devices/raspberry-pi/README.md)
