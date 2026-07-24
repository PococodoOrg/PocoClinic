# Commit strategy (evening batch)

Step-by-step plan for landing the current PocoClinic monorepo work as **reviewable, logical commits** instead of one giant diff.

**Scope:** ~400+ new files + ~68 modified tracked files (SQLite EMR, auth, patients, forms, admin, ops tooling, frontend rewrite, docs/ADRs, CI).

**Branch suggestion:** create a long-lived integration branch first, then merge to `main` when ready.

```powershell
cd D:\dev\PocoClinic
git checkout -b feature/clinic-emr-foundation
```

---

## Principles

1. **Never commit secrets or PHI** — run the secret check before every commit:
   ```powershell
   powershell -NoProfile -ExecutionPolicy Bypass -File scripts/check-no-secrets.ps1
   ```
2. **Never `git add -f`** on `backend/.env`, `backend/data/`, `loadtest/setup.env`, or `loadtest/personas.env`.
3. **Prefer `git add <paths>`** over blind `git add -A` until the tree is clean.
4. **Each commit should tell one story** — docs, persistence, backend feature, frontend feature, ops, tooling.
5. **Intermediate commits may not build in isolation** — that's OK on a feature branch; the **final tree** must pass `test-all.bat`.

---

## Dev tooling vs clinic runtime (read this first)

The repo intentionally separates **who runs a script** from **where it lives in git**:

| Audience | What | Git location | On Pi tarball? |
|----------|------|--------------|----------------|
| **Developers / CI** | Build release, run tests, lint, secret scan | [`scripts/`](../scripts/README.md), root `run-*.bat`, `test-all.bat` | No |
| **Clinic server (post-deploy)** | Install, migrate, backup, restore, cron, env wrappers | [`clinic/server/`](../clinic/server/README.md) | **Yes** → `/opt/pococlinic` |
| **Admin workstation (optional)** | Print paper binders before go-live | [`clinic/workstation/`](../clinic/workstation/README.md) | No |
| **Windows dev mirrors** | Same ops flows via `go run` for local testing | [`clinic/dev-windows/`](../clinic/dev-windows/README.md) | No |

**Rule of thumb:** if a clinic administrator or cron job runs it **after** the release is installed, it belongs under **`clinic/`** (usually `clinic/server/`). If a contributor runs it **while building** the app, it belongs under **`scripts/`** or repo-root dev launchers.

Application **source** stays in `backend/`, `frontend/`, `ops-helper/`. [`scripts/build-release.mjs`](../scripts/build-release.mjs) compiles Go binaries **and** copies `clinic/server/` into `dist/pococlinic-*-linux-*.tar.gz`.

Repo-root `backup.bat`, `migrate.bat`, etc. are **thin forwarders** to `clinic/dev-windows/` — not duplicated logic.

Decision record: [ADR-0016](../adr/0016-clinic-runtime-packaging.md).

### Documentation that explains the split

Keep these in sync when adding ops tools:

| File | Role |
|------|------|
| [`clinic/README.md`](../clinic/README.md) | Hub — layout, packaging table, dev vs server |
| [`docs/ops/tools-and-scripts.md`](../docs/ops/tools-and-scripts.md) | Full catalog (dev + deployed commands) |
| [`docs/deploy/DEPLOYMENT-BOUNDARY.md`](../docs/deploy/DEPLOYMENT-BOUNDARY.md) | Tarball contents vs microSD image-build |
| [`scripts/README.md`](../scripts/README.md) | Build/CI only; points at `clinic/server/` |
| [`scripts/release/README.md`](../scripts/release/README.md) | Stub → `clinic/server/` |
| Root [`README.md`](../README.md), [`AGENTS.md`](../AGENTS.md) | Repo map includes `clinic/` |

---

## Commit message format (Windows)

PowerShell does not support bash heredocs. For each commit below, either:

**Option A — Git Bash / WSL:** use the `$(cat <<'EOF' ... EOF)` blocks as written.

**Option B — PowerShell:** use two `-m` flags (subject + body), as shown in Commit 4.

**Option C — single line:** `git commit -m "feat(scope): short subject"` when the body is obvious.

---

## Pre-flight (once, before commit 1)

