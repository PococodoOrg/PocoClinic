# C2 — Touchscreen backup how-to

**Device binder · Raspberry Pi · Section C**  
Full guide: [touchscreen-backup.md](../touchscreen-backup.md).

## Open the kiosk

- URL on the Pi only: `http://127.0.0.1:9090/?pi=1`  
- Or autostart via `scripts/pi/start-backup-kiosk.sh`

## Backup (touch)

1. Tap **Backup**  
2. Tap **USB is plugged in**  
3. Tap **Create backup now**  
4. Copy to USB when prompted → **USB copy finished**  
5. Lock the USB away  

## Restore (touch) — disasters only

1. Staff signed out except one admin  
2. **Restore** → confirm → pick backup → type `RESTORE`  
3. Verify login + patient count + one document  
4. Fresh backup → sign [A3](./A3-pi-incident-log.md)  

Without `DOCUMENT_ENCRYPTION_KEY` from the safe vault, restored documents may not open.
