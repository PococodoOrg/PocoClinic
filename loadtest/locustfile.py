"""
PocoClinic local load / smoke tests (Locust).

Personas (see personas.env.example):
  - Admin: user mgmt, system status, audit, backups list
  - Manager: admin reports (census, activity, staff activity, performance)
  - Clinician: chart reads, updates, clinical notes, form submissions
  - Clerk: intake — create patients + generic notes (nurse role in RBAC)

Safety:
  - Loopback host only
  - No API test mode / auth bypass
  - Not wired to CI/CD

Run one scenario:
  set LOADTEST_SCENARIO=medium
  locust -f locustfile.py --host http://127.0.0.1:8080 --headless

Run all scenarios + report:
  python run_scenarios.py --host http://127.0.0.1:8080
"""

from __future__ import annotations

import os
import random
import sys

from locust import HttpUser, between, events, task
from locust.exception import StopUser

from pococlinic_common import (
    PersonaCredentials,
    active_scenario,
    active_scenario_name,
    load_dotenv,
    persona_credentials,
    pick_patient_id,
    refuse_non_loopback,
    reuse_or_login,
    unique_tag,
    validate_persona_env,
)

PERSONA_WEIGHTS: dict[str, int] = {
    "admin": 1,
    "manager": 0,
    "clinician": 1,
    "clerk": 1,
}


class PersonaMixin:
    persona_name: str = ""
    wait_time = between(2.0, 5.0)

    credentials: PersonaCredentials | None = None
    patient_ids: list[str]
    template_id: str | None
    template_field_id: str | None
    exercise_plan_id: str | None

    def on_start(self) -> None:
        load_dotenv()
        self.patient_ids = []
        self.template_id = None
        self.template_field_id = None
        self.exercise_plan_id = None
        self.credentials = persona_credentials(self.persona_name)
        if not self.credentials:
            raise StopUser(f"missing credentials for persona {self.persona_name}")
        self._login()
        self._prime_catalogs()

    def _login(self) -> None:
        cred = self.credentials
        assert cred is not None
        if cred.login == "admin":
            path = "/api/v1/auth/login/admin"
            body = {"email": cred.email, "key": cred.key, "pin": cred.pin}
            login_name = f"persona:{cred.name}:login-admin"
        else:
            path = "/api/v1/auth/login/staff"
            body = {"key": cred.key, "pin": cred.pin}
            login_name = f"persona:{cred.name}:login-staff"

        def attempt() -> tuple[int, dict[str, str], str]:
            with self.client.post(
                path, json=body, name=login_name, catch_response=True
            ) as resp:
                if resp.status_code != 200:
                    detail = resp.text[:300] if resp.text else ""
                    if resp.status_code != 429:
                        resp.failure(f"login failed: HTTP {resp.status_code} {detail}")
                    return resp.status_code, {}, detail
                user = (resp.json() or {}).get("user") or {}
                if cred.email and user.get("email", "").lower() != cred.email.lower():
                    msg = f"key logged in as {user.get('email')}, expected {cred.email}"
                    resp.failure(msg)
                    return resp.status_code, {}, msg
                if cred.api_role and user.get("role") != cred.api_role:
                    msg = f"role {user.get('role')} != expected {cred.api_role}"
                    resp.failure(msg)
                    return resp.status_code, {}, msg
                if user.get("mustChangePin"):
                    msg = "mustChangePin=true — run: python seed_personas.py --force"
                    resp.failure(msg)
                    return resp.status_code, {}, msg
                resp.success()
                return 200, self.client.cookies.get_dict(), ""

        try:
            cookies = reuse_or_login(cred.name, attempt)
        except RuntimeError as exc:
            raise StopUser(str(exc)) from exc

        self.client.cookies.clear()
        for name, value in cookies.items():
            self.client.cookies.set(name, value)

    def _prime_catalogs(self) -> None:
        with self.client.get(
            "/api/v1/patients?page=1&pageSize=25",
            name="shared:patients:list",
            catch_response=True,
        ) as resp:
            if resp.status_code == 200:
                data = resp.json() or {}
                items = data.get("patients") or data.get("items") or []
                self.patient_ids = [p["id"] for p in items if p.get("id")]
                resp.success()
            elif resp.status_code in (403, 401):
                resp.success()
            else:
                resp.failure(f"unexpected {resp.status_code}")

        with self.client.get(
            "/api/v1/form-templates",
            name="shared:forms:templates",
            catch_response=True,
        ) as resp:
            if resp.status_code == 200:
                data = resp.json() or {}
                templates = data.get("templates") if isinstance(data, dict) else data
                if isinstance(templates, list) and templates:
                    tpl = templates[0]
                    self.template_id = tpl.get("id")
                    fields = tpl.get("fields") or []
                    if fields:
                        self.template_field_id = fields[0].get("id")
                resp.success()
            elif resp.status_code in (403, 401):
                resp.success()
            else:
                resp.failure(f"unexpected {resp.status_code}")

    def _refresh_patients(self) -> None:
        with self.client.get(
            "/api/v1/patients?page=1&pageSize=25",
            name="shared:patients:list",
            catch_response=True,
        ) as resp:
            if resp.status_code == 200:
                data = resp.json() or {}
                items = data.get("patients") or data.get("items") or []
                self.patient_ids = [p["id"] for p in items if p.get("id")]
                resp.success()
            elif resp.status_code in (403, 401):
                resp.success()
            else:
                resp.failure(f"unexpected {resp.status_code}")


