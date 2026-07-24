import { Alert, Button, Group, Select, Stack, Text, Title } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { buildAuditExportUrl, downloadAuthenticatedCsv, fetchAuditLogs } from '../api/admin';
import { fetchUsers } from '../api/users';
import {
  auditArchiveDue,
  loadLastAuditArchive,
  saveAuditArchiveComplete,
} from '../components/admin/adminReminders';
import { AuditLogTable } from '../components/audit/AuditLogTable';
import { AUDIT_EVENT_LABELS } from '../types/audit';

const PAGE_SIZE = 25;

const EVENT_FILTER_OPTIONS = [
  { value: '', label: 'All events' },
  ...Object.entries(AUDIT_EVENT_LABELS).map(([value, label]) => ({ value, label })),
];

export default function AdminAuditLog() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [eventType, setEventType] = useState('');
  const [userId, setUserId] = useState('');
  const [exporting, setExporting] = useState(false);
  const [archiveDue, setArchiveDue] = useState(auditArchiveDue());
  const [lastArchived, setLastArchived] = useState<string | null>(loadLastAuditArchive());

  useEffect(() => {
    setArchiveDue(auditArchiveDue());
    setLastArchived(loadLastAuditArchive());
  }, []);

  const handleExport = async () => {
    setExporting(true);
    try {
      const url = buildAuditExportUrl({
        days: 30,
        eventType: eventType || undefined,
        userId: userId || undefined,
      });
      await downloadAuthenticatedCsv(url, `audit-log-${new Date().toISOString().slice(0, 10)}.csv`);
      notifications.show({
        color: 'teal',
        title: 'Download started',
        message: 'Save the CSV to offline storage or your ops binder, then mark archived.',
      });
    } catch {
      notifications.show({ color: 'red', message: 'Failed to download audit log.' });
    } finally {
      setExporting(false);
    }
  };

  const markArchived = () => {
    saveAuditArchiveComplete();
    setArchiveDue(auditArchiveDue());
    setLastArchived(loadLastAuditArchive());
    notifications.show({ color: 'green', message: 'Audit archive marked complete for this month.' });
  };

  const { data: staffPage } = useQuery({
    queryKey: ['users', 'audit-filter'],
    queryFn: () => fetchUsers({ page: 1, pageSize: 200 }),
  });

  const staffOptions = [
    { value: '', label: 'All staff' },
    ...(staffPage?.users ?? []).map((user) => ({ value: user.id, label: user.name })),
  ];

  const { data, isLoading, error } = useQuery({
    queryKey: ['admin-audit-logs', page, eventType, userId],
    queryFn: () =>
      fetchAuditLogs({
        page,
        pageSize: PAGE_SIZE,
        eventType: eventType || undefined,
        userId: userId || undefined,
      }),
  });

  return (
    <Stack gap="lg">
      <Group justify="space-between">
        <div>
          <Title order={2}>Audit log</Title>
          <Text c="dimmed" size="sm">
            System-wide sign-ins, patient access, staff changes, and backup operations. Admin download only.
          </Text>
        </div>
        <Group>
          <Button variant="light" loading={exporting} onClick={handleExport}>
            Download CSV (30 days)
          </Button>
          <Button variant="default" onClick={() => navigate('/admin')}>
            Back to dashboard
          </Button>
        </Group>
      </Group>

      {archiveDue ? (
        <Alert color="yellow" title="ARCHIVE AUDIT — due this month">
          <Stack gap="xs">
            <Text size="sm">
              Download the CSV and store it offline (USB, ops binder, or secure clinic archive).
              Monthly audit purge removes rows older than retention — archive first so you keep a copy.
            </Text>
            <Button size="xs" variant="light" onClick={markArchived}>
              Mark archived after saving offline
            </Button>
          </Stack>
        </Alert>
      ) : (
        <Alert color="green" title="ARCHIVE AUDIT — done this month">
          <Text size="sm">
            Audit log marked archived for {lastArchived}. Download again anytime if you need a fresh copy.
          </Text>
        </Alert>
      )}

      <Group align="flex-end" wrap="wrap">
        <Select
          label="Event type"
          data={EVENT_FILTER_OPTIONS}
          value={eventType}
          onChange={(value) => {
            setEventType(value ?? '');
            setPage(1);
          }}
          w={260}
          searchable
        />
        <Select
          label="Staff member"
          data={staffOptions}
          value={userId}
          onChange={(value) => {
            setUserId(value ?? '');
            setPage(1);
          }}
          w={260}
          searchable
        />
      </Group>

      {error && <Text c="red">Failed to load audit log.</Text>}
      {isLoading ? (
        <Text c="dimmed">Loading audit events…</Text>
      ) : (
        <AuditLogTable
          events={data?.events ?? []}
          totalCount={data?.totalCount ?? 0}
          page={page}
          pageSize={PAGE_SIZE}
          totalPages={data?.totalPages ?? 1}
          onPageChange={setPage}
          showUserColumn
          emptyMessage="No audit events match this filter."
        />
      )}

      <Text size="sm" c="dimmed">
        Audit entries support HIPAA accountability. Only administrators can download CSV copies — store them offline
        and mark archived on the dashboard so you remember before purge runs.
      </Text>
    </Stack>
  );
}
