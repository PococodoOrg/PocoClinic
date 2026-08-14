# What is PocoClinic?

PocoClinic is an open-source **Electronic Medical Records (EMR)** system for **small clinics on a private local network** — especially non-profits with limited technical staff.

## One sentence

A clinic-owned server on your LAN runs the EMR; staff sign in with **badge + PIN** from desktop browsers; **USB backups** protect data; **no internet** is required for daily care.

## Who it is for

| Good fit | Not a fit |
|----------|-----------|
| Single-location clinic on a **private LAN** | Multi-site cloud sync over the internet |
| **Desktop/laptop browsers** at front desk and exam rooms | Native **mobile apps** for staff |
| Non-technical admins who can follow **printed runbooks** | Enterprises needing SaaS support contracts |
| Clinics that want **full data control** on local hardware | Patient portals or home access over the public internet |

## What you get today

- **Authentication** — employee badge (QR) + 4-digit PIN, session timeouts, account lockout
- **Patient records** — demographics, search, chart notes, encrypted document uploads, exercise log, forms
- **Forms** — admin-built templates; staff submissions on patient charts
- **Audit log** — sign-ins, patient access, admin actions
- **Admin dashboard** — system health, backup status, staff management, reports
- **Backup & restore** — USB rotation, localhost Backup Helper, integrity verification
- **In-app help** — searchable guides on the LAN without internet

## How it is deployed

```
┌─────────────────────────────────────────────────┐
│           Clinic private LAN (no internet       │
│           required for daily care)              │
│                                                 │
│  ┌──────────────┐    ┌────────────────────────┐ │
│  │ Clinic server│    │ Staff workstations     │ │
│  │ Pi / small PC│◄──►│ Browsers → PocoClinic  │ │
│  │ DB + files   │    │ /login, /patients, …   │ │
│  │ Backup helper│    └────────────────────────┘ │
│  │ :9090 local  │                               │
│  └──────┬───────┘                               │
│         │ USB backups (labeled rotation)        │
└─────────┼───────────────────────────────────────┘
          ▼
    Locked USB drives
```

## Security model (summary)

We **assume** an isolated clinic LAN but **build defensively** anyway:

- Badge + PIN (two factors)
- Argon2id hashed credentials
- JWT sessions with 15-minute inactivity timeout
- Rate limiting and audit logging
- No cloud dependencies in production paths

Details: [NETWORK-AND-SECURITY.md](../NETWORK-AND-SECURITY.md)

## Try it (developers)

See the root [README.md](../../README.md) for local development (`run-all.bat`, SQLite file via `DATABASE_URL`). Production on a Pi: **[beginner setup](../../devices/raspberry-pi/setup.md)** · [first-time-setup](./for-administrators/first-time-setup.md) · [clinic/server/](../../clinic/server/README.md) · [deploy hub](../deploy/README.md).

## Next steps

- [Clinic network requirements](./clinic-network.md)
- [Administrator guide](./for-administrators/README.md)
- [Staff guide](./for-staff/README.md)
- [Product vision](../VISION.md)
