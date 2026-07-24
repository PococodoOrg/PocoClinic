# ADR-0011: Employee Badge Authentication (QR + PIN)

## Status
Accepted

## Context

ADR-0002 defines two-factor authentication using a **64-bit key** plus **4-digit PIN**. For clinic staff, the key must not be typed or memorized—it should be carried as a **physical employee badge** with a QR code, paired with a PIN the employee knows.

Goals:

- **Strong:** two factors, hashed at rest, rate limiting, lockout, audit trail (per ADR-0002)
- **Simple:** scan badge, enter PIN, done— suitable for front desk and clinical workflows
- **Offline:** badge validation is local; no IdP or SMS

Staff login uses badge key + PIN in the UI and API. A separate **admin bootstrap** path (email + key + PIN) remains for first-time setup, user provisioning, and badge printing before staff badges exist.

## Decision

1. **Staff-facing authentication**
   - **Factor 1 — Badge:** QR code on employee badge encodes the user's unique key material (same secret currently generated at user creation)
   - **Factor 2 — PIN:** 4 digits, entered on keyboard or touchscreen after scan
   - Routine login UI: **no email field**; identity is resolved from the scanned badge

2. **Badge lifecycle (admin)**
   - Admin creates user → system generates key → print/issue badge (QR) + communicate initial PIN securely
   - Lost badge → admin revokes/reissues badge (new key); PIN may be reset
   - Badges are clinic property; QR payload is treated as secret (not patient PHI)

3. **Technical mapping (unchanged from ADR-0002)**
   - Key → Argon2id hash + salt in `users.key_*`
   - PIN → separate Argon2id hash + salt in `users.pin_*`
   - Sessions: short JWT access token + HttpOnly refresh cookie; 15-minute inactivity timeout
   - All auth events audit-logged

4. **Login API**
   - Staff: `POST /api/v1/auth/login/staff` with `{ "key": "<from QR>", "pin": "1234" }`
   - Admin/bootstrap: `POST /api/v1/auth/login/admin` with `{ "email", "key", "pin" }` for initial setup and break-glass access
   - PIN change: `POST /api/v1/auth/change-pin` with `{ "currentPin", "newPin" }` (authenticated user)
   - Badge reissue: `POST /api/v1/auth/users/:id/reissue-badge` (admin only; returns new key once)
   - Frontend routes: `/login` (staff badge scan), `/login/admin` (administrator sign in), `/users` (staff management), `/account/pin` (PIN self-change)
   - USB scanner or paste populates badge key; focus moves to PIN field

5. **QR format**
   - v1: QR contains the raw key string (base64url) issued at user creation
   - Future: signed payload with user id + key id to support rotation without changing QR size dramatically

## Consequences

### Positive
- Aligns security model with physical clinic workflows
- Faster login than typing email and secret
- Clear story for HIPAA access control (“badge + PIN”)

### Negative
- Badge printing and reissue processes required
- Lost badge = lost factor until reissued
- Camera/scanner compatibility on clinic devices must be tested

### Mitigations
- Printable admin runbook for badge issuance
- Lockout + audit on failed PIN attempts
- Break-glass admin account documented for emergencies

## Relationship to ADR-0002

This ADR **does not replace** ADR-0002. It specifies the **operator-facing** form of the same two-factor model. Implementation should migrate UI and primary login flow to badge + PIN while keeping ADR-0002 cryptographic and session rules.
