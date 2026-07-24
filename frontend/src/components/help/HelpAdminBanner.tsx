import { Button, Group, Paper, SimpleGrid, Stack, Text, Title } from '@mantine/core';
import { IconBook, IconShieldCheck } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { getArticleById } from '../../help/content';
import { HelpArticleCard } from './HelpArticleCard';

interface HelpAdminBannerProps {
  featuredIds: string[];
}

export function HelpAdminBanner({ featuredIds }: HelpAdminBannerProps) {
  const navigate = useNavigate();

  return (
    <Paper withBorder p="md" className="help-admin-banner help-no-print">
      <Group justify="space-between" align="flex-start" mb="md" wrap="wrap">
        <Stack gap={4}>
          <Group gap="xs">
            <IconShieldCheck size={20} />
            <Title order={4}>Administrator resources</Title>
          </Group>
          <Text size="sm" c="dimmed" maw={520}>
            Setup checklists, backup procedures, and staff management — also on the Admin dashboard.
          </Text>
        </Stack>
        <Button
          leftSection={<IconBook size={16} />}
          variant="light"
          onClick={() => navigate('/admin')}
        >
          Open admin dashboard
        </Button>
      </Group>
      <SimpleGrid cols={{ base: 1, sm: 2 }}>
        {featuredIds.map((id) => {
          const article = getArticleById(id);
          if (!article) {
            return null;
          }
          return (
            <HelpArticleCard
              key={id}
              article={article}
              compact
              onClick={() => navigate(`/help/${id}`)}
            />
          );
        })}
      </SimpleGrid>
    </Paper>
  );
}
