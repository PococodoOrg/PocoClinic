import { MantineProvider } from '@mantine/core';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import { AdminRoute } from './AdminRoute';

const mockUseAuth = vi.fn();

vi.mock('../../context/AuthContext', () => ({
  useAuth: () => mockUseAuth(),
}));

function renderAdminRoute(initialPath = '/admin') {
  render(
    <MantineProvider>
      <MemoryRouter initialEntries={[initialPath]}>
        <Routes>
          <Route
            path="/admin"
            element={
              <AdminRoute>
                <div>Admin dashboard</div>
              </AdminRoute>
            }
          />
          <Route path="/login/admin" element={<div>Admin login</div>} />
        </Routes>
      </MemoryRouter>
    </MantineProvider>,
  );
}

describe('AdminRoute', () => {
  it('redirects unauthenticated users to admin login', () => {
    mockUseAuth.mockReturnValue({ isAuthenticated: false, isLoading: false, user: null });
    renderAdminRoute();
    expect(screen.getByText('Admin login')).toBeInTheDocument();
  });

  it('blocks non-admin staff', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: { role: 'nurse' },
    });
    renderAdminRoute();
    expect(screen.getByText(/administrator access required/i)).toBeInTheDocument();
  });

  it('renders admin content for admin users', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: { role: 'admin' },
    });
    renderAdminRoute();
    expect(screen.getByText('Admin dashboard')).toBeInTheDocument();
  });
});
