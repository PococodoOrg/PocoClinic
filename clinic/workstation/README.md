# Admin workstation tools (optional)

Utilities for the **admin PC** during clinic setup — not part of the EMR and not required on the Pi server.

| Tool | Purpose | Shipped on Pi? |
|------|---------|----------------|
| `print-binder.bat` | Print ops/device binders from markdown | No — run from dev machine or copy to admin PC |

Source lives in [`binder-printer/`](../../binder-printer/README.md). The image-build process may optionally bundle `binder-printer/` on an admin laptop image; it is **not** included in the server release tarball.

Repo-root `print-binder.bat` forwards here for convenience.