```powershell
cd D:\dev\PocoClinic

# Confirm ignore rules
powershell -File scripts/check-no-secrets.ps1

# Full verification baseline
.\test-all.bat

# Stage deletion of orphan lockfile (no backend package.json)
git rm backend/package-lock.json

# Optional: see what's left
git status --short | Measure-Object -Line
```

**Do not commit:** `backend/data/`, any `.env` with real values, loadtest report CSVs/JSON (except `loadtest/reports/.gitkeep`).

---

## Commit 1 — Architecture, ADRs, and repo navigation

**Story:** Move ADRs to repo root (including clinic runtime packaging), add contributor/agent orientation, tighten `.gitignore`, root README map with `clinic/`.

```powershell
git add `
  adr/ `
  docs/adr/ `
  AGENTS.md `
  README.md `
  .gitignore `
  .env.example `
  .cursor/rules/

git commit -m "$(cat <<'EOF'
docs: move ADRs to repo root and refresh monorepo navigation

Canonical architecture decisions live under adr/ (including ADR-0016 clinic
runtime packaging). Add AGENTS.md, expand README with clinic/ vs scripts/
split, and update gitignore for local PHI, SQLite journals, backups, and env.
EOF
)"
```

---

## Commit 2 — Product and ops documentation

**Story:** Docs hub, vision, features status, guides, ops binders, deploy notes — including `clinic/` references and deployed vs dev tooling.

```powershell
git add docs/

git commit -m "$(cat <<'EOF'
docs: add product vision, admin guides, and ops runbooks

Introduce docs hub (VISION, FEATURES, NETWORK-AND-SECURITY, guides, ops,
deploy boundary with clinic/server tarball layout). Document dev vs post-deploy
tooling split (ADR-0016). Remove stale Docker-compose migrate/backup examples;
align admin guides and Pi touchscreen docs with clinic/server bin wrappers.
EOF
)"
```

---

## Commit 3 — SQLite persistence and shared backend packages

**Story:** Database engine, migrations, config/dotenv, shared infra packages.

```powershell
git add `
  backend/go.mod `
  backend/go.sum `
  backend/.env.example `
  backend/internal/pkg/database/ `
  backend/internal/pkg/config/ `
  backend/internal/pkg/doccrypto/ `
  backend/internal/pkg/httperr/ `
  backend/internal/pkg/filename/ `
  backend/internal/pkg/pagination/ `
  backend/internal/pkg/logging/ `
  backend/internal/pkg/errors/ `
  backend/internal/pkg/middleware/recovery.go `
  backend/internal/pkg/middleware/security.go `
  backend/internal/pkg/security/ `
  backend/internal/pkg/staticserve/ `
  backend/cmd/migrate/ `
  docker-compose.yml

git commit -m "$(cat <<'EOF'
feat(backend): SQLite persistence, migrations, and shared infra packages

Add embedded SQLite via modernc.org/sqlite, squashed migrations, dotenv
loading, document encryption helpers, httperr, and security middleware.
Empty docker-compose documents that no DB container is required.
EOF
)"
```

**After this commit (local dev):** `cd backend && go run ./cmd/migrate` when `DATABASE_URL` is set.

**Note:** `migrate.bat` (forwarder to `clinic/dev-windows/`) lands in **Commit 10** with the rest of `clinic/`.

---

## Commit 4 — Authentication (badge + PIN, cookies, RBAC)

**Story:** Staff/admin login, session cookies, user management, SQL repos.

```powershell
git add backend/internal/features/auth/

git commit -m "feat(auth): badge+PIN login, cookie sessions, and staff RBAC" -m "Implement staff and admin login flows, refresh tokens, PIN change, user CRUD, badge reissue, SQL repositories, and auth middleware with role gates per ADR-0011."
```

---

## Commit 5 — Patients, notes, documents, exercise log, field settings

**Story:** Clinical chart backend including admin-configurable required fields.

```powershell
git add `
  backend/internal/features/patients/

git commit -m "$(cat <<'EOF'
feat(patients): chart CRUD, notes, documents, exercise log, field settings

