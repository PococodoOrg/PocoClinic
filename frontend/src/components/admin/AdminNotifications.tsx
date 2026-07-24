import { Alert, Stack, Text, Title } from '@mantine/core';

import { useNavigate } from 'react-router-dom';

import { SystemStatus } from '../../api/admin';

import {
  auditArchiveDue,
  complianceCheckDue,
  healthCheckDue,
  monthlyTestingDue,
} from './adminReminders';

import {

  daysSinceLastDrill,

  loadLastRestoreDrill,

  restoreDrillDue,

} from './maintenanceDrill';



interface AdminNotificationsProps {

  status: SystemStatus;

}



interface NotificationItem {

  id: string;

  color: string;

  title: string;

  message: string;

  actionLabel?: string;

  actionPath?: string;

}



function buildNotifications(status: SystemStatus): NotificationItem[] {

  const items: NotificationItem[] = [];



  if (status.storageMode === 'memory') {

    items.push({

      id: 'memory-storage',

      color: 'red',

      title: 'Data is not persisted',

      message: 'Set DATABASE_URL and run migrations so patient records survive restarts.',

    });

  } else if (!status.databaseConnected) {

    items.push({

      id: 'db-offline',

      color: 'red',

      title: 'Database offline',

      message: 'The clinic database file is unreachable. Check DATABASE_URL in /etc/pococlinic/env and that the SQLite file exists.',

    });

  }



  if (status.pendingMigrations.length > 0) {

    items.push({

      id: 'migrations',

      color: 'orange',

      title: `${status.pendingMigrations.length} pending migration(s)`,

      message: 'Run migrations before continuing: migrate.bat (dev) or sudo /opt/pococlinic/bin/migrate (production server).',

    });

  }



  if (status.backup.status === 'critical') {

    items.push({

      id: 'backup-critical',

      color: 'red',

      title: 'Backup overdue',

      message: 'No recent backup found. Run a USB backup today using the backup helper or Backup now below.',

      actionLabel: 'Backup guide',

      actionPath: '/help/daily-backup',

    });

  } else if (status.backup.status === 'warning') {

    items.push({

      id: 'backup-warning',

      color: 'yellow',

      title: 'Backup aging',

      message: 'The last backup is more than 24 hours old. Schedule a USB backup before end of day.',

    });

  }



  if (status.counts.lockedAccounts > 0) {

    items.push({

      id: 'locked-accounts',

      color: 'orange',

      title: `${status.counts.lockedAccounts} locked account(s)`,

      message: 'Staff accounts are locked after failed sign-in attempts. Review staff records if someone needs help.',

      actionLabel: 'View staff',

      actionPath: '/users',

    });

  }



  if (status.counts.defaultPinAccounts > 0) {

    items.push({

      id: 'default-pins',

      color: 'orange',

      title: `${status.counts.defaultPinAccounts} staff on default PIN`,

      message: 'Every staff member should choose a private PIN. Review staff records if someone needs help.',

      actionLabel: 'View staff',

      actionPath: '/users',

    });

  }



  if (status.storageMode === 'database' && status.security && !status.security.documentsStorageReady) {

    items.push({

      id: 'documents-storage',

      color: 'red',

      title: 'Document storage unavailable',

      message: 'Patient document uploads need an encrypted database connection. Check DATABASE_URL and DOCUMENT_ENCRYPTION_KEY on the server.',

    });

  }



  if (restoreDrillDue(loadLastRestoreDrill())) {

    const days = daysSinceLastDrill(loadLastRestoreDrill());

    items.push({

      id: 'restore-drill',

      color: 'yellow',

      title: 'Quarterly restore drill due',

      message: days === null

        ? 'No restore drill recorded. Test a USB backup on a spare machine when the clinic is closed.'

        : `Last drill was ${days} days ago. Verify backups still restore correctly.`,

      actionLabel: 'Drill checklist',

      actionPath: '/help/quarterly-drill',

    });

  }



  if (monthlyTestingDue()) {

    items.push({

      id: 'monthly-testing',

      color: 'yellow',

      title: 'Monthly testing checklist due',

      message: 'Complete the in-app monthly checklist — compliance, health check, backups, and staff access.',

      actionLabel: 'Open checklist',

      actionPath: '/help/compliance-checklist',

    });

  }



  if (complianceCheckDue()) {

    items.push({

      id: 'compliance-check',

      color: 'yellow',

      title: 'Run compliance checklist',

      message: 'The HIPAA compliance checklist has not been run this month. Use Overview → Run checklist.',

      actionLabel: 'Admin dashboard',

      actionPath: '/admin',

    });

  }



  if (auditArchiveDue()) {
    items.push({
      id: 'audit-archive',
      color: 'yellow',
      title: 'Archive audit log',
      message:
        'Download the audit CSV (admin only), save it to offline storage or your ops binder, then mark archived. Do this before monthly audit purge removes old rows.',
      actionLabel: 'Audit log',
      actionPath: '/admin/audit',
    });
  }

  if (healthCheckDue()) {

    items.push({

      id: 'health-check',

      color: 'yellow',

      title: 'System health check due',

      message: 'Run a health check on the Backup & maintenance tab to verify storage and document integrity.',

      actionLabel: 'Backup & maintenance',

      actionPath: '/admin?tab=backup',

    });

  }



  return items;

}



export function AdminNotifications({ status }: AdminNotificationsProps) {

  const navigate = useNavigate();

  const notifications = buildNotifications(status);



  if (notifications.length === 0) {

    return (

      <Alert color="green" title="All clear">

        No administrator actions are required right now.

      </Alert>

    );

  }



  return (

    <Stack gap="sm">

      <Title order={4}>Action needed</Title>

      {notifications.map((item) => (

        <Alert

          key={item.id}

          color={item.color}

          title={item.title}

          variant="light"

        >

          <Stack gap="xs">

            <Text size="sm">{item.message}</Text>

            {item.actionLabel && item.actionPath && (

              <Text

                size="sm"

                fw={600}

                style={{ cursor: 'pointer' }}

                onClick={() => navigate(item.actionPath!)}

              >

                {item.actionLabel} →

              </Text>

            )}

          </Stack>

        </Alert>

      ))}

    </Stack>

  );

}


