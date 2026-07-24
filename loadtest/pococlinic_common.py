"""Shared helpers for local PocoClinic Locust runs (loopback only)."""

from __future__ import annotations

import json
import os
import random
import sys
import time
import uuid
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Callable
from urllib.parse import urlparse

try:
    import gevent.lock
except ImportError:
    gevent = None  # type: ignore[assignment]

LOADTEST_DIR = Path(__file__).resolve().parent
REPO_ROOT = LOADTEST_DIR.parent
DEFAULT_BOOTSTRAP_KEY = REPO_ROOT / "backend" / "data" / "bootstrap-admin-once.txt"
SCENARIOS_FILE = LOADTEST_DIR / "scenarios.json"

# Backend login limiter: burst 5, then ~1 per 12s (see cmd/main.go loginRateLimiter).
LOGIN_MIN_INTERVAL_S = float(os.environ.get("LOADTEST_LOGIN_INTERVAL", "13"))
LOGIN_MAX_RETRIES = int(os.environ.get("LOADTEST_LOGIN_RETRIES", "8"))

_login_gate = gevent.lock.Semaphore(1) if gevent else None
_last_login_monotonic = 0.0
_persona_cookie_lock = gevent.lock.Semaphore(1) if gevent else None
_persona_cookies: dict[str, dict[str, str]] = {}


@dataclass(frozen=True)
class PersonaCredentials:
    name: str
    login: str  # "admin" | "staff"
    email: str
    key: str
    pin: str
    api_role: str  # actual RBAC role in PocoClinic


def env_bool(name: str, default: bool = False) -> bool:
    raw = os.environ.get(name, "").strip().lower()
    if raw == "":
        return default
    return raw in ("1", "true", "yes", "on")


def load_dotenv() -> None:
    for filename, override in ((".env", False), ("setup.env", False), ("personas.env", True)):
        path = LOADTEST_DIR / filename
        if not path.is_file():
            continue
        for line in path.read_text(encoding="utf-8").splitlines():
            line = line.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            key, _, value = line.partition("=")
            key = key.strip()
            value = value.strip().strip("'").strip('"')
            if not key:
                continue
            if override or key not in os.environ:
                os.environ[key] = value


def is_loopback_host(host: str) -> bool:
    if not host:
        return False
    parsed = urlparse(host if "://" in host else f"http://{host}")
    hostname = (parsed.hostname or "").lower()
    return hostname in {"localhost", "127.0.0.1", "::1"}


def read_bootstrap_key() -> str:
    path = Path(os.environ.get("SMOKE_BOOTSTRAP_KEY_FILE", str(DEFAULT_BOOTSTRAP_KEY)))
    if not path.is_file():
        return ""
    for line in path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        if "=" in line and not line.lower().startswith("http"):
            _, _, value = line.partition("=")
            value = value.strip()
            if value:
                return value
        return line
    return ""


def persona_credentials(name: str) -> PersonaCredentials | None:
    prefix = name.upper()
    email = os.environ.get(f"{prefix}_EMAIL", "").strip()
    key = os.environ.get(f"{prefix}_KEY", "").strip()
    pin = os.environ.get(f"{prefix}_PIN", "").strip()
    login = os.environ.get(f"{prefix}_LOGIN", "").strip().lower()
    api_role = os.environ.get(f"{prefix}_API_ROLE", "").strip()

    if not email and name == "admin":
        email = os.environ.get("SMOKE_EMAIL", "admin@pococlinic.local").strip()
    if not key and name == "admin" and email.lower() == "admin@pococlinic.local":
        key = read_bootstrap_key()

    if not key or not pin or len(pin) != 4:
        return None

    if not login:
        login = "admin" if email else "staff"
    if not api_role:
        defaults = {
            "admin": "admin",
            "manager": "admin",
            "clinician": "doctor",
            "clerk": "nurse",
        }
        api_role = defaults.get(name, "staff")

    return PersonaCredentials(
        name=name,
        login=login,
        email=email,
        key=key,
        pin=pin,
        api_role=api_role,
    )


def load_scenarios() -> dict[str, Any]:
    return json.loads(SCENARIOS_FILE.read_text(encoding="utf-8"))


def active_scenario_name() -> str:
    return os.environ.get("LOADTEST_SCENARIO", "tiny").strip().lower()


def active_scenario() -> dict[str, Any]:
    data = load_scenarios()
    name = active_scenario_name()
    if name not in data["scenarios"]:
        raise KeyError(f"Unknown scenario {name!r}")
    return data["scenarios"][name]


def refuse_non_loopback(host: str) -> None:
    if not is_loopback_host(host):
        sys.stderr.write(
            "\nRefusing to run: load tests only allow loopback hosts "
            "(localhost, 127.0.0.1, ::1).\n"
            f"Got host={host!r}. Local manual runs only — not CI, not clinic LAN.\n\n"
        )
        sys.exit(2)


def validate_persona_env(required: list[str]) -> None:
    missing = []
    for name in required:
        if persona_credentials(name) is None:
            missing.append(name)
    if missing:
        sys.stderr.write(
            "\nMissing persona credentials for: "
            + ", ".join(missing)
            + "\nCopy personas.env.example → personas.env and run seed_personas.py, "
            "or set {PERSONA}_KEY and {PERSONA}_PIN.\n\n"
        )
        sys.exit(2)


def unique_tag(user: Any) -> str:
    uid = getattr(user, "id", None) or random.randint(1000, 9999)
    return f"load-{uid}-{uuid.uuid4().hex[:6]}"


def pick_patient_id(user: Any) -> str | None:
    ids: list[str] = getattr(user, "patient_ids", [])
    if ids:
        return random.choice(ids)
    return getattr(user, "last_patient_id", None)


def _sleep(seconds: float) -> None:
    if gevent is not None:
        gevent.sleep(seconds)
    else:
        time.sleep(seconds)


def throttle_login_attempt() -> None:
    """Serialize login calls so we stay under the API login rate limit (per IP)."""
    global _last_login_monotonic
    if _login_gate is None:
        return
    with _login_gate:
        now = time.monotonic()
        wait = LOGIN_MIN_INTERVAL_S - (now - _last_login_monotonic)
        if wait > 0:
            _sleep(wait)
        _last_login_monotonic = time.monotonic()


def reuse_or_login(
    persona_name: str,
    login_once: Callable[[], tuple[int, dict[str, str], str]],
) -> dict[str, str]:
    """
    Return cookie dict for a persona. Many Locust users share one badge credential;
    we login once per persona and reuse cookies to avoid login rate limits.
    """
    if _persona_cookie_lock is None:
        status, cookies, err = login_once()
        if status != 200:
            raise RuntimeError(err or f"login failed HTTP {status}")
        return cookies

    with _persona_cookie_lock:
        cached = _persona_cookies.get(persona_name)
        if cached:
            return dict(cached)

        throttle_login_attempt()
        status, cookies, err = login_once()
        if status == 429:
            for _ in range(LOGIN_MAX_RETRIES):
                _sleep(LOGIN_MIN_INTERVAL_S)
                throttle_login_attempt()
                status, cookies, err = login_once()
                if status == 200:
                    break
        if status != 200:
            raise RuntimeError(err or f"login failed HTTP {status}")

        _persona_cookies[persona_name] = dict(cookies)
        return dict(cookies)

