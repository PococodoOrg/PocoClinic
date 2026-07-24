# C1 — What is PocoClinic?

**Ops binder · Section C (System help)** · Same on every workstation — do not add device-specific notes here.

## In one minute

PocoClinic is the clinic’s **electronic chart** for patients: demographics, notes, documents, forms, and exercise logs. It runs on a **server in the building**. Staff open it in a **web browser** on clinic computers or tablets.

| It is | It is not |
|-------|-----------|
| Private LAN only | Cloud / internet SaaS |
| Badge QR + PIN for staff | Passwords emailed from home |
| USB backups you control | Automatic off-site cloud backup |
| Audited access to charts | Anonymous shared logins |

## Who does what

| Role | Signs in at | Main jobs |
|------|-------------|-----------|
| **Staff** (doctor, nurse, aide, front desk) | `/login` — scan badge → PIN | Find patients, notes, documents, forms, exercise log |
| **Administrator** | `/login/admin` or admin-role badge | Staff & badges (incl. unlock), backups, health, audit, forms setup |

## Where things live

| Thing | Where |
|-------|--------|
| Patient data | SQLite file on the clinic server (`DATABASE_URL`) |
| Uploaded documents | Encrypted in the database (backed up with DB) |
| Backups | Files copied to labeled USB drives |
| Secrets (keys, encryption) | **Vault sheet in the safe** |

## Golden rules

1. **One person, one badge** — never share PIN or badge.  
2. **Sign out** on shared workstations when you leave.  
3. Keep Admin dashboard backup **green** (USB rotation) — status is on screen, not a daily paper log.  
4. **Lost badge** → tell admin immediately → reissue (old QR dies).  
5. **Server down** → open **Section A** (paper care + emergency procedures).  
6. **Do not** put the server on the public internet.

## Quick map of the screen

| Menu | Used for |
|------|----------|
| **Patients** | List, search, open chart, register |
| **Help** | How-to articles (works offline on the LAN) |
| **Change PIN** | Staff update their own PIN |
| **Dashboard** (admin) | Health, backups, setup checklist |
| **Staff** (admin) | Accounts, badges, reissue |
| **Forms** (admin) | Templates staff fill on charts |
| **Audit log** (admin) | Who did what |

## If you are stuck

1. Search **Help** in the sidebar (when the system is up).  
2. **Section C** sheets in this binder (when you want paper).  
3. **Section A** if the system is down.  
4. Call the primary administrator on the emergency contacts sheet (A3).  

**Clinic LAN address (fill in once for all devices):** `http://________________________`
