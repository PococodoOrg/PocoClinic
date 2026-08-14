# Deployment boundary — microSD / Raspberry Pi

This document defines **what PocoClinic owns** in a clinic-server deployment and **what a separate image-build process owns**. It is the contract between this repository and whoever produces the bootable microSD card.

## End goal (context)

A volunteer inserts a **microSD card** into a **Raspberry Pi** at the clinic. On power-on, the Pi becomes the **PocoClinic server** on the private LAN. Staff use **desktop browsers on clinic workstations** — not phones — to reach the EMR.

```
┌─────────────────────────────────────────────────────────────────┐
│  Separate image-build process (out of scope for this repo)        │
│  • Raspberry Pi OS (or derivative) on microSD                   │
│  • Boot → autostart clinic services                             │
│  • Optional: Pi touchscreen → Chromium kiosk → ops helper       │
│  • Hostname / static IP, firewall, time sync                    │
└───────────────────────────┬─────────────────────────────────────┘
                            │ consumes
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  PocoClinic release artifact (this repo produces)               │
│  • Pre-built ARM64 binaries                                     │
│  • Frontend static bundle                                       │
│  • Embedded SQL migrations (inside migrate binary / main)       │
│  • systemd unit templates                                       │
│  • Default config template → /etc/pococlinic/env                │
│  • Version manifest + checksums                                 │
└─────────────────────────────────────────────────────────────────┘
```

**We do not own:** Pi firmware, OS image flashing, SD card cloning workflow, HDMI/touchscreen autostart, Chromium kiosk autolaunch, or clinic router configuration. Those integrate **against** the boundaries below.

---

## Responsibility split

| Concern | Owner | Notes |
|---------|-------|-------|
| OS, kernel, microSD hardware | Image-build process | Pi OS Lite or custom image |
| First-boot autostart of services | Image-build process | Enables systemd units we ship |
| Pi touchscreen kiosk (Chromium) | Image-build process | Uses `scripts/pi/start-backup-kiosk.sh` as reference |
| PocoClinic application binaries | **This repo** | `pococlinic`, `ops-helper`, CLI tools |
| Database (SQLite file) | **Shared** | App opens file at `DATABASE_URL`; image ensures durable data dir + permissions |
| Schema migrations | **This repo** | `cmd/migrate`; run on upgrade |
| Clinical web UI (EMR) | **This repo** | React build served to LAN |
| Backup helper UI (localhost) | **This repo** | Pi touch build embedded in `ops-helper` |
| Persistent PHI data paths | **Contract** | Paths documented here; image process creates mount points |
| LAN hostname / TLS | Image-build or clinic IT | e.g. `https://pococlinic.local` (after TLS); first bring-up may use `http://<ip>:8080` — see [Pi setup](../../devices/raspberry-pi/setup.md) |
| USB backup rotation SOP | Clinic ops | Documented in runbooks, not in software |

---

## Runtime topology on the Pi

```
                    Clinic LAN (no internet required)
    ┌──────────────────────────────────────────────────────────┐
    │                                                          │
    │   Workstation browsers ──HTTP──► :8080  pococlinic       │
    │   (staff EMR)                         (API + static UI)  │
    │                                                          │
    │   ┌─────────────────────────────────────────────────┐    │
    │   │ Raspberry Pi (clinic server)                    │    │
    │   │                                                 │    │
    │   │  SQLite file  (e.g. /var/lib/pococlinic/pococlinic.db) │
    │   │         ▲                                     │    │
    │   │         │                                     │    │
    │   │  pococlinic.service ──► 0.0.0.0:8080         │    │
    │   │                                                 │    │
    │   │  pococlinic-ops-helper.service                  │    │
    │   │         └──► 127.0.0.1:9090  (localhost ONLY)   │    │
    │   │                    ▲                            │    │
    │   │                    │                            │    │
    │   │  [optional] Chromium kiosk on Pi HDMI/touch     │    │
    │   └─────────────────────────────────────────────────┘    │
    │                                                          │
    └──────────────────────────────────────────────────────────┘
```

