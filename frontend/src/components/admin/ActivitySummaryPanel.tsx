import { Badge, Card, Group, SegmentedControl, SimpleGrid, Stack, Text, Title } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { fetchActivitySummary } from '../../api/admin';

const PERIOD_OPTIONS = [
  { label: '7 days', value: '7' },
  { label: '30 days', value: '30' },
  { label: '90 days', value: '90' },
];

export function ActivitySummaryPanel() {
  const [days, setDays] = useState('7');

  const { data, isLoading, error } = useQuery({
    queryKey: ['admin-activity-summary', days],
    queryFn: () => fetchActivitySummary(Number(days)),
  });

  return (
    <Card withBorder padding="md">
      <Group justify="space-between" mb="md" align="flex-start" wrap="wrap">
        <div>
          <Title order={4}>Clinic activity</Title>
          <Text size="sm" c="dimmed">
            Sign-ins, patient work, and backups from the audit log and database.
          </Text>
        </div>
        <SegmentedControl
          value={days}
          onChange={setDays}
          data={PERIOD_OPTIONS}
          aria-label="Activity report period"
        />
      </Group>

      {error && (
        <Text c="red" size="sm">Could not load activity summary.</Text>
      )}

      {isLoading || !data ? (
        <Text c="dimmed" size="sm">Loading activity…</Text>
      ) : (
        <SimpleGrid cols={{ base: 2, sm: 3, lg: 4 }}>
          <Metric label="Staff sign-ins" value={data.signIns} color="blue" />
          <Metric label="Failed sign-ins" value={data.failedSignIns} color={data.failedSignIns > 0 ? 'orange' : 'gray'} />
          <Metric label="New patients" value={data.patientsCreated} color="green" />
          <Metric label="Chart notes" value={data.chartNotesAdded} color="gray" />
          <Metric label="Chart views" value={data.patientViews} color="gray" />
          <Metric label="Documents uploaded" value={data.documentsUploaded} color="gray" />
          <Metric label="Backups created" value={data.backupsCreated} color="teal" />
        </SimpleGrid>
      )}
    </Card>
  );
}

function Metric({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <Stack gap={4}>
      <Text size="xs" c="dimmed">{label}</Text>
      <Badge size="xl" variant="light" color={color}>{value}</Badge>
    </Stack>
  );
}
