# I4 — Hardening & TLS

**Installer binder · Section I4**

Run **after** server install (I3) and **network** (I2) are working over HTTP.

References: [Pi hardening](../../../devices/raspberry-pi/hardening.md) · [LAN TLS](../../../devices/raspberry-pi/tls-lan.md) · [ADR-0017](../../../adr/0017-raspberry-pi-hardening-and-lan-tls.md)

---

## 1. Harden the server

```bash
cd /opt/pococlinic
sudo LAN_CIDR=192.168.1.0/24 ADMIN_SSH_CIDR=192.168.1.0/24 ./scripts/harden-pi.sh
```

**Before running:** confirm SSH **key login** works — script disables password auth.

Adjust `LAN_CIDR` to match clinic Wi‑Fi subnet (Section I2).

---

## 2. Generate LAN certificates (offline)

```bash
sudo POCO_HOSTNAME=pococlinic.local POCO_LAN_IP=192.168.1.50 ./scripts/generate-lan-tls.sh
```

Output: `/etc/pococlinic/tls/` (`ca.crt`, `server.crt`, `server.key`)

---

## 3. Enable HTTPS — pick one path

### Path A — Caddy reverse proxy (recommended)

1. `sudo apt install caddy` (maintenance window)  
2. `sudo cp /opt/pococlinic/caddy/Caddyfile.example /etc/caddy/Caddyfile`  
3. Edit `/etc/pococlinic/env`:
   ```bash
   SERVER_HOST=127.0.0.1
   SERVER_PORT=8080
   ALLOWED_ORIGIN=https://pococlinic.local
   TRUSTED_PROXIES=127.0.0.1
   ```
4. `sudo systemctl enable --now caddy && sudo systemctl restart pococlinic`

### Path B — Go TLS directly

Edit `/etc/pococlinic/env`:

```bash
SERVER_TLS_CERT=/etc/pococlinic/tls/server.crt
SERVER_TLS_KEY=/etc/pococlinic/tls/server.key
SERVER_PORT=443
ALLOWED_ORIGIN=https://pococlinic.local
```

`sudo systemctl restart pococlinic`

---

## 4. Trust CA on staff devices

Copy `/etc/pococlinic/tls/ca.crt` to USB. On **each** clinic iPad/PC:

| OS | Action |
|----|--------|
| **Windows** | Import to Trusted Root Certification Authorities |
| **iPad** | Install profile → Settings → Certificate Trust → enable |
| **macOS** | Keychain → Always Trust |

Until trusted, browsers show a certificate warning.

---

## 5. Post-TLS cleanup

| Task | Done |
|------|------|
| Staff URL uses **https://** | ☐ |
| Remove UFW rule for :8080 when migration complete | ☐ |
| `OPS_HELPER_MAIN_APP_URL=https://pococlinic.local` | ☐ |
| Test login + cookie session on iPad | ☐ |

---

## TLS sign-off

**HTTPS verified from staff iPad:** ☐  
**Tester:** ____________ **Date:** ____________

---

**Next:** Section **I5** — Handoff checklist
