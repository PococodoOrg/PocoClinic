# Emergency procedures

Quick reference for clinic administrators when PocoClinic or infrastructure fails during clinic hours.

## Priority order

1. **Patient safety** — continue care using paper fallback if EMR is unavailable
2. **Contain** — sign out all staff if you suspect unauthorized access
3. **Recover** — restore from USB backup on a spare machine if server is lost
4. **Document** — log what happened in the ops binder

## Scenario: Server will not start

| Step | Action |
|------|--------|
| 1 | Confirm power and network cables |
| 2 | Restart server; wait 2 minutes |
| 3 | Confirm the SQLite file is reachable \(DATABASE_URL\) |
| 4 | Restart PocoClinic backend and frontend (or systemd services) |
| 5 | If still down, use paper intake until restored from backup |

## Scenario: Database offline (admin dashboard red)

| Step | Action |
|------|--------|
| 1 | Do **not** register new patients in production until fixed |
| 2 | Set DATABASE_URL and start PocoClinic: `set DATABASE_URL=./data/pococlinic.db` (dev) or clinic systemd unit |
| 3 | Verify `DATABASE_URL` on server |
| 4 | Run migrations — dev: `migrate.bat`; **Pi:** `sudo /opt/pococlinic/bin/migrate` |

## Scenario: Suspected unauthorized access

| Step | Action |
|------|--------|
| 1 | Note time and workstation |
| 2 | Open **Audit log** — filter failed sign-ins and patient viewed events |
| 3 | Reissue badges if a badge may be compromised |
| 4 | Change administrator PIN if admin account may be affected |

## Scenario: Ransomware or disk failure

| Step | Action |
|------|--------|
| 1 | Disconnect affected server from LAN (pull network cable) |
| 2 | Do **not** pay ransom — restore from known-good USB backup |
| 3 | Wipe and reinstall OS on clean disk |
| 4 | Restore latest **verified** backup |
| 5 | Reissue all staff badges after restore |

## Scenario: Lost USB backup drive

| Step | Action |
|------|--------|
| 1 | Treat as confidential data loss — patient records may be on drive |
| 2 | Use remaining rotation drives for daily backups |
| 3 | Follow clinic privacy incident policy |
| 4 | Take fresh backup immediately on replacement drive |

## Emergency contacts

Maintain in **Help → Clinic emergency contacts** and paper copy:

| Role | Name | Phone |
|------|------|-------|
| Primary admin | __________________ | __________________ |
| Backup admin | __________________ | __________________ |
| IT support | __________________ | __________________ |

## Related

- [Ops binder · Section A](./binder/README.md) — emergency pack + [paper fallback](./binder/A2-paper-fallback.md) + [incident log](./binder/A3-incident-restore-log.md)
- [Restore after disaster](../guide/for-administrators/backup-and-recovery.md)
- In-app help: `restore-disaster`, `troubleshooting-common`
