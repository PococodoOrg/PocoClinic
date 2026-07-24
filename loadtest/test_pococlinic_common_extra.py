"""Additional unit tests for load-test helpers."""

from __future__ import annotations

import os

import pytest

import pococlinic_common as common


def test_unique_tag_is_stable_per_call():
    class FakeUser:
        id = "abc"

    tag1 = common.unique_tag(FakeUser())
    tag2 = common.unique_tag(FakeUser())
    assert tag1.startswith("load-abc-")
    assert tag2.startswith("load-abc-")
    assert tag1 != tag2


def test_env_bool():
    assert common.env_bool("MISSING", default=True) is True
    os.environ["TEST_FLAG"] = "yes"
    try:
        assert common.env_bool("TEST_FLAG") is True
        os.environ["TEST_FLAG"] = "off"
        assert common.env_bool("TEST_FLAG") is False
    finally:
        os.environ.pop("TEST_FLAG", None)


def test_validate_persona_env_exits_when_missing(monkeypatch):
    monkeypatch.setattr(
        common,
        "persona_credentials",
        lambda name: None if name == "clerk" else common.PersonaCredentials(
            name=name,
            login="staff",
            email=f"{name}@test.local",
            key="key",
            pin="1234",
            api_role="staff",
        ),
    )
    monkeypatch.setattr(common.sys, "exit", lambda code: (_ for _ in ()).throw(SystemExit(code)))
    with pytest.raises(SystemExit) as exc:
        common.validate_persona_env(["admin", "clerk"])
    assert exc.value.code == 2


def test_active_scenario_name_default(monkeypatch):
    monkeypatch.delenv("LOADTEST_SCENARIO", raising=False)
    assert common.active_scenario_name() == "tiny"
