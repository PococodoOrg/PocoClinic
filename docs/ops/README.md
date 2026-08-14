# Operations documentation

Printable runbooks, checklists, and tooling for clinic administrators and contributors.

← [Docs hub](../README.md) · [Repo README](../../README.md)

## Physical binders

**All binders:** [**binders hub**](../binders/README.md)

| Binder | Audience |
|--------|----------|
| [**Installer**](../binders/installer/README.md) | Deploy tech — network, install, TLS, handoff |
| [**Site operations**](./binder/README.md) | Clinic desk — emergency, periodic, how-to |
| [**Pi device**](../../devices/raspberry-pi/binder/README.md) | Pi cabinet — hardware / kiosk |
| [**Safe vault**](./safe-credentials-vault.md) | Safe only — secrets |

**Print:** [**Site operations binder**](./binder/README.md) (sections A · B · C below).

Three uses only (binder stays closed otherwise):

| Section | When |
|---------|------|
| **A — Emergency** | System down, paper care, restore / incident |
| **B — Weekly / monthly** | Scheduled processes and sign-off checklists (no daily logs) |
| **C — System help** | Device-independent how-to (same on every workstation) |

Secrets: [Safe vault credentials](./safe-credentials-vault.md) → filled copy in the **safe**.

## Other ops docs

| Document | Audience | Purpose |
|----------|----------|---------|
| [Administrator runbook](./administrator-runbook.md) | Admins | Backup & restore steps (binder Section B) |
| [Emergency procedures](./emergency-procedures.md) | Admins | Outages (binder Section A) |
| [Tools & scripts](./tools-and-scripts.md) | Contributors & admins | CLI / batch / systemd catalog |
| [Clinic runtime (`clinic/`)](../../clinic/README.md) | Image builders & admins | Post-deploy install tree shipped in tarball |
| [Deployment boundary](../deploy/DEPLOYMENT-BOUNDARY.md) | Image builders | microSD / Pi packaging |
| [Pi touchscreen backup](../../devices/raspberry-pi/touchscreen-backup.md) | Admins | Localhost backup kiosk on Pi display |
| [**Pi beginner setup**](../../devices/raspberry-pi/setup.md) | First-time installers | Flash SD → install → first admin login |
| [**Interactive setup**](../../setup.bat) | Anyone after clone | `setup.bat` / `./setup.sh` — pick Raspberry Pi or local dev |
| [Devices · Raspberry Pi](../../devices/raspberry-pi/README.md) | Admins / installers | Pi hardware + device ops code |

## Related

- [Documentation guide](../guide/README.md)
- [Network & security](../NETWORK-AND-SECURITY.md)
- [ADR-0014 — Private LAN-only deployment](../../adr/0014-private-lan-only-deployment.md)
