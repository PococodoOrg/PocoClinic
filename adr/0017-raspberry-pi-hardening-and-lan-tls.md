# ADR-0017: Raspberry Pi hardening and LAN TLS

## Status

Accepted

## Context

PocoClinic runs on a **private clinic LAN** ([ADR-0014](./0014-private-lan-only-deployment.md)). Staff use browsers on clinic-owned devices; the server is often a **Raspberry Pi**. Product docs already show `https://pococlinic.local`, and [VISION.md](../docs/VISION.md) states TLS should be used on the LAN where practical.

We do **not** use public CAs or Let's Encrypt — the server is not on the internet. Clinics still need:

1. **Host hardening** — firewall, SSH, file permissions, no WAN exposure
2. **HTTPS on the LAN** — encrypt auth cookies and badge/PIN traffic; enable `Secure` cookies in production
3. **Offline-capable setup** — no cloud dependency for certificate issuance

## Decision

### Pi hardening

Ship **`clinic/server/scripts/harden-pi.sh`** in the release tarball. It applies a baseline:

- `sysctl` network hardening (no forwarding, syncookies)
- SSH key-only, no root login
- **UFW** default-deny with LAN-only HTTPS (and temporary HTTP during migration)
- Permissions on `/etc/pococlinic/env`, `/var/lib/pococlinic`, and TLS files

Document manual router/Wi‑Fi steps in [`devices/raspberry-pi/hardening.md`](../devices/raspberry-pi/hardening.md).

### LAN TLS

Ship **`clinic/server/scripts/generate-lan-tls.sh`** — offline `openssl` script that creates:

- Clinic-local **CA** (`ca.crt` / `ca.key`)
- **Server cert** with SAN for hostname + optional LAN IP

Two supported deployment paths (documented in [`devices/raspberry-pi/tls-lan.md`](../devices/raspberry-pi/tls-lan.md)):

| Path | Mechanism |
|------|-----------|
| **A (recommended)** | **Caddy** on :443 → reverse proxy to `127.0.0.1:8080`; example [`clinic/server/caddy/Caddyfile.example`](../clinic/server/caddy/Caddyfile.example) |
| **B** | **Go native TLS** via `SERVER_TLS_CERT` / `SERVER_TLS_KEY` env vars |

Staff devices **trust `ca.crt` once** (MDM, GPO, or manual) — standard practice for internal LAN services.

### Application changes

- `SERVER_TLS_CERT` and `SERVER_TLS_KEY` optional env vars; `ListenAndServeTLS` when set
- **HSTS** header only when request is HTTPS (direct TLS or `X-Forwarded-Proto: https`)
- Production `CookieSecure` already true — requires HTTPS for cookie-based auth

Ops helper remains **`127.0.0.1:9090` only** — not TLS-terminated on the LAN.

## Alternatives considered

| Alternative | Why not |
|-------------|---------|
| Let's Encrypt / public CA | Server not on public internet; violates air-gap posture |
| HTTP only on LAN | Weak against Wi‑Fi sniffing; blocks `Secure` cookies meaningfully |
| TLS only in Go, no Caddy | Binding :443 requires capabilities; reverse proxy keeps app on loopback |
| mkcert (dev tool) | Fine for dev; production uses shipped `generate-lan-tls.sh` + clinic CA |

## Consequences

- Image build or first install must run hardening + TLS scripts (or equivalent manual steps)
- Clinics must distribute and trust `ca.crt` to staff devices
- `ALLOWED_ORIGIN` must match `https://` clinic hostname after cutover
- When using Caddy: `SERVER_HOST=127.0.0.1`, `TRUSTED_PROXIES=127.0.0.1`
- Certificate rotation (~27 months default) is a documented maintenance task

## References

- [devices/raspberry-pi/hardening.md](../devices/raspberry-pi/hardening.md)
- [devices/raspberry-pi/tls-lan.md](../devices/raspberry-pi/tls-lan.md)
- [NETWORK-AND-SECURITY.md](../docs/NETWORK-AND-SECURITY.md)
- [ADR-0010 — Local air-gapped deployment](./0010-local-airgap-deployment.md)
