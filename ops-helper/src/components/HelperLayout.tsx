import { Anchor, Container, Group, Stack, Text, ThemeIcon, Title } from '@mantine/core';
import { IconHeartHandshake, IconShieldCheck } from '@tabler/icons-react';
import { Link, useLocation } from 'react-router-dom';
import { ReactNode } from 'react';

const links = [
  { to: '/', label: 'Home' },
  { to: '/backup', label: 'Daily backup' },
  { to: '/restore', label: 'Restore' },
  { to: '/checklist', label: 'Print checklist' },
];

export function HelperLayout({ children }: { children: ReactNode }) {
  const location = useLocation();

  return (
    <div className="helper-shell">
      <Container size="sm" py="xl">
        <Stack gap="xl">
          <Stack gap="xs" align="center" ta="center">
            <ThemeIcon size={64} radius="xl" color="teal" variant="light">
              <IconHeartHandshake size={34} />
            </ThemeIcon>
            <Title order={1}>PocoClinic Backup Helper</Title>
            <Text c="dimmed" maw={520}>
              A simple guide that runs only on this computer. No internet needed — just plug in your USB and follow the steps.
            </Text>
            <Group gap="xs" c="teal">
              <IconShieldCheck size={16} />
              <Text size="sm">Local only · Not visible on the clinic network</Text>
            </Group>
          </Stack>

          <Group justify="center" gap="sm" wrap="wrap" className="no-print">
            {links.map((link) => (
              <Anchor
                key={link.to}
                component={Link}
                to={link.to}
                fw={location.pathname === link.to ? 700 : 500}
                c={location.pathname === link.to ? 'teal.7' : 'dimmed'}
                underline="never"
              >
                {link.label}
              </Anchor>
            ))}
          </Group>

          {children}
        </Stack>
      </Container>
    </div>
  );
}
