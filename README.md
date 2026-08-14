# PocoClinic

**Note:** This project is in open development (including generative AI). We welcome testers and contributors.

PocoClinic is an open-source **Electronic Medical Records (EMR)** system for **small clinics on a private LAN**—especially non-profits with limited IT staff.

- Runs on a **Raspberry Pi or small PC** as the clinic server  
- Used from **clinic-owned browsers** (primarily tablets; also desktops) — **no native mobile apps**  
- Works with **no internet** for daily care  
- Backed up with **USB + printed runbooks** that non-technical admins can follow  

Full intent: [docs/VISION.md](./docs/VISION.md) · Security model: [docs/NETWORK-AND-SECURITY.md](./docs/NETWORK-AND-SECURITY.md)

---

## Choose a path

| I want to… | Start here |
|------------|------------|
| **Set up after cloning** (interactive) | **`setup.bat`** (Windows) or **`./setup.sh`** — pick Raspberry Pi, local dev, … |
| **Understand the product** (evaluate for a clinic) | [Documentation guide](./docs/guide/README.md) → [What is PocoClinic?](./docs/guide/evaluating-pococlinic.md) |
| **See why we chose X** (architecture) | [**adr/**](./adr/README.md) — especially [0014](./adr/0014-private-lan-only-deployment.md) (scope) and [0015](./adr/0015-clinic-database-engine.md) (database) |
| **Run / administer a clinic** | [For administrators](./docs/guide/for-administrators/README.md) · [Ops hub](./docs/ops/README.md) · [Binders](./docs/binders/README.md) |
| **Install on a Raspberry Pi** | **`setup.bat` → Raspberry Pi** · [written guide](./devices/raspberry-pi/setup.md) · [Installer binder](./docs/binders/installer/README.md) |
| **Print paper binders** | [binder-printer](./binder-printer/README.md) · [`clinic/workstation/`](./clinic/workstation/README.md) — not part of the EMR |
| **Packaged clinic server ops** | [`clinic/server/`](./clinic/server/README.md) — install, cron, env wrappers (shipped in tarball) |
| **Develop or contribute code** | [Getting started](#getting-started) · [AGENTS.md](./AGENTS.md) · [FEATURES.md](./docs/FEATURES.md) |
| **Browse all docs** | [docs/README.md](./docs/README.md) |

---

## Repository map

This is a **monorepo**. Top-level folders are meant to be obvious on GitHub:

| Path | What it is |
|------|------------|
| [`adr/`](./adr/README.md) | Architecture Decision Records (**why** we chose things) |
| [`docs/`](./docs/README.md) | Product vision, guides, ops binders, feature status |
| [`devices/`](./devices/README.md) | Device-specific server ops (e.g. Raspberry Pi) |
| [`frontend/`](./frontend/) | React EMR UI (clinic browsers) |
| [`backend/`](./backend/) | Go API server |
| [`ops-helper/`](./ops-helper/) | Localhost-only backup/restore UI |
| [`binder-printer/`](./binder-printer/README.md) | Standalone binder print utility |
| [`clinic/`](./clinic/README.md) | **Post-deploy ops** — server install tree (shipped in tarball), workstation tools |
| [`scripts/`](./scripts/README.md) | Build, CI, and deploy reference scripts |
| [`loadtest/`](./loadtest/README.md) | Local-only Locust load tests — run `loadtest/run.bat` |
| [`AGENTS.md`](./AGENTS.md) | Notes for coding agents / contributors |
| `llm-service/` | Placeholder only — local AI assistance is **planned**, not implemented ([FEATURES](./docs/FEATURES.md)) |

```
PocoClinic/
├── adr/                 # Decisions (GitHub-friendly)
├── docs/                # Product + ops documentation
├── devices/             # Pi / server-box ops
├── frontend/            # EMR UI
├── backend/             # EMR API
├── ops-helper/          # Backup helper (localhost)
├── binder-printer/      # Print binders (not the EMR)
├── clinic/              # Post-deploy ops (server tarball, workstation)
├── scripts/             # Build / CI (not shipped to clinics)
├── loadtest/            # Local Locust API smoke (not CI)
└── AGENTS.md
```

---

## Architecture at a glance

| Piece | Today | Notes |
|-------|--------|--------|
| Frontend | React + TypeScript + Mantine | Clinic LAN browsers; tablet-primary |
| Backend | Go modular monolith | Feature slices under `backend/internal/features/` |
| Database | **SQLite** (single file) | [ADR-0015](./adr/0015-clinic-database-engine.md) — testing in progress on clinic hardware |
| Auth | Badge (QR) + 4-digit PIN | [ADR-0011](./adr/0011-employee-badge-authentication.md) |
| Deploy | Private LAN, air-gap tarball | [ADR-0010](./adr/0010-local-airgap-deployment.md), [clinic/server/](./clinic/server/README.md) |
| Backup | USB tar.gz with SQLite file payload | [ADR-0012](./adr/0012-backup-and-recovery-ops.md) |

### Key decisions (start here)

| ADR | Why it matters |
|-----|----------------|
| [0014 — Private LAN-only](./adr/0014-private-lan-only-deployment.md) | Product scope: no cloud, no mobile apps |
| [0015 — Clinic database engine](./adr/0015-clinic-database-engine.md) | **Accepted (testing):** SQLite only |
| [0010 — Air-gapped deploy](./adr/0010-local-airgap-deployment.md) | Offline install / USB updates |
| [0012 — Backup & recovery](./adr/0012-backup-and-recovery-ops.md) | USB backups for non-technical admins |
| [0011 — Badge + PIN](./adr/0011-employee-badge-authentication.md) | Staff sign-in model |
| [0016 — Clinic runtime packaging](./adr/0016-clinic-runtime-packaging.md) | `clinic/` vs `scripts/` — what ships on the Pi |

Full index: [adr/README.md](./adr/README.md)

---

## Dev tooling vs clinic runtime

Contributors and clinic operators use **different script folders** on purpose ([ADR-0016](./adr/0016-clinic-runtime-packaging.md)):

| Folder | Who | Example |
|--------|-----|---------|
| [`scripts/`](./scripts/README.md) | Developers / CI | `build-release.mjs`, `test-all.sh`, secret scan |
| [`clinic/server/`](./clinic/server/README.md) | Clinic server after install | `install.sh`, `/opt/pococlinic/bin/backup` |
| [`clinic/dev-windows/`](./clinic/dev-windows/README.md) | Windows dev testing | `go run ./cmd/backup` via `backup.bat` forwarder |
| [`clinic/workstation/`](./clinic/workstation/README.md) | Admin PC (optional) | Print paper binders |

Repo-root `backup.bat`, `migrate.bat`, etc. forward to `clinic/dev-windows/` — they are **not** copied to production images.

Full catalog: [docs/ops/tools-and-scripts.md](./docs/ops/tools-and-scripts.md) · Deploy contract: [docs/deploy/DEPLOYMENT-BOUNDARY.md](./docs/deploy/DEPLOYMENT-BOUNDARY.md)

---

## Getting started

### One command after clone

```bat
setup.bat
```

Or on macOS/Linux: `chmod +x setup.sh && ./setup.sh`

The wizard asks what you want (**Raspberry Pi**, local development, …) and walks you through that path. Written Pi guide: [devices/raspberry-pi/setup.md](./devices/raspberry-pi/setup.md).

### Prerequisites

- Node.js 18+ (required for `setup.bat` / release builds)
- Go 1.25+ (see `backend/go.mod`) for building or local API
- No separate database server — SQLite file path via `DATABASE_URL` (see [ADR-0015](./adr/0015-clinic-database-engine.md))

### Run on Windows (quick)

```bat
run-all.bat
```

Installs deps, runs tests, starts frontend + backend when tests pass.

Or separately:

```bat
run-frontend.bat
run-backend.bat
```

### Manual

```bash
# Frontend
cd frontend && npm install && npm run dev

# Backend
cd backend && go mod tidy && go run ./cmd/main.go
```

Copy [`.env.example`](./.env.example) / [`backend/.env.example`](./backend/.env.example) as needed.

- **Local migrate / backup / restore:** `migrate.bat`, `backup.bat`, `restore.bat` (forward to [`clinic/dev-windows/`](./clinic/dev-windows/README.md))
- **Production Pi:** `setup.bat` → Raspberry Pi, or [`clinic/server/`](./clinic/server/README.md) after `build-release` — see [deployment boundary](./docs/deploy/DEPLOYMENT-BOUNDARY.md)
- **Full tool index:** [docs/ops/tools-and-scripts.md](./docs/ops/tools-and-scripts.md)

### Tests

```bat
test-all.bat
```

Or: `cd backend && go test ./...` · `cd frontend && npm test` · `python -m pytest loadtest`

Before committing: `scripts/check-no-secrets.ps1` (Windows) or `scripts/check-no-secrets.sh` (Unix).

CI: `.github/workflows/ci.yml`

---

## Contributing

1. Read [VISION.md](./docs/VISION.md) and [ADR-0014](./adr/0014-private-lan-only-deployment.md) so changes stay in product scope.  
2. Check [FEATURES.md](./docs/FEATURES.md) for current work.  
3. If you change architecture, **update an ADR** under [`adr/`](./adr/README.md) (see [AGENTS.md](./AGENTS.md)).  
4. Open an issue or PR on GitHub.

## License

MIT — see [LICENSE](./LICENSE).

## Security

Report security issues to **security@pococodo.com** (not the public issue tracker).
