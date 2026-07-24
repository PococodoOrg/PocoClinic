# Scripts

Build, CI, and release **build** automation for PocoClinic — **not** post-deploy clinic operations (those live in [`clinic/`](../clinic/README.md)).

← [Repo README](../README.md) · Clinic runtime: [`clinic/server/`](../clinic/server/README.md) · Full catalog: [docs/ops/tools-and-scripts.md](../docs/ops/tools-and-scripts.md)

| Path | Purpose | Full docs |
|------|---------|-----------|
| `build-release.mjs` | Linux ARM64 release tarball for Pi | [Tools catalog](../docs/ops/tools-and-scripts.md#shell-scripts-scripts) |
| `build-help-site.mjs` | Static HTML from `docs/guide/` | [Tools catalog](../docs/ops/tools-and-scripts.md#documentation-generators) |
| `test-all.sh` | Backend + frontend + loadtest pytest (local CI parity) | — |
| `check-no-secrets.sh` | Block commits of `.env`, DB, backup, or key material | Run before `git commit` |
| `check-no-secrets.ps1` | Same check for Windows PowerShell | Run before `git commit` |
| `deploy/` | systemd units + production install notes | [deploy/README.md](./deploy/README.md) |
| `release/` | _(moved)_ → [`clinic/server/`](../clinic/server/README.md) | Pointer only |
| `pi/start-backup-kiosk.sh` | Chromium kiosk for Pi touchscreen backup UI | [touchscreen-backup.md](../devices/raspberry-pi/touchscreen-backup.md) |

**Deployment:** [`clinic/server/`](../clinic/server/README.md) holds install/cron/env files copied into the release tarball. Build scripts live here. See [Deployment boundary](../docs/deploy/DEPLOYMENT-BOUNDARY.md).

When adding a script, document it in [docs/ops/tools-and-scripts.md](../docs/ops/tools-and-scripts.md).