class AdminPersona(PersonaMixin, HttpUser):
    persona_name = "admin"
    weight = 1

    @task(3)
    def system_status(self) -> None:
        self.client.get("/api/v1/admin/system-status", name="admin:system-status")

    @task(3)
    def list_users(self) -> None:
        self.client.get("/api/v1/auth/users?page=1&pageSize=20", name="admin:users")

    @task(2)
    def audit_logs(self) -> None:
        self.client.get("/api/v1/admin/audit-logs?page=1&pageSize=25", name="admin:audit-logs")

    @task(2)
    def list_backups(self) -> None:
        self.client.get("/api/v1/admin/backups", name="admin:backups")

    @task(1)
    def health_check(self) -> None:
        self.client.get("/api/v1/admin/health-check", name="admin:health-check")

    @task(1)
    def me(self) -> None:
        self.client.get("/api/v1/auth/me", name="shared:auth:me")

    @task(1)
    def refresh_session(self) -> None:
        with self.client.post(
            "/api/v1/auth/refresh",
            name="admin:auth:refresh",
            catch_response=True,
        ) as resp:
            if resp.status_code in (200, 429):
                resp.success()
            else:
                resp.failure(f"HTTP {resp.status_code}")

    @task(1)
    def export_audit_csv(self) -> None:
        self.client.get(
            "/api/v1/admin/audit-logs/export.csv",
            name="admin:audit-export",
        )


class ManagerPersona(PersonaMixin, HttpUser):
    persona_name = "manager"
    weight = 1

    @task(4)
    def patient_census(self) -> None:
        self.client.get("/api/v1/admin/reports/patient-census", name="manager:patient-census")

    @task(3)
    def activity_summary(self) -> None:
        self.client.get("/api/v1/admin/reports/activity-summary", name="manager:activity-summary")

    @task(3)
    def staff_activity(self) -> None:
        self.client.get("/api/v1/admin/reports/staff-activity", name="manager:staff-activity")

    @task(2)
    def performance(self) -> None:
        self.client.get("/api/v1/admin/performance", name="manager:performance")

    @task(2)
    def compliance(self) -> None:
        self.client.get("/api/v1/admin/compliance-check", name="manager:compliance")

    @task(1)
    def system_status(self) -> None:
        self.client.get("/api/v1/admin/system-status", name="manager:system-status")


class ClinicianPersona(PersonaMixin, HttpUser):
    persona_name = "clinician"
    weight = 1

    @task(4)
    def chart_round(self) -> None:
        pid = pick_patient_id(self)
        if not pid:
            self._refresh_patients()
            pid = pick_patient_id(self)
        if not pid:
            return
        self.client.get(f"/api/v1/patients/{pid}", name="clinician:patient:get")
        self.client.get(f"/api/v1/patients/{pid}/notes", name="clinician:notes:list")
        self.client.get(f"/api/v1/patients/{pid}/form-submissions", name="clinician:forms:list")
        self.client.get(f"/api/v1/patients/{pid}/documents", name="clinician:documents:list")
        self.client.get(f"/api/v1/patients/{pid}/exercise-plans", name="clinician:exercise:list")

    @task(2)
    def create_exercise_plan(self) -> None:
        pid = pick_patient_id(self)
        if not pid:
            return
        tag = unique_tag(self)
        body = {"name": f"Load plan {tag}", "description": "Load test plan"}
        with self.client.post(
            f"/api/v1/patients/{pid}/exercise-plans",
            json=body,
            name="clinician:exercise:create-plan",
            catch_response=True,
        ) as resp:
            if resp.status_code in (200, 201):
                plan = resp.json() or {}
                plan_id = plan.get("id")
                if plan_id:
                    self.exercise_plan_id = plan_id
                resp.success()
            elif resp.status_code == 429:
                resp.success()
            else:
                resp.failure(f"HTTP {resp.status_code}")

    @task(1)
    def view_form_template(self) -> None:
        if not self.template_id:
            return
        self.client.get(
            f"/api/v1/form-templates/{self.template_id}",
            name="clinician:forms:template-detail",
        )

    @task(3)
    def update_patient(self) -> None:
        pid = pick_patient_id(self)
        if not pid:
            return
        with self.client.get(
            f"/api/v1/patients/{pid}",
            name="clinician:patient:get-for-update",
            catch_response=True,
        ) as resp:
            if resp.status_code != 200:
                resp.failure(f"get before update: {resp.status_code}")
                return
            patient = resp.json() or {}
            resp.success()
        body = {
            "firstName": patient.get("firstName") or "Load",
            "lastName": patient.get("lastName") or "Patient",
            "dateOfBirth": (patient.get("dateOfBirth") or "1990-01-01")[:10],
            "gender": patient.get("gender") or "unknown",
            "email": patient.get("email") or f"{pid}@loadtest.local",
            "phoneNumber": f"555-{random.randint(1000, 9999)}",
        }
        if patient.get("middleName"):
            body["middleName"] = patient["middleName"]
        self.client.put(f"/api/v1/patients/{pid}", json=body, name="clinician:patient:update")

    @task(3)
    def clinical_note(self) -> None:
        pid = pick_patient_id(self)
        if not pid:
            return
        tag = unique_tag(self)
        body = {"body": f"Loadtest clinical note {tag}"}
        with self.client.post(
            f"/api/v1/patients/{pid}/notes",
            json=body,
            name="clinician:notes:create",
            catch_response=True,
        ) as resp:
            if resp.status_code in (201, 200):
                resp.success()
            elif resp.status_code == 403:
                resp.failure("clinician cannot create notes — check role")
            else:
                resp.failure(f"HTTP {resp.status_code}")

    @task(2)
    def submit_form(self) -> None:
        pid = pick_patient_id(self)
        if not pid or not self.template_id:
            return
        field_id = self.template_field_id or "q1"
        body = {
            "templateId": self.template_id,
            "answers": {field_id: f"loadtest-{unique_tag(self)}"},
        }
        self.client.post(
            f"/api/v1/patients/{pid}/form-submissions",
            json=body,
            name="clinician:forms:submit",
        )

    @task(1)
    def list_templates(self) -> None:
        self.client.get("/api/v1/form-templates", name="clinician:forms:templates")


