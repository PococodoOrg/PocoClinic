# Raspberry Pi — device binder

**Print this packet for the clinic Pi server.**  
Store it with the Pi (cabinet) or behind the **site operations binder** as a **device tab**.

This binder is **Pi-only**. Clinic-wide staff/admin how-to stays in the [**site operations binder**](../../../docs/ops/binder/README.md).  
Install-day network/TLS steps: [**installer binder**](../../../docs/binders/installer/README.md).  
Secrets stay on the **safe vault** sheet — not in this packet.

| Clinic | ________________________________ |
|--------|----------------------------------|
| Pi hostname / LAN IP | ________________________________ |
| Binder assembled by | ________________________________ |
| Date | ____________ |

---

## When to open this binder

| Use case | Open when… | Section |
|----------|------------|---------|
| **1. Emergency** | The **Pi** will not boot, has no network, or the backup kiosk is dead | **A** |
| **2. Weekly / monthly** | Scheduled Pi checks (storage, kiosk, cables) | **B** |
| **3. Device help** | Printed Pi how-to (same on every clinic’s Pi docs — fill IP once on the cover) | **C** |

Day-to-day charting is not done on the Pi display. Use staff browsers on the LAN.

**Easiest print path:** run the [binder-printer](../../../binder-printer/README.md) utility → choose **Raspberry Pi**.

---

## Print order

### Section A — Emergency (Pi down)

| Page | File |
|------|------|
| A1 | [Pi will not start / no network](./A1-pi-emergency.md) |
| A2 | [Backup kiosk blank or frozen](./A2-kiosk-emergency.md) |
| A3 | [Pi incident log](./A3-pi-incident-log.md) |

### Section B — Weekly / monthly (Pi)

| Page | File |
|------|------|
| B1 | [Pi periodic checks](./B1-pi-periodic.md) |

### Section C — Device help

| Page | File |
|------|------|
| C1 | [Hardware cheat sheet](./C1-hardware-cheatsheet.md) |
| C2 | [Touchscreen backup how-to](./C2-touchscreen-howto.md) |
| C3 | [Where things live on the Pi](./C3-paths-and-services.md) |

### Print checklist

- [ ] Cover (this page) — fill hostname / IP
- [ ] Section A — A1, A2, A3
- [ ] Section B — B1
- [ ] Section C — C1, C2, C3
- [ ] **Site operations binder** printed separately for clinic desk

---

## Related

- [Raspberry Pi device folder](../README.md)
- [Hardware details](../hardware.md)
- [Touchscreen backup (full)](../touchscreen-backup.md)
- [Binders hub](../../../docs/binders/README.md)
- [Site operations binder](../../../docs/ops/binder/README.md)
- [Installer binder](../../../docs/binders/installer/README.md)
- [Binder printer utility](../../../binder-printer/README.md)
