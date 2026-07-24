# B6 — Network & Wi‑Fi at the site

**Site operations binder · Section B** · For clinic admins — **not** full router setup.

**Install / major network changes:** integrator uses the separate [**installer binder**](../../binders/installer/README.md) (Section I2). This page is what to **verify** and **who to call** when Wi‑Fi misbehaves.

---

## What clinic staff need

| Fact | Your clinic |
|------|-------------|
| Staff open EMR at | `https://________________________` |
| Server lives in | ________________________________ |
| Clinic IT / router contact | ________________________________ |
| Phone | ________________________________ |

---

## Quick checks (before calling IT)

| Symptom | Try this |
|---------|----------|
| EMR won’t load on iPad | Confirm iPad is on **clinic staff Wi‑Fi** (not guest / home) |
| “Cannot connect” on all devices | Check server cabinet power; see **Section A** or Pi binder |
| Certificate warning on iPad | CA may need re-trust — call IT (installer left `ca.crt`) |
| Backup helper from tablet | **Won’t work** — helper is only on server at `127.0.0.1:9090` |

---

## Monthly (add to Section B1)

| Task | Done | Date |
|------|------|------|
| Open EMR from a **staff iPad** on Wi‑Fi (not just admin PC) | ☐ | |
| Confirm URL still uses **https://** with no warning | ☐ | |

---

## Do not ask IT to do for daily care

- Open EMR ports to the **public internet**
- Use personal phones as official chart devices
- Run backup helper from the internet

---

**Installer reference (IT only):** [I2 — Network & Wi‑Fi](../../binders/installer/I2-network-and-wifi.md)
