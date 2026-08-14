#!/usr/bin/env node
/**
 * Interactive PocoClinic setup wizard.
 *
 * Usage (from repo root):
 *   setup.bat
 *   ./setup.sh
 *   node scripts/setup.mjs
 *
 * Walks through device / target options. Raspberry Pi is the primary clinic path.
 */

import { spawnSync } from 'node:child_process';
import { createInterface } from 'node:readline/promises';
import { stdin as input, stdout as output } from 'node:process';
import crypto from 'node:crypto';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');

const rl = createInterface({ input, output });

function say(msg = '') {
  console.log(msg);
}

function blank() {
  say();
}

async function pause(message = 'Press Enter when ready…') {
  await rl.question(message);
}

async function ask(prompt, defaultValue = '') {
  const suffix = defaultValue ? ` [${defaultValue}]` : '';
  const answer = (await rl.question(`${prompt}${suffix}: `)).trim();
  return answer || defaultValue;
}

async function confirm(prompt, defaultYes = true) {
  const hint = defaultYes ? 'Y/n' : 'y/N';
  const answer = (await rl.question(`${prompt} (${hint}): `)).trim().toLowerCase();
  if (!answer) return defaultYes;
  return answer === 'y' || answer === 'yes';
}

function commandExists(cmd) {
  const probe = process.platform === 'win32' ? 'where' : 'which';
  const result = spawnSync(probe, [cmd], { encoding: 'utf8', shell: true });
  return result.status === 0;
}

function run(cmd, args, options = {}) {
  say(`\n> ${cmd} ${args.join(' ')}\n`);
  const result = spawnSync(cmd, args, {
    stdio: 'inherit',
    cwd: options.cwd || repoRoot,
    env: { ...process.env, ...options.env },
    shell: process.platform === 'win32',
  });
  return result.status === 0;
}

function banner() {
  say('');
  say('══════════════════════════════════════════════════');
  say('  PocoClinic setup');
  say('══════════════════════════════════════════════════');
  say('  Pull the repo → run one command → pick your path.');
  say('  Docs stay the source of truth; this wizard walks you.');
  say('');
}

async function checkPrereqs({ needGo = true, needNode = true } = {}) {
  say('Checking tools on this computer…');
  const nodeOk = commandExists('node');
  const npmOk = commandExists('npm');
  const goOk = commandExists('go');

  say(`  Node.js : ${nodeOk ? 'found' : 'MISSING — install from https://nodejs.org/ (LTS)'}`);
  say(`  npm     : ${npmOk ? 'found' : 'MISSING (comes with Node.js)'}`);
  if (needGo) {
    say(`  Go      : ${goOk ? 'found' : 'MISSING — install from https://go.dev/dl/ (1.25+)'}`);
  }

  if (needNode && (!nodeOk || !npmOk)) {
    say('\nInstall Node.js LTS, reopen this terminal, then run setup again.');
    return false;
  }
  if (needGo && !goOk) {
    say('\nInstall Go, reopen this terminal, then run setup again.');
    return false;
  }
  blank();
  return true;
}

function latestArm64Tarball() {
  const distDir = path.join(repoRoot, 'dist');
  if (!fs.existsSync(distDir)) return null;
  const matches = fs
    .readdirSync(distDir)
    .filter((name) => /^pococlinic-.*-linux-arm64\.tar\.gz$/.test(name))
    .map((name) => ({
      name,
      full: path.join(distDir, name),
      mtime: fs.statSync(path.join(distDir, name)).mtimeMs,
    }))
    .sort((a, b) => b.mtime - a.mtime);
  return matches[0] || null;
}

function writeSecretsHelper(piIp) {
  const access = crypto.randomBytes(48).toString('base64');
  const refresh = crypto.randomBytes(48).toString('base64');
  const docKey = crypto.randomBytes(32).toString('base64');
  const origin = `http://${piIp}:8080`;

  const distDir = path.join(repoRoot, 'dist');
  fs.mkdirSync(distDir, { recursive: true });
  const outPath = path.join(distDir, 'pi-first-boot-env-snippet.txt');

  const body = `# TEMPORARY helper for first Pi bring-up — DO NOT COMMIT
# Copy values into /etc/pococlinic/env on the Pi, then delete this file.
# Write these on the safe vault sheet. Lose DOCUMENT_ENCRYPTION_KEY ⇒ documents unreadable after restore.

JWT_ACCESS_SECRET=${access}
JWT_REFRESH_SECRET=${refresh}
DOCUMENT_ENCRYPTION_KEY=${docKey}
ALLOWED_ORIGIN=${origin}
OPS_HELPER_MAIN_APP_URL=${origin}
COOKIE_SECURE=false
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DATABASE_URL=/var/lib/pococlinic/pococlinic.db
BACKUP_DIR=/var/lib/pococlinic/backups
DOCUMENTS_DIR=/var/lib/pococlinic/documents
`;

  fs.writeFileSync(outPath, body, { encoding: 'utf8', mode: 0o600 });
  return outPath;
}

