export type BinderPage = {
  id: string;
  title: string;
  /** Path relative to binder-printer/ (starts with ../) */
  file: string;
  section: string;
  /** Optional warning shown in UI before including */
  warning?: string;
};

export type BinderPack = {
  id: string;
  name: string;
  description: string;
  /** Shown in device picker */
  deviceLabel: string;
  pages: BinderPage[];
};

export const binderPacks: BinderPack[] = [
  {
    id: 'clinic-ops',
    name: 'Clinic ops binder',
    deviceLabel: 'All clinics (device-independent)',
    description:
      'Emergency paper care, weekly/monthly checks, and staff/admin how-to. Store closed; open when needed.',
    pages: [
      { id: 'clinic-cover', title: 'Cover & print order', file: '../docs/ops/binder/README.md', section: 'Cover' },
      { id: 'clinic-a1', title: 'Emergency procedures', file: '../docs/ops/emergency-procedures.md', section: 'A — Emergency' },
      { id: 'clinic-a2', title: 'Paper care when EMR is down', file: '../docs/ops/binder/A2-paper-fallback.md', section: 'A — Emergency' },
      { id: 'clinic-a3', title: 'Incident & restore log', file: '../docs/ops/binder/A3-incident-restore-log.md', section: 'A — Emergency' },
      { id: 'clinic-b1', title: 'Weekly & monthly processes', file: '../docs/ops/binder/B1-periodic-processes.md', section: 'B — Weekly / monthly' },
      { id: 'clinic-b2', title: 'Monthly testing checklist', file: '../docs/ops/monthly-testing-checklist.md', section: 'B — Weekly / monthly' },
      { id: 'clinic-b3', title: 'Security audit checklist', file: '../docs/ops/security-audit-checklist.md', section: 'B — Weekly / monthly' },
      { id: 'clinic-b4', title: 'Backup & restore runbook', file: '../docs/ops/administrator-runbook.md', section: 'B — Weekly / monthly' },
      { id: 'clinic-b5', title: 'Physical security', file: '../docs/ops/physical-security-binder.md', section: 'B — Weekly / monthly' },
      { id: 'clinic-c1', title: 'What is PocoClinic', file: '../docs/ops/binder/C1-what-is-pococlinic.md', section: 'C — System help' },
      { id: 'clinic-c2', title: 'Staff how-to', file: '../docs/ops/binder/C2-staff-how-to.md', section: 'C — System help' },
      { id: 'clinic-c3', title: 'Administrator how-to', file: '../docs/ops/binder/C3-administrator-how-to.md', section: 'C — System help' },
    ],
  },
  {
    id: 'raspberry-pi',
    name: 'Raspberry Pi device binder',
    deviceLabel: 'Raspberry Pi server',
    description:
      'Pi-only emergency steps, periodic checks, and device help. Print in addition to the clinic ops binder.',
    pages: [
      { id: 'pi-cover', title: 'Cover & print order', file: '../devices/raspberry-pi/binder/README.md', section: 'Cover' },
      { id: 'pi-a1', title: 'Pi will not start / no network', file: '../devices/raspberry-pi/binder/A1-pi-emergency.md', section: 'A — Emergency' },
      { id: 'pi-a2', title: 'Backup kiosk blank or frozen', file: '../devices/raspberry-pi/binder/A2-kiosk-emergency.md', section: 'A — Emergency' },
      { id: 'pi-a3', title: 'Pi incident log', file: '../devices/raspberry-pi/binder/A3-pi-incident-log.md', section: 'A — Emergency' },
      { id: 'pi-b1', title: 'Pi periodic checks', file: '../devices/raspberry-pi/binder/B1-pi-periodic.md', section: 'B — Weekly / monthly' },
      { id: 'pi-c1', title: 'Hardware cheat sheet', file: '../devices/raspberry-pi/binder/C1-hardware-cheatsheet.md', section: 'C — Device help' },
      { id: 'pi-c2', title: 'Touchscreen backup how-to', file: '../devices/raspberry-pi/binder/C2-touchscreen-howto.md', section: 'C — Device help' },
      { id: 'pi-c3', title: 'Where things live on the Pi', file: '../devices/raspberry-pi/binder/C3-paths-and-services.md', section: 'C — Device help' },
    ],
  },
  {
    id: 'safe-vault',
    name: 'Safe vault credentials (blank template)',
    deviceLabel: 'Safe only (secrets)',
    description:
      'Print the blank template once, fill by hand, lock in the safe, delete digital copies. Never leave a filled copy in the printer.',
    pages: [
      {
        id: 'vault',
        title: 'Safe vault credentials sheet',
        file: '../docs/ops/safe-credentials-vault.md',
        section: 'SAFE',
        warning: 'Print blank only. Fill by hand. Lock in safe. Do not save a filled PDF.',
      },
    ],
  },
];
