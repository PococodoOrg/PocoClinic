# Database migrations

Embedded SQL migrations for **SQLite** ([ADR-0013](../../../../../adr/0013-database-migrations.md), [ADR-0015](../../../../../adr/0015-clinic-database-engine.md)).

## Layout

- `001_schema.sql` — squashed current schema (SQLite types: TEXT / BLOB / INTEGER / REAL)
- `002_patient_field_settings.sql` — `clinic_settings` table and default patient field requirements
- Applied names are recorded in `schema_migrations`

## Running

Production: set `RUN_MIGRATIONS=false` on the app and run migrations as a deploy step:

```bash
sudo /opt/pococlinic/bin/migrate
```

Local dev:

```bash
# From repo root (Windows)
migrate.bat

# Or from backend/
export DATABASE_URL=./data/pococlinic.db
go run ./cmd/migrate
```

## Rules

- Keep migrations **N-1 compatible** when expanding schema
- Prefer additive changes; avoid dropping columns in the same release that still reads them
- Do not use Postgres/Cockroach-only types (`STRING`, `TIMESTAMPTZ`, `JSONB`, `UUID` native, etc.)
