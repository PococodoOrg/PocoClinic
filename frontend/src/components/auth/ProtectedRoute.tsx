import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { Center, LoadingOverlay, Stack, Text } from '@mantine/core';
import { useAuth } from '../../context/AuthContext';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const PIN_CHANGE_EXEMPT_PREFIXES = ['/account/pin', '/help'];

function isPinChangeExempt(pathname: string): boolean {
  return PIN_CHANGE_EXEMPT_PREFIXES.some(
    (prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`),
  );
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, user } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <Center style={{ minHeight: '50vh', position: 'relative' }}>
        <LoadingOverlay visible />
        <Stack gap="xs" align="center" style={{ zIndex: 1 }}>
          <Text c="dimmed" size="sm">
            Checking your session...
          </Text>
        </Stack>
      </Center>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  if (user?.mustChangePin && !isPinChangeExempt(location.pathname)) {
    return <Navigate to="/account/pin" replace state={{ from: location.pathname }} />;
  }

  return <>{children}</>;
}
