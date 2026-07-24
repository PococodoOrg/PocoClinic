import { Card, Stack, Text, UnstyledButton } from '@mantine/core';
import { IconCheck } from '@tabler/icons-react';
import { BackupEntry } from '../api';

interface BackupPickerProps {
  backups: BackupEntry[];
  selected: string | null;
  onSelect: (filename: string) => void;
}

export function BackupPicker({ backups, selected, onSelect }: BackupPickerProps) {
  return (
    <Stack gap="sm">
      {backups.map((backup) => {
        const isSelected = selected === backup.filename;
        return (
          <UnstyledButton key={backup.filename} onClick={() => onSelect(backup.filename)} className="touch-picker">
            <Card
              withBorder
              padding="md"
              radius="md"
              className={isSelected ? 'touch-picker-card touch-picker-card--selected' : 'touch-picker-card'}
            >
              <Stack gap={4}>
                <GroupInline selected={isSelected} label={backup.friendlyAt} />
                <Text size="xs" c="dimmed">
                  {backup.filename}
                </Text>
              </Stack>
            </Card>
          </UnstyledButton>
        );
      })}
    </Stack>
  );
}

function GroupInline({ selected, label }: { selected: boolean; label: string }) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
      {selected && <IconCheck size={22} color="var(--mantine-color-teal-6)" />}
      <Text fw={700} size="lg">
        {label}
      </Text>
    </div>
  );
}