async function setupLocalDev() {
  say('── Local development (this machine) ──\n');
  if (!(await checkPrereqs({ needGo: true, needNode: true }))) return;

  if (await confirm('Install frontend dependencies (npm ci)?', true)) {
    if (!run('npm', ['ci'], { cwd: path.join(repoRoot, 'frontend') })) {
      say('Frontend install failed.');
      return;
    }
  }

  if (await confirm('Download Go modules (go mod tidy)?', true)) {
    if (!run('go', ['mod', 'tidy'], { cwd: path.join(repoRoot, 'backend') })) {
      say('go mod tidy failed.');
      return;
    }
  }

  const envExample = path.join(repoRoot, 'backend', '.env.example');
  const envFile = path.join(repoRoot, 'backend', '.env');
  if (fs.existsSync(envExample) && !fs.existsSync(envFile)) {
    if (await confirm('Create backend/.env from .env.example?', true)) {
      fs.copyFileSync(envExample, envFile);
      say('Created backend/.env — edit secrets if needed.');
    }
  }

  say('\nYou are ready to develop.');
  say('  Windows:  run-all.bat');
  say('  Or:       run-backend.bat  and  run-frontend.bat');
  say('  Migrate:  migrate.bat  (set DATABASE_URL first)');
  say('\nMore: README.md → Getting started');

  if (process.platform === 'win32' && (await confirm('Start local app now (run-all.bat)?', false))) {
    run('cmd', ['/c', 'run-all.bat'], { cwd: repoRoot });
  }
}

async function setupRaspberryPi() {
  say('── Raspberry Pi clinic server ──\n');
  say('This path builds the ARM64 package on THIS computer, then you install it on the Pi.');
  say('Full written guide (always available): devices/raspberry-pi/setup.md\n');

  if (!(await checkPrereqs({ needGo: true, needNode: true }))) return;

  say('Hardware checklist (tick mentally):');
  say('  • Raspberry Pi 4/5 (8 GB preferred)');
  say('  • Official power supply + microSD (32 GB+)');
  say('  • Ethernet recommended');
  say('  • Another PC on the same network for the browser');
  say('  • Optional USB stick to copy the package');
  blank();
  await pause('Press Enter when you have the hardware ready (or want to continue anyway)…');

  // Step: flash OS
  say('\nSTEP 1 — Flash Raspberry Pi OS');
  say('  1. Install Raspberry Pi Imager: https://www.raspberrypi.com/software/');
  say('  2. Choose your Pi model + Raspberry Pi OS (64-bit)');
  say('  3. Edit settings BEFORE writing:');
  say('       Hostname: pococlinic');
  say('       Enable SSH + set username/password');
  say('       Wi‑Fi only if you will not use Ethernet');
  say('  4. Write the image, eject, insert card into the Pi, power on');
  blank();
  await pause('Press Enter when the Pi has booted (or skip if flashing later)…');

  // Step: build tarball
  say('\nSTEP 2 — Build the release package on this PC');
  const existing = latestArm64Tarball();
  if (existing) {
    say(`  Found existing package: dist\\${existing.name}`);
  }
  const shouldBuild = await confirm(
    existing ? 'Rebuild the ARM64 release tarball?' : 'Build the ARM64 release tarball now? (several minutes)',
    !existing,
  );
  if (shouldBuild) {
    if (!run('node', ['scripts/build-release.mjs'], { cwd: repoRoot })) {
      say('Release build failed. Fix errors above, then run setup again.');
      return;
    }
  }

  const tarball = latestArm64Tarball();
  if (!tarball) {
    say('No dist/pococlinic-*-linux-arm64.tar.gz found. Build must succeed before continuing.');
    return;
  }
  say(`\n  Package ready: ${tarball.full}`);

  // Step: secrets
  say('\nSTEP 3 — Generate first-boot secrets');
  const piIp = await ask('Pi LAN IP (or hostname if you will use it in the browser)', '192.168.1.50');
  const secretsPath = writeSecretsHelper(piIp);
  say(`\n  Wrote helper snippet (not for git):\n    ${secretsPath}`);
  say('  Copy those lines into /etc/pococlinic/env on the Pi after install.sh.');
  say('  Also write them on the safe vault sheet, then delete the helper file.');
  blank();
  await pause();

  // Step: copy to Pi
  say('\nSTEP 4 — Copy the package to the Pi');
  say('  Easiest: copy the .tar.gz to a USB stick, plug into the Pi, then:');
  say('    cp /media/$USER/*/pococlinic-*-linux-arm64.tar.gz ~/');
  say('\n  Or from PowerShell on this PC:');
  say(`    scp .\\dist\\${tarball.name} clinic@pococlinic.local:~/`);
  say('  (Replace clinic with your Pi username; use the IP if .local fails.)');
  blank();
  await pause('Press Enter when the .tar.gz is on the Pi…');

  // Step: install commands
  say('\nSTEP 5 — On the Pi (SSH or keyboard), run exactly:');
  say('');
  say('    cd ~');
  say('    tar -xzf pococlinic-*-linux-arm64.tar.gz');
  say('    cd pococlinic-*-linux-arm64');
  say('    sudo ./install.sh');
  say('    sudo nano /etc/pococlinic/env');
  say('      # paste secrets from the helper file; set ALLOWED_ORIGIN to match your browser URL');
  say('      # keep COOKIE_SECURE=false for first HTTP bring-up');
  say('    sudo /opt/pococlinic/bin/migrate');
  say('    sudo systemctl enable --now pococlinic pococlinic-ops-helper');
  say('    sudo systemctl status pococlinic --no-pager');
  say('');
  await pause();

  say('\nSTEP 6 — Sign in as admin');
  say(`  On another PC:  http://${piIp}:8080/health`);
  say(`  Then:           http://${piIp}:8080/login/admin`);
  say('  On the Pi:      sudo cat /var/lib/pococlinic/bootstrap-admin-once.txt');
  say('  Use Email + Access key + PIN (0000, then change).');
  say('  After login works:  sudo rm /var/lib/pococlinic/bootstrap-admin-once.txt');
  blank();
  await pause();

  say('\nSTEP 7 — Before real patients');
  say('  Enable LAN HTTPS and remove COOKIE_SECURE=false:');
  say('    devices/raspberry-pi/tls-lan.md');
  say('  Hardening:');
  say('    devices/raspberry-pi/hardening.md');
  say('  Staff, binders, first backup:');
  say('    docs/guide/for-administrators/first-time-setup.md');
  blank();

  say('Raspberry Pi bring-up checklist complete on this PC.');
  say(`  Package:  ${tarball.full}`);
  say(`  Secrets:  ${secretsPath}  ← delete after copying to the Pi / vault`);
  say('  Guide:    devices/raspberry-pi/setup.md');
}

