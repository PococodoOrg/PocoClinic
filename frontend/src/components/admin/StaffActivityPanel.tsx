import {
  Badge,
  Card,
  Group,
  SegmentedControl,
  Table,
  Text,
  Title,
} from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { fetchStaffActivity } from '../../api/admin';

const PERIOD_OPTIONS = [
  { label: '7 days', value: '7' },
  { label: '30 days', value: '30' },
  { label: '90 days', value: '90' },
];

export function StaffActivityPanel() {
  const [days, setDays] = useState('30');

  const { data, isLoading, error } = useQuery({
    queryKey: ['admin-staff-activity', days],
    queryFn: () => fetchStaffActivity(Number(days)),
  });

  return (
    <Card withBorder padding="md">
      <Group justify="space-between" mb="md" align="flex-start" wrap="wrap">
        <div>
          <Title order={4}>Staff activity</Title>
          <Text size="sm" c="dimmed">
            Sign-ins, failed attempts, and chart views per staff member from the audit log.
          </Text>
        </div>
        <SegmentedControl
          value={days}
          onChange={setDays}
          data={PERIOD_OPTIONS}
          aria-label="Staff activity period"
        />
      </Group>

      {error && <Text c="red" size="sm">Could not load staff activity.</Text>}

      {isLoading || !data ? (
        <Text c="dimmed" size="sm">Loading staff activity…</Text>
      ) : data.staff.length === 0 ? (
        <Text c="dimmed" size="sm">No staff accounts yet.</Text>
      ) : (
        <Table withTableBorder striped highlightOnHover>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Staff</Table.Th>
              <Table.Th>Sign-ins</Table.Th>
              <Table.Th>Failed</Table.Th>
              <Table.Th>Chart views</Table.Th>
              <Table.Th>Last sign-in</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {data.staff.map((row) => (
              <Table.Tr key={row.userId}>
                <Table.Td>
                  <Text fw={600} size="sm">{row.name}</Text>
                  <Text size="xs" c="dimmed">{row.email}</Text>
                </Table.Td>
                <Table.Td>
                  <Badge variant="light" color="blue">{row.signIns}</Badge>
                </Table.Td>
                <Table.Td>
                  <Badge variant="light" color={row.failedSignIns > 0 ? 'orange' : 'gray'}>
                    {row.failedSignIns}
                  </Badge>
                </Table.Td>
                <Table.Td>{row.patientViews}</Table.Td>
                <Table.Td>
                  {row.lastSignIn ? new Date(row.lastSignIn).toLocaleString() : '—'}
                </Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      )}
    </Card>
  );
}
