# PocoClinic — Safe vault credentials sheet

**CONFIDENTIAL — clinic property**

Print this sheet **once**, fill it in **by hand** (or paste secrets only onto this printed copy), place it in the **ops binder**, and lock the binder in a **safe**. Then **delete every digital copy** of the filled values (files, screenshots, email, chat, sticky notes on the server).

| | |
|--|--|
| Clinic name | ________________________________ |
| Server hostname / LAN IP | ________________________________ |
| Sheet prepared by | ________________________________ |
| Date prepared | ____________ |
| Next review date | ____________ (recommend every 6 months) |

---

## 1. Print → lock → delete checklist

Complete in order. Initials required.

| Step | Action | Initials |
|------|--------|----------|
| 1 | Print **this blank template** from the server or admin workstation (do not email the filled sheet) | |
| 2 | Fill **Section 2** with live secrets from first-time setup (see sources below) | |
| 3 | Fill **Section 3** emergency contacts and key-holder names | |
| 4 | Read **Section 4** break-glass steps with a second administrator present | |
| 5 | Place printed sheet in the ops binder; lock binder in the safe | |
| 6 | Delete `data/bootstrap-admin-once.txt` (or equivalent) from the server | |
| 7 | Remove filled secrets from any `.env` printouts, Notepad files, Photos, Downloads | |
| 8 | Confirm no filled copy remains in email, USB scratch folders, or cloud sync | |
| 9 | Sign below | |

**We confirm the only lasting copy of these credentials is the printed sheet in the safe.**

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Primary administrator | | | |
| Witness (second admin or owner) | | | |

---

## 2. Credentials (fill by hand — never type into shared docs)

### 2A. Bootstrap / break-glass administrator

Source: first boot file `bootstrap-admin-once.txt` (usually under `./data/`), or Staff → admin user after setup.

| Item | Value (handwritten) |
|------|---------------------|
| Admin login URL | `/login/admin` on clinic LAN (e.g. `http://____:____/login/admin`) |
| Email | ________________________________________________ |
| Badge / admin key | ________________________________________________ |
| Initial PIN (change on first login) | ________ (default at first boot is often `0000`) |
| Current admin PIN hint (optional — prefer memory) | ________ or “memorized / not written” |

**After first successful admin login:** change PIN immediately. Do not leave `0000` in production.

### 2B. Application secrets (server `.env`)

These unlock sessions and decrypt patient document blobs. Without them, a restored database may be unreadable.

| Variable | Value (handwritten) | Notes |
|----------|---------------------|--------|
| `JWT_ACCESS_SECRET` | | Min 32 characters in production |
| `JWT_REFRESH_SECRET` | | Must differ from access secret |
| `DOCUMENT_ENCRYPTION_KEY` | | Base64 of 32 random bytes (`openssl rand -base64 32`) |
| `DATABASE_URL` | | Includes DB user/password if any |

Optional (if your install uses them):

| Variable | Value (handwritten) |
|----------|---------------------|
| OS / SSH password for clinic server | |
| Database admin password (if separate from URL) | |
| Wi‑Fi / router admin (clinic LAN gear) | |
| Safe combination / who knows it | |

### 2C. Where files live on the server

| Path / setting | Value (handwritten) |
|----------------|---------------------|
| App / install directory | |
| `BACKUP_DIR` | |
| `DOCUMENTS_DIR` (legacy files only) | |
| Database file (`DATABASE_URL`) | |

---

## 3. People and contacts

| Role | Name | Phone |
|------|------|-------|
| Primary administrator | | |
| Backup administrator | | |
| Clinic owner / medical director | | |
| IT / vendor support (if any) | | |

**Key holders for the safe** (names only — no combinations here unless clinic policy requires it):

1. ________________________________
2. ________________________________
3. ________________________________

---

## 4. Break-glass processes (use only when needed)

Keep this page with the credentials. Follow checklists; do not improvise on a bad day.

### 4A. Administrator locked out / lost badge

1. Open `/login/admin` on a clinic workstation.
2. Use **email + admin key + PIN** from Section 2A.
3. Open **Staff** → reissue badge for the affected user; print QR in person.
4. If the admin account itself is compromised, create a new admin, then deactivate or delete the old one after confirming access.
5. Log date, reason, and initials on the restore/incident log in the binder.

### 4B. Rebuild server / new hardware

1. Install PocoClinic and database per first-time setup guide.
2. Set `.env` values **exactly** from Section 2B (especially `DOCUMENT_ENCRYPTION_KEY` and JWT secrets).
3. Run migrations — dev: `migrate.bat`; **Pi:** `sudo /opt/pococlinic/bin/migrate` (see [clinic/server/](../../clinic/server/README.md)).
4. Restore latest verified USB backup (`RESTORE` confirmation in Backup Helper, or Admin → Restore).
5. Verify: admin sign-in, patient count, open one document, staff badge login.
6. Take a **fresh backup** immediately; update USB labels; sign restore log.

**Critical:** Restoring encrypted documents without the matching `DOCUMENT_ENCRYPTION_KEY` leaves files unreadable.

### 4C. Suspected credential leak

1. Change JWT secrets and restart the app (all sessions end).
2. Generate a **new** document encryption key only if you also re-encrypt or re-upload documents — do not rotate casually; seek guidance before changing `DOCUMENT_ENCRYPTION_KEY` on a live clinic.
3. Reissue all staff badges; require PIN changes for admins.
4. Reprint an updated vault sheet; destroy the old printed sheet (shred); update safe.
5. Note incident in ops binder; review audit log for unusual access.

### 4D. Backup reminder

- USB backup procedures: **Administrator runbook** (ops binder Section B).
- Backup Helper (server only): `http://127.0.0.1:9090`
- After any restore: verify login → patient list → one document → new backup → sign Section A incident log.

---

## 5. Procedures live in the ops binder (not only in the safe)

Keep these in the ops binder ([assembly](./binder/README.md)):

- Section A — Emergency procedures and paper fallback
- Section B — Weekly/monthly processes and backup runbook
- Section C — System help

The **safe vault sheet (this document, filled)** holds secrets only.

---

## 6. Review and re-print log

| Date | Reason (review / rotate secrets / staff change) | Prepared by | Old sheet destroyed? |
|------|--------------------------------------------------|-------------|----------------------|
| | | | ☐ Yes |
| | | | ☐ Yes |
| | | | ☐ Yes |
| | | | ☐ Yes |

---

## Digital hygiene (after printing)

- [ ] Delete `bootstrap-admin-once.txt` from the server
- [ ] Do not store filled secrets in the EMR, Help center, or browser notes
- [ ] Do not photograph the filled sheet
- [ ] Do not email or chat the filled values
- [ ] Keep blank templates in the repo/docs only — never commit a filled copy to git

**Template location in the product docs:** `docs/ops/safe-credentials-vault.md`
