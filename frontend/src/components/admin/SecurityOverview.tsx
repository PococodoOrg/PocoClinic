import { Badge, Card, Group, SimpleGrid, Text, Title } from '@mantine/core';
import { SystemStatus } from '../../api/admin';

interface SecurityOverviewProps {
  status: SystemStatus;
}

export function SecurityOverview({ status }: SecurityOverviewProps) {
  const { counts, security } = status;
  const documentsReady = security?.documentsStorageReady ?? false;

  return (
    <Card withBorder padding="md">
      <Title order={4} mb="md">Security overview</Title>
      <SimpleGrid cols={{ base: 2, sm: 4 }}>
        <Metric
          label="Active sessions"
          value={counts.activeSessions}
          color={counts.activeSessions > 0 ? 'blue' : 'gray'}
        />
        <Metric
          label="Locked accounts"
          value={counts.lockedAccounts}
          color={counts.lockedAccounts > 0 ? 'orange' : 'green'}
        />
        <Metric
          label="Default PIN"
          value={counts.defaultPinAccounts}
          color={counts.defaultPinAccounts > 0 ? 'yellow' : 'green'}
        />
        <Metric
          label="Patient documents"
          value={counts.patientDocuments}
          color="gray"
        />
      </SimpleGrid>
      <Group mt="md" gap="xs">
        <Badge color={documentsReady ? 'green' : 'red'} variant="light">
          Document storage {documentsReady ? 'ready' : 'unavailable'}
        </Badge>
        <Badge color="green" variant="light">
          LAN-only · badge + PIN
        </Badge>
      </Group>
    </Card>
  );
}

function Metric({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <div>
      <Text size="sm" c="dimmed">{label}</Text>
      <Badge size="lg" color={color} mt="xs">{value}</Badge>
    </div>
  );
}
