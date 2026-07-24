import { Paper, Stack, Text, Title } from '@mantine/core';
import { QRCodeSVG } from 'qrcode.react';
import { DEFAULT_PIN } from '../../types/user';

interface BadgePrintCardProps {
  name: string;
  email: string;
  badgeKey: string;
  pinReminder?: string;
}

export function BadgePrintCard({ name, email, badgeKey, pinReminder }: BadgePrintCardProps) {
  return (
    <Paper
      id="badge-print-card"
      withBorder
      p="lg"
      style={{ maxWidth: 320, margin: '0 auto' }}
    >
      <Stack align="center" gap="sm">
        <Title order={4}>PocoClinic</Title>
        <Text fw={600} size="lg">
          {name}
        </Text>
        <Text size="sm" c="dimmed">
          {email}
        </Text>
        <QRCodeSVG value={badgeKey} size={180} level="M" includeMargin />
        <Text size="xs" c="dimmed" ta="center">
          {pinReminder ?? `Scan this badge at staff sign-in, then enter PIN ${DEFAULT_PIN}`}
        </Text>
      </Stack>
    </Paper>
  );
}