class ClerkPersona(PersonaMixin, HttpUser):
    persona_name = "clerk"
    weight = 1

    @task(4)
    def register_patient(self) -> None:
        tag = unique_tag(self)
        body = {
            "firstName": "Load",
            "lastName": f"Patient{tag[-6:]}",
            "dateOfBirth": "1990-06-15",
            "gender": random.choice(["female", "male", "other"]),
            "email": f"{tag}@loadtest.local",
            "phoneNumber": f"555-{random.randint(1000, 9999)}",
        }
        with self.client.post(
            "/api/v1/patients",
            json=body,
            name="clerk:patients:create",
            catch_response=True,
        ) as resp:
            if resp.status_code in (200, 201):
                patient = resp.json() or {}
                pid = patient.get("id")
                if pid:
                    self.patient_ids.append(pid)
                resp.success()
            elif resp.status_code == 403:
                resp.failure("clerk persona needs nurse role for patient create")
            else:
                resp.failure(f"HTTP {resp.status_code} {resp.text[:120]}")

    @task(3)
    def intake_note(self) -> None:
        pid = pick_patient_id(self)
        if not pid:
            self._refresh_patients()
            pid = pick_patient_id(self)
        if not pid:
            return
        body = {"body": f"Front desk note — checked in {unique_tag(self)}"}
        self.client.post(
            f"/api/v1/patients/{pid}/notes",
            json=body,
            name="clerk:notes:create",
        )

    @task(2)
    def search_patients(self) -> None:
        q = random.choice(["Load", "Patient", ""])
        self.client.get(
            f"/api/v1/patients?page=1&pageSize=20&search={q}",
            name="clerk:patients:search",
        )

    @task(1)
    def list_groups(self) -> None:
        self.client.get("/api/v1/form-groups", name="clerk:forms:groups")


PERSONA_CLASSES: dict[str, type[HttpUser]] = {
    "admin": AdminPersona,
    "manager": ManagerPersona,
    "clinician": ClinicianPersona,
    "clerk": ClerkPersona,
}


@events.init.add_listener
def _configure(environment, **_kwargs):  # noqa: ANN001
    load_dotenv()
    host = (environment.host or os.environ.get("LOCUST_HOST") or "").rstrip("/")
    refuse_non_loopback(host)

    # Fresh run — drop cached persona sessions from a prior Locust process in same interpreter.
    from pococlinic_common import _persona_cookies

    _persona_cookies.clear()

    scenario = active_scenario()
    weights = scenario.get("weights") or {}
    global PERSONA_WEIGHTS
    PERSONA_WEIGHTS = {k: int(v) for k, v in weights.items()}

    for persona, cls in PERSONA_CLASSES.items():
        cls.weight = max(0, PERSONA_WEIGHTS.get(persona, 0))

    required = [p for p, w in PERSONA_WEIGHTS.items() if w > 0]
    validate_persona_env(required)

    if environment.parsed_options and getattr(environment.parsed_options, "headless", False):
        print(
            f"\nScenario: {scenario.get('label')} ({active_scenario_name()}) — "
            f"weights={PERSONA_WEIGHTS}\n",
            file=sys.stderr,
        )
