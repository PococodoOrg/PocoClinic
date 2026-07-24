# A2 — Backup kiosk blank or frozen

**Device binder · Raspberry Pi · Section A**

The Pi touchscreen is **backup/restore only** at `http://127.0.0.1:9090`. It is not used for patient charts.

## Blank screen / browser error

| Step | Action |
|------|--------|
| 1 | Confirm ops-helper service is running (Section C3) |
| 2 | On the Pi (keyboard attached), open `http://127.0.0.1:9090/?pi=1` |
| 3 | If still blank → rebuild Pi touch UI (`build-ops-helper-pi` / `npm run build:pi`) and restart service |
| 4 | Re-run kiosk script: `scripts/pi/start-backup-kiosk.sh` |

## Touch works but Backup fails

| Step | Action |
|------|--------|
| 1 | Confirm `DATABASE_URL` and SQLite are up |
| 2 | Confirm `BACKUP_DIR` exists and is writable |
| 3 | Confirm USB is mounted if the wizard expects a copy target |
| 4 | Fall back: Admin dashboard **Backup now** from a workstation on the LAN |

## Frozen UI

1. Close Chromium / reboot the Pi (after clinic hours if possible)  
2. Take a USB backup from the Admin dashboard before deeper troubleshooting  
3. Log on [A3](./A3-pi-incident-log.md)

## Related

- Full guide: [touchscreen-backup.md](../touchscreen-backup.md)
- Printed short how-to: [C2](./C2-touchscreen-howto.md)
