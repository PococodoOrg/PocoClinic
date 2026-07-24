import {
  Badge,
  Button,
  Card,
  Group,
  SimpleGrid,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useEffect, useState } from 'react';
import {
  DEFAULT_USB_DRIVES,
  loadUsbDrives,
  saveUsbDrives,
  todaysDrive,
  UsbDrive,
  weekdayLabel,
} from './usbRotation';
import './adminPrint.css';

export function UsbRotationPanel() {
  const [drives, setDrives] = useState<UsbDrive[]>(DEFAULT_USB_DRIVES);
  const today = todaysDrive(drives);

  useEffect(() => {
    setDrives(loadUsbDrives());
  }, []);

  const updateDrive = (id: string, patch: Partial<UsbDrive>) => {
    const next = drives.map((drive) => (drive.id === id ? { ...drive, ...patch } : drive));
    setDrives(next);
    saveUsbDrives(next);
  };

  const resetDefaults = () => {
    setDrives(DEFAULT_USB_DRIVES);
    saveUsbDrives(DEFAULT_USB_DRIVES);
  };

  return (
    <Card withBorder padding="md" className="usb-rotation-panel">
      <Group justify="space-between" mb="md">
        <div>
          <Title order={4}>USB rotation</Title>
          <Text size="sm" c="dimmed">
            Label each drive and use today&apos;s stick for backups.
          </Text>
        </div>
        <Group>
          <Button size="xs" variant="default" onClick={() => window.print()}>
            Print labels
          </Button>
          <Button size="xs" variant="subtle" onClick={resetDefaults}>
            Reset defaults
          </Button>
        </Group>
      </Group>

      {today && (
        <Badge size="lg" color="blue" mb="md" variant="light">
          Today: {today.label} — {today.name}
        </Badge>
      )}

      <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="md">
        {drives.map((drive) => (
          <Card
            key={drive.id}
            withBorder
            padding="md"
            className="usb-label-card"
            style={{
              borderColor: today?.id === drive.id ? 'var(--mantine-color-blue-5)' : undefined,
            }}
          >
            <Stack gap="xs">
              <Text fw={700} size="xl" ta="center" className="usb-label-text">
                {drive.label}
              </Text>
              <TextInput
                label="Drive name"
                value={drive.name}
                onChange={(event) => updateDrive(drive.id, { name: event.currentTarget.value })}
              />
              <TextInput
                label="Label text"
                value={drive.label}
                onChange={(event) => updateDrive(drive.id, { label: event.currentTarget.value.toUpperCase() })}
              />
              <Text size="xs" c="dimmed">
                Use on: {weekdayLabel(drive.weekdays)}
              </Text>
            </Stack>
          </Card>
        ))}
      </SimpleGrid>

      <Text size="sm" c="dimmed" mt="md">
        Keep at least two labeled drives in rotation. Never rely on a single USB stick.
      </Text>
    </Card>
  );
}
