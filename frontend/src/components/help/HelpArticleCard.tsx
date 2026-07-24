import { Badge, Card, Group, Stack, Text, UnstyledButton } from '@mantine/core';
import { HelpArticle } from '../../help/types';
import { getCategoryById } from '../../help/content';

interface HelpArticleCardProps {
  article: HelpArticle;
  active?: boolean;
  compact?: boolean;
  onClick: () => void;
}

export function HelpArticleCard({ article, active, compact, onClick }: HelpArticleCardProps) {
  const category = getCategoryById(article.category);

  return (
    <UnstyledButton onClick={onClick} w="100%">
      <Card
        withBorder
        padding={compact ? 'sm' : 'md'}
        style={{
          borderColor: active ? 'var(--mantine-color-blue-6)' : undefined,
          background: active ? 'var(--mantine-color-blue-light)' : undefined,
          transition: 'border-color 120ms ease, background 120ms ease',
        }}
        className="help-article-card"
      >
        <Stack gap={4}>
          <Group justify="space-between" wrap="nowrap" align="flex-start">
            <Text fw={600} size={compact ? 'sm' : 'md'}>
              {article.title}
            </Text>
            {article.audience === 'admin' && (
              <Badge size="xs" variant="light">
                Admin
              </Badge>
            )}
          </Group>
          {!compact && (
            <>
              <Text size="sm" c="dimmed" lineClamp={2}>
                {article.summary}
              </Text>
              {category && (
                <Text size="xs" c="dimmed">
                  {category.label}
                </Text>
              )}
            </>
          )}
        </Stack>
      </Card>
    </UnstyledButton>
  );
}
