# PocoClinic operations binder

**Store this binder closed** (cabinet or shelf near the admin area).  
Open it for the three uses below — not for everyday charting (that is in the EMR Help center).

Secrets stay on a **separate** sheet in the **safe**: [Safe vault credentials](../safe-credentials-vault.md).

| Clinic | ________________________________ |
|--------|----------------------------------|
| Binder assembled by | ________________________________ |
| Date | ____________ |
| Review every | 6 months or after major staff/server change |

---

## When to open this binder

| Use case | Open when… | Section |
|----------|------------|---------|
| **1. Emergency** | PocoClinic or the server is **down**, or you must restore / respond to an incident | **A** |
| **2. Weekly / monthly** | You are doing **scheduled** checks, audits, or backup drills (not daily tick sheets) | **B** |
| **3. System help** | You need **printed** how-to that is the same on every workstation (no device-specific notes) | **C** |

**Easiest print path:** [binder printer](../../../binder-printer/README.md) via [`clinic/workstation/print-binder.bat`](../../../clinic/workstation/print-binder.bat) or repo-root `print-binder.bat` → choose **Clinic ops binder**.

Day-to-day patient work stays in the browser. Daily USB backup status is on the **Admin dashboard** — do not keep a daily paper log in this binder.

---

## Print order (by section)

Print **double-sided** when possible. Use three tab dividers: **A · B · C**.

### Section A — Emergency (system down)

| Page | File |
|------|------|
| A1 | [Emergency procedures](../emergency-procedures.md) |
| A2 | [Paper care when EMR is down](./A2-paper-fallback.md) |
| A3 | [Incident & restore log](./A3-incident-restore-log.md) |
| — | Emergency contacts filled on A3 (and Help → Clinic emergency contacts) |
| SAFE | [Vault credentials](../safe-credentials-vault.md) — **filled copy in the safe**, not loose in Section A |

### Section B — Weekly / monthly processes

| Page | File |
|------|------|
| B1 | [Weekly & monthly processes](./B1-periodic-processes.md) |
| B2 | [Monthly testing checklist](../monthly-testing-checklist.md) |
| B3 | [Security audit checklist](../security-audit-checklist.md) |
| B4 | [Backup & restore runbook](../administrator-runbook.md) (how to run backups/drills — no daily log) |
| B5 | [Physical security](../physical-security-binder.md) |

### Section C — System help (device-independent)

| Page | File |
|------|------|
| C1 | [What is PocoClinic](./C1-what-is-pococlinic.md) |
| C2 | [Staff how-to](./C2-staff-how-to.md) |
| C3 | [Administrator how-to](./C3-administrator-how-to.md) |

These pages describe the product the same way on every PC, laptop, or tablet. Do **not** write workstation passwords, Wi‑Fi keys, or per-device notes here — those belong in the vault or clinic IT records.

### Print checklist

- [ ] Cover (this page)
- [ ] Section A — A1, A2, A3 (+ contacts)
- [ ] Section B — B1–B5
- [ ] Section C — C1–C3 (extra C2 copies optional for training)
- [ ] SAFE — vault sheet filled, locked, digital copies deleted
- [ ] USB rotation labels (from Admin dashboard) stored with USB drives, not required in the binder

---

## Safe vs binder

| In this binder (cabinet) | In the safe only |
|--------------------------|------------------|
| Emergency steps, paper fallback, incident logs | Filled vault credentials |
| Weekly/monthly checklists and backup **instructions** | Optional sealed spare admin badge |
| Static system help (Sections C) | Live backup USB drives (or locked drawer) |

---

## Related (do not stuff the whole repo into the binder)

- [Ops hub](../README.md)
- [Binder printer utility](../../../binder-printer/README.md) — pick clinic / Pi / vault packs
- [Raspberry Pi device binder](../../../devices/raspberry-pi/binder/README.md)
- [In-app Help](../../IN-APP-HELP.md) — searchable on the LAN when the system is up
- [First-time setup](../../guide/for-administrators/first-time-setup.md) — install once
