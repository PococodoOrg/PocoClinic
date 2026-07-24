# ADR-0014: Private LAN-Only Deployment (No Mobile, No Internet)

## Status
Accepted

## Context

PocoClinic serves **small clinics on a private local network**. There is **no product plan for mobile apps**, **no cloud dependency**, and **no clinical use over the public internet**. The clinic server and staff workstations stay on an isolated LAN; daily care must work with **zero internet connectivity**.

We still **build defensively** because a private network is not a perfect trust boundary: a compromised workstation, stolen badge, misconfigured router, or malicious insider can exist on the LAN. Security controls (auth, audit, rate limits, session timeouts) remain required—they are not optional “because we’re offline.”

## Decision

### Deployment model

1. **One clinic server** (Pi or small PC) runs the API, database, and static web UI on the **LAN only**.
2. **Staff clients** are **clinic-owned browsers** on the LAN—reached via LAN hostname or static IP (e.g. `https://poco.clinic.local`). **Primary:** iPad-class tablets in portrait (exam rooms / bedside). **Secondary:** desktop and laptop workstations at the front desk.
3. **No native mobile applications** (iOS/Android). No app-store distribution. No “companion mobile” roadmap. Tablet use is **Safari/Chrome against the local website**, not a packaged app.
4. **No patient-facing portal** or any feature that requires routing PHI to the public internet.
5. **Local-only ops tools** (e.g. backup helper on `127.0.0.1:9090`, Pi touchscreen kiosk for backup/restore) are part of **on-prem ops**, not a separate mobile product.

### Network assumptions

| Assumption | Implication |
|------------|-------------|
| Clinic LAN is **intentionally isolated** from the internet for production | No outbound clinical dependencies; releases installable from USB |
| Internet may exist elsewhere in the building | Firewall must block clinic server egress/ingress; document in runbook |
| LAN can still be **compromised or misused** | Keep auth, audit, RBAC, rate limits, secure cookies, Argon2id |

### Defensive controls (required regardless of isolation)

- Badge + PIN (two factors) on every staff session  
- JWT access tokens + HttpOnly refresh; **15-minute inactivity** timeout  
- Audit log for auth, patient access, admin actions, backup/restore  
- Admin-only gates for staff management, forms config, system restore  
- Rate limiting on auth endpoints  
- Bind production services to LAN; **never document or default to internet exposure**  
- Local backup helper **localhost-only** (not exposed on Wi‑Fi for convenience)

### Explicit non-goals

- Native mobile apps or phone-sized consumer EMR roadmap  
- Packaged tablet apps (clinic tablets use the LAN website in a browser)  
- Cloud sync, SaaS hosting, or multi-site over internet  
- Remote access features (VPN tunnels, “access from home”) as product scope  
- Features that phone home (telemetry, external AI, license servers) in production  
- Designing for public-internet threat models (CDN, WAF, OAuth with Google, etc.) as defaults  

### UI target

- **Primary:** iPad-class tablet **portrait** (~768×1024 CSS px) in a clinic browser on the LAN.  
- **Secondary:** desktop/laptop landscape (1024px and up).  
- Touch-first below ~992px (drawer nav, larger targets, card lists); denser desktop layouts from there up.  
- Phone layouts and native apps are **not** goals.

## Consequences

### Positive
- Clear scope: one network, one server, browser clients  
- Simpler architecture; no mobile build pipelines or app-store compliance  
- Aligns with volunteer-operated, air-gapped clinic reality  

### Negative
- Staff cannot use personal phones as official clients (by design)  
- Remote telehealth or home access is out of scope unless clinic adds their own VPN (unsupported)  

### Related
- [ADR-0010 — Local air-gapped deployment](./0010-local-airgap-deployment.md)  
- [ADR-0015 — Clinic database engine](./0015-clinic-database-engine.md)  
- [NETWORK-AND-SECURITY.md](../docs/NETWORK-AND-SECURITY.md)  
- [VISION.md](../docs/VISION.md)  
- [FEATURES.md](../docs/FEATURES.md)
