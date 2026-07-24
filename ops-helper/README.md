# Ops helper

Localhost-only **backup / restore UI** for clinic admins. Not the EMR charting app.

← [Repo README](../README.md) · Procedures: [docs/ops](../docs/ops/README.md) · [ADR-0012](../adr/0012-backup-and-recovery-ops.md)

## Quick start (Windows)

From the repo root (forwards to [`clinic/dev-windows/`](../clinic/dev-windows/README.md)):

```bat
run-ops-helper.bat
```

Pi / ARM builds: `build-ops-helper-pi.bat`, `run-ops-helper-pi.bat` → `clinic/dev-windows/run-ops-helper-pi.bat`.

Binds to **localhost only** by design (see product constraints and ADR-0012). On a Pi touchscreen kiosk, see [devices/raspberry-pi/touchscreen-backup.md](../devices/raspberry-pi/touchscreen-backup.md).
