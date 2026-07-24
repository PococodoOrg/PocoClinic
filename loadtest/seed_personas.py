#!/usr/bin/env python3
"""
Create or refresh local load-test persona users via the normal admin API.

Uses setup.env for your real system administrator (not the load-test Admin persona).
Writes personas.env with verified keys. Use --force to delete and recreate loadtest users.

Loopback only. No backend test hooks.
"""

from __future__ import annotations

import argparse
import os
import sys
import time
from pathlib import Path

import requests

from pococlinic_common import is_loopback_host, load_dotenv, read_bootstrap_key

LOADTEST_DIR = Path(__file__).resolve().parent
PERSONAS_ENV = LOADTEST_DIR / "personas.env"
SETUP_ENV = LOADTEST_DIR / "setup.env"

PERSONA_DEFS = [
    {
        "env_prefix": "ADMIN",
        "email": "loadtest-admin@local.test",
        "name": "Loadtest Admin",
        "role": "admin",
        "login": "admin",
    },
    {
        "env_prefix": "MANAGER",
        "email": "loadtest-manager@local.test",
        "name": "Loadtest Manager",
        "role": "admin",
        "login": "admin",
    },
    {
        "env_prefix": "CLINICIAN",
        "email": "loadtest-clinician@local.test",
        "name": "Loadtest Clinician",
        "role": "doctor",
        "login": "staff",
    },
    {
        "env_prefix": "CLERK",
        "email": "loadtest-clerk@local.test",
        "name": "Loadtest Clerk",
        "role": "nurse",
        "login": "staff",
    },
]

LOADTEST_EMAILS = {spec["email"].lower() for spec in PERSONA_DEFS}
INITIAL_PIN = "0000"


class ApiClient:
    def __init__(self, base: str) -> None:
        self.base = base.rstrip("/")
        self.session = requests.Session()

    def post(self, path: str, json: dict | None = None) -> requests.Response:
        return self.session.post(f"{self.base}{path}", json=json or {}, timeout=30)

    def get(self, path: str) -> requests.Response:
        return self.session.get(f"{self.base}{path}", timeout=30)

    def delete(self, path: str) -> requests.Response:
        return self.session.delete(f"{self.base}{path}", timeout=30)


def load_setup_admin() -> tuple[str, str, str]:
    pin = os.environ.get("LOADTEST_PIN", "1234").strip()
    email = os.environ.get("SETUP_ADMIN_EMAIL", "").strip()
    key = os.environ.get("SETUP_ADMIN_KEY", "").strip()
    setup_pin = os.environ.get("SETUP_ADMIN_PIN", "").strip()

    if SETUP_ENV.is_file():
        for line in SETUP_ENV.read_text(encoding="utf-8").splitlines():
            line = line.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            k, _, v = line.partition("=")
            k, v = k.strip(), v.strip().strip("'").strip('"')
            if k == "SETUP_ADMIN_EMAIL" and not email:
                email = v
            elif k == "SETUP_ADMIN_KEY" and not key:
                key = v
            elif k == "SETUP_ADMIN_PIN" and not setup_pin:
                setup_pin = v
            elif k == "LOADTEST_PIN" and os.environ.get("LOADTEST_PIN") is None:
                pin = v

    if not key:
        key = read_bootstrap_key()
    if not email:
        email = "admin@pococlinic.local"

    if not key or not setup_pin or len(setup_pin) != 4:
        raise RuntimeError(
            "Set SETUP_ADMIN_EMAIL, SETUP_ADMIN_KEY, SETUP_ADMIN_PIN in setup.env "
            "(copy setup.env.example). These are your real system admin credentials."
        )
    if len(pin) != 4:
        raise RuntimeError("LOADTEST_PIN must be 4 digits")
    return email, key, setup_pin, pin


def login_setup_admin(client: ApiClient, email: str, key: str, pin: str) -> None:
    client.session.cookies.clear()
    r = client.post("/api/v1/auth/login/admin", {"email": email, "key": key, "pin": pin})
    if r.status_code == 200:
        return
    if r.status_code == 401 and pin != INITIAL_PIN:
        r2 = client.post(
            "/api/v1/auth/login/admin",
            {"email": email, "key": key, "pin": INITIAL_PIN},
        )
        if r2.status_code == 200:
            client.post("/api/v1/auth/change-pin", {"currentPin": INITIAL_PIN, "newPin": pin})
            return
    raise RuntimeError(f"setup admin login failed: {r.status_code} {r.text[:250]}")


def list_users(client: ApiClient) -> list[dict]:
    r = client.get("/api/v1/auth/users?page=1&pageSize=100")
    if r.status_code != 200:
        raise RuntimeError(f"list users failed: {r.status_code}")
    data = r.json() or {}
    return data.get("users") or data.get("items") or []


def delete_loadtest_users(client: ApiClient) -> None:
    for user in list_users(client):
        email = (user.get("email") or "").lower()
        if email not in LOADTEST_EMAILS:
            continue
        uid = user.get("id")
        if not uid:
            continue
        print(f"  deleting {email} ...")
        r = client.delete(f"/api/v1/auth/users/{uid}")
        if r.status_code not in (200, 204, 404):
            raise RuntimeError(f"delete {email} failed: {r.status_code} {r.text[:200]}")
        time.sleep(0.2)


