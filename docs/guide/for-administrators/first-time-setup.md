# First-time clinic setup

Complete these steps once when PocoClinic is installed on the clinic server. Track progress in the EMR: **Admin dashboard → Administrator guide → First-time setup checklist**.

## 1. Prepare the server

- Install the OS and ensure the server has a static IP on the clinic LAN
- Create persistent directories (typical production paths under `/var/lib/pococlinic/`)
- Production env file: `/etc/pococlinic/env` (from [clinic/server/env.template](../../../clinic/server/env.template) via `install.sh`)

## 2. Configure environment

**Production:** edit `/etc/pococlinic/env` (created by [`clinic/server/install.sh`](../../../clinic/server/install.sh) from `env.template`).

**Development:** copy [`.env.example`](../../.env.example) / [`backend/.env.example`](../../../backend/.env.example).

| Variable | Purpose |
|----------|---------|
| `DATABASE_URL` | SQLite file path (e.g. `/var/lib/pococlinic/pococlinic.db`) |
| `BACKUP_DIR` | Backup tar directory (e.g. `/var/lib/pococlinic/backups`) |
| `DOCUMENT_ENCRYPTION_KEY` | Base64 32-byte AES key for document blobs (required in production) |
| `DOCUMENTS_DIR` | Legacy disk path / bootstrap parent (optional for new installs) |
| `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` | Session signing (production must be unique) |

## 3. Run migrations

**Development (Windows checkout):**

```powershell
$env:DATABASE_URL = "./data/pococlinic.db"
.\migrate.bat
```

Or from `backend/`: `go run ./cmd/migrate` with `DATABASE_URL` set.

**Production (Pi / Linux server after release tarball):**

```bash
sudo /opt/pococlinic/bin/migrate
```

See [clinic/server/install.sh](../../../clinic/server/install.sh) and [deploy hub](../../deploy/README.md).

Confirm no pending migrations appear on the admin dashboard.

## 4. Start services

**Development:**

```powershell
cd backend
go run ./cmd/main.go

cd frontend
npm run dev
```

**Production:** install from the release tarball ([clinic/server/](../../../clinic/server/README.md)), configure `/etc/pococlinic/env`, enable systemd units in [scripts/deploy/](../../../scripts/deploy/README.md):

```bash
sudo systemctl enable --now pococlinic pococlinic-ops-helper
```

Start the Backup Helper on the server for daily ops (localhost only):

```powershell
run-ops-helper.bat   # dev checkout — forwards to clinic/dev-windows/
# Opens http://127.0.0.1:9090 on the machine where it runs
```

## 5. Bootstrap administrator sign-in

1. Open `/login/admin` in a browser on the LAN
2. Sign in with administrator email, key, and PIN
3. Confirm the admin dashboard loads with **Database** storage mode (not in-memory)

## 6. Safe vault sheet (print → lock → delete)

Secrets must not stay on the server forever. Complete once at setup:

1. Print the blank [Safe vault credentials sheet](../../ops/safe-credentials-vault.md) (or Help → **Safe vault credentials** → Print)
2. Handwrite admin email/key/PIN, `JWT_*` secrets, `DOCUMENT_ENCRYPTION_KEY`, and `DATABASE_URL`
3. Add emergency contacts and safe key holders; sign with a witness
4. Place the filled sheet in the **safe** — not loose in a desk binder
5. Delete `data/bootstrap-admin-once.txt` and any digital copies of the filled values

## 6b. Assemble binders

Print with [`print-binder.bat`](../../../print-binder.bat) (forwards to [`clinic/workstation/`](../../../clinic/workstation/README.md)) or [binder-printer](../../../binder-printer/README.md). See the [binders hub](../../binders/README.md).

**Installer (deploy tech):**

- [**Installer binder**](../../binders/installer/README.md) — I1 prerequisites through I5 handoff (network, install, TLS)
- [**Raspberry Pi device binder**](../../../devices/raspberry-pi/binder/README.md) (if Pi server)

**Clinic (after handoff):**

- [**Site operations binder**](../../ops/binder/README.md) — A emergency · B weekly/monthly · C system help (admin desk)
- **Safe vault** — blank template → fill → lock in safe → delete digital copies

Store binders closed; open for those uses only.

## 7. Create staff and print badges

1. **Staff → New staff member** for each employee
2. Assign initial PIN; staff must change it on first sign-in
3. Open each staff record → **Print badge**
4. Hand badges in person; never email QR codes

## 8. Optional: form templates

**Forms → New template** — build intake, consent, or vitals forms staff will use on patient charts.

## 9. Backup readiness

1. **Admin dashboard → USB rotation** — label Mon/Wed/Fri drives; print labels
2. **Backup now** — create first backup file
3. **Verify** — confirm checksums and document integrity
4. Copy backup to labeled USB; store securely

## 10. Health check & contacts

1. **System health check → Run check** — all items green or understood
2. **Help → Clinic emergency contacts** — fill primary admin, backup admin, IT

## 11. Train staff

Point staff to:

- [Staff guide](../for-staff/README.md)
- In-app Help → Sign in with badge and PIN
- In-app Help → Find and open patients

## Verification checklist

- [ ] Admin dashboard shows Database + Connected
- [ ] No pending migrations
- [ ] Safe vault sheet printed, locked in safe, digital copies deleted
- [ ] Installer handoff signed (I5) if using deploy tech
- [ ] Site operations binder at admin desk (sections A / B / C)
- [ ] At least one staff account with printed badge
- [ ] First backup verified
- [ ] Health check passed
- [ ] Emergency contacts filled in

## Related

- [Clinic network requirements](../clinic-network.md)
- [Backup & recovery](./backup-and-recovery.md)
- [Deploy hub](../../deploy/README.md) · [clinic/server/](../../../clinic/server/README.md)
- In-app article: `admin-setup`
