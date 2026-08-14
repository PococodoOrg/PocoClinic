# Network, deployment, and security model

This document states how PocoClinic is **meant to be used** in production. For the formal decision record, see [ADR-0014](../adr/0014-private-lan-only-deployment.md). For product intent, see [VISION.md](./VISION.md).

---

## In one paragraph

PocoClinic is built for **small clinics on a private local network with no internet for clinical use**. Staff use **clinic-owned browsers on the LAN** (primarily **iPad-class tablets in portrait**, secondarily desktops/laptops)—not native mobile apps, not cloud SaaS, not access from home over the public internet. We **assume the network is isolated**, but we **still implement strong local security** (badge + PIN, audit logs, sessions, rate limits) because LANs can be misused.

---

## Production topology

```
┌─────────────────────────────────────────────────────────┐
│  Clinic private LAN (isolated — no internet required)     │
│                                                         │
│  [Exam iPad]  [Front desk PC]  [Admin laptop]           │
│       │              │               │                  │
│       └──────────────┴───────────────┘                  │
│                         │                               │
│              https://poco.clinic.local                  │
│                         │                               │
│              ┌──────────▼──────────┐                    │
│              │  Clinic server      │                    │
│              │  (Pi or small PC)   │                    │
│              │  • Go API           │                    │
│              │  • SQLite      │                    │
│              │  • Static EMR UI    │                    │
│              │  • Backup helper    │                    │
│              │    (127.0.0.1 only) │                    │
│              └─────────────────────┘                    │
└─────────────────────────────────────────────────────────┘
         ✗ No cloud    ✗ No mobile apps    ✗ No patient portal on internet
```

---

## What is in scope

| Topic | Production expectation |
|-------|------------------------|
| **Network** | Private LAN; clinic controls firewall; server not exposed to internet |
| **Clients** | Clinic-owned LAN browsers; **primary** iPad portrait (~768×1024), **secondary** desktop 1024px+ |
| **UI product** | Tablet-first touch UX in the browser; phone layouts and native apps out of scope |
| **Connectivity** | **Zero internet** required for charting, auth, backup to USB, restore |
| **Server** | Single box: API + DB + web UI on one device |
| **Updates** | Deliberate: USB or maintenance window; offline install bundle |
| **Ops UI** | Backup helper on `127.0.0.1` only; optional Pi touchscreen for backup/restore kiosk **on the server** |
| **Security** | Badge + PIN, JWT sessions, 15m inactivity timeout, audit log, admin RBAC, rate limits |

---

## What is out of scope (by design)

| Topic | Status |
|-------|--------|
| Native **iOS/Android apps** | Not planned |
| **Phone-sized** / consumer mobile EMR product | Not planned (clinic tablets in browser are in scope) |
| **Cloud** hosting or SaaS | Not planned |
| **Multi-site sync** over internet | Not planned |
| **Patient portal** on public internet | Not planned |
| **Staff remote access** (VPN, login from home) as a built-in feature | Not supported |
| External **AI/LLM** services that phone home | Not in production path |
| Telemetry, license servers, analytics SaaS | Not in production path |

Personal phones are **not** official clients. A clinic may have internet elsewhere in the building; the **PocoClinic LAN segment should not depend on it** and should block server egress/ingress (document in the ops binder).

---

## Defensive security on an isolated LAN

We do **not** treat “offline” as “trust everyone on the network.” Controls that remain required:

1. **Two-factor staff auth** — physical badge (QR) + 4-digit PIN  
2. **Session limits** — short-lived access JWT and refresh JWT in HttpOnly cookies; 15-minute inactivity logout
3. **Account lockout** — after repeated failed PIN attempts  
4. **Audit logging** — sign-in, patient view/edit, notes, admin actions, backup/restore  
5. **Role separation** — admin-only staff management, forms config, system restore  
6. **Rate limiting** — especially on auth endpoints  
7. **Localhost-only ops** — backup helper not bound to `0.0.0.0` for LAN convenience  
8. **Argon2id** — badge secrets and PINs hashed at rest  

**Threats we design for on the LAN:** stolen badge, shoulder-surfed PIN, compromised workstation, misconfigured router, curious insider, rogue device on Wi‑Fi.

**Threats we do not optimize for:** public internet DDoS, OAuth with Google, CDN/WAF, multi-tenant public cloud hosting.

---

## Development vs production

| | Development | Production clinic |
|---|-------------|-------------------|
| Internet | May use npm, Docker Hub, GitHub for builds | **Not required** for daily operation |
| `DATABASE_URL` | SQLite file on the clinic server | On-server persistent DB |
| `OPS_HELPER_HOST` | `127.0.0.1` | **`127.0.0.1` only** |
| `ALLOWED_ORIGIN` | `http://localhost:3000` | Exact clinic LAN origin only |
| `COOKIE_SECURE` | Secure in production | `false` only for temporary HTTP bring-up; HTTPS before patients ([Pi TLS](../devices/raspberry-pi/tls-lan.md)) |
| Client devices | Developer machine | Clinic-owned workstations |

Developers may use the internet to build releases; **release artifacts** must install and run on a clinic server **without ongoing internet** (see [ADR-0010](../adr/0010-local-airgap-deployment.md)).

---

## Document map

| Document | Purpose |
|----------|---------|
| [VISION.md](./VISION.md) | Product pillars and non-goals |
| [FEATURES.md](./FEATURES.md) | Feature checklist |
| [ADR-0010](../adr/0010-local-airgap-deployment.md) | Air-gapped deployment packaging |
| [ADR-0014](../adr/0014-private-lan-only-deployment.md) | LAN-only, no mobile, no internet |
| [ADR-0016](../adr/0016-clinic-runtime-packaging.md) | Post-deploy `clinic/` vs dev `scripts/` |
| [clinic/](../clinic/README.md) | Server install tree shipped in release tarball |
| [Pi beginner setup](../devices/raspberry-pi/setup.md) | Flash OS → install → first admin login |
| [DEPLOYMENT-BOUNDARY](./deploy/DEPLOYMENT-BOUNDARY.md) | Tarball vs microSD image responsibilities |
| [ADR-0015](../adr/0015-clinic-database-engine.md) | Clinic DB engine (Accepted: SQLite) |
| [administrator-runbook.md](./ops/administrator-runbook.md) | Printable ops procedures |
| [pi-touchscreen-backup.md](../devices/raspberry-pi/touchscreen-backup.md) | Server-local backup kiosk (not mobile EMR) |
| [Devices · Raspberry Pi](../devices/raspberry-pi/README.md) | Pi hardware and device ops |
| [IN-APP-HELP.md](./IN-APP-HELP.md) | In-app Help center (content authoring) |

When adding features, ask: *Does this require internet, mobile apps, or public exposure?* If yes, it is out of scope unless an ADR explicitly changes that.
