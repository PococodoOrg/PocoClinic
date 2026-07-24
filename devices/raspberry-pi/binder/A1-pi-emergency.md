# A1 — Pi will not start / no network

**Device binder · Raspberry Pi · Section A**

## Priority

1. Patient care continues on **paper** (clinic ops binder Section A)  
2. Fix or replace the Pi  
3. Restore from USB if the disk is lost (need **safe vault** encryption key after rebuild)

## Pi has no power / no lights

| Step | Action |
|------|--------|
| 1 | Confirm wall power and official / quality USB-C PSU |
| 2 | Try a known-good cable |
| 3 | If still dead → spare PSU or spare Pi; restore from latest USB backup |

## Pi powers on but no display / no boot

| Step | Action |
|------|--------|
| 1 | Reseat microSD (or NVMe); try reading the card on another computer |
| 2 | If card unreadable → flash new image / reinstall from release tarball (see deployment docs) |
| 3 | Restore database from latest **verified** USB backup |
| 4 | Re-enter `.env` secrets from the **safe vault** sheet |
| 5 | Sign [Pi incident log](./A3-pi-incident-log.md) |

## Pi boots but staff cannot open PocoClinic

| Step | Action |
|------|--------|
| 1 | Confirm Ethernet link light; cable to clinic switch |
| 2 | From a workstation, ping the Pi LAN IP written on this binder cover |
| 3 | Confirm PocoClinic / reverse-proxy service is running (Section C3) |
| 4 | Confirm workstation is on the **same LAN** (not guest Wi‑Fi) |
| 5 | If database down → start SQLite / DB service before blaming the web UI |

## Do not

- Port-forward the Pi to the internet “just to fix it remotely”  
- Email vault secrets to get help  
- Reuse an old SD image without restoring a known-good backup  

**Outage started:** ____________ **Ended:** ____________ **Initials:** ____
