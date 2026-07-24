import {
  Badge,
  Button,
  Card,
  Grid,
  Group,
  Modal,
  SimpleGrid,
  Stack,
  Table,
  Tabs,
  Text,
  Title,
} from '@mantine/core';

import {
  IconBook,
  IconDatabase,
  IconHeartRateMonitor,
  IconReportAnalytics,
  IconSettings,
} from '@tabler/icons-react';

import { useState } from 'react';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { notifications } from '@mantine/notifications';

import { useNavigate, useSearchParams } from 'react-router-dom';

import {
  createBackup,
  fetchBackups,
  fetchPatientCensus,
  fetchSystemStatus,
  restoreBackup,
  verifyBackup,
} from '../api/admin';

import { ActivitySummaryPanel } from '../components/admin/ActivitySummaryPanel';
import { AdminHelpPanel } from '../components/admin/AdminHelpPanel';

import { AdminNotifications } from '../components/admin/AdminNotifications';

import { AdminTaskBoard } from '../components/admin/AdminTaskBoard';

import { BackupSticker } from '../components/admin/BackupSticker';

import { CompliancePanel } from '../components/admin/CompliancePanel';

import { HealthCheckPanel } from '../components/admin/HealthCheckPanel';

import { PerformancePanel } from '../components/admin/PerformancePanel';
import { PatientFieldSettingsPanel } from '../components/admin/PatientFieldSettingsPanel';

import { MaintenanceReminders } from '../components/admin/MaintenanceReminders';

import { MaintenanceSchedulePanel } from '../components/admin/MaintenanceSchedulePanel';

import { SecurityOverview } from '../components/admin/SecurityOverview';

import { StaffActivityPanel } from '../components/admin/StaffActivityPanel';

import { UsbRotationPanel } from '../components/admin/UsbRotationPanel';

import { getErrorMessage } from '../utils/apiError';

function formatUptime(seconds: number): string {
  const hours = Math.floor(seconds / 3600);

  const minutes = Math.floor((seconds % 3600) / 60);

  if (hours > 0) {
    return `${hours}h ${minutes}m`;
  }

  return `${minutes}m`;
}

