import { describe, expect, it } from 'vitest';
import { filterHelpArticles, searchHelpArticles } from './search';

describe('help search', () => {
  it('finds articles by keyword for staff', () => {
    const results = searchHelpArticles('badge pin', false);
    expect(results.some((article) => article.id === 'sign-in')).toBe(true);
  });

  it('hides admin-only articles from staff', () => {
    const results = filterHelpArticles({ query: 'backup', isAdmin: false });
    expect(results.every((article) => article.audience !== 'admin')).toBe(true);
  });

  it('includes admin backup articles for administrators', () => {
    const results = filterHelpArticles({ query: 'backup', isAdmin: true });
    expect(results.some((article) => article.id === 'daily-backup')).toBe(true);
  });
});
