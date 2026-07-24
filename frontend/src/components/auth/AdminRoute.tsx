import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { Center, LoadingOverlay, Text } from '@mantine/core';
import { useAuth } from '../../context/AuthContext';

interface AdminRouteProps {
  children: React.ReactNode;
}

export function AdminRoute({ children }: AdminRouteProps) {
  const { user, isAuthenticated, isLoading } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <Center style={{ minHeight: '50vh', position: 'relative' }}>
        <LoadingOverlay visible />
      </Center>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login/admin" replace state={{ from: location.pathname }} />;
  }

  if (user?.role !== 'admin') {
    return (
      <Center style={{ minHeight: '50vh' }}>
        <Text c="dimmed">Administrator access required.</Text>
      </Center>
    );
  }

  return <>{children}</>;
}
