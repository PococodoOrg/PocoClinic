import { HelpArticle, HelpAudience } from './types';
import { helpArticles } from './content';

function normalize(text: string): string {
  return text.toLowerCase();
}

function articleMatchesQuery(article: HelpArticle, query: string): boolean {
  const q = normalize(query.trim());
  if (!q) {
    return true;
  }

  const haystack = [
    article.title,
    article.summary,
    article.category,
    ...article.tags,
    ...article.sections.flatMap((section) => [
      section.heading ?? '',
      ...(section.paragraphs ?? []),
      ...(section.list ?? []),
      ...(section.steps ?? []),
    ]),
  ]
    .join(' ')
    .toLowerCase();

  return q.split(/\s+/).every((term) => haystack.includes(term));
}

function audienceAllowed(article: HelpArticle, isAdmin: boolean): boolean {
  if (article.audience === 'all') {
    return true;
  }
  return isAdmin && article.audience === 'admin';
}

export function filterHelpArticles(options: {
  query?: string;
  categoryId?: string;
  isAdmin: boolean;
  audience?: HelpAudience;
}): HelpArticle[] {
  const { query = '', categoryId, isAdmin } = options;

  return helpArticles.filter((article) => {
    if (!audienceAllowed(article, isAdmin)) {
      return false;
    }
    if (categoryId && article.category !== categoryId) {
      return false;
    }
    return articleMatchesQuery(article, query);
  });
}

export function searchHelpArticles(query: string, isAdmin: boolean): HelpArticle[] {
  if (!query.trim()) {
    return [];
  }

  return filterHelpArticles({ query, isAdmin }).slice(0, 12);
}
