import { Badge, Box, Group, Stack, Text, UnstyledButton } from '@mantine/core';
import {
  IconDeviceUsb,
  IconHome,
  IconRefresh,
  IconShield,
} from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { ReactNode } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { fetchStatus, statusHeadline } from '../api';

const navItems = [
  { to: '/', label: 'Status', icon: IconHome },
  { to: '/backup', label: 'Backup', icon: IconDeviceUsb },
  { to: '/restore', label: 'Restore', icon: IconRefresh },
];

export function PiTouchLayout({ children }: { children: ReactNode }) {
  const location = useLocation();
  const { data: status } = useQuery({
    queryKey: ['helper-status'],
    queryFn: fetchStatus,
    refetchInterval: 30_000,
  });

  const headline = status ? statusHeadline(status) : null;

  return (
    <Box className="pi-shell">
      <header className="pi-header">
        <Group justify="space-between" wrap="nowrap">
          <Stack gap={0}>
            <Text fw={800} size="lg">
              PocoClinic Backup
            </Text>
            <Text size="xs" c="dimmed">
              Touch to back up or restore
            </Text>
          </Stack>
          {headline && (
            <Badge color={headline.color} size="lg" variant="filled" leftSection={<IconShield size={14} />}>
              {headline.title}
            </Badge>
          )}
        </Group>
      </header>

      <main className="pi-main">{children}</main>

      <nav className="pi-nav" aria-label="Backup and restore">
        {navItems.map((item) => {
          const active = location.pathname === item.to;
          const Icon = item.icon;
          return (
            <UnstyledButton
              key={item.to}
              component={NavLink}
              to={item.to}
              className={`pi-nav-item${active ? ' pi-nav-item--active' : ''}`}
            >
              <Icon size={28} stroke={active ? 2.2 : 1.6} />
              <Text fw={active ? 800 : 600} size="sm">
                {item.label}
              </Text>
            </UnstyledButton>
          );
        })}
      </nav>
    </Box>
  );
}
