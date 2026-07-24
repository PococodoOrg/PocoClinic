# Clinic server install tree

Files in this directory are **copied into the release tarball** and installed under `/opt/pococlinic` on the clinic server (Raspberry Pi or small PC).

## On the deployed system

```
/opt/pococlinic/
├── pococlinic          # EMR binary (LAN :8080)
├── ops-helper          # localhost backup UI (:9090)
├── migrate             # migration CLI binary
├── backup              # backup CLI binary
├── restore             # restore CLI binary
├── audit-purge         # audit retention CLI binary
├── static/             # production EMR UI
├── bin/                # env-aware wrappers (prefer these in cron)
│   ├── migrate
│   ├── backup
│   ├── restore
│   └── audit-purge
├── cron/               # example /etc/cron.d drop-ins
├── scripts/            # harden-pi.sh, generate-lan-tls.sh
├── caddy/              # Caddyfile.example for LAN HTTPS
├── install.sh
└── env.template
```

Persistent data stays under `/var/lib/pococlinic/` and is **not** replaced on upgrade.

## Operator commands

Always prefer `bin/` wrappers — they source `/etc/pococlinic/env`:

```bash
sudo /opt/pococlinic/bin/migrate
sudo /opt/pococlinic/bin/backup
sudo /opt/pococlinic/bin/restore --confirm
sudo /opt/pococlinic/bin/audit-purge --days 365 --dry-run
```

## Install / upgrade

See [install.sh](./install.sh). Full flow: [docs/deploy/DEPLOYMENT-BOUNDARY.md](../../docs/deploy/DEPLOYMENT-BOUNDARY.md).
