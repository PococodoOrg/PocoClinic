# Clinic runtime tooling

**Post-deployment** helpers for the clinic server and admin workstation — **not** repository development scripts.

← [Repo README](../README.md) · Deploy contract: [docs/deploy/DEPLOYMENT-BOUNDARY.md](../docs/deploy/DEPLOYMENT-BOUNDARY.md) · Dev/build scripts: [scripts/README.md](../scripts/README.md)

## Why this folder exists

The monorepo mixes two audiences:

| Audience | Examples | Where it lives |
|----------|----------|----------------|
| **Developers / CI** | `test-all.bat`, `build-release.mjs`, `check-no-secrets` | [`scripts/`](../scripts/README.md), repo root dev launchers |
| **Clinic server (Pi)** | `backup`, `restore`, `migrate`, cron, install | **`clinic/server/`** → packaged into `/opt/pococlinic` |
| **Admin workstation** | Print paper binders before go-live | **`clinic/workstation/`** (optional on image) |

Application **source** (Go/React) stays in `backend/`, `frontend/`, and `ops-helper/`. This folder holds **how those binaries are operated** after install.

## Where does a new script go?

```
Does a clinic admin or cron run it on the deployed Pi?
├── YES → backend/cmd/<tool> (Go binary) + clinic/server/bin/<tool> (env wrapper)
│         Also update clinic/server/cron/ if scheduled
└── NO
    ├── Building/testing the app locally on Windows?
    │   └── clinic/dev-windows/<tool>.bat (+ repo-root forwarder optional)
    ├── Admin PC only (not the server)?
    │   └── clinic/workstation/
    └── CI, release build, lint, docs generation?
        └── scripts/
```

When in doubt, read [ADR-0016](../adr/0016-clinic-runtime-packaging.md) and add a row to [tools-and-scripts.md](../docs/ops/tools-and-scripts.md).

## Layout

```
clinic/
├── server/           # Shipped in the release tarball → /opt/pococlinic
│   ├── bin/          # Env-aware wrappers (cron, SSH)
│   ├── cron/         # Example cron drop-ins
│   ├── install.sh    # Pi install/upgrade
│   └── env.template  # First-boot env file template
├── workstation/      # Optional admin-PC tools (not the EMR)
│   └── print-binder.bat
└── dev-windows/      # Local Windows mirrors (go run) — NOT in tarball
    ├── backup.bat
    ├── restore.bat
    └── …
```

## Release packaging

`node scripts/build-release.mjs` copies:

| Tarball path | Source |
|--------------|--------|
| `pococlinic`, `backup`, `restore`, … | `go build` from `backend/cmd/` |
| `bin/*` | `clinic/server/bin/` |
| `cron/` | `clinic/server/cron/` |
| `install.sh`, `env.template` | `clinic/server/` |
| `systemd/` | `scripts/deploy/*.service` |

After install on the Pi, operators use:

```bash
sudo /opt/pococlinic/bin/backup
sudo /opt/pococlinic/bin/migrate
```

Wrappers load `/etc/pococlinic/env` then exec the matching binary in `/opt/pococlinic/`.

## Local development (Windows)

Repo-root `.bat` files are **thin forwarders** to `clinic/dev-windows/` so volunteers can test backup/restore flows without hunting paths. They use `go run ./cmd/…` and are **never** copied to production images.

For day-to-day coding, use `run-backend.bat`, `run-frontend.bat`, and `test-all.bat` at the repo root.

## Related

- [ADR-0016 — Clinic runtime packaging](../adr/0016-clinic-runtime-packaging.md)
- [Tools & scripts catalog](../docs/ops/tools-and-scripts.md)
- [Administrator runbook](../docs/ops/administrator-runbook.md)
