# Local API load tests (Locust)

**Everything lives in `loadtest/`.** Local only — not CI/CD.

## Run (Windows)

Start backend (`cd backend` → `go run ./cmd`), then:

```bat
cd loadtest
copy setup.env.example setup.env
REM edit setup.env — system admin badge + PIN
run.bat
```

Or directly:

```bat
python run.py
python run.py tiny
python run.py ui
python run.py --seed
```

`run.py` handles venv, health check, persona verify/seed, all scenarios, and writes `reports/summary-latest.md`.

## What was fixed

- **No more bat logic** — `run.bat` / `run.sh` call `run.py` only
- **Subprocess deadlock** — Locust output streams live (Windows pipe buffer hang)
- **Progress** — prints `Progress: 2/6` per scenario + elapsed time
- **Timeouts** — each scenario capped so a stuck Locust run cannot block the rest
- **Login reuse** — one login per persona under a lock (fewer 429s)

## Scenarios

| ID | Users | Duration |
|----|------:|----------|
| smoke | 4 | 60s |
| tiny | 4 | 90s |
| small | 5 | 120s |
| medium | 10 | 180s |
| large | 17 | 240s |
| xlarge | 35 | 300s |
| stress100 | 100 | 300s |

Full suite takes ~25–35 minutes (stress100 is slow at spawn rate 1/s).

## Credentials

```bat
python run.py --seed
```

Uses `setup.env` (system admin) to recreate `personas.env`. Do not hand-edit badge keys.

## Safety

- Loopback host only
- Not in CI
- `setup.env`, `personas.env`, `reports/*` gitignored
