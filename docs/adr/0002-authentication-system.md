# ADR-0002: Authentication System Design

## Status
Accepted

## Context
The PocoClinic EMR system requires a secure but user-friendly authentication system. As a medical system handling sensitive patient data, security is paramount, but it must also be practical for healthcare workers to use in their daily work.

Key requirements:
- Must be HIPAA compliant
- Must be secure against common attack vectors
- Must be user-friendly for healthcare workers
- Must support audit logging of authentication attempts
- Must allow for quick access in emergency situations
- Must support session management and timeout

## Decision
We will implement a two-factor authentication system consisting of:

1. **Primary Authentication**: 64-bit key
   - Generated uniquely for each user
   - High entropy for security
   - Can be stored securely in password managers
   - Used for initial authentication

2. **Secondary Authentication**: 4-digit PIN
   - Personal to each user
   - Quick to enter in emergency situations
   - Changed periodically
   - Rate-limited to prevent brute force attacks

3. **Implementation Details**:
   - Keys will be stored using Argon2id for password hashing
   - PINs will be stored separately with their own salt
   - Sessions will be managed using JWTs with short expiration
   - Refresh tokens will be stored in a secure HTTP-only cookie
   - Failed attempts will be rate-limited and logged
   - Automatic session timeout after 15 minutes of inactivity
   - All authentication attempts will be audit logged

4. **Security Measures**:
   - Rate limiting on all authentication endpoints
   - Automatic account locking after multiple failed attempts
   - Secure session management
   - HTTPS required for all authentication requests
   - Audit logging of all authentication events
   - Regular security audits and penetration testing

## Consequences

### Positive
- Strong security with two-factor authentication
- User-friendly for daily healthcare operations
- Compliant with HIPAA requirements
- Clear audit trail for security events
- Support for emergency quick access
- Protection against common attack vectors

### Negative
- Additional complexity in the authentication flow
- Need to manage two separate credentials
- Requires secure delivery method for initial 64-bit keys
- Additional overhead in user management
- Need for periodic PIN changes

### Mitigations
- Clear documentation for users and administrators
- Simple self-service portal for PIN management
- Automated key rotation procedures
- Clear recovery procedures for lost credentials
- Regular security training for users 