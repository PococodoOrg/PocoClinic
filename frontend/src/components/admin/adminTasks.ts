import { SystemStatus } from '../../api/admin';
import {
  auditArchiveDue,
  complianceCheckDue,
  healthCheckDue,
  monthlyTestingDue,
} from './adminReminders';
import { loadLastRestoreDrill, restoreDrillDue } from './maintenanceDrill';

export type AdminTaskStatus = 'ok' | 'due' | 'urgent';

export interface AdminTask {
  id: string;
  label: string;
  status: AdminTaskStatus;
  hint: string;
  path?: string;
}

const STATUS_ORDER: Record<AdminTaskStatus, number> = {
  urgent: 0,
  due: 1,
  ok: 2,
};

export function buildAdminTasks(status: SystemStatus): AdminTask[] {
  const tasks: AdminTask[] = [];

  const backupStatus =
    status.backup.status === 'critical'
      ? 'urgent'
      : status.backup.status === 'warning'
        ? 'due'
        : 'ok';
  tasks.push({
    id: 'backup',
    label: 'BACKUP',
    status: backupStatus,
    hint:
      backupStatus === 'urgent'
        ? 'No recent backup — run USB backup today'
        : backupStatus === 'due'
          ? 'Last backup is aging — schedule one before end of day'
          : 'Backup is up to date',
    path: '/admin?tab=backup',
  });

  tasks.push({
    id: 'audit-archive',
    label: 'ARCHIVE AUDIT',
    status: auditArchiveDue() ? 'due' : 'ok',
    hint: auditArchiveDue()
      ? 'Download audit CSV and store offline before monthly purge'
      : 'Audit log archived this month',
    path: '/admin/audit',
  });

  tasks.push({
    id: 'restore-drill',
    label: 'RESTORE DRILL',
    status: restoreDrillDue(loadLastRestoreDrill()) ? 'due' : 'ok',
    hint: restoreDrillDue(loadLastRestoreDrill())
      ? 'Test that a USB backup still restores correctly'
      : 'Restore drill recorded recently',
    path: '/help/quarterly-drill',
  });

  tasks.push({
    id: 'health-check',
    label: 'HEALTH CHECK',
    status: healthCheckDue() ? 'due' : 'ok',
    hint: healthCheckDue()
      ? 'Run system health check on the admin dashboard'
      : 'Health check run this month',
    path: '/admin?tab=backup',
  });

  tasks.push({
    id: 'compliance',
    label: 'COMPLIANCE',
    status: complianceCheckDue() ? 'due' : 'ok',
    hint: complianceCheckDue()
      ? 'Run the HIPAA compliance checklist'
      : 'Compliance checklist done this month',
    path: '/admin?tab=overview',
  });

  tasks.push({
    id: 'monthly-test',
    label: 'MONTHLY TEST',
    status: monthlyTestingDue() ? 'due' : 'ok',
    hint: monthlyTestingDue()
      ? 'Walk through the monthly testing checklist'
      : 'Monthly testing marked complete',
    path: '/help/compliance-checklist',
  });

  if (status.counts.lockedAccounts > 0) {
    tasks.push({
      id: 'locked-accounts',
      label: 'UNLOCK STAFF',
      status: 'due',
      hint: `${status.counts.lockedAccounts} staff locked after failed sign-ins — review and unlock from staff records`,
      path: '/users',
    });
  }

  if (status.counts.defaultPinAccounts > 0) {
    tasks.push({
      id: 'default-pins',
      label: 'DEFAULT PINS',
      status: 'due',
      hint: `${status.counts.defaultPinAccounts} staff still on factory PIN — review staff records`,
      path: '/users',
    });
  }

  return tasks.sort(
    (a, b) => STATUS_ORDER[a.status] - STATUS_ORDER[b.status] || a.label.localeCompare(b.label),
  );
}

export function adminTasksNeedAttention(tasks: AdminTask[]): boolean {
  return tasks.some((task) => task.status !== 'ok');
}
