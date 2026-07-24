# Clinic network requirements

PocoClinic is designed for a **private, isolated local area network (LAN)**. Daily clinical work does not require internet access.

## Minimum hardware

| Role | Suggestion |
|------|------------|
| **Server** | Raspberry Pi 4/5 (4–8 GB RAM) or small PC with SSD — see [Devices · Raspberry Pi](../../devices/raspberry-pi/hardware.md) |
| **Database** | SQLite on the same server (or dedicated small PC) |
| **Staff clients** | Clinic-owned **tablets** (primary) and desktops/laptops on the LAN |
| **Backup media** | At least two labeled USB drives (rotation e.g. Mon/Wed/Fri) |

## Network layout

- Server and workstations on the **same private subnet**
- **No port forwarding** of PocoClinic to the public internet
- **No mobile apps** — browsers on clinic PCs only
- Backup Helper at `http://127.0.0.1:9090` — **localhost on the server only**, not exposed on Wi‑Fi

## What needs to work on the LAN

| Service | Port (default) | Who connects |
|---------|----------------|--------------|
| PocoClinic web UI + API | 443 (HTTPS, recommended) or 8080 (HTTP migrate-only) | Staff browsers on the LAN |
| SQLite | _(file on disk — no port)_ | Server process only |
| Backup Helper | 9090 | Server browser only (localhost) |

## Internet

| Need internet? | When |
|----------------|------|
| **No** | Sign-in, patient charts, forms, notes, documents, daily USB backup |
| **Optional** | Initial software install, security patches during maintenance windows |

Clinics that require full egress blocking should use the **release tarball** built with `scripts/build-release.mjs` (copies [`clinic/server/`](../../clinic/server/README.md)) — see [DEPLOYMENT-BOUNDARY](../deploy/DEPLOYMENT-BOUNDARY.md) and [ADR-0010](../../adr/0010-local-airgap-deployment.md).

**HTTPS on the LAN:** [Installer binder I2](../../binders/installer/I2-network-and-wifi.md) (print checklist) · [Pi TLS guide](../../devices/raspberry-pi/tls-lan.md) · [Pi hardening](../../devices/raspberry-pi/hardening.md) · [ADR-0017](../../adr/0017-raspberry-pi-hardening-and-lan-tls.md)

## Security expectations

Even on a private LAN:

- Treat badges like physical keys; reissue when lost
- Staff choose private PINs; do not share
- Review audit log for unusual access
- Keep USB backups in a locked location
- Run quarterly restore drills

Full model: [NETWORK-AND-SECURITY.md](../NETWORK-AND-SECURITY.md) · [ADR-0014](../../adr/0014-private-lan-only-deployment.md)

## Related

- [Binders hub](../../binders/README.md) · [Installer I2 — Network & Wi‑Fi](../../binders/installer/I2-network-and-wifi.md)
- [First-time setup](./for-administrators/first-time-setup.md)
- [Deploy hub](../deploy/README.md) · [clinic/server/](../../clinic/server/README.md)
- [In-app help: How it works](../../frontend/src/help/content.ts) (article `how-it-works`)
