import {
  Badge,
  Button,
  Card,
  Group,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { fetchPerformance } from '../../api/admin';

export function PerformancePanel() {
  const { data, isFetching, refetch, isError } = useQuery({
    queryKey: ['admin-performance'],
    queryFn: fetchPerformance,
    enabled: false,
  });

  return (
    <Card withBorder padding="md">
      <Group justify="space-between" mb="md">
        <div>
          <Title order={4}>Performance snapshot</Title>
          <Text size="sm" c="dimmed">
            Memory, goroutines, and database pool usage — run during slow periods or after upgrades.
          </Text>
        </div>
        <Button loading={isFetching} onClick={() => refetch()} variant="light">
          Sample now
        </Button>
      </Group>

      {isError && (
        <Text c="red" size="sm">Could not load performance metrics.</Text>
      )}

      {data && (
        <Stack gap="sm">
          <SimpleGrid cols={{ base: 2, sm: 4 }}>
            <Metric label="Goroutines" value={String(data.goroutines)} />
            <Metric label="Memory (alloc)" value={`${data.memoryAllocMB.toFixed(1)} MB`} />
            <Metric label="Memory (sys)" value={`${data.memorySysMB.toFixed(1)} MB`} />
            {data.databasePool ? (
              <Metric
                label="DB pool in use"
                value={`${data.databasePool.inUseConns}/${data.databasePool.maxConns}`}
              />
            ) : (
              <Metric label="DB pool" value="n/a" />
            )}
          </SimpleGrid>
          <Text size="xs" c="dimmed">
            Sampled {new Date(data.collectedAt).toLocaleString()}
          </Text>
        </Stack>
      )}
    </Card>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <Stack gap={4}>
      <Text size="xs" c="dimmed">{label}</Text>
      <Badge size="lg" variant="light">{value}</Badge>
    </Stack>
  );
}
