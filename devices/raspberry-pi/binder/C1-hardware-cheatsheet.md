# C1 — Hardware cheat sheet

**Device binder · Raspberry Pi · Section C**

Fill clinic-specific blanks on the **binder cover**. Do not write passwords here.

## Role

| | |
|--|--|
| This device | Clinic **server** (API, DB, static UI) |
| Staff charting | Browsers on clinic PCs / tablets on the **same LAN** |
| Pi touchscreen | Optional **backup/restore kiosk only** |

## Typical kit

| Item | Clinic note |
|------|-------------|
| Board | Pi 5 or Pi 4, prefer **8 GB** |
| Boot media | Endurance microSD and/or NVMe |
| Data / backups | Prefer USB SSD for DB + `BACKUP_DIR` |
| Network | Wired Ethernet to clinic switch |
| PSU | Official / quality USB-C |

Full detail: [hardware.md](../hardware.md).

## Golden rules

1. No public internet port-forward to this Pi  
2. Backup USB sticks are not the live database disk  
3. Vault secrets live in the **safe**, not taped to the Pi  
4. Lost SD → reinstall + restore USB + vault keys  

**LAN address (same as cover):** `http://________________________`
