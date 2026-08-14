# C3 — Where things live on the Pi

**Device binder · Raspberry Pi · Section C** · Paths are typical — adjust if your image differs.

## Common paths (fill clinic actuals)

| What | Typical path | This clinic |
|------|--------------|-------------|
| App install | `/opt/pococlinic` | ________________ |
| Ops wrappers (cron/SSH) | `/opt/pococlinic/bin/backup`, `bin/migrate`, … | ________________ |
| Env / secrets file | `/etc/pococlinic/env` | ________________ |
| Backups (`BACKUP_DIR`) | `/var/lib/pococlinic/backups` | ________________ |
| Database file | `/var/lib/pococlinic/pococlinic.db` (typical) | ________________ |
| First-boot admin credentials | `/var/lib/pococlinic/bootstrap-admin-once.txt` (delete after setup) | ________________ |
| Ops helper | `ops-helper` on port **9090** localhost | |

## Services (names may vary)

| Service | Purpose | How you start it (clinic note) |
|---------|---------|--------------------------------|
| Database | SQLite / SQLite | ________________ |
| PocoClinic API + UI | Main EMR | ________________ |
| Ops helper | Backup kiosk | ________________ |

## Useful commands (admin with keyboard)

```bash
# Prefer bin/ wrappers (load /etc/pococlinic/env)
sudo /opt/pococlinic/bin/backup
sudo /opt/pococlinic/bin/migrate

# Service status
sudo systemctl status pococlinic
sudo systemctl status pococlinic-ops-helper
curl -s http://127.0.0.1:9090/ | head
```

## Build reminders (dev / upgrade)

- **First install walkthrough:** [`devices/raspberry-pi/setup.md`](../setup.md)
- Pi touch UI: `build-ops-helper-pi.bat` or `ops-helper` → `npm run build:pi`  
- Release tarball: `build-release.bat` (linux-arm64) — includes [`clinic/server/`](../../../clinic/server/README.md)  
- Contract: [`docs/deploy/DEPLOYMENT-BOUNDARY.md`](../../../docs/deploy/DEPLOYMENT-BOUNDARY.md)
