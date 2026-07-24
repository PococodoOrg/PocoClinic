# ADR-0010: Local Air-Gapped Deployment

## Status
Accepted

## Context

PocoClinic is intended for small clinics—often non-profits with no IT staff—that need an EMR on a **private local network with no internet for clinical use**. Typical deployment is a Raspberry Pi or small PC on the clinic LAN, serving **staff browsers on clinic workstations** on the same network. There is **no mobile app** and **no cloud** in scope.

See also: [ADR-0014 — Private LAN-only deployment](./0014-private-lan-only-deployment.md).

Constraints:

- Day-to-day operation must not depend on cloud APIs, DNS to external services, or SaaS auth
- Hardware is modest (ARM64, 4–8 GB RAM, SSD preferred)
- Updates are intentional (USB or scheduled maintenance), not silent background pulls
- Data must remain on clinic-controlled storage

## Decision

We will target **single-node, offline-capable deployment** as the primary production shape:

1. **Topology**
   - One clinic server runs: optional reverse proxy, Go API + static frontend build, and a **SQLite file** at `DATABASE_URL` ([ADR-0015](./0015-clinic-database-engine.md))
   - Client devices reach the server via LAN hostname or static IP (e.g. `https://poco.clinic.local`)
   - No inbound or outbound internet requirement for clinical use

2. **Packaging (target)**
   - Release artifact: versioned tarball for **microSD / Raspberry Pi** with pre-built ARM64 binaries and migration bundle — see [DEPLOYMENT-BOUNDARY.md](../docs/deploy/DEPLOYMENT-BOUNDARY.md)
   - Install script or flashable image for Raspberry Pi OS / Debian (image-build process is separate from application tarball)
   - `systemd` units for API and ops helper with restart on failure (no separate database daemon)

3. **Configuration defaults**
   - Bind services to LAN interface; document firewall expectation (clinic LAN only)
   - `ENV=production`: secure cookies, no dev secrets
   - No telemetry or external license checks

4. **Development vs production**
   - Developers may use local SQLite files and `.env` for builds; **releases** must document offline install from USB
   - CI produces checksum-signed release bundles

5. **Updates**
   - Documented upgrade path: stop services → backup → **`/opt/pococlinic/bin/migrate`** (or `go run ./cmd/migrate` in dev) → deploy new app → start services
   - Migrations must stay N-1 compatible; breaking schema changes use expand/contract across two releases ([ADR-0013](./0013-database-migrations.md))
   - No auto-update channel in v1

## Consequences

### Positive
- Matches clinic privacy and reliability needs
- Predictable costs (hardware only)
- Works during ISP outages

### Negative
- We own the entire release and update pipeline
- Security patches require deliberate clinic maintenance
- ARM64 and resource limits may constrain database choice

### Mitigations
- Simple health endpoint and admin “system OK” dashboard
- Backup-before-upgrade enforced in runbook
- Keep modular monolith (ADR-0001) to limit moving parts
- Release layout and microSD boundary: [DEPLOYMENT-BOUNDARY.md](../docs/deploy/DEPLOYMENT-BOUNDARY.md)
