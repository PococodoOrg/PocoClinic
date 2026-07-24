# Raspberry Pi

PocoClinic’s primary **clinic server** target: a Raspberry Pi on the private LAN serving the EMR to staff browsers.

← [Devices](../README.md) · [Repo README](../../README.md) · [Docs hub](../../docs/README.md)

**This folder** holds Pi-specific documentation and (over time) Pi-only ops code — kiosk launchers, hardware notes.  
Server install/cron/env wrappers ship in [`clinic/server/`](../../clinic/server/README.md).  
The application release contract (tarball vs OS image) stays in [`docs/deploy/DEPLOYMENT-BOUNDARY.md`](../../docs/deploy/DEPLOYMENT-BOUNDARY.md).

## In this folder

| Path | Contents |
|------|----------|
| [binder/](./binder/README.md) | **Print-me** Pi device binder (A emergency · B periodic · C help) |
| [hardening.md](./hardening.md) | Firewall, SSH, permissions — run `scripts/harden-pi.sh` |
| [tls-lan.md](./tls-lan.md) | **HTTPS on the LAN** — local CA, Caddy or Go TLS |
| [touchscreen-backup.md](./touchscreen-backup.md) | Full touchscreen kiosk guide |
| [hardware.md](./hardware.md) | Recommended models, RAM, storage, display |

**Print the Pi binder:** [`clinic/workstation/print-binder.bat`](../../clinic/workstation/print-binder.bat) (or repo-root `print-binder.bat`) → choose **Raspberry Pi**, or follow [binder/README.md](./binder/README.md).

**Install on Pi:** extract release tarball and run [`clinic/server/install.sh`](../../clinic/server/install.sh) — see [deployment boundary](../../docs/deploy/DEPLOYMENT-BOUNDARY.md).

## Quick facts

| Item | Guidance |
|------|----------|
| Role | Clinic **server** only — not a staff charting tablet product |
| Clients | Staff use browsers on clinic PCs / tablets on the **same LAN** |
| Ops UI | Optional Pi touchscreen → backup/restore at `http://127.0.0.1:9090` only |
| Release | Linux **ARM64** tarball from `build-release.bat` / `scripts/build-release.mjs` |
| Server ops (on Pi) | [`clinic/server/`](../../clinic/server/README.md) → `/opt/pococlinic/bin/*` |

| Path | Purpose |
|------|---------|
| [`clinic/server/`](../../clinic/server/README.md) | Install script, env template, cron, `bin/` wrappers (in tarball) |
| [`scripts/pi/start-backup-kiosk.sh`](../../scripts/pi/start-backup-kiosk.sh) | Chromium full-screen kiosk (may move here later) |
| `build-ops-helper-pi.bat` / `ops-helper` `build:pi` | Touch layout for backup helper |
| [`clinic/dev-windows/run-ops-helper-pi.bat`](../../clinic/dev-windows/run-ops-helper-pi.bat) | Dev launch for Pi UI layout |
| `scripts/build-release.mjs` | ARM64 production tarball (copies `clinic/server/`) |

## Related docs

- [Devices hub](../README.md)
- [Ops hub](../../docs/ops/README.md)
- [ADR-0010 — Local air-gapped deployment](../../adr/0010-local-airgap-deployment.md)
- [ADR-0014 — Private LAN-only](../../adr/0014-private-lan-only-deployment.md)
- [ADR-0015 — Clinic database engine](../../adr/0015-clinic-database-engine.md)
