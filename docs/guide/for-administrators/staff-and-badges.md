# Staff accounts and badges

## Create a staff account

1. **Staff → New staff member**
2. Enter name, email, and assign a role (`admin`, `doctor`, `nurse`, or `staff`)
3. Set an initial PIN (staff must change it on first sign-in)
4. Print the badge QR immediately and give it to the employee

## Daily sign-in (staff)

1. Staff opens `/login`
2. Scan badge QR (or paste key) → enter PIN
3. After inactivity, the session ends and they sign in again

Administrators can also use `/login/admin` (email + key + PIN) for bootstrap / break-glass.

## Lost or compromised badge

1. Open the staff detail page
2. **Reissue badge** — old QR stops working
3. Print the new badge and destroy the old printout

## Locked accounts

After too many failed PIN attempts, the account locks for about 15 minutes.

- Admin can review locked count on **Security overview**
- On the staff detail page, use **Unlock account** to clear the lock immediately (or wait for the timeout / help them reset PIN)

## Roles

Assignable clinic roles (permissions differ slightly; admins manage staff and backups):

| Role | Typical access |
|------|----------------|
| `staff` | Chart read, document upload, own PIN change |
| `nurse` / `doctor` | Chart work including notes, forms, demographics edits (doctor can delete documents) |
| `admin` | Staff & badges, forms admin, audit, backup, unlock, full clinic ops |

Sign-in path is still badge + PIN at `/login` for all of these (except bootstrap at `/login/admin`).

## Related

- [Daily operations](./daily-operations.md)
- [Sign-in and sessions (staff)](../for-staff/sign-in-and-sessions.md)
