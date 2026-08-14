# PocoClinic Features

## How to Use This Document

1. **Feature Updates**
   - Update feature status using the defined status emojis
   - Check off completed sub-tasks within features
   - Add new sub-tasks as needed
   - Document any blockers or dependencies

2. **When Implementing**
   - Update the feature status
   - Add necessary tests
   - Update related documentation
   - Create or update relevant ADRs
   - Ensure HIPAA compliance
   - Add to changelog

## System Philosophy

PocoClinic is designed with a **"Simple but Secure"** philosophy. Full product intent is in **[VISION.md](./VISION.md)**.

Three pillars:

1. **Local & offline** — runs on a Raspberry Pi or small PC on an **isolated clinic LAN**; **no internet** for daily use; **no mobile apps**  
2. **Badge + PIN auth** — employee QR badge (primary) + 4-digit PIN (secondary); strong but fast for staff  
3. **Friendly ops** — USB backup/restore and printed runbooks so non-technical admins can protect and recover the system  

**Product scope** ([ADR-0014](../adr/0014-private-lan-only-deployment.md), [NETWORK-AND-SECURITY.md](./NETWORK-AND-SECURITY.md)):

- **Assume** production is always a **private, isolated LAN** (small clinic, no cloud).  
- **Build defensively anyway** — auth, audit, sessions, rate limits — because LANs can be compromised or misused.  
- **Staff clients** = clinic-owned **browsers** — primarily **tablets**, also desktops/laptops; not phones or native apps.

Design filters:

- Easy-to-follow processes for non-technical administrators  
- Physical documentation and backup procedures  
- Clear, friendly user interfaces  
- Robust but straightforward security measures  

## Status Definitions

| Status | Description |
|--------|-------------|
| 🚀 Live | Feature is implemented and deployed |
| ✅ Complete | Feature is implemented and tested |
| 🏗️ In Progress | Feature is currently being developed |
| 📝 Planned | Feature is planned but not started |
| 🔄 Review | Feature needs review or revision |
| ⏸️ Paused | Development temporarily paused |

## Feature Overview

| Feature | Status | Description | Dependencies |
|---------|--------|-------------|--------------|
| System Administration | ✅ | Admin dashboard, backup management, system health monitoring | Documentation system, USB management |
| Authentication | ✅ | Badge (QR key) + PIN, sessions, cookie JWTs | SQLite, JWT cookies |
| Patient Management | ✅ | Patient records, notes, documents, exercise log | Auth system |
| User Interface | ✅ | React-based UI with error handling and navigation | Mantine UI, React Router |
| Backup System | ✅ | USB tar.gz with SQLite file payload (docs + exercise included) | USB management, Documentation |
| API Layer | ✅ | REST endpoints for the EMR and ops helper | Go backend, Auth system |
| Audit Logging | ✅ | HIPAA-oriented action and event tracking | Database, Auth system |
| Reporting | ✅ | Census, activity, staff, form CSV (no ad-hoc builder) | Data access layer |
| Documentation | ✅ | Physical admin guide and system documentation | Documentation generator |
| AI Assistance | 📝 | Lightweight task assistance and guidance | Local LLM, Task templates |

## Current Development Focus

Aligned with [VISION.md](./VISION.md) roadmap:

- Badge + PIN login UX (QR scanner) ✅  
- Local production deploy (Pi / microSD / systemd) 🏗️ — [release build](../../scripts/build-release.mjs), [`clinic/server/`](../../clinic/server/README.md), [deployment boundary](./deploy/DEPLOYMENT-BOUNDARY.md)  
- USB backup & restore v1 ✅  
- Auth, patient CRUD, audit logging, persistence ✅  

## Detailed Feature Specifications

### System Administration
**Status**: ✅ Complete
- Physical Administrator's Guide
  - [x] Step-by-step setup instructions (`docs/guide/for-administrators/first-time-setup.md`, in-app admin-setup)
  - [x] Troubleshooting guides (`docs/guide/for-administrators/troubleshooting.md`, in-app troubleshooting-common)
  - [x] Emergency procedures (`docs/ops/emergency-procedures.md`)
  - [x] Contact information forms (Help → Clinic emergency contacts + runbook template)
