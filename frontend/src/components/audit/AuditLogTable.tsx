import { Badge, Group, Pagination, Stack, Table, Text } from '@mantine/core';
import { AuditEntry, formatAuditEvent } from '../../types/audit';

interface AuditLogTableProps {
  events: AuditEntry[];
  totalCount: number;
  page: number;
  pageSize: number;
  totalPages: number;
  onPageChange?: (page: number) => void;
  showUserColumn?: boolean;
  emptyMessage?: string;
}

function formatDetails(entry: AuditEntry): string {
  const details = entry.details ?? {};
  const parts: string[] = [];

  if (details.filename) {
    parts.push(details.filename);
  }
  if (details.email) {
    parts.push(details.email);
  }
  if (details.role) {
    parts.push(`role: ${details.role}`);
  }
  if (details.reason) {
    parts.push(details.reason.replace(/_/g, ' '));
  }
  if (details.error) {
    parts.push(details.error);
  }
  if (entry.resourceType === 'patient' && entry.resourceId) {
    parts.push(`patient ${entry.resourceId.slice(0, 8)}…`);
  }

  return parts.join(' · ');
}

export function AuditLogTable({
  events,
  totalCount,
  page,
  pageSize,
  totalPages,
  onPageChange,
  showUserColumn = false,
  emptyMessage = 'No activity recorded yet.',
}: AuditLogTableProps) {
  return (
    <Stack gap="md">
      {events.length === 0 ? (
        <Text c="dimmed" ta="center" py="md">
          {emptyMessage}
        </Text>
      ) : (
        <Table striped highlightOnHover withTableBorder>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>When</Table.Th>
              <Table.Th>Event</Table.Th>
              {showUserColumn && <Table.Th>User</Table.Th>}
              <Table.Th>Details</Table.Th>
              <Table.Th>IP address</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {events.map((entry) => (
              <Table.Tr key={entry.id}>
                <Table.Td>{new Date(entry.createdAt).toLocaleString()}</Table.Td>
                <Table.Td>
                  <Group gap="xs">
                    <Text size="sm">{formatAuditEvent(entry)}</Text>
                    {!entry.success && (
                      <Badge color="red" size="xs" variant="light">
                        Failed
                      </Badge>
                    )}
                  </Group>
                </Table.Td>
                {showUserColumn && (
                  <Table.Td>
                    <Text size="sm" c="dimmed">
                      {entry.userId ? `${entry.userId.slice(0, 8)}…` : '—'}
                    </Text>
                  </Table.Td>
                )}
                <Table.Td>
                  <Text size="sm" c="dimmed">
                    {formatDetails(entry) || '—'}
                  </Text>
                </Table.Td>
                <Table.Td>
                  <Text size="sm" c="dimmed">
                    {entry.ipAddress || '—'}
                  </Text>
                </Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      )}

      {totalCount > pageSize && onPageChange && (
        <Group justify="space-between">
          <Text size="sm" c="dimmed">
            Showing page {page} of {totalPages} ({totalCount} events)
          </Text>
          <Pagination value={page} onChange={onPageChange} total={totalPages} />
        </Group>
      )}
    </Stack>
  );
}
