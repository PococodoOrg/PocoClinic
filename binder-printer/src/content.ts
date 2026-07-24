/** Raw markdown loaded at build time (paths outside binder-printer/). */
import instCover from '../../docs/binders/installer/README.md?raw';
import instI1 from '../../docs/binders/installer/I1-prerequisites.md?raw';
import instI2 from '../../docs/binders/installer/I2-network-and-wifi.md?raw';
import instI3 from '../../docs/binders/installer/I3-server-install.md?raw';
import instI4 from '../../docs/binders/installer/I4-hardening-and-tls.md?raw';
import instI5 from '../../docs/binders/installer/I5-handoff-checklist.md?raw';
import clinicCover from '../../docs/ops/binder/README.md?raw';
import clinicA1 from '../../docs/ops/emergency-procedures.md?raw';
import clinicA2 from '../../docs/ops/binder/A2-paper-fallback.md?raw';
import clinicA3 from '../../docs/ops/binder/A3-incident-restore-log.md?raw';
import clinicB1 from '../../docs/ops/binder/B1-periodic-processes.md?raw';
import clinicB2 from '../../docs/ops/monthly-testing-checklist.md?raw';
import clinicB3 from '../../docs/ops/security-audit-checklist.md?raw';
import clinicB4 from '../../docs/ops/administrator-runbook.md?raw';
import clinicB5 from '../../docs/ops/physical-security-binder.md?raw';
import clinicB6 from '../../docs/ops/binder/B6-network-and-wifi-at-site.md?raw';
import clinicC1 from '../../docs/ops/binder/C1-what-is-pococlinic.md?raw';
import clinicC2 from '../../docs/ops/binder/C2-staff-how-to.md?raw';
import clinicC3 from '../../docs/ops/binder/C3-administrator-how-to.md?raw';
import piCover from '../../devices/raspberry-pi/binder/README.md?raw';
import piA1 from '../../devices/raspberry-pi/binder/A1-pi-emergency.md?raw';
import piA2 from '../../devices/raspberry-pi/binder/A2-kiosk-emergency.md?raw';
import piA3 from '../../devices/raspberry-pi/binder/A3-pi-incident-log.md?raw';
import piB1 from '../../devices/raspberry-pi/binder/B1-pi-periodic.md?raw';
import piC1 from '../../devices/raspberry-pi/binder/C1-hardware-cheatsheet.md?raw';
import piC2 from '../../devices/raspberry-pi/binder/C2-touchscreen-howto.md?raw';
import piC3 from '../../devices/raspberry-pi/binder/C3-paths-and-services.md?raw';
import safeVault from '../../docs/ops/safe-credentials-vault.md?raw';

/** Map of path keys used in manifest.ts `file` fields → markdown source */
export const markdownByFile: Record<string, string> = {
  '../docs/binders/installer/README.md': instCover,
  '../docs/binders/installer/I1-prerequisites.md': instI1,
  '../docs/binders/installer/I2-network-and-wifi.md': instI2,
  '../docs/binders/installer/I3-server-install.md': instI3,
  '../docs/binders/installer/I4-hardening-and-tls.md': instI4,
  '../docs/binders/installer/I5-handoff-checklist.md': instI5,
  '../docs/ops/binder/README.md': clinicCover,
  '../docs/ops/emergency-procedures.md': clinicA1,
  '../docs/ops/binder/A2-paper-fallback.md': clinicA2,
  '../docs/ops/binder/A3-incident-restore-log.md': clinicA3,
  '../docs/ops/binder/B1-periodic-processes.md': clinicB1,
  '../docs/ops/monthly-testing-checklist.md': clinicB2,
  '../docs/ops/security-audit-checklist.md': clinicB3,
  '../docs/ops/administrator-runbook.md': clinicB4,
  '../docs/ops/physical-security-binder.md': clinicB5,
  '../docs/ops/binder/B6-network-and-wifi-at-site.md': clinicB6,
  '../docs/ops/binder/C1-what-is-pococlinic.md': clinicC1,
  '../docs/ops/binder/C2-staff-how-to.md': clinicC2,
  '../docs/ops/binder/C3-administrator-how-to.md': clinicC3,
  '../devices/raspberry-pi/binder/README.md': piCover,
  '../devices/raspberry-pi/binder/A1-pi-emergency.md': piA1,
  '../devices/raspberry-pi/binder/A2-kiosk-emergency.md': piA2,
  '../devices/raspberry-pi/binder/A3-pi-incident-log.md': piA3,
  '../devices/raspberry-pi/binder/B1-pi-periodic.md': piB1,
  '../devices/raspberry-pi/binder/C1-hardware-cheatsheet.md': piC1,
  '../devices/raspberry-pi/binder/C2-touchscreen-howto.md': piC2,
  '../devices/raspberry-pi/binder/C3-paths-and-services.md': piC3,
  '../docs/ops/safe-credentials-vault.md': safeVault,
};
