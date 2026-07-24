# Moved: clinic server install tree

Production install scripts, env template, cron examples, and `/opt/pococlinic/bin/` wrappers now live under **[`clinic/server/`](../../clinic/server/README.md)**.

This folder is kept as a **pointer only** so older links to `scripts/release/` still resolve. Do not add new files here.

Release packaging copies from `clinic/server/` via [`scripts/build-release.mjs`](../build-release.mjs).