### Network rules (non-negotiable)

| Listener | Bind address | Reachable from | Purpose |
|----------|--------------|----------------|---------|
| EMR API + UI | `0.0.0.0:8080` (or reverse proxy) | Clinic LAN only | Staff charting |
| Ops helper | `127.0.0.1:9090` | Pi console only | Backup/restore kiosk |
| SQLite | _(file on disk — no port)_ | App process only | Database |

The image-build process must **not** port-forward 9090 to the LAN “for convenience.”

---

## Target filesystem layout

Paths assume a standard FHS layout. The image-build process creates users, directories, and permissions; the PocoClinic tarball installs into ` /opt`.

```
/opt/pococlinic/                    # Read-only application (upgrade replaces this tree)
├── pococlinic                      # main server binary
├── ops-helper                      # localhost backup UI binary
├── migrate                         # migration CLI (also runnable via go run in dev)
├── backup                          # cron backup CLI
├── restore                         # disaster recovery CLI
├── audit-purge                     # optional retention CLI
├── bin/                            # env-aware wrappers (prefer in cron/SSH)
│   ├── migrate
│   ├── backup
│   ├── restore
│   └── audit-purge
├── cron/                           # example /etc/cron.d drop-ins
├── install.sh
├── env.template
└── static/                         # frontend production build

/etc/pococlinic/
└── env                             # secrets + paths (not in tarball; created at install)

/var/lib/pococlinic/                # Persistent PHI — survives app upgrades
├── backups/                        # BACKUP_DIR — tarball archives
├── documents/                      # DOCUMENTS_DIR — legacy files only (new uploads are DB blobs)
└── pococlinic.db                   # SQLite file (typical; path from DATABASE_URL)

/var/log/pococlinic/                # Optional stdout from cron jobs
```

### microSD vs durable storage

microSD cards wear with heavy writes. **Recommended boundary:**

| Data | Preferred location |
|------|-------------------|
| OS + PocoClinic binaries | microSD (read-mostly after install) |
| SQLite data | USB SSD or NAS mount if available |
| Document files + backups | Same durable volume as DB when possible |

If everything must live on one microSD for v1, document backup frequency and SD replacement cadence in the ops binder. The software boundary (`BACKUP_DIR`, `DOCUMENTS_DIR`, DB data dir) stays the same — only mount points change.

---

## Release artifact (what we ship)

**Target shape:** versioned tarball, e.g. `pococlinic-1.0.0-linux-arm64.tar.gz`

Build on a dev machine:

```bash
node scripts/build-release.mjs
# or: build-release.bat
```

Output: `dist/pococlinic-<version>-linux-arm64.tar.gz` with `MANIFEST.json` checksums.

| Contents | Source in repo |
|----------|----------------|
| `pococlinic`, `ops-helper`, CLI binaries | `go build` from `backend/cmd/` |
| `bin/*` wrappers | `clinic/server/bin/` — load `/etc/pococlinic/env` |
| `install.sh`, `env.template`, `cron/` | `clinic/server/` |
| Frontend static files | `frontend` build → tarball `static/` |
| Ops-helper static (Pi build) | `build-ops-helper-pi.bat` output baked into ops-helper embed dir at build time |
| systemd units | `scripts/deploy/*.service` |
| `MANIFEST.json` + SHA256 checksums | Generated at release |
| Operator docs | Subset copied from `docs/ops/`, `docs/deploy/` |

**Not in the tarball:** `DATABASE_URL` secrets, JWT secrets, or existing PHI. Those are created on first clinic setup.

---

## Process boundaries (systemd)

Units in `scripts/deploy/` define the **application** boundary:

| Unit | Starts | Depends on |
|------|--------|------------|
| `pococlinic.service` | Main EMR (embeds SQLite file access) | network |
| `pococlinic-ops-helper.service` | Backup UI | network |

The image-build process is responsible for:

```bash
systemctl enable pococlinic pococlinic-ops-helper
```

Optional cron (documented, not systemd):

| Schedule | Command |
|----------|---------|
| Daily ~6 PM | `/opt/pococlinic/bin/backup` |
| Monthly | `/opt/pococlinic/bin/audit-purge` |

---

## Install & upgrade flow

Boundary between **first install** and **version upgrade**:

### First install (clinic go-live)

1. Image-build process flashes microSD and boots Pi on LAN.
2. Image creates `pococlinic` user, `/var/lib/pococlinic/*`, `/etc/pococlinic/env`.
3. Extract tarball and run **`sudo ./install.sh`** (or manual extract → `/opt/pococlinic/`).
4. Run `/opt/pococlinic/bin/migrate` (or `sudo ./install.sh` then migrate per install output).
5. `systemctl start pococlinic pococlinic-ops-helper`.
6. Administrator completes in-app setup (bootstrap admin, staff, first backup) — see admin guide.

### Upgrade (new release tarball)

1. **Stop** `pococlinic` and `pococlinic-ops-helper` (keep DB running).
2. **Backup** via CLI or admin dashboard.
3. Replace `/opt/pococlinic/` application tree (not `/var/lib/pococlinic/`).
4. Run **`sudo /opt/pococlinic/bin/migrate`**.
5. **Start** services.
6. Admin runs health check + compliance checklist.

Rollback boundary: restore previous `/opt/pococlinic/` binary tree **and** DB restore from pre-upgrade backup if migrations are not reversible.

---

## Build-time vs run-time boundary

| Phase | Where it runs | Output |
|-------|---------------|--------|
| **Build** | Developer machine or CI (may use internet) | arm64 tarball, checksums |
| **Transfer** | USB stick to clinic | Tarball copied to Pi |
| **Install** | Pi (offline OK) | Extract to `/opt`, configure env |
| **Run** | Pi 24/7 on LAN | systemd services, no outbound clinical deps |

Development shortcuts (`run-all.bat`, Vite dev server) **do not cross** into production — they exist only on developer workstations.

---

## Remaining validation (not blocking packaging)

| Item | Status |
|------|--------|
| Pi 4GB soak test | SQLite under `/var/lib/pococlinic/` — accepted per [ADR-0015](../../adr/0015-clinic-database-engine.md); hardware soak still open |
| Evaluator docs on image | Optional copy of `docs/guide-site/` — generate via `build-help-site.mjs` |

**Resolved (for reference):**

| Item | State |
|------|--------|
| Frontend serving | Main binary serves `static/` when `ENV=production` |
| Release build | `scripts/build-release.mjs` + `clinic/server/` install tree |
| Clinic runtime packaging | [ADR-0016](../../adr/0016-clinic-runtime-packaging.md) — `clinic/server/bin/*` wrappers |
| Embedded migrations | `go:embed` in `cmd/migrate` / main |

---

## Related documents

- **[Pi beginner setup](../../devices/raspberry-pi/setup.md)** — flash OS + install + first login (for humans, not image builders)
- [Deploy documentation hub](./README.md)
- [Tools & scripts catalog](../ops/tools-and-scripts.md) — every helper referenced above
- [Clinic runtime packaging (ADR-0016)](../../adr/0016-clinic-runtime-packaging.md)
- [clinic/server/](../../clinic/server/README.md) — install tree source
- [scripts/deploy/README.md](../../scripts/deploy/README.md) — systemd install steps
- [ADR-0010 — Local air-gapped deployment](../../adr/0010-local-airgap-deployment.md)
- [ADR-0014 — Private LAN-only deployment](../../adr/0014-private-lan-only-deployment.md)
- [ADR-0015 — Clinic database engine](../../adr/0015-clinic-database-engine.md)
- [Pi touchscreen backup](../../devices/raspberry-pi/touchscreen-backup.md)
- [Devices · Raspberry Pi](../../devices/raspberry-pi/README.md)
