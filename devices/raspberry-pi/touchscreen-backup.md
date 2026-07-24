# Raspberry Pi — touchscreen backup kiosk

Use this guide when the clinic server is a Raspberry Pi with an official (or compatible) touchscreen shield.

> **Scope:** This is a **local backup/restore kiosk on the server** — not a mobile EMR app. Clinical charting is done from **staff browsers on clinic workstations** on the LAN. There is **no mobile product** and **no internet** requirement. See [NETWORK-AND-SECURITY.md](../../docs/NETWORK-AND-SECURITY.md).

The UI is **backup and restore only** — large buttons, bottom navigation, no keyboard required except typing `RESTORE` during recovery.

## What you get

- **Status** — backup health at a glance
- **Backup** — step-by-step wizard (USB → create file → copy to USB → done)
- **Restore** — guided recovery with tap-to-select backup + type `RESTORE`
- Runs at **http://127.0.0.1:9090** (**localhost only** on the server — never exposed on the clinic LAN for remote access)

## One-time setup on the Pi

### 1. Build the Pi touch UI (on dev PC or on the Pi)

```bash
# From project root
./build-ops-helper-pi.bat   # Windows
# or on Linux/Pi:
cd ops-helper && npm install && npm run build:pi
cp -r dist/* ../backend/cmd/ops-helper/static/
```

The Pi build sets `VITE_PI_TOUCH=true` so the touch layout is always used.

### 2. Configure environment

Production values live in **`/etc/pococlinic/env`** (created from [clinic/server/env.template](../../clinic/server/env.template) on install). At minimum:

```bash
DATABASE_URL=/var/lib/pococlinic/pococlinic.db
BACKUP_DIR=/var/lib/pococlinic/backups
OPS_HELPER_PORT=9090
```

### 3. Start the helper service

**Development (Windows checkout):**

```powershell
run-ops-helper-pi.bat
```

**Production:** use the shipped systemd unit (do not hand-edit paths unless you know why):

```bash
sudo cp /opt/pococlinic/systemd/pococlinic-ops-helper.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now pococlinic-ops-helper
```

The unit binds **`127.0.0.1:9090` only** and loads `/etc/pococlinic/env`. Source: [`scripts/deploy/pococlinic-ops-helper.service`](../../scripts/deploy/pococlinic-ops-helper.service).

### 4. Full-screen kiosk (recommended for touchscreen shield)

```bash
chmod +x scripts/pi/start-backup-kiosk.sh
./scripts/pi/start-backup-kiosk.sh
```

This opens Chromium edge-to-edge on the Pi display. The URL includes `?pi=1` to lock touch mode.

## Autostart on boot

Use the **shipped** unit from the release tarball — do not copy the outdated inline example below from older docs.

```bash
sudo cp /opt/pococlinic/systemd/pococlinic-ops-helper.service /etc/systemd/system/
sudo systemctl enable --now pococlinic-ops-helper
```

Reference copy in the repo: [`scripts/deploy/pococlinic-ops-helper.service`](../../scripts/deploy/pococlinic-ops-helper.service) (`EnvironmentFile=/etc/pococlinic/env`, `ExecStart=/opt/pococlinic/ops-helper`, `WorkingDirectory=/opt/pococlinic`).

For kiosk autologin + Chromium, configure your Pi OS desktop session to run `scripts/pi/start-backup-kiosk.sh` after login (Desktop Autostart or `.config/autostart/`).

## Display sizes

Tested layouts target common Pi official displays:

| Display | Resolution | Notes |
|---------|------------|--------|
| 7" official | 800 × 480 | Primary target |
| 5" / 3.5" shields | 480 × 320+ | Vertical stepper, bottom nav |

Touch targets are at least **56px** tall; backup dates use full-width tap cards.

## Daily workflow (touch)

1. Tap **Backup** in the bottom bar
2. Tap **USB is plugged in**
3. Tap **Create backup now**
4. Copy file to USB when prompted, tap **USB copy finished**
5. Store USB in locked drawer

## Restore workflow (touch)

1. Ensure all staff signed out
2. Tap **Restore** → **I understand**
3. Tap the backup date to select it → **Use this backup**
4. Type `RESTORE` on the on-screen keyboard → **Restore clinic data**

## Troubleshooting

| Issue | Fix |
|-------|-----|
| Blank screen | Check `systemctl status pococlinic-ops-helper`; rebuild Pi UI with `build-ops-helper-pi.bat` if `static/` missing |
| Desktop layout on Pi | Rebuild with `npm run build:pi`, not regular `build` |
| Can't tap small links | Open with `?pi=1` or rebuild Pi bundle |
| Backup button errors | Verify `/etc/pococlinic/env` has `DATABASE_URL` and `BACKUP_DIR`; run health check in EMR admin |

## Related

- [Raspberry Pi index](./README.md)
- [Hardware](./hardware.md)
- [NETWORK-AND-SECURITY.md](../../docs/NETWORK-AND-SECURITY.md)
- [Administrator runbook](../../docs/ops/administrator-runbook.md)
- [Tools & scripts](../../docs/ops/tools-and-scripts.md)
- [clinic/server/](../../clinic/server/README.md) — install tree on deployed Pi
- [Deployment boundary](../../docs/deploy/DEPLOYMENT-BOUNDARY.md)
- [scripts/deploy/README.md](../../scripts/deploy/README.md)
- [ADR-0014](../../adr/0014-private-lan-only-deployment.md)
- [ADR-0012](../../adr/0012-backup-and-recovery-ops.md)
