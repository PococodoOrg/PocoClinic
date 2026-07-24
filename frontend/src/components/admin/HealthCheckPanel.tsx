import {
  Badge,
  Button,
  Card,
  Group,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { fetchHealthCheck } from '../../api/admin';
import { saveHealthCheckComplete } from './adminReminders';

function statusColor(status: string): string {
  switch (status) {
    case 'ok':
      return 'green';
    case 'warning':
      return 'yellow';
    case 'critical':
      return 'red';
    default:
      return 'gray';
  }
}

export function HealthCheckPanel() {
  const { data, isFetching, refetch, isError } = useQuery({
    queryKey: ['admin-health-check'],
    queryFn: fetchHealthCheck,
    enabled: false,
  });

  useEffect(() => {
    if (data) {
      saveHealthCheckComplete();
    }
  }, [data]);

  return (
    <Card withBorder padding="md">
      <Group justify="space-between" mb="md">
        <div>
          <Title order={4}>System health check</Title>
          <Text size="sm" c="dimmed">
            Verify database, storage paths, migrations, and document file integrity.
          </Text>
        </div>
        <Button loading={isFetching} onClick={() => refetch()}>
          Run check
        </Button>
      </Group>

      {isError && (
        <Text c="red" size="sm">Health check failed. Ensure you are signed in as an administrator.</Text>
      )}

      {data && (
        <Stack gap="sm">
          <Badge size="lg" color={statusColor(data.status)} variant="light">
            Overall: {data.status}
          </Badge>
          {data.checks.map((check) => (
            <Group key={check.id} justify="space-between" align="flex-start" wrap="nowrap">
              <div>
                <Text fw={600} size="sm">{check.label}</Text>
                <Text size="sm" c="dimmed">{check.message}</Text>
              </div>
              <Badge color={statusColor(check.status)}>{check.status}</Badge>
            </Group>
          ))}
          <Text size="xs" c="dimmed">
            Checked {new Date(data.checkedAt).toLocaleString()}
          </Text>
        </Stack>
      )}
    </Card>
  );
}
