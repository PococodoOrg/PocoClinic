# Printable binders

PocoClinic uses **paper binders** for work that happens when the screen is wrong — outages, scheduled checks, or install day. Each binder has a clear audience. Do not merge them into one thick manual.

← [Docs hub](../README.md) · Print utility: [binder-printer](../../binder-printer/README.md)

---

## Which binder do I need?

| Binder | Who keeps it | When to open it |
|--------|--------------|-----------------|
| [**Installer**](installer/README.md) | Deploy tech / integrator (until handoff) | New clinic, server replace, major network change |
| [**Site operations**](../ops/binder/README.md) | Clinic admin desk (locked cabinet) | Emergency, weekly/monthly checks, printed how-to |
| [**Raspberry Pi device**](../../devices/raspberry-pi/binder/README.md) | With the Pi server (cabinet tab) | Pi hardware / kiosk problems |
| [**Safe vault**](../ops/safe-credentials-vault.md) | **Physical safe only** | Never open for routine work — secrets sheet |

**Rule:** Installers finish with the **handoff checklist** (installer Section I5), then leave the **site operations binder** for clinic staff.

---

## Print all packs

From the repo root:

```bat
print-binder.bat
```

Choose **Installer**, **Site operations**, **Raspberry Pi**, and/or **Safe vault (blank)** in [binder-printer](../../binder-printer/README.md).

---

## What lives outside binders

| Topic | Where |
|-------|--------|
| Searchable help when EMR is up | In-app **Help** (`/help`) |
| Long-form guides | [docs/guide/](../guide/README.md) |
| **First-time Pi install** | [devices/raspberry-pi/setup.md](../../devices/raspberry-pi/setup.md) |
| Architecture decisions | [adr/](../../adr/README.md) |
| Release tarball install tree | [clinic/server/](../../clinic/server/README.md) |

---

## Related

- [Ops hub](../ops/README.md)
- [First-time setup](../guide/for-administrators/first-time-setup.md)
- [Clinic network requirements](../guide/clinic-network.md)
- [ADR-0017 — Pi hardening & LAN TLS](../../adr/0017-raspberry-pi-hardening-and-lan-tls.md)