async function setupSmallPc() {
  say('── Small PC / Linux x86 ──\n');
  say('Same release flow as the Pi, but build with --arch amd64:');
  say('  node scripts/build-release.mjs --arch amd64');
  say('Then follow clinic/server/install.sh on the Linux box.');
  say('\nA dedicated beginner walkthrough is not finished yet.');
  say('For now, prefer Raspberry Pi (option 2) or ask in the repo issues.');
  blank();
  await pause();
}

async function openDocsMenu() {
  say('── Documentation ──\n');
  say('  1) devices/raspberry-pi/setup.md     — Pi beginner install');
  say('  2) docs/guide/for-administrators/first-time-setup.md');
  say('  3) docs/VISION.md');
  say('  4) adr/README.md');
  say('  5) Back');
  const choice = await ask('Open which path in your editor/notes (number)', '5');
  const map = {
    1: 'devices/raspberry-pi/setup.md',
    2: 'docs/guide/for-administrators/first-time-setup.md',
    3: 'docs/VISION.md',
    4: 'adr/README.md',
  };
  const rel = map[choice];
  if (!rel) return;
  const full = path.join(repoRoot, rel);
  say(`\n  ${full}`);
  if (process.platform === 'win32') {
    spawnSync('cmd', ['/c', 'start', '', full], { shell: true, cwd: repoRoot });
  } else if (process.platform === 'darwin') {
    spawnSync('open', [full]);
  } else {
    spawnSync('xdg-open', [full], { stdio: 'ignore' });
  }
}

async function mainMenu() {
  for (;;) {
    banner();
    say('What do you want to set up?');
    say('');
    say('  1) Develop locally on this computer');
    say('  2) Raspberry Pi clinic server   ← primary clinic path');
    say('  3) Small PC / Linux x86 (preview)');
    say('  4) Open documentation');
    say('  q) Quit');
    say('');
    const choice = (await ask('Choice', '2')).toLowerCase();

    blank();
    if (choice === '1') await setupLocalDev();
    else if (choice === '2') await setupRaspberryPi();
    else if (choice === '3') await setupSmallPc();
    else if (choice === '4') await openDocsMenu();
    else if (choice === 'q' || choice === 'quit' || choice === 'exit') break;
    else {
      say('Unknown choice.');
      continue;
    }

    blank();
    if (!(await confirm('Return to the main menu?', true))) break;
  }
}

async function main() {
  try {
    process.chdir(repoRoot);
    await mainMenu();
    say('\nDone. Run setup again anytime:  setup.bat   or   node scripts/setup.mjs\n');
  } finally {
    rl.close();
  }
}

main().catch((err) => {
  console.error(err);
  rl.close();
  process.exit(1);
});
