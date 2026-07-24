import { Alert, Button, Card, Group, Stack, Text, Title } from '@mantine/core';

import { useState } from 'react';

import { useNavigate } from 'react-router-dom';

import { SystemStatus } from '../../api/admin';

import {
  auditArchiveDue,
  complianceCheckDue,
  healthCheckDue,
  loadLastAuditArchive,
  loadLastComplianceCheck,
  loadLastMonthlyTesting,
  monthlyTestingDue,
  saveAuditArchiveComplete,
  saveComplianceCheckComplete,
  saveMonthlyTestingComplete,
} from './adminReminders';

import {

  daysSinceLastDrill,

  loadLastRestoreDrill,

  restoreDrillDue,

  RESTORE_DRILL_INTERVAL_DAYS,

  saveLastRestoreDrill,

} from './maintenanceDrill';



interface MaintenanceRemindersProps {

  status: SystemStatus;

}



export function MaintenanceReminders({ status }: MaintenanceRemindersProps) {

  const navigate = useNavigate();

  const [lastDrill, setLastDrill] = useState<Date | null>(() => loadLastRestoreDrill());

  const [monthlyDone, setMonthlyDone] = useState<string | null>(() => loadLastMonthlyTesting());

  const [complianceDone, setComplianceDone] = useState<string | null>(() => loadLastComplianceCheck());
  const [auditArchiveDone, setAuditArchiveDone] = useState<string | null>(() => loadLastAuditArchive());



  const drillDue = restoreDrillDue(lastDrill);

  const daysSince = daysSinceLastDrill(lastDrill);

  const monthlyDue = monthlyTestingDue();

  const complianceDue = complianceCheckDue();
  const healthDue = healthCheckDue();
  const auditArchiveDueNow = auditArchiveDue();



  const markDrillComplete = () => {

    saveLastRestoreDrill();

    setLastDrill(new Date());

  };



  const markMonthlyComplete = () => {

    saveMonthlyTestingComplete();

    setMonthlyDone(loadLastMonthlyTesting());

  };



  const markComplianceComplete = () => {
    saveComplianceCheckComplete();
    setComplianceDone(loadLastComplianceCheck());
  };

  const markAuditArchiveComplete = () => {
    saveAuditArchiveComplete();
    setAuditArchiveDone(loadLastAuditArchive());
  };

  const backupDue = status.backup.status === 'critical' || status.backup.status === 'warning';

  return (

    <Card withBorder padding="md">

      <Title order={4} mb="md">Maintenance reminders</Title>

      <Stack gap="sm">

        {status.pendingMigrations.length > 0 && (

          <Alert color="orange" title="Pending migrations">

            Run migrations:{' '}
            <Text span ff="monospace">
              migrate.bat
            </Text>{' '}
            (dev) or{' '}
            <Text span ff="monospace">
              sudo /opt/pococlinic/bin/migrate
            </Text>{' '}
            (production).

          </Alert>

        )}



        {!status.security?.documentsStorageReady && status.storageMode === 'database' && (

          <Alert color="red" title="Document storage unavailable">

            Patient document uploads may fail. Check DATABASE_URL and DOCUMENT_ENCRYPTION_KEY on the server.

          </Alert>

        )}



        {backupDue ? (
          <Alert
            color={status.backup.status === 'critical' ? 'red' : 'yellow'}
            title={status.backup.status === 'critical' ? 'BACKUP overdue' : 'BACKUP aging'}
          >
            <Stack gap="xs">
              <Text size="sm">
                {status.backup.status === 'critical'
                  ? 'No recent backup found. Run a USB backup today before end of clinic.'
                  : 'Last backup is more than 24 hours old. Swap in today\'s labeled USB drive and run backup now.'}
              </Text>
              <Group>
                <Button size="xs" variant="light" onClick={() => navigate('/help/daily-backup')}>
                  Backup guide
                </Button>
                <Button size="xs" onClick={() => navigate('/admin?tab=backup')}>
                  Go to backup tab
                </Button>
              </Group>
            </Stack>
          </Alert>
        ) : (
          <Alert color="green" title="BACKUP up to date">
            Last backup looks current. Keep rotating labeled USB drives daily.
          </Alert>
        )}

        {auditArchiveDueNow ? (
          <Alert color="yellow" title="ARCHIVE AUDIT due">
            <Stack gap="xs">
              <Text size="sm">
                Download the audit CSV (admin only), save it offline or in your ops binder, then mark archived.
                {auditArchiveDone ? ` Last archived: ${auditArchiveDone}.` : ''}
              </Text>
              <Group>
                <Button size="xs" variant="light" onClick={() => navigate('/admin/audit')}>
                  Open audit log
                </Button>
                <Button size="xs" onClick={markAuditArchiveComplete}>
                  Mark archived
                </Button>
              </Group>
            </Stack>
          </Alert>
        ) : (
          <Alert color="green" title="ARCHIVE AUDIT done">
            Audit log marked archived for {auditArchiveDone}.
          </Alert>
        )}

        {drillDue ? (

          <Alert color="yellow" title="Quarterly restore drill due">

            <Stack gap="xs">

              <Text size="sm">

                {lastDrill

                  ? `Last drill was ${daysSince} days ago. Schedule a test restore on a spare machine (not during clinic hours).`

                  : `No restore drill recorded. Test a USB backup at least every ${RESTORE_DRILL_INTERVAL_DAYS} days.`}

              </Text>

              <Group>

                <Button size="xs" variant="light" onClick={() => navigate('/help/quarterly-drill')}>

                  Drill checklist

                </Button>

                <Button size="xs" onClick={markDrillComplete}>

                  Mark drill completed

                </Button>

              </Group>

            </Stack>

          </Alert>

        ) : (

          <Alert color="green" title="Restore drill up to date">

            Last completed {lastDrill?.toLocaleDateString()} ({daysSince} days ago).

            <Button size="xs" variant="subtle" ml="sm" onClick={markDrillComplete}>

              Log again

            </Button>

          </Alert>

        )}



        {monthlyDue ? (

          <Alert color="yellow" title="Monthly testing checklist due">

            <Stack gap="xs">

              <Text size="sm">

                Walk through the monthly testing checklist — compliance, health check, backups, and staff access.

                {monthlyDone ? ` Last completed: ${monthlyDone}.` : ''}

              </Text>

              <Group>

                <Button size="xs" variant="light" onClick={() => navigate('/help/compliance-checklist')}>

                  Open checklist

                </Button>

                <Button size="xs" onClick={markMonthlyComplete}>

                  Mark month complete

                </Button>

              </Group>

            </Stack>

          </Alert>

        ) : (

          <Alert color="green" title="Monthly testing recorded">

            Checklist marked complete for {monthlyDone}.

          </Alert>

        )}



        {complianceDue ? (

          <Alert color="yellow" title="Compliance checklist not run this month">

            <Stack gap="xs">

              <Text size="sm">Run the HIPAA compliance checklist on the Overview tab.</Text>

              <Group>

                <Button size="xs" variant="light" onClick={() => navigate('/admin')}>

                  Go to Overview

                </Button>

                <Button size="xs" onClick={markComplianceComplete}>

                  Mark completed

                </Button>

              </Group>

            </Stack>

          </Alert>

        ) : (

          <Alert color="green" title="Compliance check recorded this month">

            Last marked {complianceDone}.

          </Alert>

        )}



        {healthDue && (

          <Alert color="yellow" title="System health check due">

            <Text size="sm">

              Run a health check on the Backup &amp; maintenance tab to verify database, storage, and document integrity.

            </Text>

          </Alert>

        )}



        {status.counts.defaultPinAccounts > 0 && (

          <Alert color="orange" title={`${status.counts.defaultPinAccounts} staff on default PIN`}>

            These accounts must change PIN on next sign-in. Review staff records if someone needs help.

          </Alert>

        )}

      </Stack>

    </Card>

  );

}


