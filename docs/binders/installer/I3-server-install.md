# I3 — Server install

**Installer binder · Section I3**

Install PocoClinic from the **release tarball** on Linux (Raspberry Pi or small PC).

Full reference: [clinic/server/](../../../clinic/server/README.md) · [Deploy hub](../../deploy/README.md)

---

## 1. Extract and install

```bash
tar -xzf pococlinic-*-linux-arm64.tar.gz
cd pococlinic-*-linux-arm64
sudo ./install.sh
```

Installs to `/opt/pococlinic`, creates `pococlinic` user, copies systemd units.

| Path | Purpose |
|------|---------|
| `/opt/pococlinic/` | Binaries + static UI |
| `/etc/pococlinic/env` | Secrets and config |
| `/var/lib/pococlinic/` | Database, backups, documents |

---

## 2. Configure environment

Edit `/etc/pococlinic/env` (from `env.template`):

| Variable | Example |
|----------|---------|
| `DATABASE_URL` | `/var/lib/pococlinic/pococlinic.db` |
| `BACKUP_DIR` | `/var/lib/pococlinic/backups` |
| `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` | Unique 32+ chars each |
| `DOCUMENT_ENCRYPTION_KEY` | `openssl rand -base64 32` |
| `ALLOWED_ORIGIN` | `https://pococlinic.local` (after TLS) |

Generate secrets **on the server**; copy to safe vault at handoff.

---

## 3. Migrate database

```bash
sudo /opt/pococlinic/bin/migrate
```

---

## 4. Enable services

```bash
sudo systemctl enable --now pococlinic pococlinic-ops-helper
sudo systemctl status pococlinic pococlinic-ops-helper
```

Ops helper: **http://127.0.0.1:9090** on the server only.

---

## 5. Bootstrap admin

1. On LAN browser: `/login/admin`  
2. Sign in with bootstrap credentials (delete `bootstrap-admin-once.txt` after)  
3. Admin dashboard → **Database connected**, no pending migrations  

---

## 6. Optional cron

```bash
sudo cp /opt/pococlinic/cron/pococlinic-backup /etc/cron.d/
sudo cp /opt/pococlinic/cron/pococlinic-audit-purge /etc/cron.d/
```

---

## Install sign-off

| Check | Done |
|-------|------|
| `systemctl status pococlinic` active | ☐ |
| `/health` returns OK on LAN (HTTP OK before TLS) | ☐ |
| First staff user created + badge printed (or scheduled at handoff) | ☐ |

---

**Next:** Section **I4** — Hardening & TLS