function backupColor(status: string): string {
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

function backupLabel(status: string): string {
  switch (status) {
    case 'ok':
      return 'Up to date';

    case 'warning':
      return 'Backup aging';

    case 'critical':
      return 'Backup overdue';

    default:
      return 'No backup found';
  }
}

export default function AdminDashboard() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const activeTab = searchParams.get('tab') ?? 'overview';

  const queryClient = useQueryClient();

  const [restoreTarget, setRestoreTarget] = useState<string | null>(null);

  const {
    data: status,
    isLoading: statusLoading,
    error: statusError,
  } = useQuery({
    queryKey: ['admin-system-status'],

    queryFn: fetchSystemStatus,

    refetchInterval: 60_000,
  });

  const {
    data: census,
    isLoading: censusLoading,
    error: censusError,
  } = useQuery({
    queryKey: ['admin-patient-census'],

    queryFn: fetchPatientCensus,
  });

  const { data: backups = [] } = useQuery({
    queryKey: ['admin-backups'],

    queryFn: fetchBackups,

    enabled: status?.storageMode === 'database',
  });

  const backupMutation = useMutation({
    mutationFn: createBackup,

    onSuccess: async (result) => {
      await queryClient.invalidateQueries({
        queryKey: ['admin-system-status'],
      });

      await queryClient.invalidateQueries({ queryKey: ['admin-backups'] });

      notifications.show({
        title: 'Backup created',

        message: result.filename,

        color: 'green',
      });
    },

    onError: (error: Error) => {
      notifications.show({
        title: 'Backup failed',

        message: getErrorMessage(error, 'Backup failed'),

        color: 'red',
      });
    },
  });

  const verifyMutation = useMutation({
    mutationFn: verifyBackup,

    onSuccess: (result) => {
      const detail = result.summary
        ? ` ${result.summary.tableRows?.patients ?? 0} patients, ${result.summary.documentFiles ?? 0} document files.`
        : '';

      notifications.show({
        title: result.valid ? 'Backup verified' : 'Backup verification failed',

        message: (result.message ?? result.filename) + detail,

        color: result.valid ? 'green' : 'red',
      });
    },

    onError: (error: Error) => {
      notifications.show({
        title: 'Verification failed',

        message: getErrorMessage(error, 'Backup verification failed'),

        color: 'red',
      });
    },
  });

  const restoreMutation = useMutation({
    mutationFn: (filename: string) => restoreBackup({ filename, confirm: true }),

    onSuccess: async (_result, filename) => {
      setRestoreTarget(null);

      await queryClient.invalidateQueries({
        queryKey: ['admin-system-status'],
      });

      await queryClient.invalidateQueries({ queryKey: ['admin-backups'] });

      await queryClient.invalidateQueries({ queryKey: ['admin-audit-logs'] });

      notifications.show({
        title: 'Restore completed',

        message: filename,

        color: 'green',
      });
    },

    onError: (error: Error) => {
      notifications.show({
        title: 'Restore failed',

        message: getErrorMessage(error, 'Restore failed'),

        color: 'red',
      });
    },
  });

  if (statusError || censusError) {
    return (
      <Text c="red">
        Failed to load admin dashboard. Ensure you are signed in as an administrator.
      </Text>
    );
  }

  return (
    <Stack gap="lg" className="admin-dashboard">
      <Group justify="space-between" align="flex-start" wrap="wrap">
        <div>
          <Title order={2}>Admin dashboard</Title>

          <Text c="dimmed" size="sm">
            System health, administrator guide, backups, and patient census.
          </Text>
        </div>

        <Group>
          <Button variant="default" onClick={() => navigate('/admin/audit')}>
            Audit log
          </Button>

          <Button
            variant="light"

            component="a"

            href="http://127.0.0.1:9090"

            target="_blank"

            rel="noreferrer"
          >
            Backup helper
          </Button>

          <Button
            variant="light"
            leftSection={<IconBook size={16} />}
            onClick={() => navigate('/help?category=administration')}
          >
            Help articles
          </Button>
        </Group>
      </Group>

      <Tabs
        value={activeTab}
        onChange={(value) => setSearchParams(value && value !== 'overview' ? { tab: value } : {})}
        keepMounted={false}
      >
        <Tabs.List>
          <Tabs.Tab value="overview" leftSection={<IconHeartRateMonitor size={16} />}>
            Overview
          </Tabs.Tab>

          <Tabs.Tab value="guide" leftSection={<IconBook size={16} />}>
            Administrator guide
          </Tabs.Tab>

          <Tabs.Tab value="backup" leftSection={<IconDatabase size={16} />}>
            Backup & maintenance
          </Tabs.Tab>

          <Tabs.Tab value="reports" leftSection={<IconReportAnalytics size={16} />}>
            Reports
          </Tabs.Tab>

          <Tabs.Tab value="settings" leftSection={<IconSettings size={16} />}>
            Clinic settings
          </Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="overview" pt="md">
          <Stack gap="lg">
            {status && !statusLoading && (
              <>
                <AdminTaskBoard status={status} />

                <Card withBorder padding="md">
                  <AdminNotifications status={status} />
                </Card>
              </>
            )}

            {statusLoading || !status ? (
              <Text c="dimmed">Loading system status…</Text>
            ) : (
              <>
                <SimpleGrid cols={{ base: 1, sm: 2, lg: 4 }}>
                  <Card withBorder padding="md">
                    <Text size="sm" c="dimmed">
                      Storage
                    </Text>

                    <Group gap="xs" mt="xs">
                      <Badge color={status.storageMode === 'database' ? 'green' : 'yellow'}>
                        {status.storageMode === 'database' ? 'Database' : 'In-memory'}
                      </Badge>

                      {status.storageMode === 'database' && (
                        <Badge color={status.databaseConnected ? 'green' : 'red'} variant="light">
                          {status.databaseConnected ? 'Connected' : 'Offline'}
                        </Badge>
                      )}
                    </Group>
                  </Card>

                  <Card withBorder padding="md">
                    <Text size="sm" c="dimmed">
                      Backup
                    </Text>

                    <Badge color={backupColor(status.backup.status)} mt="xs">
                      {backupLabel(status.backup.status)}
                    </Badge>

                    {status.backup.lastBackupAt && (
                      <Text size="xs" c="dimmed" mt="xs">
                        Last: {new Date(status.backup.lastBackupAt).toLocaleString()}
                      </Text>
                    )}
                  </Card>

                  <Card withBorder padding="md">
                    <Text size="sm" c="dimmed">
                      Uptime
                    </Text>

                    <Title order={3} mt="xs">
                      {formatUptime(status.uptimeSeconds)}
                    </Title>

                    <Text size="xs" c="dimmed">
                      {status.environment} · v{status.appVersion}
                    </Text>
                  </Card>

                  <Card withBorder padding="md">
                    <Text size="sm" c="dimmed">
                      Migrations
                    </Text>

                    {status.pendingMigrations.length === 0 ? (
                      <Badge color="green" mt="xs">
                        Up to date
                      </Badge>
                    ) : (
                      <Badge color="red" mt="xs">
                        {status.pendingMigrations.length} pending
                      </Badge>
                    )}

                    {status.latestMigration && (
                      <Text size="xs" c="dimmed" mt="xs">
                        Latest: {status.latestMigration}
                      </Text>
                    )}
                  </Card>
                </SimpleGrid>

                <SimpleGrid cols={{ base: 2, sm: 3, lg: 6 }}>
                  <StatCard label="Patients" value={status.counts.patients} />

                  <StatCard label="Staff" value={status.counts.users} />

                  <StatCard label="Forms" value={status.counts.formTemplates} />

                  <StatCard label="Submissions" value={status.counts.formSubmissions} />

                  <StatCard label="Sessions" value={status.counts.activeSessions} />

                  <StatCard label="Documents" value={status.counts.patientDocuments} />
                </SimpleGrid>

                <SecurityOverview status={status} />

                <CompliancePanel />

                <PerformancePanel />
              </>
            )}
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="guide" pt="md">
          <AdminHelpPanel />
        </Tabs.Panel>

        <Tabs.Panel value="backup" pt="md">
          <Stack gap="lg">
            <Grid gap="lg">
              <Grid.Col span={{ base: 12, lg: 6 }}>
                <HealthCheckPanel />
              </Grid.Col>

              <Grid.Col span={{ base: 12, lg: 6 }}>
                {status && !statusLoading && <MaintenanceReminders status={status} />}
              </Grid.Col>
            </Grid>

            <MaintenanceSchedulePanel />

            <UsbRotationPanel />

            {status && !statusLoading && <BackupSticker backup={status.backup} />}

            <Card withBorder padding="md">
              <Group justify="space-between" mb="md">
                <Title order={4}>Backups</Title>

                <Button
                  loading={backupMutation.isPending}

                  disabled={status?.storageMode !== 'database'}

                  onClick={() => backupMutation.mutate()}
                >
                  Backup now
                </Button>
              </Group>

              {status?.storageMode !== 'database' ? (
                <Text size="sm" c="dimmed">
                  Backups require a persistent database. Set DATABASE_URL and restart the server.
                </Text>
              ) : backups.length === 0 ? (
                <Text size="sm" c="dimmed">
                  No backups yet. Use Backup now or run{' '}
                  <Text span ff="monospace">
                    go run ./cmd/backup
                  </Text>
                  .
                </Text>
              ) : (
                <Table withTableBorder>
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th>File</Table.Th>

                      <Table.Th>Created</Table.Th>

                      <Table.Th>Version</Table.Th>

                      <Table.Th />
                    </Table.Tr>
                  </Table.Thead>

                  <Table.Tbody>
                    {[...backups]
                      .reverse()
                      .slice(0, 5)
                      .map((backup) => (
                        <Table.Tr key={backup.filename}>
                          <Table.Td>{backup.filename}</Table.Td>

                          <Table.Td>{new Date(backup.createdAt).toLocaleString()}</Table.Td>

                          <Table.Td>{backup.appVersion || '—'}</Table.Td>

                          <Table.Td>
                            <Group gap="xs" justify="flex-end">
                              <Button
                                size="xs"

                                variant="light"

                                loading={verifyMutation.isPending}

                                onClick={() => verifyMutation.mutate(backup.filename)}
                              >
                                Verify
                              </Button>

                              <Button
                                size="xs"

                                color="red"

                                variant="light"

                                onClick={() => setRestoreTarget(backup.filename)}
                              >
                                Restore
                              </Button>
                            </Group>
                          </Table.Td>
                        </Table.Tr>
                      ))}
                  </Table.Tbody>
                </Table>
              )}

              <Text size="sm" c="dimmed" mt="md">
                Restore replaces all database contents (including encrypted patient documents and
                exercise logs). Sign out all staff first.
              </Text>
            </Card>
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="reports" pt="md">
          <Stack gap="lg">
            <ActivitySummaryPanel />

            <StaffActivityPanel />

            <Card withBorder padding="md">
              <Group justify="space-between" mb="md">
                <Title order={4}>Patient census</Title>

                <Button variant="light" onClick={() => navigate('/forms/reports')}>
                  Form reports
                </Button>
              </Group>

              {censusLoading || !census ? (
                <Text c="dimmed">Loading census…</Text>
              ) : (
                <Grid>
                  <Grid.Col span={{ base: 12, md: 4 }}>
                    <Stack gap="xs">
                      <Text size="sm" c="dimmed">
                        Total patients
                      </Text>

                      <Title order={2}>{census.totalPatients}</Title>

                      <Text size="sm" c="dimmed">
                        {census.addedLast30Days} added in the last 30 days
                      </Text>
                    </Stack>
                  </Grid.Col>

                  <Grid.Col span={{ base: 12, md: 8 }}>
                    <Table withTableBorder>
                      <Table.Thead>
                        <Table.Tr>
                          <Table.Th>Gender</Table.Th>

                          <Table.Th>Count</Table.Th>
                        </Table.Tr>
                      </Table.Thead>

                      <Table.Tbody>
                        {census.byGender.length === 0 ? (
                          <Table.Tr>
                            <Table.Td colSpan={2}>
                              <Text c="dimmed">No patient records yet.</Text>
                            </Table.Td>
                          </Table.Tr>
                        ) : (
                          census.byGender.map((row) => (
                            <Table.Tr key={row.gender}>
                              <Table.Td>{row.gender}</Table.Td>

                              <Table.Td>{row.count}</Table.Td>
                            </Table.Tr>
                          ))
                        )}
                      </Table.Tbody>
                    </Table>
                  </Grid.Col>
                </Grid>
              )}
            </Card>
          </Stack>
        </Tabs.Panel>

        <Tabs.Panel value="settings" pt="md">
          <PatientFieldSettingsPanel />
        </Tabs.Panel>
      </Tabs>

      <Modal
        opened={restoreTarget !== null}

        onClose={() => setRestoreTarget(null)}

        title="Confirm database restore"

        centered
      >
        <Stack gap="md">
          <Text size="sm">
            This will replace <strong>all</strong> patients, forms, staff accounts, audit history,
            and patient document files with the contents of:
          </Text>

          <Text ff="monospace" size="sm">
            {restoreTarget}
          </Text>

          <Text size="sm" c="red">
            Active staff sessions may become invalid. Ensure everyone is signed out before
            continuing.
          </Text>

          <Group justify="flex-end">
            <Button variant="default" onClick={() => setRestoreTarget(null)}>
              Cancel
            </Button>

            <Button
              color="red"

              loading={restoreMutation.isPending}

              onClick={() => restoreTarget && restoreMutation.mutate(restoreTarget)}
            >
              Restore database
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Stack>
  );
}

function StatCard({ label, value }: { label: string; value: number }) {
  return (
    <Card withBorder padding="md">
      <Text size="sm" c="dimmed">
        {label}
      </Text>

      <Title order={3} mt="xs">
        {value}
      </Title>
    </Card>
  );
}
