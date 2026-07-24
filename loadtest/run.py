#!/usr/bin/env python3
"""
PocoClinic load-test runner (local only, not CI).

Replaces run.bat with one Python entry point: venv, health check, personas, scenarios.

  python run.py              # all scenarios + report
  python run.py smoke        # quick endpoint sweep (~60s)
  python run.py tiny         # one scenario
  python run.py ui           # Locust web UI
  python run.py --help
"""

from __future__ import annotations

import argparse
import os
import shutil
import subprocess
import sys
import time
import urllib.error
import urllib.request
from pathlib import Path

LOADTEST_DIR = Path(__file__).resolve().parent
VENV_DIR = LOADTEST_DIR / ".venv"
REPORTS_DIR = LOADTEST_DIR / "reports"

if str(LOADTEST_DIR) not in sys.path:
    sys.path.insert(0, str(LOADTEST_DIR))

from pococlinic_common import is_loopback_host, load_dotenv  # noqa: E402


def venv_python() -> Path:
    if os.name == "nt":
        return VENV_DIR / "Scripts" / "python.exe"
    return VENV_DIR / "bin" / "python"


def venv_locust() -> Path:
    if os.name == "nt":
        return VENV_DIR / "Scripts" / "locust.exe"
    return VENV_DIR / "bin" / "locust"


def ensure_venv() -> None:
    py = venv_python()
    if not py.is_file():
        print("[1/5] Creating virtual environment...", flush=True)
        subprocess.run([sys.executable, "-m", "venv", str(VENV_DIR)], check=True, cwd=LOADTEST_DIR)
    else:
        print("[1/5] Virtual environment OK", flush=True)

    print("[2/5] Installing dependencies...", flush=True)
    subprocess.run(
        [str(py), "-m", "pip", "install", "-r", "requirements.txt", "-q"],
        check=True,
        cwd=LOADTEST_DIR,
    )


def check_health(host: str) -> None:
    print(f"[3/5] Checking API at {host} ...", flush=True)
    url = f"{host.rstrip('/')}/health"
    try:
        with urllib.request.urlopen(url, timeout=5) as resp:
            if resp.status != 200:
                raise RuntimeError(f"HTTP {resp.status}")
    except (urllib.error.URLError, TimeoutError, RuntimeError) as exc:
        print(f"\nERROR: Backend not reachable at {host}\n  ({exc})", flush=True)
        print("Start it: cd backend && go run ./cmd\n", flush=True)
        raise SystemExit(1)
    print("      API reachable.", flush=True)


def ensure_personas(host: str, force_seed: bool) -> None:
    setup_env = LOADTEST_DIR / "setup.env"
    personas_env = LOADTEST_DIR / "personas.env"
    py = str(venv_python())

    if not setup_env.is_file():
        print("\nMissing setup.env — copy setup.env.example and add system admin credentials.\n", flush=True)
        raise SystemExit(1)

    if force_seed or not personas_env.is_file():
        print("[4/5] Seeding load-test personas...", flush=True)
        subprocess.run([py, "seed_personas.py", "--host", host, "--force"], check=True, cwd=LOADTEST_DIR)
        return

    print("[4/5] Verifying persona logins...", flush=True)
    proc = subprocess.run([py, "verify_personas.py", "--host", host], cwd=LOADTEST_DIR)
    if proc.returncode != 0:
        print("\nPersona verification failed. Try: python run.py --seed\n", flush=True)
        raise SystemExit(1)


def run_ui(host: str, extra: list[str]) -> int:
    locust = venv_locust()
    if not locust.is_file():
        locust = shutil.which("locust") or "locust"  # type: ignore[assignment]

    scenario = os.environ.get("LOADTEST_SCENARIO", "small")
    os.environ["LOADTEST_SCENARIO"] = scenario
    print(f"[5/5] Locust UI (hint scenario={scenario}) — http://localhost:8089", flush=True)
    cmd = [str(locust), "-f", "locustfile.py", "--host", host, *extra]
    return subprocess.run(cmd, cwd=LOADTEST_DIR).returncode


def run_scenarios(host: str, scenario_ids: list[str] | None) -> int:
    py = str(venv_python())
    locust = str(venv_locust())
    cmd = [py, "run_scenarios.py", "--host", host, "--locust", locust]
    if scenario_ids:
        cmd.extend(["--scenarios", *scenario_ids])

    print("[5/5] Running load scenarios (output streams below)...", flush=True)
    proc = subprocess.run(cmd, cwd=LOADTEST_DIR)
    return proc.returncode


def main() -> int:
    load_dotenv()

    parser = argparse.ArgumentParser(description="PocoClinic local load-test runner")
    parser.add_argument(
        "mode",
        nargs="?",
        default="all",
        help="all | ui | smoke | tiny | small | medium | large | xlarge | stress100",
    )
    parser.add_argument("--host", default=os.environ.get("LOADTEST_HOST", "http://127.0.0.1:8080"))
    parser.add_argument("--seed", action="store_true", help="Force re-seed personas before run")
    parser.add_argument("--skip-setup", action="store_true", help="Skip venv/pip (already prepared)")
    parser.add_argument("extra", nargs=argparse.REMAINDER, help="Extra args for Locust UI mode")
    args = parser.parse_args()

    if not is_loopback_host(args.host):
        print("Refusing non-loopback host.", file=sys.stderr)
        return 2

    print("\n=== PocoClinic load test (local only) ===", flush=True)
    print(f"Folder: {LOADTEST_DIR}", flush=True)
    if args.mode.lower() in ("all",):
        print("Full suite: 6 scenarios, ~25-35 min. Use 'python run.py tiny' for a quick smoke.\n", flush=True)
    else:
        print(flush=True)

    if not args.skip_setup:
        ensure_venv()
    check_health(args.host)
    ensure_personas(args.host, args.seed)

    mode = args.mode.lower()
    if mode in ("help", "-h", "--help"):
        parser.print_help()
        return 0
    if mode == "ui":
        extra = [a for a in args.extra if a != "--"]
        return run_ui(args.host, extra)
    if mode == "all":
        return run_scenarios(args.host, None)

    return run_scenarios(args.host, [mode])


if __name__ == "__main__":
    raise SystemExit(main())
