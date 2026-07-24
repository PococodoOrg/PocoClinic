import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { AuthProvider, useAuth } from './AuthContext';

vi.mock('../api/auth', () => ({
  clearLegacyStoredToken: vi.fn(),
  getCurrentUser: vi.fn(),
  loginStaff: vi.fn(),
  loginAdmin: vi.fn(),
  logoutRequest: vi.fn(),
  refreshAccessToken: vi.fn(),
}));

vi.mock('../queryClient', () => ({
  queryClient: { clear: vi.fn() },
}));

import {
  clearLegacyStoredToken,
  getCurrentUser,
  loginStaff,
  logoutRequest,
  refreshAccessToken,
} from '../api/auth';

function AuthProbe() {
  const { user, isAuthenticated, isLoading, authError } = useAuth();
  if (isLoading) {
    return <div>loading</div>;
  }
  return (
    <div>
      <span data-testid="authenticated">{String(isAuthenticated)}</span>
      <span data-testid="email">{user?.email ?? 'none'}</span>
      <span data-testid="error">{authError ?? 'none'}</span>
    </div>
  );
}

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getCurrentUser).mockRejectedValue(new Error('no session'));
    vi.mocked(refreshAccessToken).mockResolvedValue(null);
  });

  it('restores an existing session on mount', async () => {
    vi.mocked(getCurrentUser).mockResolvedValue({
      id: 'user-1',
      email: 'staff@clinic.test',
      name: 'Staff',
      role: 'staff',
      mustChangePin: false,
    });

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId('authenticated')).toHaveTextContent('true');
    });
    expect(screen.getByTestId('email')).toHaveTextContent('staff@clinic.test');
    expect(clearLegacyStoredToken).toHaveBeenCalled();
  });

  it('falls back to refresh when getCurrentUser fails', async () => {
    vi.mocked(refreshAccessToken).mockResolvedValue({
      user: {
        id: 'user-2',
        email: 'doctor@clinic.test',
        name: 'Doctor',
        role: 'doctor',
        mustChangePin: false,
      },
    });

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId('email')).toHaveTextContent('doctor@clinic.test');
    });
    expect(refreshAccessToken).toHaveBeenCalled();
  });

  it('loginStaff stores the returned user', async () => {
    vi.mocked(loginStaff).mockResolvedValue({
      user: {
        id: 'user-3',
        email: 'nurse@clinic.test',
        name: 'Nurse',
        role: 'nurse',
        mustChangePin: false,
      },
    });

    function LoginProbe() {
      const { loginStaff: doLogin, user, isLoading } = useAuth();
      if (isLoading) {
        return <div>loading</div>;
      }
      return (
        <div>
          <button type="button" onClick={() => doLogin({ key: 'badge', pin: '1234' })}>
            login
          </button>
          <span data-testid="email">{user?.email ?? 'none'}</span>
        </div>
      );
    }

    render(
      <AuthProvider>
        <LoginProbe />
      </AuthProvider>,
    );

    await waitFor(() => {
      expect(screen.queryByText('loading')).not.toBeInTheDocument();
    });

    screen.getByRole('button', { name: 'login' }).click();

    await waitFor(() => {
      expect(screen.getByTestId('email')).toHaveTextContent('nurse@clinic.test');
    });
  });

  it('logout clears the user', async () => {
    vi.mocked(getCurrentUser).mockResolvedValue({
      id: 'user-4',
      email: 'admin@clinic.test',
      name: 'Admin',
      role: 'admin',
      mustChangePin: false,
    });
    vi.mocked(logoutRequest).mockResolvedValue(undefined);

    function LogoutProbe() {
      const { logout, user, isLoading } = useAuth();
      if (isLoading) {
        return <div>loading</div>;
      }
      return (
        <div>
          <button type="button" onClick={() => logout()}>
            logout
          </button>
          <span data-testid="email">{user?.email ?? 'none'}</span>
        </div>
      );
    }

    render(
      <AuthProvider>
        <LogoutProbe />
      </AuthProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId('email')).toHaveTextContent('admin@clinic.test');
    });

    screen.getByRole('button', { name: 'logout' }).click();

    await waitFor(() => {
      expect(screen.getByTestId('email')).toHaveTextContent('none');
    });
  });

  it('throws when useAuth is used outside provider', () => {
    expect(() => render(<AuthProbe />)).toThrow(/AuthProvider/);
  });
});
