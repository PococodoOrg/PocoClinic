import { MantineProvider } from '@mantine/core';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import Login from './Login';

const mockNavigate = vi.fn();
const mockLoginStaff = vi.fn();
const mockClearAuthError = vi.fn();

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

vi.mock('../context/AuthContext', () => ({
  useAuth: () => ({
    loginStaff: mockLoginStaff,
    authError: null,
    clearAuthError: mockClearAuthError,
  }),
}));

function renderLogin(state?: { from?: string; reason?: string }) {
  render(
    <MantineProvider>
      <MemoryRouter initialEntries={[{ pathname: '/login', state }]}>
        <Login />
      </MemoryRouter>
    </MantineProvider>,
  );
}

describe('Login page', () => {
  it('shows session expired message', () => {
    renderLogin({ reason: 'session-expired' });
    expect(screen.getByText(/your session ended/i)).toBeInTheDocument();
  });

  it('renders badge entry and staff sign-in heading', () => {
    renderLogin();
    expect(screen.getByRole('heading', { name: /staff sign in/i })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: /badge/i })).toBeInTheDocument();
    expect(screen.getByLabelText(/four-digit pin/i)).toBeInTheDocument();
  });

  it('keeps sign in disabled until badge and PIN are complete', async () => {
    const user = userEvent.setup();
    renderLogin();

    const signIn = screen.getByRole('button', { name: /sign in/i });
    expect(signIn).toBeDisabled();

    await user.type(screen.getByRole('textbox', { name: /badge/i }), 'BADGE-12345678');
    expect(signIn).toBeDisabled();
  });
});