Add patient list filters, encrypted document storage, clinical notes,
exercise plans/sessions, and clinic-wide configurable required patient
fields (clinic_settings + validation on create/update).
EOF
)"
```

---

## Commit 6 — Forms, audit, and admin operations

**Story:** Form builder/submissions, audit trail, admin dashboard API, backup hooks.

```powershell
git add `
  backend/internal/features/forms/ `
  backend/internal/features/audit/ `
  backend/internal/features/admin/ `
  backend/internal/pkg/auditretention/ `
  backend/internal/pkg/backup/ `
  backend/internal/pkg/opshelper/ `
  backend/cmd/backup/ `
  backend/cmd/restore/ `
  backend/cmd/audit-purge/ `
  backend/cmd/ops-helper/ `
  build-ops-helper.bat `
  build-ops-helper-pi.bat

git commit -m "$(cat <<'EOF'
feat(admin): forms, audit, backups, and clinic admin API

Add form templates/groups/submissions, audit logging and retention, admin
status/health/compliance endpoints, SQLite backup/restore commands, and
localhost ops-helper UI wiring.
EOF
)"
```

**Note:** `backup.bat`, `restore.bat`, `run-ops-helper*.bat` forwarders land in **Commit 10** under `clinic/dev-windows/`.

---

## Commit 7 — Wire API server (`cmd/main`)

**Story:** Single entrypoint registering all features (depends on commits 3–6).

```powershell
git add `
  backend/cmd/main.go `
  run-backend.bat

git commit -m "$(cat <<'EOF'
feat(backend): wire modular monolith entrypoint and routes

Register auth, patients, forms, audit, admin, exercise log, and static
frontend handlers in cmd/main with SQLite-backed repositories and
in-memory fallbacks when DATABASE_URL is unset.
EOF
)"
```

---

## Commit 8 — Frontend platform (Vite, Vitest, ESLint, API client, auth)

**Story:** Replace Jest/MSW skeleton with Vitest, cookie auth client, routing shell.

```powershell
git add `
  frontend/package.json `
  frontend/package-lock.json `
  frontend/vite.config.ts `
  frontend/tsconfig.json `
  frontend/tsconfig.node.json `
  frontend/eslint.config.js `
  frontend/.prettierrc.json `
  frontend/.prettierignore `
  frontend/index.html `
  frontend/src/main.tsx `
  frontend/src/App.tsx `
  frontend/src/config.ts `
  frontend/src/queryClient.ts `
  frontend/src/setupTests.ts `
  frontend/src/env.d.ts `
  frontend/src/vite-env.d.ts `
  frontend/src/api/ `
  frontend/src/context/ `
  frontend/src/utils/ `
  frontend/src/components/ErrorBoundary.tsx `
  frontend/src/components/auth/ `
  frontend/src/components/layout/ `
  frontend/src/layout/ `
  frontend/jest.config.js `
  frontend/src/__mocks__/ `
  frontend/src/mocks/ `
  frontend/src/types/api.ts

git commit -m "$(cat <<'EOF'
feat(frontend): Vitest, ESLint, cookie-auth API client, and app shell

Migrate from Jest/MSW stub to Vitest, add ESLint flat config, axios client
with refresh handling, AuthContext, protected routes, and lazy-loaded
routing layout for clinic LAN browsers.
EOF
)"
```

Note: `jest.config.js`, `__mocks__/`, `mocks/`, and `types/api.ts` are **deletions** — include them so the commit records the removal.

---

## Commit 9 — Frontend clinical and admin UI

**Story:** Pages and components staff/admins use daily.

```powershell
git add `
  frontend/src/pages/ `
  frontend/src/components/patients/ `
  frontend/src/components/forms/ `
  frontend/src/components/users/ `
  frontend/src/components/admin/ `
  frontend/src/components/account/ `
  frontend/src/components/audit/ `
  frontend/src/components/help/ `
  frontend/src/help/ `
  frontend/src/styles/ `
  frontend/src/types/ `
  frontend/src/__tests__/ `
  run-frontend.bat

git commit -m "$(cat <<'EOF'
feat(frontend): patient chart, admin dashboard, forms, and staff UI

