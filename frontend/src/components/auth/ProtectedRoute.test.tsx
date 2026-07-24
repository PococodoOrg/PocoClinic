import { MantineProvider } from '@mantine/core';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import { ProtectedRoute } from './ProtectedRoute';

const mockUseAuth = vi.fn();

vi.mock('../../context/AuthContext', () => ({
  useAuth: () => mockUseAuth(),
}));

function renderProtected(initialPath = '/patients') {
  return render(
    <MantineProvider>
      <MemoryRouter initialEntries={[initialPath]}>
        <Routes>
          <Route
            path="/patients"
            element={
              <ProtectedRoute>
                <div>Patient area</div>
              </ProtectedRoute>
            }
          />
          <Route
            path="/account/pin"
            element={
              <ProtectedRoute>
                <div>PIN change</div>
              </ProtectedRoute>
            }
          />
          <Route
            path="/help"
            element={
              <ProtectedRoute>
                <div>Help</div>
              </ProtectedRoute>
            }
          />
          <Route path="/login" element={<div>Login page</div>} />
        </Routes>
      </MemoryRouter>
    </MantineProvider>,
  );
}

describe('ProtectedRoute', () => {
  it('shows loading state while session restores', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: true,
      user: null,
    });

    renderProtected();
    expect(screen.getByText('Checking your session...')).toBeInTheDocument();
  });

  it('redirects unauthenticated users to login', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
      user: null,
    });

    renderProtected();
    expect(screen.getByText('Login page')).toBeInTheDocument();
  });

  it('renders children for authenticated users', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: { mustChangePin: false },
    });

    renderProtected();
    expect(screen.getByText('Patient area')).toBeInTheDocument();
  });

  it('redirects to PIN change when mustChangePin is set', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: { mustChangePin: true },
    });

    renderProtected('/patients');
    expect(screen.queryByText('Patient area')).not.toBeInTheDocument();
    expect(screen.getByText('PIN change')).toBeInTheDocument();
  });

  it('allows help while PIN change is required', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: { mustChangePin: true },
    });

    renderProtected('/help');
    expect(screen.getByText('Help')).toBeInTheDocument();
  });
});
