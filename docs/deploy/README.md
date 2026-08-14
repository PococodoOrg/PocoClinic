# Deploy documentation

How PocoClinic is **packaged**, **installed**, and **upgraded** on the clinic server — separate from day-to-day admin runbooks in [`docs/ops/`](../ops/README.md).

← [Docs hub](../README.md) · [Repo README](../../README.md)

## Start here

| Document | Audience | Purpose |
|----------|----------|---------|
| **`setup.bat` / `scripts/setup.mjs`** | First-time installers | Interactive wizard — pick Raspberry Pi, local dev, … |
| **[Pi beginner setup](../../devices/raspberry-pi/setup.md)** | First-time installers | Flash SD → install → first admin login (step-by-step written) |
| [Installer binder](../binders/installer/README.md) | Deploy techs | Printable I1–I5 checklists |
| [DEPLOYMENT-BOUNDARY.md](./DEPLOYMENT-BOUNDARY.md) | Image builders, contributors | Tarball vs microSD image responsibilities |
| [clinic/server/](../../clinic/server/README.md) | Installers | Install script, env template, cron, `bin/` wrappers |
| [ADR-0016](../../adr/0016-clinic-runtime-packaging.md) | Contributors | Why `clinic/` is separate from `scripts/` |
| [scripts/deploy/README.md](../../scripts/deploy/README.md) | Sysadmins | systemd units, manual install notes |
| [scripts/build-release.mjs](../../scripts/build-release.mjs) | Contributors | Builds `dist/pococlinic-*-linux-*.tar.gz` |

## Quick reference

**First time?** From the repo root run **`setup.bat`** (Windows) or **`./setup.sh`**, then choose **Raspberry Pi**. Written fallback: [devices/raspberry-pi/setup.md](../../devices/raspberry-pi/setup.md).

**Build** (developer machine):

```bash
node scripts/build-release.mjs
# or: build-release.bat
```

**Install** (clinic server — after OS is ready):

```bash
tar -xzf pococlinic-1.0.0-linux-arm64.tar.gz
cd pococlinic-1.0.0-linux-arm64
sudo ./install.sh
sudo nano /etc/pococlinic/env   # JWT_*, DOCUMENT_ENCRYPTION_KEY, ALLOWED_ORIGIN
sudo /opt/pococlinic/bin/migrate
sudo systemctl enable --now pococlinic pococlinic-ops-helper
```

**Not in the tarball:** secrets, existing PHI, Windows dev batch files under `clinic/dev-windows/`.

See [tools & scripts catalog](../ops/tools-and-scripts.md) for the full dev vs deployed command index.
