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

Set `DATABASE_URL` before use, e.g. copy `backend/.env.example` to `backend/.env` with `DATABASE_URL=./data/pococlinic.db`, then run `migrate.bat`.

`run-ops-helper.bat` checks for Go and Node (on first UI build), creates `data/` and `backups/` if missing, waits for the server before opening the browser, and requires `backend/.env` when `DATABASE_URL` is not set in the shell.
