# PocoClinic Administrator Runbook

Print this guide and keep it in the ops binder **Section B** (weekly/monthly processes).  
Do **not** keep a daily paper backup log — use the Admin dashboard status and USB rotation.

**Binder assembly:** [ops binder cover & print order](./binder/README.md)

**Contributor reference:** [Tools & scripts catalog](./tools-and-scripts.md) · [clinic/server/](../../clinic/server/README.md) · [Deployment boundary](../deploy/DEPLOYMENT-BOUNDARY.md)

## Clinic network (read first)

PocoClinic is designed for a **private local network with no internet for daily care**:

- The **server** (Raspberry Pi or small PC) and **staff workstations** stay on the **clinic LAN only**.
- Staff use **clinic-owned desktop or laptop browsers** — not mobile apps, not access from home over the internet.
- The clinic network should **not depend on internet** for sign-in, charting, or USB backup.
- **Security still matters on the LAN** (badge + PIN, audit log, session timeouts) — see [NETWORK-AND-SECURITY.md](../NETWORK-AND-SECURITY.md).

**Firewall expectation:** Do not port-forward the PocoClinic server to the public internet. Block outbound internet from the server if your clinic policy requires full isolation.

**Backup helper:** Runs at **http://127.0.0.1:9090** on the server PC only — not exposed on the Wi‑Fi for convenience.

## Daily backup checklist

| Step | Action | Initials / date |
|------|--------|-----------------|
| 1 | Confirm PocoClinic is running and staff can sign in | |
| 2 | Insert today's labeled USB drive | |
| 3 | Run backup (Admin → **Backup now**, Backup Helper, or `sudo /opt/pococlinic/bin/backup` on the server) | |
| 4 | Confirm green success / new file in backup list | |
| 5 | Remove USB and store in locked location | |
| 6 | Confirm dashboard backup status is green | |

**Rotation:** use at least two USB drives (e.g. Mon/Wed/Fri). Never rely on a single drive.

**Warning thresholds** (shown on Admin dashboard):

- Green: backup within 24 hours
- Yellow: 24–48 hours
- Red: over 48 hours — take a backup before end of day

## Backup from the admin UI

1. Sign in as administrator (`/login/admin` if needed).
2. Open **Admin** in the navigation bar.
3. Click **Backup now**.
4. Verify the new file appears in the backup table.

## Backup Helper (recommended for daily use)

A friendly guide runs **only on the clinic computer** — not on the network:

1. On the server PC, open **http://127.0.0.1:9090** (or run `run-ops-helper.bat`).
2. Follow the **Daily backup** wizard step by step.
3. Print the checklist from **Print checklist** and keep it in the ops binder.

Start the helper:

```powershell
build-ops-helper.bat   # first time only
run-ops-helper.bat
```

## Backup from the command line

**Development (Windows checkout):**

```powershell
$env:DATABASE_URL = "./data/pococlinic.db"
.\backup.bat
```

**Production (Pi / Linux server):**

```bash
sudo /opt/pococlinic/bin/backup
```

Backups are written to `BACKUP_DIR` (default `./backups` in dev; `/var/lib/pococlinic/backups` on Pi).

## Restore (disaster recovery)

**Warning:** Restore replaces all database contents — patients, forms, staff, sessions, and audit history.

### From the Backup Helper (recommended)

1. Sign out all staff except one administrator.
2. On **this computer**, open **http://127.0.0.1:9090**.
3. Open **Restore** and follow the guided steps (type `RESTORE` to confirm).
4. After restore, verify login and spot-check patient count.

### From the admin UI

1. Sign out all staff except one administrator.
2. Open **Admin** → backup table → **Restore** on the chosen file.
3. Read the confirmation dialog carefully and confirm.
4. After restore, verify login and spot-check patient count.

### From the command line

**Development:**

```powershell
.\restore.bat
# or a specific file:
.\restore.bat pococlinic-backup-YYYYMMDD-HHMMSS.tar.gz
```

**Production (Pi):**

```bash
sudo /opt/pococlinic/bin/restore --confirm
# or: sudo /opt/pococlinic/bin/restore --confirm --file pococlinic-backup-YYYYMMDD-HHMMSS.tar.gz
```

## After restore verification

- [ ] Admin can sign in
- [ ] Patient list loads and count looks reasonable
- [ ] Staff badges still work (reissue if restore was from an old backup)
- [ ] Take a fresh backup immediately after successful restore
- [ ] Note date, backup file used, and initials on restore log sheet

## Quarterly restore drill

1. Use a **test PC** or spare machine — not the production clinic PC during hours.
2. Install from a release tarball (see [clinic/server/](../../clinic/server/README.md)) or dev checkout with a copy of a USB backup.
3. Restore the latest USB backup with `--confirm`.
4. Verify login and patient count.
5. Sign drill log in binder.

## Common issues

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| Data disappears after restart | `DATABASE_URL` not set | Set URL; dev: `migrate.bat`; Pi: `sudo /opt/pococlinic/bin/migrate`; restart service |
| Backup button disabled | In-memory mode | Configure persistent database |
| Restore fails checksum | Corrupted or edited backup file | Use a different USB / earlier backup |
| Staff locked out | 5 failed PIN attempts | Wait 15 minutes or admin resets via staff management |

## Emergency contacts

| Role | Name | Phone |
|------|------|-------|
| Primary admin | __________________ | __________________ |
| Backup admin | __________________ | __________________ |
| IT support | __________________ | __________________ |

## Related documentation

- **Safe vault (secrets)** — [safe-credentials-vault.md](./safe-credentials-vault.md) — print, fill, lock in safe, delete digital copies
- **Documentation guide** — [docs/guide/README.md](../guide/README.md) (evaluating, admin, staff — web-publishable)
- **Physical security** — [physical-security-binder.md](./physical-security-binder.md)
- **Emergency procedures** — [emergency-procedures.md](./emergency-procedures.md)
- **Security audit checklist** — [security-audit-checklist.md](./security-audit-checklist.md)
- **Production deploy** — [clinic/server/](../../clinic/server/README.md) · [scripts/deploy/README.md](../../scripts/deploy/README.md) · [deploy hub](../deploy/README.md)
- **In-app Help center** — Help in the sidebar or **?** in the header (search, backup checklists, contacts, troubleshooting)
- **Administrator guide** — Admin dashboard → Administrator guide (setup checklist + all admin tasks)
- [NETWORK-AND-SECURITY.md](../NETWORK-AND-SECURITY.md) — LAN-only scope, no mobile, defensive security
- [VISION.md](../VISION.md) — product intent
- [FEATURES.md](../FEATURES.md) — product feature status
- [ADR-0014 — Private LAN-only deployment](../../adr/0014-private-lan-only-deployment.md)
- [ADR-0012](../../adr/0012-backup-and-recovery-ops.md) — backup architecture decision
- [ADR-0015](../../adr/0015-clinic-database-engine.md) — clinic DB engine (proposed)
- [.env.example](../../.env.example) — development environment variables
