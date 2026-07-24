# I2 — Network & Wi‑Fi

**Installer binder · Section I2**

Configure **clinic router / Wi‑Fi** so staff browsers can reach the server at a **stable HTTPS URL**. PocoClinic does not configure your router — this page is the checklist.

---

## Must configure

| Task | Done | Record here |
|------|------|-------------|
| Pi/PC on **same subnet** as staff tablets/PCs (not guest Wi‑Fi) | ☐ | Subnet: ____________ |
| **DHCP reservation** or static IP for server | ☐ | IP: ____________ |
| **Local DNS** or hosts: `pococlinic.local` → server IP | ☐ | |
| **No port forwarding** from internet to 443 / 8080 / 9090 | ☐ | Verified with clinic IT |
| **Client isolation OFF** on staff Wi‑Fi (AP isolation blocks LAN access) | ☐ | |
| Staff Wi‑Fi **WPA2/WPA3**; guest network **isolated** from clinic subnet | ☐ | |

---

## Staff URL (write on cover)

After TLS (Section I4), staff open:

```text
https://pococlinic.local
```

(or clinic-chosen hostname — must match TLS cert and `ALLOWED_ORIGIN`)

---

## Ports (LAN only)

| Port | Service | Exposed to Wi‑Fi? |
|------|---------|-------------------|
| **443** | EMR (HTTPS) | Yes — staff browsers |
| **8080** | EMR (HTTP) | Migrate only — **close after HTTPS cutover** |
| **9090** | Backup helper | **No** — localhost on server only |
| **22** | SSH admin | Restrict to admin subnet if possible |

---

## Recommended

| Task | Done |
|------|------|
| **Wired Ethernet** from Pi/PC to switch (Wi‑Fi for server is fallback) | ☐ |
| Document router admin login in **safe vault** (not this binder) | ☐ |
| Test from **staff iPad** on clinic Wi‑Fi before handoff | ☐ |

---

## Common failures

| Symptom | Likely cause |
|---------|----------------|
| iPad cannot reach server | Guest Wi‑Fi, AP isolation, wrong subnet |
| Certificate warning forever | CA not trusted on device (Section I4) |
| Login works on IP but not name | DNS / `.local` not configured |
| Backup helper “not found” from tablet | Expected — helper is **server localhost only** |

---

## Verification (from staff iPad on clinic Wi‑Fi)

1. `https://pococlinic.local/health` → OK JSON  
2. Padlock clean after CA trust  
3. `/login/admin` loads  

**Tested by:** ____________ **Date:** ____________

---

**Next:** Section **I3** — Server install
