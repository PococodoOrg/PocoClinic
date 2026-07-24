"""Unit tests for load-test helpers (no live API required)."""

from __future__ import annotations

import os
import sys

import pytest

import pococlinic_common as common


def test_is_loopback_host_accepts_localhost_variants():
    assert common.is_loopback_host("http://127.0.0.1:8080")
    assert common.is_loopback_host("http://localhost:8080")
    assert common.is_loopback_host("http://[::1]:8080")


def test_is_loopback_host_rejects_remote():
    assert not common.is_loopback_host("http://192.168.1.10:8080")
    assert not common.is_loopback_host("")


def test_refuse_non_loopback_exits(monkeypatch):
    monkeypatch.setattr(sys, "exit", lambda code: (_ for _ in ()).throw(SystemExit(code)))
    with pytest.raises(SystemExit) as exc:
        common.refuse_non_loopback("http://10.0.0.5")
    assert exc.value.code == 2


def test_persona_credentials_requires_key_and_pin(monkeypatch):
    monkeypatch.delenv("CLINICIAN_KEY", raising=False)
    monkeypatch.delenv("CLINICIAN_PIN", raising=False)
    assert common.persona_credentials("clinician") is None


def test_persona_credentials_builds_defaults(monkeypatch):
    monkeypatch.setenv("CLINICIAN_EMAIL", "doc@test.local")
    monkeypatch.setenv("CLINICIAN_KEY", "badge-key")
    monkeypatch.setenv("CLINICIAN_PIN", "1234")
    monkeypatch.setenv("CLINICIAN_LOGIN", "staff")
    cred = common.persona_credentials("clinician")
    assert cred is not None
    assert cred.login == "staff"
    assert cred.api_role == "doctor"


def test_active_scenario_unknown_raises(monkeypatch):
    monkeypatch.setenv("LOADTEST_SCENARIO", "not-a-scenario")
    with pytest.raises(KeyError):
        common.active_scenario()


def test_load_dotenv_reads_personas_env(tmp_path, monkeypatch):
    env_file = tmp_path / "personas.env"
    env_file.write_text("CLINICIAN_PIN=5678\n", encoding="utf-8")
    monkeypatch.setattr(common, "LOADTEST_DIR", tmp_path)
    monkeypatch.delenv("CLINICIAN_PIN", raising=False)
    common.load_dotenv()
    assert os.environ["CLINICIAN_PIN"] == "5678"
