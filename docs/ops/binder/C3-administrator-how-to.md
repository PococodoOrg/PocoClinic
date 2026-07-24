# C3 — Administrator how-to

**Ops binder · Section C (System help)** · Device-independent.  
Scheduled checks: **Section B**. Outages: **Section A**. Secrets: **safe vault**.

**Admin login:** `/login/admin` · **URL:** `https://________________________`

---

## When you open the dashboard

1. Read **Action needed** — fix red/yellow first  
2. Confirm **Backup** is green (within 24 hours) — if not, run USB backup (Section B / runbook)  
3. Glance at **Security** — locked accounts, default PINs  

---

## Staff accounts and badges

**New hire:** Staff → New → name, email, role, initial PIN → **Print badge** in person → they **Change PIN**.  
**Lost badge:** Reissue → old QR dies immediately → print new.  
**Lockout:** Wait 15 minutes or help reset PIN — never share admin PIN.  
**Leaving:** Deactivate/delete per policy.

---

## Forms, audit, health

| Task | Where |
|------|--------|
| Form templates | **Forms** |
| Who did what | **Audit log** (weekly skim; monthly CSV — Section B) |
| Health / compliance | Dashboard checks |
| Migrations alert | Run `migrate.bat` (dev) or `sudo /opt/pococlinic/bin/migrate` (server) |

---

## Backups (summary)

Instructions: **Section B4** — [Administrator runbook](../administrator-runbook.md).  
Do **not** keep a daily paper log in this binder; use the dashboard color and USB rotation.

| Step | Reminder |
|------|----------|
| 1 | Today’s labeled USB |
| 2 | **Backup now** or server-only helper `http://127.0.0.1:9090` |
| 3 | Confirm file; optionally **Verify** |
| 4 | USB → locked storage |

---

## Restore (disasters)

1. Sign out all staff except one admin  
2. Backup Helper on server → confirm `RESTORE`  
3. Verify login, patient count, one document  
4. Fresh backup → sign **Section A3**  

Need **DOCUMENT_ENCRYPTION_KEY** from the safe vault sheet after rebuild.

---

## Where to look

| Need | Section |
|------|---------|
| System down / paper care | **A** |
| Weekly / monthly checklists | **B** |
| How the product works | **C** (this section) |
| Break-glass keys | **Safe** |
