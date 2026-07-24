# Administrator troubleshooting

| Symptom | Likely cause | What to try |
|---------|--------------|-------------|
| Patient data gone after restart | `DATABASE_URL` not set | Set URL; dev: `migrate.bat`; Pi: `sudo /opt/pococlinic/bin/migrate`; restart service |
| Dashboard shows in-memory mode | No database connection | Check `DATABASE_URL` |
| Backup button disabled | In-memory mode | Fix database configuration |
| Document upload fails | Wrong/missing encryption key or DB error | Check `DOCUMENT_ENCRYPTION_KEY`; run health check |
| Restore fails verify | Corrupted USB file | Try earlier backup drive |
| Staff locked out | 5 failed PINs | **Unlock account** on staff detail, or wait ~15 minutes |
| Session keeps ending | 15 min inactivity | Normal — badge + PIN again |
| Pending migrations alert | Schema not applied | Dev: `migrate.bat`; Pi: `sudo /opt/pococlinic/bin/migrate` |
| Document integrity warning | Missing encrypted blob or legacy disk file | Restore from backup or re-upload |

## Health check first

**Admin dashboard → System health check → Run check**

Interpreting results:

- **Database critical** — fix `DATABASE_URL` and SQLite first
- **Document integrity critical** — DB records point to missing files; may need restore
- **Document integrity warning** — missing encrypted content or orphan legacy files; investigate before deleting

## When to restore

Only when:

- Server disk failure
- Database corruption
- Catastrophic operator error

Always sign out staff first. Prefer the most recent **verified** backup.

## Getting more help

- In-app **Help → Troubleshooting → Common problems**
- **Help → Clinic emergency contacts**
- Repository [administrator runbook](../ops/administrator-runbook.md)

## Related

- In-app article: `troubleshooting-common`
