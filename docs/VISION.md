# PocoClinic — Product Vision & Design Intent

This document is the north star for what PocoClinic is **for**, who it serves, and what trade-offs we accept. Technical details live in [ADRs](../adr/README.md); feature tracking lives in [FEATURES.md](./FEATURES.md).

---

## One-sentence intent

**PocoClinic is an offline-first EMR that runs on a small clinic-owned device on a private LAN, lets staff sign in with a badge and PIN from clinic browsers (primarily tablets, also desktops), and can be backed up and restored by non-technical people using simple, repeatable steps— with no internet and no native mobile apps.**

---

## Who uses it (and who does not)

| In scope | Out of scope |
|----------|--------------|
| Clinic-owned **browsers on the private LAN** — **primary** iPad-class tablets (portrait), **secondary** desktops/laptops | Native **mobile apps** (iOS/Android) or phone-sized consumer EMR |
| Browser-based EMR at bedside / exam rooms and front desk | Patient portals or staff access over the **public internet** |
| One **local server** (Pi or small PC) | Cloud hosting, SaaS, multi-site sync over internet |
| **Local-only** backup/restore (USB, localhost helper, Pi touchscreen for ops) | Packaged tablet apps (use the LAN website instead) |
| Deliberate **offline** updates from USB at maintenance time | Silent auto-updates that require internet |

**Network assumption:** Production clinics treat the server and workstations as an **isolated private network** with **no internet required or expected** for daily care.

**Security assumption:** We **design for isolation** but **implement defensively** anyway—badge + PIN, audit logs, session timeouts, rate limits—because LANs can still be misused (stolen badge, compromised PC, misconfigured router, insider threat). We do **not** optimize for internet-facing deployment.

**Related ADR:** [0014 — Private LAN-only deployment](../adr/0014-private-lan-only-deployment.md)

---

## The three pillars

### 1. Local, air-gapped clinic server

**Intent:** The system runs entirely on hardware the clinic controls—typically a Raspberry Pi or similar small PC—on the **local network only**. It must work **without internet access** for day-to-day care.

| Principle | What it means |
|-----------|---------------|
| No cloud dependency | Patient data and auth never require an external service to function |
| Single-box deploy | One device serves the UI and API to **clinic LAN browsers** (tablets and workstations—not native apps) |
| Modest hardware | Target: 4–8 GB RAM, SSD storage, ARM64 or x86 |
| Offline updates | Software and security updates are deliberate (USB or admin window), not silent auto-updates |

**Implications for design:**

- No third-party SaaS, analytics, or AI that phones home by default
- Dependencies and images must be installable/buildable without ongoing internet on the clinic device (or bundled in a release artifact)
- TLS is still used on the LAN where practical; “air-gapped” means **no egress to the public internet**, not “no encryption”
- Database and files live on local disk with clear backup boundaries
- **Defensive security stays on** (auth, audit, RBAC) even though we assume an isolated LAN—see [ADR-0014](../adr/0014-private-lan-only-deployment.md)

**Related ADR:** [0010 — Local air-gapped deployment](../adr/0010-local-airgap-deployment.md)

---

### 2. Strong but simple authentication (badge + PIN)

**Intent:** Clinic staff authenticate quickly during a busy day, without memorizing passwords, while still meeting a **two-factor** bar appropriate for PHI access.

**Chosen model:**

1. **Employee badge (QR code)** — primary factor  
   - Each staff member receives a physical badge with a QR code  
   - The QR encodes a **unique high-entropy credential** (64-bit key material, represented as a string the system hashes at rest)  
   - Scanning the badge identifies *who* is signing in; the badge should be treated like a physical key

2. **4-digit PIN** — secondary factor  
   - Known only to the employee  
   - Fast to enter at the front desk or in an exam room  
   - Rate-limited; account locks after repeated failures  

**UX target:** Staff sign in at `/login` → scan badge → enter PIN. Administrators use `/login/admin` (email + key + PIN) for first-time setup and badge printing before staff badges exist.

**Security properties (already aligned in code direction):**

- Badge secret and PIN stored separately with Argon2id  
- Short-lived access tokens; refresh via HttpOnly cookie  
- Session timeout after **15 minutes of inactivity**  
- All login, refresh, and logout events audit-logged  

**Related ADR:** [0011 — Employee badge authentication](../adr/0011-employee-badge-authentication.md) (clarifies [ADR-0002](../adr/0002-authentication-system.md))

---

### 3. Hardy operations — backup & restore for non-technical admins

**Intent:** A clinic volunteer or office manager—not a sysadmin—must be able to **protect** and **recover** the entire system using a printed runbook and labeled USB drives.

