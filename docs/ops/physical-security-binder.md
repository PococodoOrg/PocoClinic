# Physical security — ops binder insert

Print this page for ops binder **Section B** (periodic processes). Secrets: [Safe vault credentials](./safe-credentials-vault.md).

## Server and network

| Control | Clinic responsibility |
|---------|----------------------|
| Server location | Locked room or cabinet; limited key holders |
| Network | Private LAN only — no port forwarding to the internet |
| Wi‑Fi | Clinic staff network separate from guest Wi‑Fi when possible |
| Backup helper | http://127.0.0.1:9090 — localhost on server only |

## Workstations

| Control | Clinic responsibility |
|---------|----------------------|
| Front desk PCs | Clinic-owned; auto-lock OS when unattended |
| Shared login | Each staff member uses **own badge + PIN** — no shared accounts |
| Sign out | Staff use **Sign out** when leaving a shared workstation |
| Badges | Treat like keys; report loss immediately for reissue |

## USB backups

| Control | Clinic responsibility |
|---------|----------------------|
| Labeled drives | Mon / Wed / Fri rotation (or clinic schedule) |
| Storage | Locked drawer or safe when not in use |
| Credentials sheet | [Safe vault sheet](./safe-credentials-vault.md) — filled copy in safe only |
| Off-site | At least one copy stored away from the server room periodically |

## Badge and PIN policy (summary)

- Administrator prints badge in person — never email QR codes
- Reissue lost badges immediately (old QR stops working)
- Staff choose a private PIN on first sign-in
- After 5 failed PIN attempts, account locks for 15 minutes

## Quarterly review (administrator)

- [ ] Run system health check on admin dashboard
- [ ] Verify latest USB backup
- [ ] Complete restore drill on spare PC
- [ ] Review audit log for unusual failed sign-ins
- [ ] Confirm emergency contacts in Help center are current

## Related

- [Emergency procedures](./emergency-procedures.md)
- [Security audit checklist](./security-audit-checklist.md)
- [Administrator runbook](./administrator-runbook.md)
- [Safe vault credentials sheet](./safe-credentials-vault.md)
- [Ops binder assembly](./binder/README.md)
