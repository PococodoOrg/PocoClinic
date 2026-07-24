# Daily administrator operations

A typical week for a clinic administrator using PocoClinic.

## Each morning (5 minutes)

1. Open **Admin dashboard**
2. Read **Action needed** alerts — resolve anything red or yellow first
3. Confirm **Backup** status is green (within 24 hours)
4. Glance at **Security overview** — locked accounts or default PIN counts

## USB backup (daily or per rotation schedule)

1. Insert **today's labeled USB drive** (see USB rotation panel)
2. **Backup now** on the dashboard, or use Backup Helper at `http://127.0.0.1:9090` on the server
3. Confirm new file in backup table
4. **Verify** after copying to USB (recommended weekly at minimum)
5. Remove USB; store locked

## Staff support

| Situation | Action |
|-----------|--------|
| Forgot PIN | Reset via Staff → edit account; if locked, use **Unlock account** on the staff detail page (or wait ~15 min) |
| Lost badge | Staff detail → Reissue badge → print immediately |
| Locked out | Staff detail → **Unlock account**, or wait for lockout to expire |
| New hire | Create account → print badge → confirm PIN change on first login |
| Leaving employee | Deactivate or delete staff record per clinic policy |

## Weekly

- Review **Audit log** for unusual patient access or failed sign-ins
- Confirm **Maintenance reminders** — restore drill not overdue
- Export form reports if clinic needs CSV for quality review

## Monthly

- Run **System health check** after any server maintenance
- Re-read backup age trends; adjust USB rotation if backups are slipping

## Related

- [Ops binder · Administrator how-to](../../ops/binder/C3-administrator-how-to.md) — printable Section C
- [Ops binder · Weekly/monthly](../../ops/binder/B1-periodic-processes.md) — Section B
- [Backup & recovery](./backup-and-recovery.md)
- [Staff & badges](./staff-and-badges.md)
- In-app article: `admin-dashboard`
