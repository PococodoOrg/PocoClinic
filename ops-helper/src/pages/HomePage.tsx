import { ReactNode } from 'react';
import {
  Badge,
  Button,
  Card,
  Group,
  SimpleGrid,
  Stack,
  Text,
  ThemeIcon,
  Title,
  UnstyledButton,
} from '@mantine/core';
import {
  IconChevronRight,
  IconDatabase,
  IconDeviceUsb,
  IconHeart,
  IconRefresh,
  IconShield,
  IconUsers,
} from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { fetchStatus, statusHeadline } from '../api';
import { usePiTouch } from '../context/PiTouchContext';

export function HomePage() {
  const { isPiTouch } = usePiTouch();
  const { data: status, isLoading, error } = useQuery({
    queryKey: ['helper-status'],
    queryFn: fetchStatus,
    refetchInterval: 30_000,
  });

  const headline = status ? statusHeadline(status) : null;

  if (isPiTouch) {
    return (
      <Stack gap="md">
        <Card withBorder padding="lg" radius="md" className="pi-status-card">
          {isLoading && <Text>Checking backup status…</Text>}
          {error && <Text c="red">Helper service not running. Ask your installer to start the backup helper.</Text>}
          {status && headline && (
            <Stack gap="sm">
              <Badge color={headline.color} size="xl" variant="light" w="fit-content">
                {headline.title}
              </Badge>
              <Text size="lg" fw={600}>
                {headline.detail}
              </Text>
              <SimpleGrid cols={2}>
                <MiniStat label="Patients" value={String(status.patientCount)} />
                <MiniStat label="Backups" value={String(status.backupCount)} />
              </SimpleGrid>
              {status.lastBackupAt && (
                <Text size="sm" c="dimmed">
                  Last backup: {new Date(status.lastBackupAt).toLocaleString()}
                </Text>
              )}
            </Stack>
          )}
        </Card>

        <Stack gap="sm">
          <Text fw={700} size="lg">
            What would you like to do?
          </Text>
          <UnstyledButton component={Link} to="/backup" className="touch-picker">
            <Card withBorder padding="lg" radius="md" className="pi-action-tile pi-action-tile--backup">
              <Group justify="space-between" wrap="nowrap">
                <Group gap="md" wrap="nowrap">
                  <ThemeIcon size={56} radius="md" color="teal" variant="light">
                    <IconDeviceUsb size={32} />
                  </ThemeIcon>
                  <Stack gap={2}>
                    <Text fw={800} size="xl">
                      Back up clinic
                    </Text>
                    <Text c="dimmed" size="sm">
                      Save today&apos;s data, then copy to USB
                    </Text>
                  </Stack>
                </Group>
                <IconChevronRight size={28} />
              </Group>
            </Card>
          </UnstyledButton>

          <UnstyledButton component={Link} to="/restore" className="touch-picker">
            <Card withBorder padding="lg" radius="md" className="pi-action-tile pi-action-tile--restore">
              <Group justify="space-between" wrap="nowrap">
                <Group gap="md" wrap="nowrap">
                  <ThemeIcon size={56} radius="md" color="orange" variant="light">
                    <IconRefresh size={32} />
                  </ThemeIcon>
                  <Stack gap={2}>
                    <Text fw={800} size="xl">
                      Restore clinic
                    </Text>
                    <Text c="dimmed" size="sm">
                      Only when recovering from a problem
                    </Text>
                  </Stack>
                </Group>
                <IconChevronRight size={28} />
              </Group>
            </Card>
          </UnstyledButton>
        </Stack>
      </Stack>
    );
  }

  return (
    <Stack gap="lg">
      <Card withBorder padding="xl" radius="lg" shadow="sm">
        {isLoading && <Text c="dimmed">Checking your backup status…</Text>}
        {error && <Text c="red">Could not reach the helper service. Is it running on this computer?</Text>}
        {status && headline && (
          <Stack gap="md">
            <Group justify="space-between" align="flex-start">
              <div>
                <Badge color={headline.color} size="lg" variant="light" mb="sm">
                  {headline.title}
                </Badge>
                <Title order={3}>{headline.detail}</Title>
              </div>
              <ThemeIcon size={52} radius="md" color={headline.color} variant="light">
                <IconShield size={28} />
              </ThemeIcon>
            </Group>
            <SimpleGrid cols={{ base: 1, xs: 3 }}>
              <Stat label="Patients protected" value={String(status.patientCount)} icon={<IconUsers size={18} />} />
              <Stat label="Backups saved" value={String(status.backupCount)} icon={<IconDatabase size={18} />} />
              <Stat
                label="Database"
                value={status.databaseConnected ? 'Connected' : 'Not ready'}
                icon={<IconHeart size={18} />}
              />
            </SimpleGrid>
            {status.lastBackupAt && (
              <Text size="sm" c="dimmed">
                Last backup: {new Date(status.lastBackupAt).toLocaleString()}
                {status.lastBackupFile ? ` (${status.lastBackupFile})` : ''}
              </Text>
            )}
          </Stack>
        )}
      </Card>

      <SimpleGrid cols={{ base: 1, sm: 2 }}>
        <ActionCard
          title="Daily backup"
          description="Step-by-step guide to save today's clinic data to this computer, then copy to USB."
          icon={<IconDeviceUsb size={28} />}
          to="/backup"
          buttonLabel="Start backup guide"
          color="teal"
        />
        <ActionCard
          title="Restore clinic data"
          description="Use only when recovering from a problem. We will walk you through every confirmation."
          icon={<IconRefresh size={28} />}
          to="/restore"
          buttonLabel="Start restore guide"
          color="orange"
        />
      </SimpleGrid>

      {status?.mainAppUrl && (
        <Card withBorder padding="md" radius="lg" className="no-print">
          <Group justify="space-between">
            <Text size="sm">Ready to chart patients?</Text>
            <Button component="a" href={status.mainAppUrl} variant="light" color="gray">
              Open PocoClinic
            </Button>
          </Group>
        </Card>
      )}
    </Stack>
  );
}

function MiniStat({ label, value }: { label: string; value: string }) {
  return (
    <Card withBorder padding="sm" radius="md">
      <Text size="xs" c="dimmed" tt="uppercase" fw={700}>
        {label}
      </Text>
      <Text fw={800} size="xl">
        {value}
      </Text>
    </Card>
  );
}

function Stat({ label, value, icon }: { label: string; value: string; icon: ReactNode }) {
  return (
    <Card withBorder padding="md" radius="md">
      <Group gap="xs" mb={4}>
        {icon}
        <Text size="xs" c="dimmed" tt="uppercase" fw={600}>
          {label}
        </Text>
      </Group>
      <Text fw={700} size="lg">
        {value}
      </Text>
    </Card>
  );
}

function ActionCard({
  title,
  description,
  icon,
  to,
  buttonLabel,
  color,
}: {
  title: string;
  description: string;
  icon: React.ReactNode;
  to: string;
  buttonLabel: string;
  color: string;
}) {
  return (
    <Card withBorder padding="xl" radius="lg" className="helper-card">
      <Stack gap="md">
        <ThemeIcon size={52} radius="md" color={color} variant="light">
          {icon}
        </ThemeIcon>
        <Title order={3}>{title}</Title>
        <Text c="dimmed">{description}</Text>
        <Button component={Link} to={to} color={color} size="md">
          {buttonLabel}
        </Button>
      </Stack>
    </Card>
  );
}
