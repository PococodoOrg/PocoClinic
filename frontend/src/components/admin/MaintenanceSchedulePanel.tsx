import { Button, Card, Code, CopyButton, Group, Stack, Text, Title, Tooltip } from '@mantine/core';

const BACKUP_CRON = `# /etc/cron.d/pococlinic-backup
0 18 * * * root set -a && . /etc/pococlinic/env && set +a && /opt/pococlinic/backup >> /var/log/pococlinic/backup.log 2>&1`;

const AUDIT_CRON = `# /etc/cron.d/pococlinic-audit-purge (requires AUDIT_RETENTION_DAYS)
0 2 1 * * root set -a && . /etc/pococlinic/env && set +a && /opt/pococlinic/audit-purge >> /var/log/pococlinic/audit-purge.log 2>&1`;

function CronBlock({ title, description, snippet, installHint }: {
  title: string;
  description: string;
  snippet: string;
  installHint: string;
}) {
  return (
    <Stack gap="xs">
      <div>
        <Text fw={600} size="sm">{title}</Text>
        <Text size="sm" c="dimmed">{description}</Text>
      </div>
      <Code block style={{ whiteSpace: 'pre-wrap', fontSize: 12 }}>{snippet}</Code>
      <Group justify="space-between">
        <Text size="xs" c="dimmed">{installHint}</Text>
        <CopyButton value={snippet}>
          {({ copied, copy }) => (
            <Tooltip label={copied ? 'Copied' : 'Copy cron snippet'}>
              <Button size="xs" variant="light" onClick={copy}>
                {copied ? 'Copied' : 'Copy'}
              </Button>
            </Tooltip>
          )}
        </CopyButton>
      </Group>
    </Stack>
  );
}

export function MaintenanceSchedulePanel() {
  return (
    <Card withBorder padding="md">
      <Title order={4} mb="xs">Scheduled maintenance (cron)</Title>
      <Text size="sm" c="dimmed" mb="md">
        PocoClinic has no background scheduler — use system cron on the Pi. Snippets ship in the release tarball under <Code>cron/</Code>.
      </Text>
      <Stack gap="lg">
        <CronBlock
          title="Daily backup — 6:00 PM"
          description="Creates a tarball in BACKUP_DIR. Pair with USB rotation on the admin dashboard."
          snippet={BACKUP_CRON}
          installHint="Release: sudo cp /opt/pococlinic/cron/pococlinic-backup /etc/cron.d/"
        />
        <CronBlock
          title="Monthly audit purge — 1st at 2:00 AM"
          description="Removes audit rows older than AUDIT_RETENTION_DAYS. Run --dry-run first when tuning retention."
          snippet={AUDIT_CRON}
          installHint="Release: sudo cp /opt/pococlinic/cron/pococlinic-audit-purge /etc/cron.d/"
        />
      </Stack>
    </Card>
  );
}