**Ops philosophy:** *Simple, physical, verifiable.*

| Practice | Purpose |
|----------|---------|
| **Daily USB backup** | Copy a verified backup bundle to a labeled drive |
| **Drive rotation** | e.g. Mon/Wed/Fri drives, kept in a locked drawer |
| **Printed checklist** | Laminated steps: “green light = OK”, “red = call support” |
| **Quarterly restore drill** | Restore to a spare device or VM to prove backups work |
| **Administrator’s binder** | Setup, backup, restore, and emergency contacts in one place |

**Backup bundle should include (target state):**

- Database dump (patients, users, sessions, audit logs)  
- Configuration (non-secret settings)  
- Backup manifest with date, version, and checksum  
- Optional: printable “last good backup” sticker for the binder  

**Restore target:** From bare metal (or fresh SD card) + USB → clinic operational within a documented time budget (goal: under 1 hour with runbook).

**Related ADR:** [0012 — Backup and recovery operations](../adr/0012-backup-and-recovery-ops.md)

---

## System philosophy (unchanged, refined)

From [FEATURES.md](./FEATURES.md):

> **Simple but Secure** — easy processes for non-technical administrators, physical documentation, friendly UI, robust but straightforward security.

Everything we build should pass this filter:

- *Can a volunteer follow the runbook?*  
- *Does it work with **no internet** on the clinic network?*  
- *Is the secure path also the easy path?*  
- *Would this still be safe if someone rogue is on the LAN?*

---

## Current implementation vs intent

Honest snapshot — see [FEATURES.md](./FEATURES.md) for live checklist detail.

| Area | Intent | Today | Remaining gap |
|------|--------|-------|----------------|
| **Local deploy** | Pi / small PC, LAN-only | Release tarball (`build-release.mjs`), [`clinic/server/`](../clinic/server/README.md), systemd, `/etc/pococlinic/env` | microSD **image-build** (separate from app tarball); Pi soak test |
| **Air-gap** | No internet required for care | No clinical cloud deps; offline tarball install after USB transfer | Fully offline **OS** flash workflow (image-build process) |
| **Auth UX** | Badge QR + PIN | Staff/admin login, badge reissue, PIN change, staff CRUD | — |
| **Auth security** | Badge+PIN, sessions, audit | Argon2id, HttpOnly cookies, lockout, audit log | — |
| **Patient records** | Core EMR | SQLite ([ADR-0015](../adr/0015-clinic-database-engine.md)), notes, documents, exercise log | — |
| **Backup / restore** | USB + runbook | CLI, admin UI, ops helper, printed binders | Drill habit (process) |
| **Ops binder** | Physical admin guide | Markdown + binder-printer + in-app help | Optional PDF export automation |

---

## Explicit non-goals (for now)

To protect the three pillars, we **defer** unless a clinic explicitly opts in (outside supported product scope):

- Native **mobile applications** (iOS/Android) or mobile-first EMR  
- Multi-site cloud sync  
- Patient-facing portals requiring public internet  
- Staff **remote access** (VPN, “login from home”) as a built-in feature  
- Always-on external AI / LLM services  
- Complex role hierarchies beyond small-clinic needs  
- Features that require a dedicated IT team to operate  
- Public-internet deployment patterns (OAuth with Google, CDN-hosted PHI, etc.)

---

## Suggested roadmap order

1. ~~**Badge login UX**~~ ✅  
2. **Production local deploy** — tarball + `clinic/server/` ✅ in repo; **Pi microSD image** + soak test 🏗️  
3. ~~**Backup / restore v1**~~ ✅  
4. ~~**Admin health dashboard**~~ ✅  
5. **Ops binder polish** — print via [`clinic/workstation/`](../clinic/workstation/README.md) ✅; optional static help on image  
6. **AI assistance** — local LLM placeholder 📝 ([FEATURES](./FEATURES.md))

---

## Document map

| Document | Role |
|----------|------|
| [VISION.md](./VISION.md) (this file) | Why we exist; pillars; gaps |
| [NETWORK-AND-SECURITY.md](./NETWORK-AND-SECURITY.md) | LAN-only scope, no mobile, defensive security |
| [FEATURES.md](./FEATURES.md) | Feature checklist and status |
| [adr/](../adr/README.md) | Durable technical decisions (repo root — GitHub-friendly) |
| [clinic/](../clinic/README.md) | Post-deploy server install tree (tarball) vs dev scripts |
| [README.md](../README.md) | Developer onboarding |

When intent and code diverge, **update this file and FEATURES.md first**, then ADRs, then implement.
