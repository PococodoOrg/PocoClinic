import { Button, Card, Group, Stack, Text, Title } from '@mantine/core';
import { IconPrinter } from '@tabler/icons-react';
import { BackupStatus } from '../../api/admin';
import './adminPrint.css';

interface BackupStickerProps {
  backup: BackupStatus;
  clinicName?: string;
}

function formatBackupAge(hours?: number): string {
  if (hours == null) {
    return 'Unknown';
  }
  if (hours < 24) {
    return `${Math.round(hours)} hours ago`;
  }
  return `${Math.round(hours / 24)} days ago`;
}

export function BackupSticker({ backup, clinicName = 'PocoClinic' }: BackupStickerProps) {
  const handlePrint = () => {
    window.print();
  };

  return (
    <Card withBorder padding="md" className="backup-sticker">
      <Group justify="space-between" mb="sm">
        <Title order={4}>Backup label</Title>
        <Button leftSection={<IconPrinter size={16} />} variant="light" size="xs" onClick={handlePrint}>
          Print label
        </Button>
      </Group>
      <Stack gap={4} className="backup-sticker-body">
        <Text fw={700} size="lg">{clinicName}</Text>
        <Text size="sm">Last backup: {backup.lastBackupFile ?? 'None recorded'}</Text>
        <Text size="sm">
          Created: {backup.lastBackupAt ? new Date(backup.lastBackupAt).toLocaleString() : '—'}
        </Text>
        <Text size="sm">Age: {formatBackupAge(backup.backupAgeHours)}</Text>
        <Text size="xs" c="dimmed" mt="sm">
          Affix to USB rotation drive. Verify backup integrity monthly.
        </Text>
      </Stack>
    </Card>
  );
}
