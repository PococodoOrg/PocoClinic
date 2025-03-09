import React from 'react';
import { AppShell, Box, Text, Group, Button, Container } from '@mantine/core';
import { useNavigate, useLocation } from 'react-router-dom';

interface AppLayoutProps {
  children: React.ReactNode;
}

export function AppLayout({ children }: AppLayoutProps) {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <AppShell
      padding="md"
      header={{ height: 60 }}
    >
      <Box component="header" p="xs">
        <Container size="xl">
          <Group justify="space-between">
            <Group>
              <Text size="xl" fw={700} onClick={() => navigate('/')} style={{ cursor: 'pointer' }}>
                PocoClinic EMR
              </Text>
            </Group>
            <Group>
              <Button
                variant={location.pathname === '/patients' ? 'filled' : 'light'}
                onClick={() => navigate('/patients')}
              >
                Patients
              </Button>
              {/* Add more navigation buttons as needed */}
            </Group>
          </Group>
        </Container>
      </Box>
      <Container size="xl">
        {children}
      </Container>
    </AppShell>
  );
} 