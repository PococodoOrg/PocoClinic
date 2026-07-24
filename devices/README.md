# Devices

Device-specific **ops code and documentation** for running PocoClinic as the clinic server.

This is a **top-level** folder (not under `docs/`). Put Raspberry Pi scripts, unit files, and install helpers next to the docs that describe them. Shared product docs (vision, ADRs, ops binder) stay in [`docs/`](../docs/).

Staff charting still uses **clinic browsers** on the LAN. Device folders cover the **server box** (and optional on-server ops displays), not mobile apps.

← [Repo README](../README.md) · [Docs hub](../docs/README.md)

## Targets

| Folder | Role | Status |
|--------|------|--------|
| [raspberry-pi/](./raspberry-pi/README.md) | Primary clinic server (ARM64) + optional touchscreen backup kiosk | Active |
| *(future)* | Small x86 PC / NUC | Not started |

**Print device binders:** [`binder-printer/`](../binder-printer/README.md) · [`clinic/workstation/`](../clinic/workstation/README.md) (not part of the EMR).

**Clinic server install tree (tarball):** [`clinic/server/`](../clinic/server/README.md) — migrate, backup, cron ([ADR-0016](../adr/0016-clinic-runtime-packaging.md)).

## Shared product docs (not duplicated here)

| Topic | Where |
|-------|--------|
| microSD / release contract | [`docs/deploy/DEPLOYMENT-BOUNDARY.md`](../docs/deploy/DEPLOYMENT-BOUNDARY.md) |
| LAN layout | [`docs/guide/clinic-network.md`](../docs/guide/clinic-network.md) |
| Security model | [`docs/NETWORK-AND-SECURITY.md`](../docs/NETWORK-AND-SECURITY.md) |
| Printed ops binder (clinic-wide) | [`docs/ops/binder/`](../docs/ops/binder/README.md) |
| Cross-platform scripts catalog | [`docs/ops/tools-and-scripts.md`](../docs/ops/tools-and-scripts.md) |

When adding a device type: create `devices/<name>/` with a `README.md`, then add device-local code under that tree.
