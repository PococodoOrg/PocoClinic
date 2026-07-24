# I1 — Install prerequisites

**Installer binder · Section I1**

Complete before network or server work.

---

## Clinic confirms

| Item | Done | Notes |
|------|------|-------|
| Private LAN only — **no** public internet requirement for daily care | ☐ | |
| Clinic-owned staff devices (tablets + desktops) on same LAN as server | ☐ | |
| Locked space for server (cabinet / closet) | ☐ | |
| Two+ labeled USB drives for backup rotation | ☐ | |
| Safe or locked drawer for vault sheet | ☐ | |
| Primary admin available for badge print + training window | ☐ | |

---

## Installer brings

| Item | Done |
|------|------|
| Release tarball (`pococlinic-*-linux-arm64.tar.gz`) or pre-flashed image + tarball | ☐ |
| Ethernet cable (recommended for Pi server) | ☐ |
| USB keyboard/monitor **or** SSH access for first boot | ☐ |
| Laptop on clinic LAN for testing | ☐ |
| Printed **installer** + **site operations** binder packs | ☐ |
| Blank [safe vault](../../ops/safe-credentials-vault.md) template | ☐ |

---

## Server hardware (minimum)

| Role | Guidance |
|------|----------|
| **Pi server** | Pi 4/5, **8 GB RAM** preferred; USB SSD for database strongly preferred |
| **Small PC** | SSD, 8 GB+ RAM, wired Ethernet |
| **Not supported** | Pi Zero / 3 for production; consumer Wi‑Fi-only without Ethernet fallback |

Details: [Pi hardware](../../../devices/raspberry-pi/hardware.md)

---

## Software versions (record on cover)

| Item | Version |
|------|---------|
| PocoClinic release tarball | ____________ |
| OS image (if Pi) | ____________ |
| OpenSSL / Caddy (if used) | ____________ |

---

## Out of scope for PocoClinic installer

- Public cloud hosting or DNS
- Staff personal phones as official clients
- Patient portal over the internet
- Port forwarding EMR ports to WAN

---

**Next:** Section **I2** — Network & Wi‑Fi