def register_user(client: ApiClient, spec: dict) -> str:
    r = client.post(
        "/api/v1/auth/register",
        {"email": spec["email"], "name": spec["name"], "role": spec["role"]},
    )
    if r.status_code != 201:
        raise RuntimeError(f"register {spec['email']} failed: {r.status_code} {r.text[:250]}")
    key = (r.json() or {}).get("key") or ""
    if not key:
        raise RuntimeError(f"no badge key returned for {spec['email']}")
    return key


def set_persona_pin(spec: dict, key: str, target_pin: str, base: str) -> None:
    session = requests.Session()
    if spec["login"] == "admin":
        r = session.post(
            f"{base}/api/v1/auth/login/admin",
            json={"email": spec["email"], "key": key, "pin": INITIAL_PIN},
            timeout=15,
        )
    else:
        r = session.post(
            f"{base}/api/v1/auth/login/staff",
            json={"key": key, "pin": INITIAL_PIN},
            timeout=15,
        )
    if r.status_code == 429:
        time.sleep(13)
        return set_persona_pin(spec, key, target_pin, base)
    if r.status_code != 200:
        raise RuntimeError(
            f"initial login for {spec['email']} failed: {r.status_code} {r.text[:200]}"
        )
    r2 = session.post(
        f"{base}/api/v1/auth/change-pin",
        json={"currentPin": INITIAL_PIN, "newPin": target_pin},
        timeout=15,
    )
    if r2.status_code not in (200, 204):
        raise RuntimeError(f"change-pin for {spec['email']} failed: {r2.status_code}")
    me = session.get(f"{base}/api/v1/auth/me", timeout=15)
    if me.status_code != 200:
        raise RuntimeError(f"/auth/me for {spec['email']} failed: {me.status_code}")
    user = me.json() or {}
    if user.get("email", "").lower() != spec["email"].lower():
        raise RuntimeError(f"key for {spec['email']} matched {user.get('email')}")
    if user.get("role") != spec["role"]:
        raise RuntimeError(f"{spec['email']} role {user.get('role')} != {spec['role']}")
    if user.get("mustChangePin"):
        raise RuntimeError(f"{spec['email']} still mustChangePin")
    time.sleep(1)


def write_personas_env(rows: list[tuple[dict, str, str]]) -> None:
    pin_line = rows[0][2] if rows else "1234"
    lines = [
        "# Generated by seed_personas.py — local load tests only",
        f"# Verified logins. Shared PIN: {pin_line}",
        "",
    ]
    for spec, key, pin in rows:
        prefix = spec["env_prefix"]
        lines.extend(
            [
                f"# {spec['name']} ({spec['role']})",
                f"{prefix}_EMAIL={spec['email']}",
                f"{prefix}_KEY={key}",
                f"{prefix}_PIN={pin}",
                f"{prefix}_LOGIN={spec['login']}",
                f"{prefix}_API_ROLE={spec['role']}",
                "",
            ]
        )
    PERSONAS_ENV.write_text("\n".join(lines), encoding="utf-8")
    print(f"Wrote {PERSONAS_ENV}")


def main() -> int:
    load_dotenv()
    if SETUP_ENV.is_file():
        for line in SETUP_ENV.read_text(encoding="utf-8").splitlines():
            line = line.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            k, _, v = line.partition("=")
            k, v = k.strip(), v.strip().strip("'").strip('"')
            if k and k not in os.environ:
                os.environ[k] = v

    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="http://127.0.0.1:8080")
    parser.add_argument(
        "--force",
        action="store_true",
        help="Delete existing loadtest-* users before recreating",
    )
    args = parser.parse_args()

    if not is_loopback_host(args.host):
        print("Refusing non-loopback host.", file=sys.stderr)
        return 2

    try:
        setup_email, setup_key, setup_pin, loadtest_pin = load_setup_admin()
    except RuntimeError as exc:
        print(exc, file=sys.stderr)
        return 1

    client = ApiClient(args.host)
    print(f"Logging in as setup admin ({setup_email})...")
    login_setup_admin(client, setup_email, setup_key, setup_pin)

    if args.force:
        print("Removing old loadtest persona users...")
        delete_loadtest_users(client)

    created: list[tuple[dict, str, str]] = []
    for spec in PERSONA_DEFS:
        existing = next(
            (u for u in list_users(client) if (u.get("email") or "").lower() == spec["email"].lower()),
            None,
        )
        if existing and not args.force:
            raise RuntimeError(
                f"{spec['email']} already exists. Run with --force to recreate, "
                "or delete loadtest users manually."
            )

        print(f"Creating {spec['email']} ({spec['role']})...")
        key = register_user(client, spec)
        set_persona_pin(spec, key, loadtest_pin, args.host)
        created.append((spec, key, loadtest_pin))

    write_personas_env(created)
    print("\nVerified all persona logins. Next: python verify_personas.py")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
