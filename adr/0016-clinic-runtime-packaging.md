# ADR-0016: Clinic runtime packaging

## Status

Accepted

## Context

The monorepo contains two kinds of scripts:

1. **Repository / CI tooling** — build release tarballs, run tests, lint, generate docs (`scripts/`, root `run-backend.bat`, `test-all.bat`).
2. **Post-deployment clinic operations** — backup, restore, migrate, audit purge, cron drop-ins, install script, env template, and optional admin-workstation helpers (print binders).

Mixing these under `scripts/` and the repo root made it unclear which files ship on the clinic server image versus which exist only for developers. Operators and the separate microSD image-build process need an explicit contract for what lands in `/opt/pococlinic` after install.

Related: [Deployment boundary](../docs/deploy/DEPLOYMENT-BOUNDARY.md), [ADR-0010](./0010-local-airgap-deployment.md), [ADR-0012](./0012-backup-and-recovery-ops.md).

## Decision

Introduce a top-level **`clinic/`** folder for post-deployment runtime tooling:

| Subfolder | Audience | In release tarball? |
|-----------|----------|---------------------|
| `clinic/server/` | Raspberry Pi / clinic server | **Yes** — install.sh, env.template, cron, `bin/` wrappers |
| `clinic/workstation/` | Admin PC (optional) | **No** — e.g. print-binder launcher |
| `clinic/dev-windows/` | Windows dev mirrors (`go run`) | **No** — repo-root `.bat` forwarders only |

Application **source** stays in `backend/`, `frontend/`, `ops-helper/`. `scripts/build-release.mjs` compiles Go binaries and copies `clinic/server/` into the tarball.

On the deployed system, operators prefer **`/opt/pococlinic/bin/*`** wrappers that load `/etc/pococlinic/env` before exec'ing the matching binary.

## Alternatives considered

| Option | Rejected because |
|--------|------------------|
| Keep ops scripts in `scripts/release/` | Name implies "build scripts"; still mixed with dev tooling |
| Put everything under `devices/` | Device docs are Pi-specific; server install tree applies to any Linux clinic host |
| Only document paths, no folder move | Drift — new helpers would land in the wrong place again |

## Consequences

**Easier**

- Clear split: `scripts/` = build/CI, `clinic/` = what clinics run after deploy
- Image-build process has one folder to read for install/cron/env contract
- Cron and SSH docs can standardize on `/opt/pococlinic/bin/backup`

**Harder**

- Doc links must point at `clinic/` (repo-root `.bat` forwarders preserve old entry points)
- Two places to update when adding a new production CLI (Go `cmd/` + optional `clinic/server/bin/` wrapper)

## References

- [`clinic/README.md`](../clinic/README.md)
- [`clinic/server/README.md`](../clinic/server/README.md)
