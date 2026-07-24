import { Divider, Stack, Title } from '@mantine/core';
import { HelpArticleCard } from './HelpArticleCard';
import { HelpArticle } from '../../help/types';

interface HelpRelatedArticlesProps {
  articles: HelpArticle[];
  onSelect: (id: string) => void;
}

export function HelpRelatedArticles({ articles, onSelect }: HelpRelatedArticlesProps) {
  if (articles.length === 0) {
    return null;
  }

  return (
    <Stack gap="sm" className="help-no-print">
      <Divider />
      <Title order={5}>Related articles</Title>
      <Stack gap="xs">
        {articles.map((article) => (
          <HelpArticleCard
            key={article.id}
            article={article}
            compact
            onClick={() => onSelect(article.id)}
          />
        ))}
      </Stack>
    </Stack>
  );
}
