# PocoClinic documentation guide

Clinic-facing guides for people evaluating PocoClinic, administrators setting it up, and staff learning daily workflows. For the full docs map (vision, ops, ADRs), see **[docs/README.md](../README.md)**. Repo front door: **[README.md](../../README.md)**.

This folder can be published as static web help later (GitHub Pages, MkDocs, etc.) without changing the content.

## Three ways to get help

| Layer | Audience | Where |
|-------|----------|--------|
| **In-app Help center** | Staff & admins at the clinic | Sidebar → Help & Support, or **?** in the header |
| **Administrator guide panel** | Clinic administrators | Admin dashboard → Administrator guide |
| **Repository guides (this folder)** | Evaluators, new admins, contributors | Markdown in `docs/guide/` |

In-app help works **offline on the LAN**. These markdown files are for reading in the repo, printing, or future web publishing.

## Static HTML site (evaluators)

Generate an offline-browsable mirror for prospective clinics:

```bash
node scripts/build-help-site.mjs
```

Output lands in `docs/guide-site/` — open `index.html` in any browser on the LAN.

## Start here

### Thinking about PocoClinic?

- [What is PocoClinic?](./evaluating-pococlinic.md) — who it is for, what it is not
- [Clinic network requirements](./clinic-network.md) — LAN, hardware, and security assumptions

### Running a clinic

- [For administrators](./for-administrators/README.md) — setup, daily ops, backup
- [For staff](./for-staff/README.md) — sign-in, patients, forms

### Technical depth

- [Product vision](../VISION.md)
- [Devices](../../devices/README.md) — hardware targets + device ops ([Raspberry Pi](../../devices/raspberry-pi/README.md))
- [Network & security](../NETWORK-AND-SECURITY.md)
- [Operations hub](../ops/README.md) — runbooks, checklists, **tools catalog**
- [Tools & scripts](../ops/tools-and-scripts.md) — every helper, CLI, and batch file
- [Deployment boundary](../deploy/DEPLOYMENT-BOUNDARY.md) — microSD / Pi packaging contract · [deploy hub](../deploy/README.md)
- [Clinic runtime (`clinic/`)](../../clinic/README.md) — post-deploy install tree ([ADR-0016](../../adr/0016-clinic-runtime-packaging.md))
- [Ops binder](../ops/binder/README.md) — A emergency · B weekly/monthly · C system help
- [Administrator runbook](../ops/administrator-runbook.md) — printable backup procedures
- [Safe vault credentials](../ops/safe-credentials-vault.md) — print, fill, lock in safe, delete digital copies
- [Architecture decisions](../../adr/README.md) — repo-root ADRs (GitHub-friendly)
- [Feature status](../FEATURES.md)

## Keeping docs aligned

When procedures change:

1. Update the **in-app article** in `frontend/src/help/content.ts`
2. Update the matching **admin guide task** in `frontend/src/help/adminGuide.ts` if needed
3. Update the **markdown guide** in this folder
4. Update `docs/ops/administrator-runbook.md` for printable binder content

See also [IN-APP-HELP.md](../IN-APP-HELP.md) for developer notes on the help system.
