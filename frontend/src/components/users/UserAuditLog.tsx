import { Stack, Text } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { fetchUserAuditLogs } from '../../api/users';
import { AuditLogTable } from '../audit/AuditLogTable';

interface UserAuditLogProps {
  userId: string;
}

export function UserAuditLog({ userId }: UserAuditLogProps) {
  const { data, isLoading, error } = useQuery({
    queryKey: ['user-audit', userId],
    queryFn: () => fetchUserAuditLogs(userId, { page: 1, pageSize: 50 }),
  });

  if (isLoading) {
    return (
      <Text c="dimmed" ta="center" py="md">
        Loading activity history...
      </Text>
    );
  }

  if (error) {
    return <Text c="red">Could not load activity history.</Text>;
  }

  return (
    <Stack gap="md">
      <Text fw={500}>Activity history</Text>
      <Text size="sm" c="dimmed">
        Sign-ins, PIN changes, badge events, and patient actions performed by this staff member.
      </Text>

      <AuditLogTable
        events={data?.events ?? []}
        totalCount={data?.totalCount ?? 0}
        page={1}
        pageSize={50}
        totalPages={data?.totalPages ?? 1}
        emptyMessage="No activity recorded yet."
      />

      {data && data.totalCount > data.events.length && (
        <Text size="sm" c="dimmed">
          Showing latest {data.events.length} of {data.totalCount} events.
        </Text>
      )}
    </Stack>
  );
}
