#!/usr/bin/env node
/**
 * Build a PocoClinic release tarball for Linux ARM64 (Raspberry Pi).
 *
 * Usage:
 *   node scripts/build-release.mjs
 *   node scripts/build-release.mjs --version 1.0.0
 *   node scripts/build-release.mjs --arch amd64   # dev/testing on x86 Linux
 *
 * Output: dist/pococlinic-<version>-linux-<arch>.tar.gz
 */

import { createHash } from 'node:crypto';
import { spawnSync } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');

function parseArgs() {
  const args = process.argv.slice(2);
  let version = process.env.APP_VERSION || '';
  let arch = 'arm64';
  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--version' && args[i + 1]) {
      version = args[++i];
    } else if (args[i] === '--arch' && args[i + 1]) {
      arch = args[++i];
    }
  }
  if (!version) {
    const pkg = JSON.parse(fs.readFileSync(path.join(repoRoot, 'frontend/package.json'), 'utf8'));
    version = pkg.version || '1.0.0';
  }
  return { version, arch, goArch: arch === 'amd64' ? 'amd64' : 'arm64' };
}

function run(cmd, args, options = {}) {
  console.log(`> ${cmd} ${args.join(' ')}`);
  const result = spawnSync(cmd, args, {
    stdio: 'inherit',
    cwd: options.cwd || repoRoot,
    env: { ...process.env, ...options.env },
    shell: process.platform === 'win32',
  });
  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

function sha256File(filePath) {
  const hash = createHash('sha256');
  hash.update(fs.readFileSync(filePath));
  return hash.digest('hex');
}

function copyRecursive(src, dest) {
  fs.mkdirSync(dest, { recursive: true });
  for (const entry of fs.readdirSync(src, { withFileTypes: true })) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);
    if (entry.isDirectory()) {
      copyRecursive(srcPath, destPath);
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

function main() {
  const { version, arch, goArch } = parseArgs();
  const releaseName = `pococlinic-${version}-linux-${arch}`;
  const stagingRoot = path.join(repoRoot, 'dist', releaseName);
  const tarballPath = path.join(repoRoot, 'dist', `${releaseName}.tar.gz`);

  console.log(`Building ${releaseName}…`);

  if (fs.existsSync(stagingRoot)) {
    fs.rmSync(stagingRoot, { recursive: true, force: true });
  }
  fs.mkdirSync(stagingRoot, { recursive: true });

  // 1. Frontend EMR (same-origin API in production)
  run('npm', ['ci'], { cwd: path.join(repoRoot, 'frontend') });
  run('npm', ['run', 'build'], {
    cwd: path.join(repoRoot, 'frontend'),
    env: { VITE_API_BASE_URL: '/api/v1' },
  });
  copyRecursive(path.join(repoRoot, 'frontend/dist'), path.join(stagingRoot, 'static'));

  // 2. Ops helper Pi touch UI → embed dir, then build binary
  run('npm', ['ci'], { cwd: path.join(repoRoot, 'ops-helper') });
  run('npm', ['run', 'build:pi'], { cwd: path.join(repoRoot, 'ops-helper') });
  const opsStatic = path.join(repoRoot, 'backend/cmd/ops-helper/static');
  fs.rmSync(opsStatic, { recursive: true, force: true });
  copyRecursive(path.join(repoRoot, 'ops-helper/dist'), opsStatic);
  copyRecursive(path.join(repoRoot, 'ops-helper/dist'), path.join(stagingRoot, 'ops-helper-static'));

  // 3. Go binaries (linux cross-compile)
  const goEnv = {
    GOOS: 'linux',
    GOARCH: goArch,
    CGO_ENABLED: '0',
  };
  const backendDir = path.join(repoRoot, 'backend');
  const bins = [
    { src: './cmd/main.go', out: 'pococlinic' },
    { src: './cmd/ops-helper', out: 'ops-helper' },
    { src: './cmd/migrate', out: 'migrate' },
    { src: './cmd/backup', out: 'backup' },
    { src: './cmd/restore', out: 'restore' },
    { src: './cmd/audit-purge', out: 'audit-purge' },
  ];
  for (const bin of bins) {
    run(
      'go',
      ['build', '-o', path.join(stagingRoot, bin.out), bin.src],
      { cwd: backendDir, env: goEnv },
    );
  }

  // 4. systemd units, install script, env template, cron, bin wrappers
  fs.mkdirSync(path.join(stagingRoot, 'systemd'), { recursive: true });
  for (const unit of ['pococlinic.service', 'pococlinic-ops-helper.service']) {
    fs.copyFileSync(
      path.join(repoRoot, 'scripts/deploy', unit),
      path.join(stagingRoot, 'systemd', unit),
    );
  }
  fs.copyFileSync(
    path.join(repoRoot, 'clinic/server/install.sh'),
    path.join(stagingRoot, 'install.sh'),
  );
  fs.copyFileSync(
    path.join(repoRoot, 'clinic/server/env.template'),
    path.join(stagingRoot, 'env.template'),
  );
  copyRecursive(path.join(repoRoot, 'clinic/server/cron'), path.join(stagingRoot, 'cron'));
  copyRecursive(path.join(repoRoot, 'clinic/server/bin'), path.join(stagingRoot, 'bin'));
  copyRecursive(path.join(repoRoot, 'clinic/server/scripts'), path.join(stagingRoot, 'scripts'));
  copyRecursive(path.join(repoRoot, 'clinic/server/caddy'), path.join(stagingRoot, 'caddy'));

  // 5. Operator docs (subset)
  fs.mkdirSync(path.join(stagingRoot, 'docs'), { recursive: true });
  for (const doc of [
    'docs/ops/administrator-runbook.md',
    'docs/ops/tools-and-scripts.md',
    'docs/ops/monthly-testing-checklist.md',
    'docs/deploy/DEPLOYMENT-BOUNDARY.md',
  ]) {
    const base = path.basename(doc);
    fs.copyFileSync(path.join(repoRoot, doc), path.join(stagingRoot, 'docs', base));
  }

  // 6. MANIFEST.json + checksums
  const manifest = {
    name: 'pococlinic',
    version,
    platform: `linux-${arch}`,
    builtAt: new Date().toISOString(),
    files: [],
  };

  function walkManifest(dir, prefix = '') {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      const rel = prefix ? `${prefix}/${entry.name}` : entry.name;
      const full = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        walkManifest(full, rel);
      } else {
        manifest.files.push({
          path: rel.replace(/\\/g, '/'),
          sha256: sha256File(full),
          bytes: fs.statSync(full).size,
        });
      }
    }
  }
  walkManifest(stagingRoot);
  manifest.files.sort((a, b) => a.path.localeCompare(b.path));
  fs.writeFileSync(
    path.join(stagingRoot, 'MANIFEST.json'),
    `${JSON.stringify(manifest, null, 2)}\n`,
  );

  // 7. Tarball
  fs.mkdirSync(path.join(repoRoot, 'dist'), { recursive: true });
  if (fs.existsSync(tarballPath)) {
    fs.unlinkSync(tarballPath);
  }
  run('tar', ['-czf', tarballPath, '-C', path.join(repoRoot, 'dist'), releaseName]);

  console.log('');
  console.log(`Release ready: ${tarballPath}`);
  console.log(`Staging dir:   ${stagingRoot}`);
  console.log(`Files:         ${manifest.files.length}`);
}

main();
