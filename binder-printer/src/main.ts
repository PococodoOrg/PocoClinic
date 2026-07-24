import { marked } from 'marked';
import { binderPacks, type BinderPack, type BinderPage } from './manifest';
import { markdownByFile } from './content';

marked.use({ async: false });

function loadMarkdown(file: string): string {
  const body = markdownByFile[file];
  if (body) {
    return body;
  }
  return `# Missing file\n\nCould not load \`${file}\`. Add it to \`src/content.ts\` and \`src/manifest.ts\`.`;
}

type State = {
  selectedPackIds: Set<string>;
  selectedPageIds: Set<string>;
  step: 'pick' | 'pages' | 'preview';
};

const state: State = {
  selectedPackIds: new Set(),
  selectedPageIds: new Set(),
  step: 'pick',
};

function pagesForSelection(): BinderPage[] {
  const pages: BinderPage[] = [];
  for (const pack of binderPacks) {
    if (!state.selectedPackIds.has(pack.id)) continue;
    for (const page of pack.pages) {
      if (state.selectedPageIds.has(page.id)) {
        pages.push(page);
      }
    }
  }
  return pages;
}

function selectAllPagesForPacks(packs: BinderPack[]) {
  for (const pack of packs) {
    for (const page of pack.pages) {
      state.selectedPageIds.add(page.id);
    }
  }
}

function renderApp() {
  const app = document.getElementById('app');
  if (!app) return;

  if (state.step === 'pick') {
    app.innerHTML = `
      <div class="shell no-print">
        <header class="hero">
          <p class="eyebrow">PocoClinic · not part of the EMR</p>
          <h1>Binder printer</h1>
          <p class="lede">Choose which paper packs to print. Preview, then use your browser print dialog.</p>
        </header>
        <section class="cards">
          ${binderPacks
            .map(
              (pack) => `
            <label class="card ${state.selectedPackIds.has(pack.id) ? 'card--on' : ''}">
              <input type="checkbox" data-pack="${pack.id}" ${state.selectedPackIds.has(pack.id) ? 'checked' : ''} />
              <div>
                <strong>${escapeHtml(pack.name)}</strong>
                <span class="meta">${escapeHtml(pack.deviceLabel)}</span>
                <p>${escapeHtml(pack.description)}</p>
              </div>
            </label>`,
            )
            .join('')}
        </section>
        <footer class="actions">
          <button type="button" class="primary" id="to-pages" ${state.selectedPackIds.size === 0 ? 'disabled' : ''}>
            Next: choose pages
          </button>
        </footer>
      </div>`;

    app.querySelectorAll<HTMLInputElement>('input[data-pack]').forEach((input) => {
      input.addEventListener('change', () => {
        const id = input.dataset.pack!;
        if (input.checked) state.selectedPackIds.add(id);
        else state.selectedPackIds.delete(id);
        renderApp();
      });
    });
    app.querySelector('#to-pages')?.addEventListener('click', () => {
      const packs = binderPacks.filter((p) => state.selectedPackIds.has(p.id));
      state.selectedPageIds = new Set();
      selectAllPagesForPacks(packs);
      state.step = 'pages';
      renderApp();
    });
    return;
  }

  if (state.step === 'pages') {
    const packs = binderPacks.filter((p) => state.selectedPackIds.has(p.id));
    app.innerHTML = `
      <div class="shell no-print">
        <header class="hero">
          <p class="eyebrow">Step 2 of 3</p>
          <h1>Choose pages</h1>
          <p class="lede">Everything is selected by default. Uncheck pages you do not need.</p>
        </header>
        ${packs
          .map((pack) => {
            const bySection = groupBySection(pack.pages);
            return `
              <section class="pack-block">
                <h2>${escapeHtml(pack.name)}</h2>
                ${Object.entries(bySection)
                  .map(
                    ([section, pages]) => `
                  <h3>${escapeHtml(section)}</h3>
                  <ul class="page-list">
                    ${pages
                      .map(
                        (page) => `
                      <li>
                        <label>
                          <input type="checkbox" data-page="${page.id}" ${state.selectedPageIds.has(page.id) ? 'checked' : ''} />
                          <span>${escapeHtml(page.title)}</span>
                        </label>
                        ${page.warning ? `<p class="warn">${escapeHtml(page.warning)}</p>` : ''}
                      </li>`,
                      )
                      .join('')}
                  </ul>`,
                  )
                  .join('')}
              </section>`;
          })
          .join('')}
        <footer class="actions">
          <button type="button" class="ghost" id="back-pick">Back</button>
          <button type="button" class="primary" id="to-preview" ${state.selectedPageIds.size === 0 ? 'disabled' : ''}>
            Preview &amp; print
          </button>
        </footer>
      </div>`;

    app.querySelectorAll<HTMLInputElement>('input[data-page]').forEach((input) => {
      input.addEventListener('change', () => {
        const id = input.dataset.page!;
        if (input.checked) state.selectedPageIds.add(id);
        else state.selectedPageIds.delete(id);
        renderApp();
      });
    });
    app.querySelector('#back-pick')?.addEventListener('click', () => {
      state.step = 'pick';
      renderApp();
    });
    app.querySelector('#to-preview')?.addEventListener('click', () => {
      state.step = 'preview';
      renderApp();
    });
    return;
  }

  // preview
  const pages = pagesForSelection();
  const articles = pages
    .map((page) => {
      const md = loadMarkdown(page.file);
      const html = marked.parse(md) as string;
      return `
        <article class="print-page" data-page="${escapeHtml(page.id)}">
          <header class="print-banner no-print">
            <span>${escapeHtml(page.section)}</span>
            <strong>${escapeHtml(page.title)}</strong>
            ${page.warning ? `<em class="warn">${escapeHtml(page.warning)}</em>` : ''}
          </header>
          <div class="md-body">${html}</div>
        </article>`;
    })
    .join('');

  app.innerHTML = `
    <div class="shell">
      <header class="hero no-print">
        <p class="eyebrow">Step 3 of 3</p>
        <h1>Print preview</h1>
        <p class="lede">${pages.length} page(s) ready. Use Print → PDF or your clinic printer.</p>
        <div class="actions">
          <button type="button" class="ghost" id="back-pages">Back</button>
          <button type="button" class="primary" id="do-print">Print</button>
        </div>
      </header>
      <div class="preview-stack">${articles}</div>
    </div>`;

  app.querySelector('#back-pages')?.addEventListener('click', () => {
    state.step = 'pages';
    renderApp();
  });
  app.querySelector('#do-print')?.addEventListener('click', () => window.print());
}

function groupBySection(pages: BinderPage[]): Record<string, BinderPage[]> {
  const out: Record<string, BinderPage[]> = {};
  for (const page of pages) {
    (out[page.section] ??= []).push(page);
  }
  return out;
}

function escapeHtml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;');
}

renderApp();
