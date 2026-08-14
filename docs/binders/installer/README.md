# PocoClinic installer binder

**For deploy technicians and integrators** — not for day-to-day clinic staff.

Complete this binder during **install or server replacement**, then **hand off** using Section **I5**. The clinic keeps the separate [**site operations binder**](../../ops/binder/README.md).

Secrets: print the blank [Safe vault](../../ops/safe-credentials-vault.md) → fill → **safe** (not loose in this binder).

| Clinic | ________________________________ |
|--------|----------------------------------|
| Installer / company | ________________________________ |
| Server type | ☐ Raspberry Pi ☐ Small PC |
| LAN IP (Pi/PC) | ________________________________ |
| Staff URL | `https://________________________` |
| Install date | ____________ |

---

## When to use this binder

| Step | Section |
|------|---------|
| Before you touch hardware | **I1** Prerequisites |
| Router / Wi‑Fi / DNS | **I2** Network & Wi‑Fi |
| Tarball, migrate, systemd | **I3** Server install |
| Firewall, HTTPS, CA trust | **I4** Hardening & TLS |
| Sign-off to clinic admin | **I5** Handoff checklist |

If the server is a **Raspberry Pi**, also print the [Pi device binder](../../../devices/raspberry-pi/binder/README.md) and leave it in the cabinet.

**First Pi install?** Screen-by-screen: [devices/raspberry-pi/setup.md](../../../devices/raspberry-pi/setup.md).

**Print path:** [binder-printer](../../../binder-printer/README.md) → **Installer binder**.

---

## Print order

Use tab dividers **I1 · I2 · I3 · I4 · I5**.

| Page | File |
|------|------|
| Cover | This page |
| I1 | [Prerequisites](./I1-prerequisites.md) |
| I2 | [Network & Wi‑Fi](./I2-network-and-wifi.md) |
| I3 | [Server install](./I3-server-install.md) |
| I4 | [Hardening & TLS](./I4-hardening-and-tls.md) |
| I5 | [Handoff checklist](./I5-handoff-checklist.md) |

### Print checklist

- [ ] Cover — fill clinic name, IP, URL, date
- [ ] I1–I5
- [ ] Pi device binder (if applicable)
- [ ] Blank safe vault sheet → fill at handoff → lock in safe
- [ ] Site operations binder printed for clinic desk (separate pack)

---

## After handoff

Give the clinic:

1. Filled **site operations binder** (store closed at admin desk)
2. **Pi device binder** in the server cabinet (Pi installs)
3. **USB backup drives** (labeled) + first verified backup
4. **CA certificate file** (`ca.crt`) on USB for trusting iPads/PCs — or MDM profile
5. Router/DNS notes on cover sheet (I5 copy for clinic IT)

Installer binder may be **archived** by the integrator — clinic does not need it for daily work.

---

## Related

- [Binders hub](../README.md)
- [Deployment boundary](../../deploy/DEPLOYMENT-BOUNDARY.md)
- [clinic/server/](../../../clinic/server/README.md)
