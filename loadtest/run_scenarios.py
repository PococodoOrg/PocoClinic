#!/usr/bin/env python3
"""
Run PocoClinic Locust scenarios locally and write a comparison report.

Streams Locust output live (no subprocess capture — avoids Windows pipe deadlocks).
"""

from __future__ import annotations

import argparse
import csv
import json
import os
import re
import subprocess
import sys
import time
from datetime import datetime, timezone
from pathlib import Path

from pococlinic_common import is_loopback_host, load_dotenv, load_scenarios

LOADTEST_DIR = Path(__file__).resolve().parent
REPORTS_DIR = LOADTEST_DIR / "reports"
DURATION_RE = re.compile(r"^(?P<n>\d+)(?P<u>s|m|h)?$", re.I)


def parse_duration_seconds(raw: str) -> int:
    raw = raw.strip()
    m = DURATION_RE.match(raw)
    if not m:
        return 120
    n = int(m.group("n"))
    unit = (m.group("u") or "s").lower()
    if unit == "m":
        return n * 60
    if unit == "h":
        return n * 3600
    return n


def scenario_timeout(scenario: dict) -> int:
    duration = parse_duration_seconds(str(scenario.get("duration", "120s")))
    users = int(scenario.get("users", 1))
    spawn_rate = max(1, int(scenario.get("spawn_rate", 1)))
    spawn_seconds = users / spawn_rate
    # Login throttle + SQLite teardown on Windows can be slow.
    return int(duration + spawn_seconds + 180)


def parse_stats_csv(path: Path) -> dict:
    if not path.is_file():
        return {}
    rows = list(csv.DictReader(path.open(encoding="utf-8")))
    agg = next((r for r in rows if r.get("Name") == "Aggregated"), None)
    if not agg:
        return {}
    return {
        "requests": int(float(agg.get("Request Count") or 0)),
        "failures": int(float(agg.get("Failure Count") or 0)),
        "rps": float(agg.get("Requests/s") or 0),
        "median_ms": float(agg.get("Median Response Time") or 0),
        "avg_ms": float(agg.get("Average Response Time") or 0),
        "p95_ms": float(agg.get("95%") or 0),
        "p99_ms": float(agg.get("99%") or 0),
        "max_ms": float(agg.get("Max Response Time") or 0),
    }


def persona_breakdown(stats_path: Path) -> dict[str, dict]:
    if not stats_path.is_file():
        return {}
    out: dict[str, dict] = {}
    for row in csv.DictReader(stats_path.open(encoding="utf-8")):
        name = row.get("Name") or ""
        if name in ("Aggregated", "Total"):
            continue
        if ":" not in name:
            continue
        prefix = name.split(":", 1)[0]
        if prefix not in ("admin", "manager", "clinician", "clerk", "persona", "shared"):
            continue
        bucket = prefix if prefix != "persona" else name.split(":", 2)[1]
        entry = out.setdefault(bucket, {"requests": 0, "failures": 0, "p95_ms": 0.0})
        entry["requests"] += int(float(row.get("Request Count") or 0))
        entry["failures"] += int(float(row.get("Failure Count") or 0))
        p95 = float(row.get("95%") or 0)
        if p95 > entry["p95_ms"]:
            entry["p95_ms"] = p95
    return out


def run_scenario(
    host: str,
    scenario_id: str,
    scenario: dict,
    locust_bin: str,
) -> dict:
    REPORTS_DIR.mkdir(parents=True, exist_ok=True)
    prefix = REPORTS_DIR / scenario_id
    env = os.environ.copy()
    env["LOADTEST_SCENARIO"] = scenario_id
    env["PYTHONUNBUFFERED"] = "1"

    cmd = [
        locust_bin,
        "-f",
        str(LOADTEST_DIR / "locustfile.py"),
        "--host",
        host,
        "--headless",
        "-u",
        str(scenario["users"]),
        "-r",
        str(scenario["spawn_rate"]),
        "-t",
        str(scenario["duration"]),
        "--csv",
        str(prefix),
        "--only-summary",
        "--exit-code-on-error",
        "0",
    ]

    timeout = scenario_timeout(scenario)
    print(f"\n{'=' * 60}", flush=True)
    print(f"Scenario: {scenario_id} — {scenario.get('label')}", flush=True)
    print(f"Users: {scenario['users']}  spawn: {scenario['spawn_rate']}/s  duration: {scenario['duration']}", flush=True)
    print(f"Timeout: {timeout}s", flush=True)
    print(f"Command: {' '.join(cmd)}", flush=True)
    print(f"{'=' * 60}\n", flush=True)

    started = time.monotonic()
    try:
        proc = subprocess.run(
            cmd,
            cwd=LOADTEST_DIR,
            env=env,
            timeout=timeout,
        )
        exit_code = proc.returncode
    except subprocess.TimeoutExpired:
        print(f"\nTIMEOUT: {scenario_id} exceeded {timeout}s — continuing.\n", flush=True)
        exit_code = 124
    elapsed = round(time.monotonic() - started, 1)

    stats_path = Path(f"{prefix}_stats.csv")
    stats = parse_stats_csv(stats_path)
    stats["scenario"] = scenario_id
    stats["label"] = scenario.get("label", scenario_id)
    stats["users"] = scenario["users"]
    stats["duration"] = scenario["duration"]
    stats["exit_code"] = exit_code
    stats["elapsed_s"] = elapsed
    stats["personas"] = persona_breakdown(stats_path)
    stats["fail_pct"] = (
        round(100.0 * stats["failures"] / stats["requests"], 2)
        if stats.get("requests")
        else 0.0
    )

    status = "OK" if exit_code == 0 and stats.get("requests", 0) > 0 else "WARN"
    print(
        f"\nFinished {scenario_id} in {elapsed}s — {status} "
        f"(exit={exit_code}, requests={stats.get('requests', 0)}, fail={stats.get('fail_pct', 0)}%)\n",
        flush=True,
    )
    return stats


