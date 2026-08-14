# Set up PocoClinic on a Raspberry Pi

**Audience:** anyone installing PocoClinic for the first time — no Linux experience required.

**Goal:** by the end, you open a browser on another computer on the same network, reach the clinic server, and sign in as administrator.

← [Raspberry Pi index](./README.md) · [Devices hub](../README.md)

---

## Fastest start (interactive)

From a clone of this repo on your Windows/Mac/Linux PC:

```bat
setup.bat
```

Or: `./setup.sh` · `node scripts/setup.mjs`

Choose **Raspberry Pi clinic server**. The wizard checks tools, can build the ARM64 package, generates first-boot secrets, and pauses on each install step. This page is the same path written out in full (use it if you prefer reading, or if the wizard is unavailable).

---

## What you will have when finished

| Thing | What it means |
|-------|----------------|
| Raspberry Pi running 24/7 | The **clinic server** (not a staff tablet) |
| Staff URL (first bring-up) | `http://YOUR-PI-IP:8080` |
| Admin sign-in | Email + badge key + PIN from a one-time file on the Pi |
| Database | SQLite file on the Pi under `/var/lib/pococlinic/` |

Staff charting happens on **other** clinic devices (iPad / desktop browsers) on the same Wi‑Fi or Ethernet. The Pi is the server box.

**Time:** about 1–2 hours for a first install (plus download/build time).

**Important:** production login cookies are normally **HTTPS-only**. This guide uses a temporary `COOKIE_SECURE=false` setting so you can verify the install over HTTP. **Turn on LAN TLS and remove that setting before real patient use** ([tls-lan.md](./tls-lan.md)).

---

## Checklist (print or tick as you go)

- [ ] **0** — Have the release package (`.tar.gz` file)
- [ ] **1** — Flash microSD with Raspberry Pi OS
- [ ] **2** — Boot the Pi and log in (keyboard or SSH)
- [ ] **3** — Copy the package onto the Pi
- [ ] **4** — Run `install.sh`
- [ ] **5** — Edit secrets in `/etc/pococlinic/env`
- [ ] **6** — Migrate database + start services
- [ ] **7** — Open the EMR in a browser and sign in as admin
- [ ] **8** — (Before patients) Enable HTTPS and remove `COOKIE_SECURE=false`

---

## What you need

### Hardware

| Item | Notes |
|------|--------|
| **Raspberry Pi 5** (preferred) or **Pi 4** | **8 GB RAM** preferred; 4 GB works for trying it out |
| Official power supply | Cheap chargers cause SD corruption |
| microSD card | 32 GB+; A2 / endurance-rated if you can |
| Ethernet cable | Strongly recommended for the server |
| Monitor + HDMI + USB keyboard | For first boot **or** use SSH from another PC |
| Another computer on the same network | Windows/Mac/Linux — to open the EMR in a browser |
| USB stick (optional) | Easiest way to copy the package to the Pi |

Optional later: USB SSD for the database, official touchscreen for backup-only kiosk.

### Software on your Windows PC (to build the package)

