# HTTPS on the clinic LAN (Raspberry Pi)

PocoClinic runs on a **private LAN** — not the public internet. Staff still benefit from **HTTPS** on Wi‑Fi (encrypted badge/PIN traffic, HttpOnly cookies with `Secure`, no browser “Not secure” warnings).

Public certificate authorities (Let’s Encrypt) are **out of scope**. Use a **clinic-local CA** generated on the server with `openssl` (works **offline**).

← [Pi index](./README.md) · [Hardening](./hardening.md) · [NETWORK-AND-SECURITY](../../docs/NETWORK-AND-SECURITY.md)

---

## Choose a path

| Path | Best for | Port staff use |
|------|----------|----------------|
| **A — Caddy reverse proxy** (recommended) | Pi production; keeps Go on loopback | `https://pococlinic.local` (:443) |
| **B — Go TLS directly** | Simple single-binary deploy | `https://pococlinic.local:443` |

Both use the same certificates from `generate-lan-tls.sh`.

---

## Step 1 — Generate certificates (both paths)

On the Pi after install:

```bash
cd /opt/pococlinic
sudo POCO_HOSTNAME=pococlinic.local POCO_LAN_IP=192.168.1.50 ./scripts/generate-lan-tls.sh
```

Files land in `/etc/pococlinic/tls/`:

| File | Purpose |
|------|---------|
| `ca.crt` | Install on **each staff iPad/PC** as a trusted root |
| `ca.key` | **Offline vault only** — can mint new server certs |
| `server.crt` / `server.key` | Server identity for HTTPS |

---

## Step 2A — Caddy reverse proxy (recommended)

1. Install Caddy on the Pi OS image (maintenance window; not required daily):
   ```bash
   sudo apt install caddy
   ```
2. Copy the example config:
   ```bash
   sudo cp /opt/pococlinic/caddy/Caddyfile.example /etc/caddy/Caddyfile
   sudo systemctl enable --now caddy
   ```
3. Edit `/etc/pococlinic/env`:
   ```bash
   SERVER_HOST=127.0.0.1
   SERVER_PORT=8080
   ALLOWED_ORIGIN=https://pococlinic.local
   TRUSTED_PROXIES=127.0.0.1
   OPS_HELPER_MAIN_APP_URL=https://pococlinic.local
   ```
4. Restart: `sudo systemctl restart pococlinic`
5. Open **`https://pococlinic.local`** from a staff browser.

The ops helper stays **`http://127.0.0.1:9090`** only (localhost).

---

## Step 2B — Go TLS directly

Edit `/etc/pococlinic/env`:

```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=443
SERVER_TLS_CERT=/etc/pococlinic/tls/server.crt
SERVER_TLS_KEY=/etc/pococlinic/tls/server.key
ALLOWED_ORIGIN=https://pococlinic.local
```

Restart: `sudo systemctl restart pococlinic`

Binding to port 443 may require `CAP_NET_BIND_SERVICE` on the binary or running behind a reverse proxy — **Path A avoids this**.

---

## Step 3 — Trust the CA on clinic devices

Until staff devices trust `ca.crt`, browsers show a certificate warning.

### Windows

1. Copy `ca.crt` from the server (USB or admin share).
2. `certmgr.msc` → Trusted Root Certification Authorities → Import `ca.crt`.

### iPad / iPhone (supervised or MDM preferred)

1. AirDrop / email / Apple Configurator profile with `ca.crt`.
2. Settings → General → About → Certificate Trust Settings → enable full trust for the PocoClinic CA.

### macOS

Double-click `ca.crt` → add to System keychain → always trust.

Record the rollout in the ops binder. Re-trust after CA rotation.

---

## Verification checklist

- [ ] Staff URL is **`https://`** (not plain HTTP)
- [ ] Padlock shows no warning after CA trust
- [ ] Sign-in works (cookies are `Secure` in production)
- [ ] `curl -k https://pococlinic.local/health` returns OK from a workstation
- [ ] Backup helper still **only** on `127.0.0.1:9090`
- [ ] Router has **no** port forward to the Pi

---

## Related

- [Pi hardening](./hardening.md)
- [Clinic network requirements](../../docs/guide/clinic-network.md)
- [ADR-0017 — Pi hardening and LAN TLS](../../adr/0017-raspberry-pi-hardening-and-lan-tls.md)
