# Agent notes (PocoClinic)

For **coding agents and contributors**. Clinic admins: start at [docs/guide/](./docs/guide/README.md) or the [repo README](./README.md).

## Monorepo map

| Path | What it is |
|------|------------|
| [adr/](./adr/README.md) | Architecture Decision Records — **why** we chose things |
| [docs/](./docs/README.md) | Product, ops binders ([hub](./docs/binders/README.md)), vision, features, network |
| [devices/](./devices/README.md) | Per-device ops (start: [Raspberry Pi](./devices/raspberry-pi/README.md)) |
| [binder-printer/](./binder-printer/README.md) | Print binders utility (not the EMR) |
| [clinic/](./clinic/README.md) | Post-deploy ops — server install tree, cron, env (see [ADR-0016](./adr/0016-clinic-runtime-packaging.md)) |
| [backend/](./backend/) | Go API |
| [frontend/](./frontend/) | React EMR UI |
| [ops-helper/](./ops-helper/) | Localhost backup/restore UI |
| [clinic/](./clinic/README.md) | Post-deploy ops — server install tree ([ADR-0016](./adr/0016-clinic-runtime-packaging.md)) |
| [scripts/](./scripts/README.md) | Build / CI helpers (**not** clinic runtime) |
| [loadtest/](./loadtest/README.md) | Local Locust load tests — `loadtest/run.bat` (Windows) |
| `llm-service/` | Placeholder — AI assistance planned only |

Root [README.md](./README.md) is the GitHub front door — keep its **Choose a path** table and repo map current when you add hubs.

## Dev vs clinic runtime (do not mix)

| Path | Put here when… |
|------|----------------|
| `scripts/` | Building, testing, or CI on a dev machine |
| `clinic/server/` | Tool ships on the Pi tarball (`/opt/pococlinic`) |
| `clinic/dev-windows/` | Windows dev mirror of server ops (`go run`) |
| `clinic/workstation/` | Admin-PC utility, not the EMR server |

New production CLI: `backend/cmd/<name>` + `clinic/server/bin/<name>` wrapper + `build-release.mjs` (already copies `clinic/server/bin/`). See [clinic/README.md](./clinic/README.md).

## When you change architecture

1. Read [adr/README.md](./adr/README.md).
2. Update or add an ADR under `adr/` (never under `docs/adr/` — stub only).
3. Refresh the ADR index and root README links if the decision is product-visible.
4. Cursor rules in `.cursor/rules/` (especially `monorepo-adr-docs.mdc` and `product-constraints.mdc`) encode the same habits.

## High-signal ADRs

- [0014 — Private LAN-only](./adr/0014-private-lan-only-deployment.md)
- [0015 — Clinic database engine](./adr/0015-clinic-database-engine.md) (Accepted — testing: SQLite only)
- [0012 — Backup & recovery](./adr/0012-backup-and-recovery-ops.md)
- [0010 — Air-gapped deploy](./adr/0010-local-airgap-deployment.md)
- [0016 — Clinic runtime packaging](./adr/0016-clinic-runtime-packaging.md)
- [0017 — Pi hardening & LAN TLS](./adr/0017-raspberry-pi-hardening-and-lan-tls.md)