Add patient list/chart, notes/documents/exercise log sections, form
renderer, staff user management with badge QR, admin dashboard (backup,
compliance, clinic field settings), and in-app help content.
EOF
)"
```

---

## Commit 10 — Clinic runtime, ops apps, devices, loadtest, scripts, CI

**Story:** Post-deploy tooling (`clinic/`), ops-helper and binder-printer apps, Pi docs, build/CI scripts, and repo-root forwarders. This commit is where **dev vs deployed** boundaries become concrete in the tree.

```powershell
git add `
  clinic/ `
  ops-helper/ `
  binder-printer/ `
  devices/ `
  loadtest/ `
  scripts/ `
  .github/ `
  test-all.bat `
  run-all.bat `
  build-release.bat `
  backup.bat `
  restore.bat `
  migrate.bat `
  run-ops-helper.bat `
  run-ops-helper-pi.bat `
  print-binder.bat `
  COMMIT-STRATEGY.md

git commit -m "$(cat <<'EOF'
chore: clinic runtime packaging, ops apps, Pi docs, loadtest, CI

Add clinic/ (server install tree for tarball, dev-windows mirrors,
workstation print-binder). Wire build-release.mjs to copy clinic/server/.
Add ops-helper, binder-printer, Raspberry Pi ops docs, Locust loadtest,
secret-scan scripts, GitHub Actions, and repo-root .bat forwarders.
EOF
)"
```

**What lands in `clinic/`:**

| Path | Shipped on Pi? |
|------|----------------|
| `clinic/server/` (install.sh, env.template, cron, bin/) | Yes |
| `clinic/dev-windows/` | No — local `go run` only |
| `clinic/workstation/` | No — admin PC print utility |

---

## Final verification (after commit 10)

```powershell
powershell -File scripts/check-no-secrets.ps1
.\test-all.bat
git status   # should be clean (or only COMMIT-STRATEGY.md if you defer it)
git log --oneline -10
```

---

## Optional adjustments

### Fewer commits (squash pairs)

| Merge | Result |
|-------|--------|
| 1 + 2 | Single **docs** commit |
| 3 + 7 | **backend infra + main** together (large) |
| 8 + 9 | Single **frontend** commit |
| 9 + 10 | **UI + clinic runtime/ops** together |

Minimum sensible set: **5 commits** — docs → backend → frontend → clinic runtime + ops → CI/tooling.

### Single commit (fastest)

If you just want it on the branch tonight:

```powershell
powershell -File scripts/check-no-secrets.ps1
git rm backend/package-lock.json
git add -A
git commit -m "$(cat <<'EOF'
feat: clinic LAN EMR monorepo (SQLite, auth, patients, admin, frontend)

Private-LAN EMR with badge+PIN auth, SQLite persistence, patient chart
(notes/documents/exercise log), form builder, admin backup dashboard,
configurable patient required fields, ops helpers, and expanded docs/ADRs.
EOF
)"
```

Use this only if you don't need reviewable history.

---

## Push and PR (when ready)

```powershell
git push -u origin feature/clinic-emr-foundation
gh pr create --title "Clinic LAN EMR foundation" --body-file .github/PULL_REQUEST_TEMPLATE.md
```

If no PR template exists, paste summary from `docs/FEATURES.md` and the test plan from `test-all.bat`.

---

## Post-merge clinic checklist (not part of git)

- [ ] Build tarball: `node scripts/build-release.mjs` (or `build-release.bat`)
- [ ] On Pi: extract tarball, `sudo ./install.sh`, edit `/etc/pococlinic/env`
- [ ] On Pi: extract tarball, `sudo ./install.sh`, edit `/etc/pococlinic/env`
- [ ] Run `sudo /opt/pococlinic/bin/migrate`
- [ ] Enable cron: `sudo cp /opt/pococlinic/cron/pococlinic-backup /etc/cron.d/`
- [ ] Set strong `JWT_*` secrets and `DOCUMENT_ENCRYPTION_KEY` in production env
- [ ] Delete `bootstrap-admin-once.txt` after admin setup
- [ ] Confirm USB backup drill on a test stick

---

## Delete this file?

`COMMIT-STRATEGY.md` is a **one-time operator guide**. Either:

- Commit it in **Commit 10** (helps future you), or
- Keep it local and delete after pushing.

Do **not** leave real credentials written in this file.
