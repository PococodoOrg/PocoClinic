#!/usr/bin/env python3
"""Verify every persona can log in and matches expected email/role."""

from __future__ import annotations

import argparse
import sys
import time

import requests

from pococlinic_common import PersonaCredentials, load_dotenv, persona_credentials, refuse_non_loopback


def try_login(base: str, cred: PersonaCredentials) -> tuple[int, dict | None, str]:
    session = requests.Session()
    if cred.login == "admin":
        path = "/api/v1/auth/login/admin"
        body = {"email": cred.email, "key": cred.key, "pin": cred.pin}
    else:
        path = "/api/v1/auth/login/staff"
        body = {"key": cred.key, "pin": cred.pin}

    for attempt in range(3):
        resp = session.post(f"{base.rstrip('/')}{path}", json=body, timeout=15)
        if resp.status_code == 429 and attempt < 2:
            time.sleep(13)
            continue
        break
    if resp.status_code != 200:
        return resp.status_code, None, resp.text[:300]

    user = (resp.json() or {}).get("user") or {}
    me = session.get(f"{base.rstrip('/')}/api/v1/auth/me", timeout=15)
    if me.status_code != 200:
        return me.status_code, user, f"/auth/me failed: {me.text[:200]}"
    return 200, me.json() or user, ""


def main() -> int:
    load_dotenv()
    parser = argparse.ArgumentParser(description="Verify load-test persona logins")
    parser.add_argument("--host", default="http://127.0.0.1:8080")
    parser.add_argument(
        "personas",
        nargs="*",
        default=["admin", "manager", "clinician", "clerk"],
        help="Persona names to verify (default: all four)",
    )
    args = parser.parse_args()
    refuse_non_loopback(args.host)

    failed = False
    print(f"Verifying personas against {args.host}\n")
    for name in args.personas:
        cred = persona_credentials(name)
        if not cred:
            print(f"[FAIL] {name}: missing {name.upper()}_KEY / _PIN in personas.env")
            failed = True
            continue

        status, user, err = try_login(args.host, cred)
        if status != 200 or not user:
            print(f"[FAIL] {name}: HTTP {status} — {err or 'login failed'}")
            print(f"       login={cred.login} email={cred.email or '(staff)'} pin=****")
            failed = True
            continue

        if user.get("mustChangePin"):
            print(f"[FAIL] {name}: mustChangePin still true — run seed_personas.py --force")
            failed = True
            continue

        email = user.get("email", "")
        role = user.get("role", "")
        if cred.email and email.lower() != cred.email.lower():
            print(f"[FAIL] {name}: key belongs to {email!r}, expected {cred.email!r}")
            print("       Re-run: python seed_personas.py --force")
            failed = True
            continue
        if cred.api_role and role != cred.api_role:
            print(f"[FAIL] {name}: role is {role!r}, expected {cred.api_role!r}")
            failed = True
            continue

        print(f"[OK]   {name}: {email} ({role}) via {cred.login} login")

    print()
    if failed:
        print("Fix: fill setup.env, then  python seed_personas.py --force")
        return 1
    print("All persona logins OK.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