- Backup System
  - [x] CLI backup command (`cmd/backup`; dev: `clinic/dev-windows/backup.bat`) with manifest + SQLite file (`VACUUM INTO`)
  - [x] Admin UI “Backup now” (`POST /admin/backups`)
  - [x] Admin UI restore with confirmation (`POST /admin/restore`)
  - [x] Local Backup Helper site (`cmd/ops-helper`, http://127.0.0.1:9090, localhost only)
  - [x] Raspberry Pi touchscreen mode (large targets, bottom nav, backup/restore only)
  - [x] Pi kiosk launcher (`scripts/pi/start-backup-kiosk.sh`, `build-ops-helper-pi.bat`)
  - [x] Friendly step-by-step backup and restore wizards
  - [x] Printable daily checklist in helper UI
  - [x] Daily USB backup reminders (dashboard shows age warnings)
  - [x] Backup age warnings on admin dashboard (ok / warning / critical)
  - [x] Admin task board (BACKUP, ARCHIVE AUDIT, RESTORE DRILL, HEALTH CHECK, COMPLIANCE, MONTHLY TEST badges)
  - [x] Admin task notifications (backup overdue, DB offline, locked accounts, pending migrations, monthly testing, compliance, health check, audit archive)
  - [x] CLI restore command (`cmd/restore`; dev: `clinic/dev-windows/restore.bat`) with checksum verification
  - [x] Labeled USB rotation system (admin dashboard, printable labels, localStorage)
  - [x] Backup verification process (admin Verify button + checksum manifest)
  - [x] Recovery testing tracker (quarterly restore drill reminders on admin dashboard)
- System Health Dashboard
  - [x] Simple status indicators (storage mode, DB connectivity, migrations)
  - [x] Resource counts (patients, staff, forms, sessions)
  - [x] Backup status tracking
  - [x] Maintenance reminders (BACKUP, ARCHIVE AUDIT, restore drill, monthly testing, compliance, health check, default PIN alerts)
  - [x] Security status overview (sessions, locks, document storage health)

### Authentication System
**Status**: ✅ Complete

Target UX: **scan employee badge (QR) → enter PIN** at `/login`; administrators use `/login/admin` (email + key + PIN) for bootstrap setup ([ADR-0011](../adr/0011-employee-badge-authentication.md)).

- [x] Basic user model with hashed key + PIN (Argon2id)
- [x] Session management (JWT access + refresh in HttpOnly cookies; 15m inactivity)
- [x] Staff login (`POST /auth/login/staff`) and admin bootstrap login (`POST /auth/login/admin`); protected patient routes
- [x] Audit logging for auth events
- [x] Staff badge login UI at `/login` (USB scanner or paste key → PIN)
- [x] Admin bootstrap login UI at `/login/admin`
- [x] Staff management UI at `/users` (admin only): list, create, edit, delete, badge QR print
- [x] Badge reissue from staff detail page (admin only)
- [x] PIN self-change at `/account/pin`
- [x] Per-user activity history (sign-ins, PIN changes, badge events, patient actions)
- [x] Basic form builder (admin) and patient-scoped form submissions
- [x] Account locking after failed attempts (5 failures → 15 minute lock)
- [x] Admin unlock for locked staff accounts (`POST /auth/users/:id/unlock`)
- [x] Staff deactivation (inactive accounts cannot sign in; retained for audit history)
- [x] Session inactivity timeout UI (15-minute warning modal + friendly sign-out)
- [x] Must-change-PIN gate (blocks clinical routes until personal PIN is set)
- [x] Account locking covered by login handler tests
- [x] Physical security documentation in ops binder (`docs/ops/physical-security-binder.md`)
  - [x] Safe vault credentials sheet — print once, fill secrets, lock in safe, delete digital copies (`docs/ops/safe-credentials-vault.md` + Help article)

### Patient Management
**Status**: ✅ Complete
- [x] Patient registration
- [x] Demographics management
- [x] Search and filters (name/email search, gender, DOB range, registered-since date)
- [x] Patient detail and edit views
- [x] Clinical chart notes (timeline; author edits/deletes own notes; admin can delete any; audit logged)
- [x] Document uploads (PDF, images, text — AES-GCM encrypted in DB, cascade with patient, audit logged)
- [x] Exercise log (PT plans + session entries with sets/reps/resistance/difficulty; audit logged)
- [x] Audit logging (create, view, update, delete)

### User Interface
**Status**: ✅ Complete
- [x] Basic layout
- [x] Error boundary implementation
- [x] Sidebar navigation with grouped sections (Clinical / Admin / Account)
- [x] Breadcrumb trail on desktop
- [x] Dark/Light theme toggle (follows system by default)
- [x] Collapsible sidebar for smaller clinic displays (not a mobile product)
- [x] Clinic workstation layout polish (1024px+ LAN browsers) — patient list, chart two-column, sticky headers
- [x] Accessibility compliance (keyboard/screen reader at front desk) — skip link, focus styles, row keyboard nav, ARIA labels
- In-app Help System
  - [x] Help center hub with categories and search (`/help`)
  - [x] Article pages with steps, tables, and alerts (`/help/:articleId`)
  - [x] Contextual “Help for this page” drawer (header ? button)
  - [x] Role-aware content (staff vs administrator articles)
  - [x] Printable backup/restore checklists
  - [x] Editable clinic emergency contacts (local browser storage)
  - [x] Contextual inline tips on key EMR screens (dismissible, links to guides)
  - [x] Admin setup, health check, security, and form report articles
  - [x] Keyboard and accessibility tips article
- Admin Dashboard
  - [x] System status overview (`/admin`)
  - [x] Backup age warnings
  - [x] Task notifications (action-needed alerts on admin dashboard)
  - [x] Simple action buttons (backup trigger from UI)
  - [x] Administrator guide panel (setup checklist + all admin tasks)
  - [x] Configurable required patient chart fields (admin Clinic settings tab; staff form + API validation)
- Repository documentation (`docs/guide/`)
  - [x] Documentation hub and evaluating guide
  - [x] Clinic network requirements
  - [x] Administrator guides (setup, daily ops, staff, backup, troubleshooting)
  - [x] Staff guides (sign-in, patient care)
  - [x] Three-layer documentation model documented in IN-APP-HELP.md

### Backup and Recovery
**Status**: ✅ Complete
- USB Backup System
  - [x] Backup command writes `pococlinic-backup-*.tar.gz` with manifest
  - [x] Restore command with `-confirm` safety gate
  - [x] Admin UI “Backup now”
  - [x] Admin UI restore with confirmation (`POST /admin/restore`)
  - [x] Local Backup Helper (`run-ops-helper.bat`, localhost :9090) with guided wizards
  - [x] Patient document encrypted blobs included in the SQLite backup file; legacy disk files still optional under `documents/`
  - [x] Exercise plans and log entries included in the SQLite backup file
  - [x] Restore reloads encrypted document rows with the database; legacy `documents/` extracted to DOCUMENTS_DIR when present
  - [x] Backup verification (checksums + encrypted-content or legacy disk integrity)
  - [x] Admin UI backup checksum verification (`POST /admin/backups/verify`)
  - [x] Recovery testing (quarterly drill reminders)
- Physical Tracking
  - [x] Printable daily checklist (helper UI + runbook)
  - [x] USB drive labels (admin dashboard rotation panel + print)
- Recovery Procedures
  - [x] Step-by-step recovery guide (helper restore wizard + runbook)
  - [x] Data integrity verification (backup manifest summary + live document checks)
  - [x] System health checks (`GET /admin/health-check`, admin dashboard)

### API Layer
**Status**: ✅ Complete

Conventions for agents and contributors: `.cursor/rules/backend-api-conventions.mdc`

- [x] RESTful endpoints under `/api/v1`
- [x] Centralized error handling (`pkg/httperr` + panic recovery middleware)
- [x] Authentication required on all business routes (patients, forms, staff admin)
- [x] Role-based admin gates (`RequireRole(RoleAdmin)`)
- [x] Rate limiting and security headers
- [x] Auth middleware responses migrated to `httperr`

### Audit Logging
**Status**: ✅ Complete
- [x] User action tracking (auth, patients, staff, badge events, chart notes create/edit/delete, exercise log)
- [x] Admin system-wide audit log viewer with filters
- [x] Backup and restore operations logged
- [x] System event logging (health check, audit purge events in audit log)
- [x] HIPAA compliance checks (automated checklist on admin dashboard)
- [x] Log rotation (`cmd/audit-purge`, `AUDIT_RETENTION_DAYS`, ops checklist)
- [x] Log analysis tools (admin activity summary report, audit log filters, admin-only audit CSV download + archive reminder)

### Reporting
**Status**: ✅ Complete (custom ad-hoc builder deferred)
- [x] Form template submission report (paginated, date filter, export all matching rows as CSV)
- [x] Basic patient census report (totals, gender breakdown, 30-day registrations — in-app only)
- [x] Statistical analysis (clinic activity summary — sign-ins, patient work, backups by period)
- [x] Staff activity report (per-user sign-ins, failures, chart views)
- [x] Scheduled maintenance templates (backup + audit purge cron snippets in release + admin dashboard panel)
- [ ] Custom report builder (adhoc query UI — deferred / out of scope for v1)

### Physical Documentation
**Status**: ✅ Complete
- Repository guides (`docs/guide/`)
  - [x] Evaluating PocoClinic (prospective clinics)
  - [x] Clinic network requirements
  - [x] Administrator setup and operations guides
  - [x] Staff daily workflow guides
  - [x] Ops binder inserts (physical security, safe vault credentials, emergency procedures, security audit checklist, monthly testing checklist)
  - [x] Printable ops binder (`docs/ops/binder/`) — A emergency, B weekly/monthly, C system help (no daily logs)
  - [x] Production deploy systemd units (`scripts/deploy/`)
  - [x] Clinic server install tree (`clinic/server/` — install.sh, cron, bin wrappers; [ADR-0016](../adr/0016-clinic-runtime-packaging.md))
  - [x] Release build script (`scripts/build-release.mjs`, `build-release.bat`) — copies `clinic/server/` into tarball
  - [x] Production static UI served from main binary (`ENV=production`, `STATIC_DIR`)
  - [x] Tools & scripts catalog (`docs/ops/tools-and-scripts.md`)
  - [x] microSD / Pi deployment boundary (`docs/deploy/DEPLOYMENT-BOUNDARY.md`)
  - [x] Devices hub at repo root (`devices/`) with Raspberry Pi folder (docs + future ops code)
  - [x] Interactive setup wizard (`setup.bat` / `scripts/setup.mjs`) — Raspberry Pi + local dev paths
  - [x] Raspberry Pi beginner setup walkthrough (`devices/raspberry-pi/setup.md`)
  - [x] Raspberry Pi print-me binder (`devices/raspberry-pi/binder/`)
  - [x] Standalone binder printer utility (`binder-printer/`, `clinic/workstation/`) — not part of EMR

  - [x] Web-published help site (static HTML from markdown — `node scripts/build-help-site.mjs`)

### AI Assistance
**Status**: 📝 Planned — **local-only on clinic hardware**; no cloud LLM; deferred until core EMR and ops are complete
- Local LLM Integration
  - [ ] Lightweight model selection
    - Primary Option: Llama-2-7b-chat-q4 (GGUF format)
      - ~4GB RAM usage
      - ~4GB disk space
      - CPU-only operation possible
      - Good balance of capability vs resource usage
    - Backup Option: GPT4All-J-6B (GGML format)
      - ~3GB RAM usage
      - ~3.7GB disk space
      - Optimized for CPU
    - Minimum System Requirements:
      - 8GB RAM total
      - 10GB free disk space
      - x86_64 CPU with AVX2 support
  - [ ] Offline-first operation
    - [ ] Local model file management
    - [ ] Versioned model updates
    - [ ] Fallback to rule-based responses
  - [ ] Resource usage monitoring
    - [ ] RAM usage limits
    - [ ] CPU usage throttling
    - [ ] Disk space monitoring
  - [ ] Model updates management
    - [ ] Manual update process
    - [ ] Integrity verification
    - [ ] Rollback capability
- Task Templates
  - [ ] Common procedure guidance
    - [ ] Pre-defined prompt templates
    - [ ] Context-aware responses
    - [ ] Step-by-step instructions
  - [ ] Form filling assistance
    - [ ] Field explanation
    - [ ] Data validation suggestions
    - [ ] Common value recommendations
  - [ ] Documentation lookup
    - [ ] Natural language queries
    - [ ] Context-based search
    - [ ] Quick reference generation
  - [ ] Simple report generation
    - [ ] Template-based outputs
    - [ ] Data summarization
    - [ ] Format consistency
- System Integration
  - [ ] Context-aware help
    - [x] Route-based suggestions in help drawer
    - [x] Inline tips on patient list, chart, admin, PIN, and staff pages
  - [ ] Natural language search
    - [ ] Query optimization
    - [ ] Result ranking
    - [ ] Search scope control
  - [ ] Task completion suggestions
    - [ ] Next step recommendations
    - [ ] Common patterns recognition
    - [ ] Error prevention hints
  - [ ] Error explanation assistance
    - [ ] Plain language translations
    - [ ] Resolution suggestions
    - [ ] Prevention tips
- Privacy & Security
  - [ ] Local-only processing
    - [ ] Network isolation verification
    - [ ] Data flow monitoring
    - [ ] Cache management
  - [ ] PHI/PII awareness
    - [ ] Pattern recognition
    - [ ] Data masking
    - [ ] Sanitization rules
  - [ ] Audit logging of AI usage
    - [ ] Query logging
    - [ ] Response tracking
    - [ ] Usage patterns
  - [ ] Configurable usage limits
    - [ ] Rate limiting
    - [ ] Token quotas
    - [ ] Access controls

## Quality Assurance
**Status**: ✅ Complete
- [x] Monthly testing procedures — checklist + in-app monthly reminder (localStorage tracker on admin dashboard)
- [x] Backup verification (admin Verify + quarterly drill tracker)
- [x] Security audit checklist (`docs/ops/security-audit-checklist.md`)
- [x] Performance review (admin performance snapshot — memory, goroutines, DB pool)
- [x] Documentation updates (three-layer help + `docs/guide/`) 