# Security audit checklist

Use monthly or quarterly — print and sign in the ops binder.

**Clinic:** __________________ **Date:** __________________ **Initials:** ____

## Authentication & access

- [ ] All active staff have individual accounts (no shared logins)
- [ ] No staff remain on default PIN (check Security overview on admin dashboard)
- [ ] Lost badges reissued and old QRs invalidated
- [ ] Locked account count reviewed; failures investigated in audit log

## Backup & recovery

- [ ] Last backup within 24 hours (green on dashboard)
- [ ] Latest backup **Verified** successfully
- [ ] USB rotation labels current; drives stored locked
- [ ] Restore drill completed within 90 days (Maintenance reminders)

## System health

- [ ] Admin dashboard: Database connected, migrations up to date
- [ ] System health check: all green or documented exceptions
- [ ] Document storage writable; document integrity OK

## Network & physical

- [ ] Server in locked location; access log reviewed if applicable
- [ ] No port forwarding of PocoClinic to internet
- [ ] Workstations auto-lock when unattended
- [ ] Emergency contacts updated in Help center

## Audit log hygiene

- [ ] Audit log reviewed for unusual patient access or failed sign-ins
- [ ] Audit retention job scheduled if using `AUDIT_RETENTION_DAYS` (see below)

## Audit log rotation (optional)

If the clinic policy requires limiting audit history retention:

**Production (Pi / Linux server):**

```bash
sudo /opt/pococlinic/bin/audit-purge --days 365 --dry-run
sudo /opt/pococlinic/bin/audit-purge --days 365
```

**Development (Windows checkout):**

```powershell
cd backend
$env:AUDIT_RETENTION_DAYS = "365"
go run ./cmd/audit-purge --days 365 --dry-run
go run ./cmd/audit-purge --days 365
```

Schedule via cron on the server (`clinic/server/cron/pococlinic-audit-purge`) or Windows Task Scheduler during off-hours.

## Notes / follow-up actions

_____________________________________________________________________________

_____________________________________________________________________________

## Related

- [Physical security binder](./physical-security-binder.md)
- [Emergency procedures](./emergency-procedures.md)
