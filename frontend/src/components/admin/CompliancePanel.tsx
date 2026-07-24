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
import { fetchComplianceCheck } from '../../api/admin';
import { saveComplianceCheckComplete } from './adminReminders';

function overallStatusColor(status: string): string {
  switch (status) {
    case 'ok':
      return 'green';
    case 'warn':
      return 'yellow';
    case 'fail':
      return 'red';
    default:
      return 'gray';
  }
}

function overallStatusLabel(status: string): string {
  switch (status) {
    case 'ok':
      return 'Pass';
    case 'warn':
      return 'Review needed';
    case 'fail':
      return 'Action required';
    default:
      return status;
  }
}

function statusColor(status: string): string {
  switch (status) {
    case 'pass':
      return 'green';
    case 'warn':
      return 'yellow';
    case 'fail':
      return 'red';
    default:
      return 'gray';
  }
}

function statusLabel(status: string): string {
  switch (status) {
    case 'pass':
      return 'Pass';
    case 'warn':
      return 'Review';
    case 'fail':
      return 'Fail';
    default:
      return status;
  }
}

export function CompliancePanel() {
  const { data, isFetching, refetch, isError } = useQuery({
    queryKey: ['admin-compliance-check'],
    queryFn: fetchComplianceCheck,
    enabled: false,
  });

  useEffect(() => {
    if (data) {
      saveComplianceCheckComplete();
    }
  }, [data]);

  return (
    <Card withBorder padding="md">
      <Group justify="space-between" mb="md">
        <div>
          <Title order={4}>HIPAA compliance checklist</Title>
          <Text size="sm" c="dimmed">
            Automated operational controls for backups, authentication, audit logging, and storage.
          </Text>
        </div>
        <Button loading={isFetching} onClick={() => refetch()}>
          Run checklist
        </Button>
      </Group>

      {isError && (
        <Text c="red" size="sm">Compliance check failed. Ensure you are signed in as an administrator.</Text>
      )}

      {data && (
        <Stack gap="sm">
          <Badge size="lg" color={overallStatusColor(data.status)} variant="light">
            Overall: {overallStatusLabel(data.status)}
          </Badge>
          {data.checks.map((check) => (
            <Group key={check.id} justify="space-between" align="flex-start" wrap="nowrap">
              <div>
                <Text fw={600} size="sm">{check.label}</Text>
                <Text size="sm" c="dimmed">{check.message}</Text>
                {check.control && (
                  <Text size="xs" c="dimmed" mt={4}>{check.control}</Text>
                )}
              </div>
              <Badge color={statusColor(check.status)}>{statusLabel(check.status)}</Badge>
            </Group>
          ))}
          <Text size="xs" c="dimmed">
            Checked {new Date(data.checkedAt).toLocaleString()} — review any failed items before end of month.
          </Text>
        </Stack>
      )}
    </Card>
  );
}
