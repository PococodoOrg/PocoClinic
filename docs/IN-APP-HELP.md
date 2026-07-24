# Documentation (three layers)

PocoClinic documentation is intentionally split so it reaches the right audience in the right place.

## 1. In-app Help center (staff & admins on the LAN)

Works **offline** — content ships in the frontend bundle.

| Entry | Location |
|-------|----------|
| **Help center** | Sidebar → Help & Support, or `/help` |
| **Quick help drawer** | **?** button in the header |
| **Deep links** | `/help/:articleId` |
| **Print** | Printable articles (backup/restore) → Print button |

**Source files:** `frontend/src/help/content.ts`, `contextual.ts`, `search.ts`

## 2. Administrator guide panel (admins in the EMR)

On **Admin dashboard → Administrator guide**:

- First-time **setup checklist** (progress saved in browser localStorage)
- Accordion of **all admin tasks** with Open / Read guide links
- Links to in-app help articles for every setup and ops task

**Source files:** `frontend/src/help/adminGuide.ts`, `frontend/src/components/admin/AdminHelpPanel.tsx`

## 3. Repository markdown guides (evaluators & binders)

For people **thinking about PocoClinic**, new administrators reading ahead, or future **web-based help**:

| Folder | Audience |
|--------|----------|
| [adr/README.md](../adr/README.md) | **Architecture Decision Records** (why we chose X) |
| [docs/guide/](./guide/README.md) | Hub — evaluating, network, admin, staff |
| [docs/guide/evaluating-pococlinic.md](./guide/evaluating-pococlinic.md) | Prospective clinics |
| [docs/guide/for-administrators/](./guide/for-administrators/README.md) | Setup and ops |
| [docs/guide/for-staff/](./guide/for-staff/README.md) | Daily staff workflows |
| [docs/ops/administrator-runbook.md](./ops/administrator-runbook.md) | Printable binder (daily backup) |
| [docs/ops/binder/README.md](./ops/binder/README.md) | **Physical ops binder** — A emergency · B weekly/monthly · C system help |
| [devices/raspberry-pi/binder/README.md](../devices/raspberry-pi/binder/README.md) | **Pi device binder** — print-me A/B/C for the server |
| [binder-printer/README.md](../binder-printer/README.md) | Standalone utility to print binders — [`clinic/workstation/`](../clinic/workstation/README.md) |
| [docs/ops/safe-credentials-vault.md](./ops/safe-credentials-vault.md) | **Print → fill → lock in safe → delete digital copies** |
| [docs/ops/tools-and-scripts.md](./ops/tools-and-scripts.md) | Every CLI, batch file, and script |
| [devices/README.md](../devices/README.md) | Device ops + docs (Raspberry Pi, …) |
| [devices/raspberry-pi/README.md](../devices/raspberry-pi/README.md) | Pi hardware, kiosk, and device binder |
| [docs/deploy/DEPLOYMENT-BOUNDARY.md](./deploy/DEPLOYMENT-BOUNDARY.md) | microSD / Pi packaging contract |

These can be published as static HTML later (GitHub Pages, MkDocs, etc.) without changing the EMR.

## Adding or editing in-app articles

1. Add an entry to `helpArticles` in `frontend/src/help/content.ts` with a unique `id`.
2. If the task belongs on the admin panel, add or update a task in `frontend/src/help/adminGuide.ts`.
3. Optionally add the `id` to `getContextualArticleIds()` in `contextual.ts`.
4. Update the matching markdown file in `docs/guide/` when procedures change.
5. Run `npm test run src/help/search.test.ts`.

## Clinic emergency contacts

Article `clinic-contacts` lets administrators save names and phones in **local browser storage** on the clinic network.

## Keeping docs aligned

When procedures change, update **all three layers** that cover the topic:

1. `frontend/src/help/content.ts`
2. `frontend/src/help/adminGuide.ts` (if admin-facing)
3. Matching file under `docs/guide/` and/or `docs/ops/`
4. If you add a CLI or script, document it in `docs/ops/tools-and-scripts.md`

## Related

- [FEATURES.md](./FEATURES.md) — documentation checklist
- [VISION.md](./VISION.md) — product intent
- [Tools & scripts catalog](./ops/tools-and-scripts.md)
- [Deployment boundary (microSD / Pi)](./deploy/DEPLOYMENT-BOUNDARY.md)