If someone already gave you a file named like `pococlinic-1.0.0-linux-arm64.tar.gz`, skip this list and go to [Step 0B](#0b--you-already-have-the-tarball).

Otherwise install on the PC that has this git repo:

| Tool | Why |
|------|-----|
| [Node.js](https://nodejs.org/) (LTS) | Builds the web UI and the release script |
| [Go](https://go.dev/dl/) (1.22+) | Builds the server binaries |
| Git | You already have the repo |

---

## Step 0 — Get the release package

You need one file:

```text
pococlinic-<version>-linux-arm64.tar.gz
```

### 0A — Build it on your Windows PC

1. Open a terminal in the **repo root** (`D:\dev\PocoClinic` or wherever you cloned it).
2. Run:

```bat
build-release.bat
```

3. Wait until it finishes (several minutes the first time).
4. Find the file under:

```text
dist\pococlinic-<version>-linux-arm64.tar.gz
```

Copy that file to a USB stick, or keep the path handy for `scp` later.

### 0B — You already have the tarball

Put it on a USB stick, or note the path on your PC. Continue to Step 1.

---

## Step 1 — Flash Raspberry Pi OS onto the microSD

1. On your Windows PC, download and install **[Raspberry Pi Imager](https://www.raspberrypi.com/software/)**.
2. Insert the microSD card (USB adapter is fine).
3. Open Imager and choose:
   - **Raspberry Pi Device** — your model (e.g. Pi 5)
   - **Operating System** — **Raspberry Pi OS (64-bit)**  
     Lite is fine if you are comfortable with SSH only. Desktop is easier with a monitor.
   - **Storage** — your microSD card
4. Click the **gear** (or **Edit settings**) and set **all** of these **before** writing:

| Setting | Suggested value |
|---------|-----------------|
| Hostname | `pococlinic` (other PCs can often reach it as `pococlinic.local`) |
| Enable SSH | **On** (password authentication is OK for first setup) |
| Username | pick one you will remember (example: `clinic`) |
| Password | strong password — write it down |
| Wi‑Fi | Only if you will not use Ethernet; use the **clinic** network, not guest |
| Locale / keyboard | Your country |

5. Write the image. When Imager says it is done, eject the card safely.

---

## Step 2 — First boot and find the Pi

1. Insert the microSD into the Pi.
2. Plug in **Ethernet** (preferred), then power.
3. Wait about 1–2 minutes for first boot.

### Option A — Monitor + keyboard on the Pi

Log in with the username/password you set in Imager. You are done with this step when you see a terminal prompt (Desktop: open **Terminal**).

### Option B — SSH from Windows (no monitor)

1. On Windows, open **PowerShell**.
2. Try:

```powershell
ssh clinic@pococlinic.local
```

Replace `clinic` with the username you chose. Accept the fingerprint (`yes`), then enter the password.

3. If `pococlinic.local` does not resolve, find the Pi’s IP:
   - Check your router’s “connected devices” list for hostname `pococlinic`, **or**
   - Temporarily plug in a monitor and run: `hostname -I`

Then connect with the IP:

```powershell
ssh clinic@192.168.1.50
```

(Use the real IP you found.)

**Write down the Pi’s LAN IP.** You need it in Step 5 and Step 7. Example: `192.168.1.50`.

On the Pi, confirm the IP anytime with:

```bash
hostname -I
```

---

## Step 3 — Copy the package onto the Pi

Pick one method.

### Method A — USB stick (easiest)

1. On Windows, copy `pococlinic-*-linux-arm64.tar.gz` onto a USB stick.
2. Plug the stick into the Pi.
3. On the Pi, find the USB mount (Desktop often auto-mounts under `/media/clinic/…`). In Terminal:

```bash
ls /media/$USER/
```

4. Copy the tarball to your home folder (adjust the path to match what you see):

```bash
cp /media/$USER/*/pococlinic-*-linux-arm64.tar.gz ~/
cd ~
ls *.tar.gz
```

You should see the `.tar.gz` file listed.

### Method B — Copy over the network (`scp`)

From **PowerShell on Windows** (not on the Pi), from the folder that contains the tarball:

```powershell
scp .\dist\pococlinic-*-linux-arm64.tar.gz clinic@pococlinic.local:~/
```

Use your username and IP if `.local` fails.

---

## Step 4 — Install PocoClinic

On the Pi (SSH or Terminal):

```bash
cd ~
tar -xzf pococlinic-*-linux-arm64.tar.gz
cd pococlinic-*-linux-arm64
sudo ./install.sh
```

**What “good” looks like:** the script ends with `Install complete.` and prints **Next steps**.

This creates:

| Path | Purpose |
|------|---------|
| `/opt/pococlinic/` | Application |
| `/etc/pococlinic/env` | Config + secrets (you edit this next) |
| `/var/lib/pococlinic/` | Database and backups |

---

## Step 5 — Configure secrets (required — do not skip)

The stock `REPLACE-ME` placeholders will **not** start the app. You must set real secrets.

1. Open the config file:

```bash
sudo nano /etc/pococlinic/env
```

2. Generate three secret values (run these on the Pi, one at a time). **Copy each output** onto paper — you need them in the file and later in the safe vault:

```bash
openssl rand -base64 48
openssl rand -base64 48
openssl rand -base64 32
```

3. In `nano`, replace these lines (paste your generated values; keep the two JWT values **different**):

```bash
JWT_ACCESS_SECRET=paste-first-openssl-output-here
JWT_REFRESH_SECRET=paste-second-openssl-output-here
DOCUMENT_ENCRYPTION_KEY=paste-third-openssl-output-here
```

4. Set `ALLOWED_ORIGIN` and `OPS_HELPER_MAIN_APP_URL` to the **exact** URL you will type in the browser (use your Pi IP):

```bash
ALLOWED_ORIGIN=http://192.168.1.50:8080
OPS_HELPER_MAIN_APP_URL=http://192.168.1.50:8080
```

Replace `192.168.1.50` with **your** Pi IP from Step 2.  
If you will use the hostname instead, use that exact origin, e.g. `http://pococlinic.local:8080`.

5. Keep temporary HTTP login working (required for this first-bring-up path):

```bash
COOKIE_SECURE=false
```

(The template already includes this. Leave it until Step 8 / TLS.)

6. Leave these as-is for a first install:

```bash
DATABASE_URL=/var/lib/pococlinic/pococlinic.db
BACKUP_DIR=/var/lib/pococlinic/backups
DOCUMENTS_DIR=/var/lib/pococlinic/documents
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

7. Save and exit nano: `Ctrl+O`, Enter, then `Ctrl+X`.

**Important:** write the three secrets on paper (or print the [safe vault sheet](../../docs/ops/safe-credentials-vault.md) later). If you lose `DOCUMENT_ENCRYPTION_KEY`, uploaded documents cannot be decrypted after a restore.

---

## Step 6 — Database migrate and start services

Still on the Pi:

```bash
sudo /opt/pococlinic/bin/migrate
sudo systemctl enable --now pococlinic pococlinic-ops-helper
sudo systemctl status pococlinic --no-pager
```

**What “good” looks like:** status shows `active (running)` in green.

If it failed, jump to [Troubleshooting](#troubleshooting).

Quick health check from the Pi:

```bash
curl -s http://127.0.0.1:8080/health
```

You should see JSON that looks healthy (not a connection error).

Optional daily backup cron:

```bash
sudo cp /opt/pococlinic/cron/pococlinic-backup /etc/cron.d/
```

---

## Step 7 — Open the EMR and sign in as admin

1. On **another computer** on the same network, open a browser.
2. Go to:

```text
http://YOUR-PI-IP:8080/health
```

You should see a health response. Then open:

```text
http://YOUR-PI-IP:8080/login/admin
```

3. On the Pi, read the one-time admin credentials:

```bash
sudo cat /var/lib/pococlinic/bootstrap-admin-once.txt
```

You will see:

- **Email** (usually `admin@pococlinic.local`)
- **Temporary PIN** (`0000` — you must change it on first login)
- **Badge key** (long string — treat like a password)

4. On the admin login page fill in:
   - **Email**
   - **Access key** (the badge key from the file)
   - **PIN** (`0000`, then change when prompted)
5. Confirm the admin dashboard loads and storage shows **Database** (not in-memory).
6. After you can sign in reliably:

```bash
sudo rm /var/lib/pococlinic/bootstrap-admin-once.txt
```

---

## You are done (first bring-up)

PocoClinic is running on the LAN over HTTP for setup and testing.

### Step 8 — Before real patients (required)

| When | What | Guide |
|------|------|--------|
| **Before real patients** | HTTPS + trust CA on staff devices; **delete** `COOKIE_SECURE=false` (or set `true`) | [tls-lan.md](./tls-lan.md) |
| Before real patients | Firewall / SSH hardening | [hardening.md](./hardening.md) |
| Before go-live | Safe vault sheet + binders + staff badges | [First-time setup](../../docs/guide/for-administrators/first-time-setup.md) |
| Optional | Touchscreen backup kiosk | [touchscreen-backup.md](./touchscreen-backup.md) |

After HTTPS is on, set `ALLOWED_ORIGIN` (and `OPS_HELPER_MAIN_APP_URL`) to match, for example `https://pococlinic.local`, then:

```bash
sudo systemctl restart pococlinic
```

---

## Troubleshooting

| Problem | What to try |
|---------|-------------|
| Cannot SSH / `pococlinic.local` fails | Use the numeric IP from the router; confirm Ethernet link lights; re-check Imager SSH settings and re-flash if needed |
| `install.sh` says “Run as root” | Use `sudo ./install.sh` |
| Service fails immediately | `sudo journalctl -u pococlinic -n 50 --no-pager` — almost always `REPLACE-ME` secrets still in place or bad `DOCUMENT_ENCRYPTION_KEY` |
| Browser cannot connect | Same Wi‑Fi/Ethernet subnet? Guest Wi‑Fi often blocks device-to-device. Ping the Pi IP from the PC |
| Page loads but login fails / session lost | Confirm `COOKIE_SECURE=false` for HTTP; after TLS, use `https://` and remove that override. Also check `ALLOWED_ORIGIN` matches the browser address **exactly** |
| `/health` works on Pi but not from PC | Firewall or wrong IP; confirm `SERVER_HOST=0.0.0.0` |
| No bootstrap file | Service may not have started cleanly; fix env, restart, then `sudo ls -la /var/lib/pococlinic/` |
| `migrate` errors | Confirm `DATABASE_URL` path and that `/var/lib/pococlinic` exists and is owned by `pococlinic` |
| Go TLS on port 443 fails to bind | Reinstall units from a fresh tarball (`install.sh` copies them) — units include `CAP_NET_BIND_SERVICE`. Or use Caddy ([tls-lan.md](./tls-lan.md) Path A) |

Useful commands:

```bash
sudo systemctl status pococlinic pococlinic-ops-helper --no-pager
sudo journalctl -u pococlinic -n 80 --no-pager
sudo nano /etc/pococlinic/env
sudo systemctl restart pococlinic
```

Ops helper (backup UI) is **only** on the Pi itself: `http://127.0.0.1:9090` — it will not open from a tablet on Wi‑Fi. That is intentional.

---

## Related docs (after this guide)

| Doc | Role |
|-----|------|
| [hardware.md](./hardware.md) | Model / storage recommendations |
| [tls-lan.md](./tls-lan.md) | HTTPS before real patient use |
| [Installer binder](../../docs/binders/installer/README.md) | Printable I1–I5 for formal clinic installs |
| [clinic/server/](../../clinic/server/README.md) | What `install.sh` installs |
| [DEPLOYMENT-BOUNDARY.md](../../docs/deploy/DEPLOYMENT-BOUNDARY.md) | Tarball vs full microSD image product split |
| [ADR-0014](../../adr/0014-private-lan-only-deployment.md) | Private LAN product rules |
