# Raspberry Pi hardening

Baseline security for the **clinic server** Pi on a private LAN. This complements application auth (badge + PIN) — it does **not** replace it.

← [Pi index](./README.md) · [LAN TLS setup](./tls-lan.md) · [Physical security](../../docs/ops/physical-security-binder.md)

---

## Goals

| Goal | How |
|------|-----|
| No public internet exposure | Router: **no port forwarding** to the Pi |
| Limit LAN attack surface | Firewall: SSH + HTTPS (+ temporary HTTP 8080) from clinic subnet only |
| Protect secrets | `/etc/pococlinic/env` mode `640`; TLS keys mode `600` |
| Safer remote admin | SSH: key-only, no root login |
| Encrypted staff traffic | [LAN TLS](./tls-lan.md) with clinic-local CA |

---

## Automated script

Ships in the release tarball at `/opt/pococlinic/scripts/harden-pi.sh`.

```bash
cd /opt/pococlinic
sudo LAN_CIDR=192.168.1.0/24 ADMIN_SSH_CIDR=192.168.1.0/24 ./scripts/harden-pi.sh
```

**Before running:** ensure SSH key login works — the script disables password authentication.

### What the script does

- Writes `sysctl` hardening (no IP forwarding, syncookies)
- SSH: `PermitRootLogin no`, `PasswordAuthentication no`, `MaxAuthTries 3`
- **UFW**: default deny inbound; allow SSH from admin CIDR; allow 443 (and 8080 during TLS migration) from LAN
- Fixes ownership on `/var/lib/pococlinic` and TLS directory permissions
- Enables **fail2ban** if installed

### What it does not do

- Configure your clinic router or Wi‑Fi AP
- Install OS security updates (schedule separately in a maintenance window)
- Replace physical access controls (locked cabinet)

---

## Manual checklist (ops binder)

| Item | Action |
|------|--------|
| Hostname / DNS | `pococlinic.local` → Pi LAN IP (router DNS or local DNS) |
| Static IP | DHCP reservation or static config — record in safe vault |
| Default accounts | Change/disable Pi OS default password; use dedicated `pococlinic` service user (install.sh) |
| SSH keys | Admin laptops use keys; disable password auth after verify |
| HTTPS | Run [generate-lan-tls.sh](../../clinic/server/scripts/generate-lan-tls.sh); trust CA on staff devices |
| Bind app behind proxy | `SERVER_HOST=127.0.0.1` when using Caddy |
| Ops helper | Confirm `OPS_HELPER_HOST=127.0.0.1` — never `0.0.0.0` |
| USB backups | Encrypted clinic process; drives stored locked |
| Power | Official PSU; UPS optional for brownout-prone sites |
| Display kiosk | Chromium kiosk only for backup UI — not general browsing |

---

## Ongoing maintenance

- **Monthly:** verify HTTPS cert expiry (`openssl x509 -in /etc/pococlinic/tls/server.crt -noout -dates`)
- **Quarterly:** restore drill + review audit log
- **As needed:** OS patches during a scheduled window (internet optional at clinic discretion)

---

## Related

- [LAN TLS](./tls-lan.md)
- [Hardware notes](./hardware.md)
- [NETWORK-AND-SECURITY.md](../../docs/NETWORK-AND-SECURITY.md)
- [ADR-0017 — Pi hardening and LAN TLS](../../adr/0017-raspberry-pi-hardening-and-lan-tls.md)
