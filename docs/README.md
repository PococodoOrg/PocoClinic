# Documentation

Human-facing docs for PocoClinic. Architecture **decisions** live in root [`adr/`](../adr/README.md) (not here). Device/server ops code lives in [`devices/`](../devices/README.md).

← Back to [repo README](../README.md)

## Start here by role

| Role | Go to |
|------|--------|
| Evaluating the product | [guide/](./guide/README.md) → [evaluating-pococlinic.md](./guide/evaluating-pococlinic.md) |
| Clinic administrator | [guide/for-administrators/](./guide/for-administrators/README.md) · [ops/](./ops/README.md) · [deploy/](./deploy/README.md) |
| Clinic staff (workflows) | [guide/for-staff/](./guide/for-staff/README.md) |
| Contributor / developer | [FEATURES.md](./FEATURES.md) · [adr/](../adr/README.md) · [clinic/](../clinic/README.md) · [AGENTS.md](../AGENTS.md) |

## What’s in this folder

| Path | Purpose |
|------|---------|
| [VISION.md](./VISION.md) | Product north star — who it serves, what we refuse |
| [FEATURES.md](./FEATURES.md) | Feature status / roadmap checklist |
| [NETWORK-AND-SECURITY.md](./NETWORK-AND-SECURITY.md) | LAN assumptions + defensive security |
| [guide/](./guide/README.md) | Clinic & staff guides (also feed in-app / static help) |
| [ops/](./ops/README.md) | Runbooks, checklists, tools catalog, **printed binder** |
| [ops/binder/](./ops/binder/README.md) | Clinic-wide paper binder (A emergency · B periodic · C help) |
| [deploy/](./deploy/README.md) | Release tarball vs microSD image · [DEPLOYMENT-BOUNDARY](./deploy/DEPLOYMENT-BOUNDARY.md) |
| [IN-APP-HELP.md](./IN-APP-HELP.md) | How in-app help content is maintained |
| [adr/](./adr/README.md) | **Stub only** → real ADRs are at [`/adr`](../adr/README.md) |

## Related (outside `docs/`)

| Path | Purpose |
|------|---------|
| [adr/](../adr/README.md) | Architecture Decision Records |
| [devices/](../devices/README.md) | Raspberry Pi (and future) server boxes |
| [binder-printer/](../binder-printer/README.md) | Print binders from markdown |
| [ops-helper/](../ops-helper/) | Localhost backup/restore UI (source; binary ships in tarball) |
| [clinic/](../clinic/README.md) | Post-deploy server install tree + workstation tools |
| [scripts/](../scripts/README.md) | Build, CI, and release **build** scripts (not clinic runtime) |
