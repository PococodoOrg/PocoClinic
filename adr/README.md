# Architecture Decision Records (ADRs)

**Canonical location:** this folder at the **repo root** (`/adr`).  
Read these on GitHub like any other markdown — no special tooling required.

ADRs capture **why** we chose something, what we rejected, and what that implies for a single-clinic, private-LAN EMR (often on a Raspberry Pi).

## Start here

| Priority | ADR | One-line why |
|----------|-----|----------------|
| 1 | [0014 — Private LAN-only](./0014-private-lan-only-deployment.md) | No cloud, no mobile apps, clinic browsers on the LAN only |
| 2 | [0015 — Clinic database engine](./0015-clinic-database-engine.md) | **Accepted (testing):** SQLite file DB — no Cockroach |
| 3 | [0010 — Local air-gapped deploy](./0010-local-airgap-deployment.md) | Offline-capable install; USB updates |
| 4 | [0012 — Backup & recovery](./0012-backup-and-recovery-ops.md) | USB backups non-technical admins can run |
| 5 | [0011 — Badge + PIN](./0011-employee-badge-authentication.md) | Staff auth shaped for front desk, not SSO |
| 6 | [0016 — Clinic runtime packaging](./0016-clinic-runtime-packaging.md) | `clinic/` vs `scripts/` — what ships on the Pi |
| 7 | [0017 — Pi hardening & LAN TLS](./0017-raspberry-pi-hardening-and-lan-tls.md) | Firewall, local CA, HTTPS on the LAN |

Related (not an ADR): [Deployment boundary](../docs/deploy/DEPLOYMENT-BOUNDARY.md) · [clinic/](../clinic/README.md) — tarball install tree vs microSD image.

## Full index

| ADR | Title | Status | Topic |
|-----|-------|--------|-------|
| [0001](./0001-modular-monolith-architecture.md) | Modular monolith | Accepted | Backend structure |
| [0002](./0002-authentication-system.md) | Authentication system | Accepted | Auth overview |
| [0003](./0003-react-frontend-framework.md) | React frontend | Accepted | UI stack |
| [0004](./0004-golang-backend.md) | Go backend | Accepted | API server |
| [0005](./0005-use-vite-for-frontend.md) | Vite for frontend | Accepted | Build tooling |
| [0006](./0006-mantine-ui-library.md) | Mantine UI | Accepted | Components |
| [0007](./0007-react-query-state-management.md) | React Query | Accepted | Client state |
| [0008](./0008-testing-strategy.md) | Testing strategy | Accepted | QA |
| [0009](./0009-rest-api-design.md) | REST API design | Accepted | HTTP API |
| [0010](./0010-local-airgap-deployment.md) | Local air-gapped deployment | Accepted | Offline / Pi |
| [0011](./0011-employee-badge-authentication.md) | Employee badge authentication | Accepted | Badge + PIN |
| [0012](./0012-backup-and-recovery-ops.md) | Backup and recovery ops | Accepted | USB backup |
| [0013](./0013-database-migrations.md) | Database migrations | Accepted | Schema changes |
| [0014](./0014-private-lan-only-deployment.md) | Private LAN-only deployment | Accepted | Product scope |
| [0015](./0015-clinic-database-engine.md) | Clinic database engine | **Accepted — testing** | SQLite only |
| [0016](./0016-clinic-runtime-packaging.md) | Clinic runtime packaging | Accepted | `clinic/` vs `scripts/` split |
| [0017](./0017-raspberry-pi-hardening-and-lan-tls.md) | Pi hardening and LAN TLS | Accepted | Firewall, local CA, HTTPS |

## Format (keep it GitHub-readable)

```markdown
# ADR-NNNN: Title

## Status
Proposed | Accepted | Deprecated | Superseded by ADR-XXXX

## Context
What forces the decision? (constraints, clinic reality)

## Decision
What we chose.

## Alternatives considered
What we rejected and why (short table is fine).

## Consequences
What gets easier / harder.
```

## Naming

`NNNN-short-title-with-hyphens.md` — sequential numbers, never reuse.

## Mono-repo note

Product guides stay under [`docs/`](../docs/). Device ops under [`devices/`](../devices/). **Decisions that bind the whole repo live here** so GitHub’s root tree shows `adr/` next to `backend/`, `frontend/`, and `README.md`.

Agent/Cursor habit: see [`AGENTS.md`](../AGENTS.md) and `.cursor/rules/monorepo-adr-docs.mdc` — update ADRs when architecture changes; keep root README links handy for GitHub navigation.
