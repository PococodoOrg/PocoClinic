# Binder printer

Standalone **print utility** for PocoClinic paper binders.  
**Not part of the EMR** — run it on an admin PC when assembling binders.

← [Repo README](../README.md) · Clinic binder source: [docs/ops/binder](../docs/ops/binder/README.md) · Pi binder: [devices/raspberry-pi/binder](../devices/raspberry-pi/binder/README.md)

## What it does

1. Asks which binder(s) you need (clinic-wide ops, Raspberry Pi device, optional safe vault template)
2. Lets you tick sections/pages
3. Opens a print-ready preview → use the browser **Print** dialog (PDF or paper)

## Quick start (Windows)

From the repo root (forwards to [`clinic/workstation/print-binder.bat`](../clinic/workstation/print-binder.bat)):

```bat
print-binder.bat
```

Or:

```bat
cd binder-printer
npm install
npm run dev
```

Then open the URL shown (default http://localhost:5199).

## Binders included

| Pack | Source |
|------|--------|
| Clinic ops | `docs/ops/binder/` (+ linked runbooks/checklists) |
| Raspberry Pi | `devices/raspberry-pi/binder/` |
| Safe vault (blank template) | `docs/ops/safe-credentials-vault.md` — print blank only; fill by hand |

## Notes

- Requires Node.js 18+
- Reads Markdown from the repo over the Vite dev server (keep this folder inside the clone)
- Do not commit filled vault secrets into the repo
- Production clinics can print once, then store paper; no need to leave this utility running

## Related

- [Clinic ops binder](../docs/ops/binder/README.md)
- [Raspberry Pi binder](../devices/raspberry-pi/binder/README.md)
- [Devices hub](../devices/README.md)
