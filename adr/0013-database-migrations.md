# ADR-0013: Database Migrations and N-1 Compatibility

## Status
Accepted

## Context

PocoClinic uses embedded SQL migrations applied to SQLite. Deployments may be:

- **Single-node** (Raspberry Pi, one VM) — one app process, migrate then start
- **Multi-instance** (future cluster, blue/green) — many app replicas sharing one database

Running migrations inside every app instance on startup causes races and ties schema changes to process boot. Rolling deploys also require the **previous app version** to keep working while the schema moves forward.

## Decision

### 1. Migrations as a deployment step

- **`cmd/migrate`** applies pending migrations once, then exits successfully or with error.
- **App (`cmd/main`)** connects with `database.OpenPool` and does **not** migrate when `RUN_MIGRATIONS=false`.
- **`RUN_MIGRATIONS` defaults:**
  - `true` when `ENV != production` (local dev convenience)
  - `false` when `ENV=production` (cluster / Pi production script runs migrate first)

Production upgrade sequence:

```
stop traffic (or drain) → backup → /opt/pococlinic/bin/migrate → deploy new app → start
```

(Dev checkout: `migrate.bat` or `go run ./cmd/migrate`.)

### 2. N-1 schema compatibility (expand / contract)

Each migration in a release must be **backward compatible** with the app version currently in production.

| Phase | Deploy | Schema | App |
|-------|--------|--------|-----|
| **Expand** | 1 | Add new columns/tables; keep old ones | Roll out new app; old app still safe |
| **Contract** | 2 | Drop/rename deprecated schema | All instances already on new app |

Breaking changes **always** take two deployments. Never drop or rename columns still read by N-1 code in the same release.

Examples:

- **One deploy:** add nullable column; add new table; add column with `DEFAULT`
- **Two deploys:** add `email_v2` → deploy app using it → drop `email_v1`

### 3. Migration tracking

- Applied files recorded in `schema_migrations` table
- One transaction per migration file
- Files ordered lexicographically by name (`001_`, `002_`, …)

## Consequences

### Positive

- Safe rolling deploys and cluster rollouts
- Clear operator step for Pi and LAN deploy (`cmd/migrate`)
- Documented rule for agents and contributors (`migrations/README.md`)

### Negative

- Contributors must think in two phases for destructive changes
- Production scripts must run migrate before app (documented in ADR-0010)

## References

- [ADR-0010: Local air-gapped deployment](./0010-local-airgap-deployment.md)
- [ADR-0014: Private LAN-only deployment](./0014-private-lan-only-deployment.md)
- [migrations/README.md](../backend/internal/pkg/database/migrations/README.md)
