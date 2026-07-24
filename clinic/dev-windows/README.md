# Windows dev mirrors (not shipped)

These batch files run **Go CLI tools via `go run`** against a local checkout. They mirror production ops flows for developers and evaluators on Windows.

They are **never** copied into the release tarball. Production uses [`clinic/server/`](../server/README.md) wrappers installed under `/opt/pococlinic/bin/`.

Repo-root `.bat` files (`backup.bat`, etc.) forward here so paths stay stable.

## Scripts

| Script | Maps to |
|--------|---------|
| `migrate.bat` | `go run ./cmd/migrate` |
| `backup.bat` | `go run ./cmd/backup` |
| `restore.bat` | `go run ./cmd/restore --confirm` |
| `run-ops-helper.bat` | Desktop ops-helper UI |
| `run-ops-helper-pi.bat` | Pi-touch ops-helper layout |

Set `DATABASE_URL` before use, e.g. `set DATABASE_URL=./data/pococlinic.db` from `backend/`.
