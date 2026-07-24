# Raspberry Pi — hardware

Recommendations for the clinic **server** box.

Staff charting still happens on **clinic workstations** (desktop, laptop, or tablet browsers on the LAN). The Pi is the always-on server — optional touchscreen is for **backup/restore ops only**.

## Recommended models

| Model | Notes |
|-------|--------|
| **Raspberry Pi 5** (4–8 GB) | Preferred for new installs |
| **Raspberry Pi 4** (4–8 GB) | Acceptable; prefer 8 GB if database + app share one board |
| Pi 3 / Zero | **Not recommended** for production clinic load |

Exact SQLite performance on 4 GB boards varies — prefer **8 GB** or move the database data directory to USB SSD / NVMe when possible (see storage below).

## Storage

| Use | Recommendation |
|-----|----------------|
| OS + binaries | microSD (A2 / endurance-rated) or NVMe (Pi 5) |
| Database + backups + documents | Prefer **USB SSD** or dedicated durable volume — microSD wears under heavy writes |
| USB backup rotation | Separate labeled sticks (clinic ops); not the same drive as live DB |

Path contract (`DATABASE_URL` SQLite file, `BACKUP_DIR`, legacy `DOCUMENTS_DIR`) is defined in [`docs/deploy/DEPLOYMENT-BOUNDARY.md`](../../docs/deploy/DEPLOYMENT-BOUNDARY.md) and [`clinic/server/env.template`](../../clinic/server/env.template).

## Network

- Wired Ethernet to the clinic LAN when possible (more reliable than Wi‑Fi for the server)
- Static IP or reserved DHCP lease — record on the safe vault / ops contacts sheet
- **No** port forwarding to the public internet

See [`docs/guide/clinic-network.md`](../../docs/guide/clinic-network.md).

## Display (optional)

| Display | Use |
|---------|-----|
| Official 7" (800×480) | Primary target for [touchscreen backup kiosk](./touchscreen-backup.md) |
| Smaller official / compatible shields | Supported with large touch targets; test before clinic go-live |
| No display | Fine — use Backup Helper from a keyboard/monitor temporarily, or Admin dashboard from a workstation |

## Power and placement

- Locked room or cabinet; limited key holders ([Physical security](../../docs/ops/physical-security-binder.md))
- Quality USB-C / official PSU — brownouts corrupt SD cards
- Label the Pi with clinic hostname / LAN IP

## Related

- [Raspberry Pi index](./README.md)
- [Touchscreen backup](./touchscreen-backup.md)
- [Deployment boundary](../../docs/deploy/DEPLOYMENT-BOUNDARY.md)
