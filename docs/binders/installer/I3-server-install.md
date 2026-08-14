# I3 — Server install

**Installer binder · Section I3**

Install PocoClinic from the **release tarball** on Linux (Raspberry Pi or small PC).

**First time / need hand-holding?** Use the full walkthrough (flash OS → login):  
[devices/raspberry-pi/setup.md](../../../devices/raspberry-pi/setup.md)

Short reference: [clinic/server/](../../../clinic/server/README.md) · [Deploy hub](../../deploy/README.md)

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

Edit `/etc/pococlinic/env` (from `env.template`). Stock placeholders will **not** start — replace them.

On the Pi, generate values:

```bash
openssl rand -base64 48   # JWT_ACCESS_SECRET
openssl rand -base64 48   # JWT_REFRESH_SECRET (must differ)
openssl rand -base64 32   # DOCUMENT_ENCRYPTION_KEY
```

| Variable | Example |
|----------|---------|
| `DATABASE_URL` | `/var/lib/pococlinic/pococlinic.db` |
| `BACKUP_DIR` | `/var/lib/pococlinic/backups` |
| `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` | From openssl above (32+ chars each, different) |
| `DOCUMENT_ENCRYPTION_KEY` | From openssl above (required in production) |
| `ALLOWED_ORIGIN` | First boot HTTP: `http://<pi-ip>:8080` — after TLS: `https://pococlinic.local` |
| `COOKIE_SECURE` | `false` for temporary HTTP bring-up only — **remove before patients** (see [tls-lan.md](../../../devices/raspberry-pi/tls-lan.md)) |

`ALLOWED_ORIGIN` must match the browser URL **exactly**. Copy secrets to the safe vault at handoff.

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

1. On LAN browser: `http://<pi-ip>:8080/login/admin` (HTTPS after Section I4)  
2. Read credentials: `sudo cat /var/lib/pococlinic/bootstrap-admin-once.txt`  
3. Sign in (PIN `0000` must be changed); then delete the bootstrap file  
4. Admin dashboard → **Database connected**, no pending migrations  

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
