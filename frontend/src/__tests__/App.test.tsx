import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MantineProvider } from '@mantine/core';
import App from '../App';

vi.mock('../api/auth', () => ({
  clearLegacyStoredToken: vi.fn(),
  logoutRequest: vi.fn(),
  refreshAccessToken: vi.fn(async () => null),
  refreshAccessTokenOnce: vi.fn(async () => null),
  getCurrentUser: vi.fn(async () => ({
    id: '1',
    name: 'Test User',
    email: 'test@example.com',
    role: 'admin',
  })),
  loginStaff: vi.fn(),
  loginAdmin: vi.fn(),
  attachAuthInterceptor: vi.fn(),
}));

vi.mock('../context/AuthContext', () => ({
  AuthProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  useAuth: () => ({
    user: { id: '1', name: 'Test User', email: 'test@example.com', role: 'admin' },
    isAuthenticated: true,
    isLoading: false,
    loginStaff: vi.fn(),
    loginAdmin: vi.fn(),
    logout: vi.fn(async () => undefined),
  }),
}));

vi.mock('../components/patients/PatientList', () => ({
  PatientList: () => <div>Patient List</div>,
}));

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: false },
  },
});

const TestWrapper = ({ children }: { children: React.ReactNode }) => (
  <QueryClientProvider client={queryClient}>
    <MantineProvider>{children}</MantineProvider>
  </QueryClientProvider>
);

describe('App Component', () => {
  it('renders without crashing', async () => {
    render(<App />, { wrapper: TestWrapper });
    expect(await screen.findByRole('button', { name: /PocoClinic home/i })).toBeInTheDocument();
  });

  it('contains header elements', async () => {
    render(<App />, { wrapper: TestWrapper });
    expect(await screen.findByRole('button', { name: /PocoClinic home/i })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Patients' })).toBeInTheDocument();
  });
});
