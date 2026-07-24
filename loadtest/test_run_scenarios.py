"""Unit tests for scenario runner helpers."""

from __future__ import annotations

from pathlib import Path

import run_scenarios


def test_parse_duration_seconds():
    assert run_scenarios.parse_duration_seconds("90s") == 90
    assert run_scenarios.parse_duration_seconds("2m") == 120
    assert run_scenarios.parse_duration_seconds("1h") == 3600
    assert run_scenarios.parse_duration_seconds("bad") == 120


def test_scenario_timeout_includes_spawn_and_buffer():
    timeout = run_scenarios.scenario_timeout(
        {"duration": "90s", "users": 4, "spawn_rate": 1},
    )
    assert timeout >= 90 + 4 + 180


def test_persona_breakdown_groups_stats(tmp_path: Path):
    csv_path = tmp_path / "tiny_stats.csv"
    csv_path.write_text(
        "Type,Name,Request Count,Failure Count,Median Response Time,95%\n"
        "GET,admin:users,10,0,5,8\n"
        "GET,clinician:patient:get,20,1,12,30\n"
        "GET,shared:patients:list,15,0,3,6\n"
        "None,Aggregated,45,1,8,30\n",
        encoding="utf-8",
    )

    breakdown = run_scenarios.persona_breakdown(csv_path)
    assert breakdown["admin"]["requests"] == 10
    assert breakdown["clinician"]["failures"] == 1
    assert breakdown["shared"]["requests"] == 15
