#!/usr/bin/env node
/**
 * Build a static HTML help site from docs/guide markdown for offline evaluator browsing.
 *
 * Usage: node scripts/build-help-site.mjs
 * Output: docs/guide-site/
 */

import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');
const guideDir = path.join(repoRoot, 'docs', 'guide');
const outDir = path.join(repoRoot, 'docs', 'guide-site');

function escapeHtml(text) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function markdownToHtml(markdown) {
  const lines = markdown.replace(/\r\n/g, '\n').split('\n');
  const html = [];
  let inCode = false;
  let inList = false;

  const closeList = () => {
    if (inList) {
      html.push('</ul>');
      inList = false;
    }
  };

  for (const line of lines) {
    if (line.startsWith('```')) {
      closeList();
      if (!inCode) {
        html.push('<pre><code>');
        inCode = true;
      } else {
        html.push('</code></pre>');
        inCode = false;
      }
      continue;
    }
    if (inCode) {
      html.push(`${escapeHtml(line)}\n`);
      continue;
    }

    if (line.startsWith('# ')) {
      closeList();
      html.push(`<h1>${inlineFormat(line.slice(2))}</h1>`);
      continue;
    }
    if (line.startsWith('## ')) {
      closeList();
      html.push(`<h2>${inlineFormat(line.slice(3))}</h2>`);
      continue;
    }
    if (line.startsWith('### ')) {
      closeList();
      html.push(`<h3>${inlineFormat(line.slice(4))}</h3>`);
      continue;
    }
    if (line.startsWith('- ')) {
      if (!inList) {
        html.push('<ul>');
        inList = true;
      }
      html.push(`<li>${inlineFormat(line.slice(2))}</li>`);
      continue;
    }
    if (line.trim() === '') {
      closeList();
      continue;
    }
    closeList();
    html.push(`<p>${inlineFormat(line)}</p>`);
  }
  closeList();
  if (inCode) {
    html.push('</code></pre>');
  }
  return html.join('\n');
}

function inlineFormat(text) {
  let out = escapeHtml(text);
  out = out.replace(/`([^`]+)`/g, '<code>$1</code>');
  out = out.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
  out = out.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>');
  return out;
}

function collectMarkdownFiles(dir, base = dir) {
  const entries = fs.readdirSync(dir, { withFileTypes: true });
  const files = [];
  for (const entry of entries) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...collectMarkdownFiles(full, base));
    } else if (entry.name.endsWith('.md')) {
      files.push({
        abs: full,
        rel: path.relative(base, full).replace(/\\/g, '/'),
      });
    }
  }
  return files.sort((a, b) => a.rel.localeCompare(b.rel));
}

function pageTemplate(title, bodyHtml, navLinks) {
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>${escapeHtml(title)} — PocoClinic Guide</title>
  <link rel="stylesheet" href="${navLinks ? '../site.css' : 'site.css'}" />
</head>
<body>
  <header class="site-header">
    <a href="${navLinks ? '../index.html' : 'index.html'}" class="site-title">PocoClinic Guide</a>
    <p class="site-tagline">Offline documentation for clinic evaluators and administrators</p>
  </header>
  <div class="layout">
    ${navLinks ? `<nav class="sidebar"><ul>${navLinks}</ul></nav>` : ''}
    <main class="content">${bodyHtml}</main>
  </div>
</body>
</html>`;
}

const css = `body { font-family: system-ui, sans-serif; line-height: 1.6; margin: 0; color: #1a1a1a; }
.site-header { background: #0b3d5c; color: #fff; padding: 1.25rem 2rem; }
.site-title { color: #fff; font-weight: 700; text-decoration: none; font-size: 1.25rem; }
.site-tagline { margin: 0.25rem 0 0; opacity: 0.85; font-size: 0.9rem; }
.layout { display: flex; max-width: 1100px; margin: 0 auto; padding: 1.5rem; gap: 2rem; }
.sidebar { flex: 0 0 220px; font-size: 0.9rem; }
.sidebar ul { list-style: none; padding: 0; margin: 0; }
.sidebar a { color: #0b3d5c; text-decoration: none; }
.sidebar a:hover { text-decoration: underline; }
.content { flex: 1; min-width: 0; }
.content h1 { margin-top: 0; }
.content pre { background: #f4f4f4; padding: 1rem; overflow-x: auto; border-radius: 4px; }
.content code { background: #f0f0f0; padding: 0.1em 0.35em; border-radius: 3px; font-size: 0.9em; }
.index-list { list-style: none; padding: 0; }
.index-list li { margin: 0.5rem 0; }
.index-list a { font-size: 1.05rem; }
`;

function main() {
  if (!fs.existsSync(guideDir)) {
    console.error('Guide directory not found:', guideDir);
    process.exit(1);
  }

  fs.rmSync(outDir, { recursive: true, force: true });
  fs.mkdirSync(outDir, { recursive: true });
  fs.writeFileSync(path.join(outDir, 'site.css'), css);

  const files = collectMarkdownFiles(guideDir);
  const pages = [];

  for (const file of files) {
    const markdown = fs.readFileSync(file.abs, 'utf8');
    const titleMatch = markdown.match(/^#\s+(.+)/m);
    const title = titleMatch ? titleMatch[1].trim() : file.rel.replace(/\.md$/, '');
    const htmlName = file.rel.replace(/\.md$/, '.html');
    const outPath = path.join(outDir, htmlName);
    fs.mkdirSync(path.dirname(outPath), { recursive: true });

    const depth = htmlName.split('/').length - 1;
    const prefix = depth > 0 ? '../'.repeat(depth) : '';
    const navLinks = files
      .map((f) => {
        const name = f.rel.replace(/\.md$/, '.html');
        const linkPrefix = depth > 0 ? '../'.repeat(depth) : '';
        const label = f.rel.replace(/\.md$/, '').replace(/\//g, ' › ');
        return `<li><a href="${linkPrefix}${name}">${escapeHtml(label)}</a></li>`;
      })
      .join('');

    const body = markdownToHtml(markdown);
    fs.writeFileSync(outPath, pageTemplate(title, body, navLinks.replace(prefix, prefix)));
    pages.push({ title, href: htmlName });
  }

  const indexBody = `<h1>PocoClinic documentation</h1>
<p>Static mirror of <code>docs/guide/</code> for prospective clinics evaluating the system on a LAN without the app running.</p>
<ul class="index-list">
${pages.map((p) => `<li><a href="${p.href}">${escapeHtml(p.title)}</a> <span style="color:#666">(${escapeHtml(p.href)})</span></li>`).join('\n')}
</ul>
<p>Open <code>index.html</code> from this folder in any browser. Regenerate after guide updates with <code>node scripts/build-help-site.mjs</code>.</p>`;

  fs.writeFileSync(path.join(outDir, 'index.html'), pageTemplate('Home', indexBody, null));

  console.log(`Built ${pages.length} pages → ${outDir}`);
}

main();