def write_report(results: list[dict], host: str) -> None:
    REPORTS_DIR.mkdir(parents=True, exist_ok=True)
    ts = datetime.now(timezone.utc).strftime("%Y%m%d-%H%M%S")
    json_path = REPORTS_DIR / f"summary-{ts}.json"
    md_path = REPORTS_DIR / f"summary-{ts}.md"
    latest_json = REPORTS_DIR / "summary-latest.json"
    latest_md = REPORTS_DIR / "summary-latest.md"

    payload = {
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "host": host,
        "results": results,
    }
    text = json.dumps(payload, indent=2)
    json_path.write_text(text, encoding="utf-8")
    latest_json.write_text(text, encoding="utf-8")

    lines = [
        "# PocoClinic load test summary",
        "",
        f"Generated: {payload['generatedAt']}",
        f"Host: `{host}`",
        "",
        "| Scenario | Users | Duration | Elapsed | Requests | Fail % | RPS | p95 ms |",
        "|----------|------:|----------|--------:|--------:|-------:|----:|-------:|",
    ]
    for r in results:
        lines.append(
            f"| {r.get('label', r['scenario'])} | {r.get('users', '')} | {r.get('duration', '')} | "
            f"{r.get('elapsed_s', '?')}s | {r.get('requests', 0)} | {r.get('fail_pct', 0)}% | "
            f"{r.get('rps', 0):.1f} | {r.get('p95_ms', 0):.0f} |"
        )

    lines.extend(["", "## Persona breakdown", ""])
    for r in results:
        lines.append(f"### {r.get('label', r['scenario'])}")
        personas = r.get("personas") or {}
        if not personas:
            lines.append("_No persona-tagged stats_")
        else:
            lines.append("| Persona | Requests | Failures | p95 ms |")
            lines.append("|---------|--------:|---------:|-------:|")
            for persona, pstats in sorted(personas.items()):
                lines.append(
                    f"| {persona} | {pstats.get('requests', 0)} | "
                    f"{pstats.get('failures', 0)} | {pstats.get('p95_ms', 0):.0f} |"
                )
        lines.append("")

    md_text = "\n".join(lines)
    md_path.write_text(md_text, encoding="utf-8")
    latest_md.write_text(md_text, encoding="utf-8")
    print(f"\nWrote {latest_md} and {latest_json}\n", flush=True)


def main() -> int:
    load_dotenv()
    parser = argparse.ArgumentParser(description="Run local Locust scenarios")
    parser.add_argument("--host", default="http://127.0.0.1:8080")
    parser.add_argument("--scenarios", nargs="*", help="Scenario ids (default: all)")
    parser.add_argument("--locust", default="locust")
    args = parser.parse_args()

    if not is_loopback_host(args.host):
        print("Refusing non-loopback host.", file=sys.stderr)
        return 2

    if not Path(args.locust).is_file() and args.locust == "locust":
        # Allow bare name on PATH
        pass
    elif not Path(args.locust).is_file():
        print(f"Locust not found: {args.locust}", file=sys.stderr)
        return 1

    data = load_scenarios()
    ids = args.scenarios or list(data["scenarios"].keys())
    results: list[dict] = []
    failures = 0

    print(f"Running {len(ids)} scenario(s): {', '.join(ids)}", flush=True)

    est_seconds = sum(scenario_timeout(data["scenarios"][sid]) for sid in ids)
    est_min = max(1, round(est_seconds / 60))
    print(
        f"Estimated wall time: ~{est_min} min (max timeout budget {est_seconds}s). "
        "Output streams live — long pauses during stress100 are normal.\n",
        flush=True,
    )

    for i, sid in enumerate(ids, start=1):
        if sid not in data["scenarios"]:
            print(f"Unknown scenario {sid!r}", file=sys.stderr)
            return 1
        print(f"\n>>> Progress: {i}/{len(ids)}", flush=True)
        result = run_scenario(args.host, sid, data["scenarios"][sid], args.locust)
        results.append(result)
        if result.get("exit_code") not in (0, None) or result.get("requests", 0) == 0:
            failures += 1

    write_report(results, args.host.rstrip("/"))

    if failures:
        print(f"Completed with {failures} scenario(s) needing review.", flush=True)
    else:
        print("All scenarios finished.", flush=True)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
