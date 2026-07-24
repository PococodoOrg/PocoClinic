import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { PatientList } from './PatientList';

const mockNavigate = vi.fn();

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

vi.mock('@mantine/hooks', () => ({
  useMediaQuery: () => false,
}));

vi.mock('../../api/patients', () => ({
  fetchPatients: vi.fn(),
}));

import { fetchPatients } from '../../api/patients';

const patientId = '11111111-1111-1111-1111-111111111111';

function renderList() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <MemoryRouter>
          <PatientList />
        </MemoryRouter>
      </MantineProvider>
    </QueryClientProvider>,
  );
}

describe('PatientList', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockNavigate.mockReset();
    vi.mocked(fetchPatients).mockResolvedValue({
      patients: [
        {
          id: patientId,
          firstName: 'Load',
          lastName: 'Patient',
          dateOfBirth: '1990-01-01',
          gender: 'female',
          email: 'load@example.com',
          phoneNumber: '555-1234',
          createdAt: '2026-01-01T00:00:00.000Z',
          updatedAt: '2026-01-01T00:00:00.000Z',
        },
      ],
      totalCount: 1,
      currentPage: 1,
      pageSize: 10,
      totalPages: 1,
    });
  });

  it('renders patients from the API', async () => {
    renderList();

    await waitFor(() => {
      expect(screen.getByText('Patient')).toBeInTheDocument();
      expect(screen.getByText('Load')).toBeInTheDocument();
    });
    expect(fetchPatients).toHaveBeenCalled();
  });

  it('shows an error alert when loading fails', async () => {
    vi.mocked(fetchPatients).mockRejectedValue(new Error('network'));

    renderList();

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(/error loading patients/i);
    });
  });

  it('navigates to patient details when a row is opened', async () => {
    const user = userEvent.setup();
    renderList();

    await waitFor(() => {
      expect(screen.getByText('Patient')).toBeInTheDocument();
      expect(screen.getByText('Load')).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: /view load patient/i }));

    expect(mockNavigate).toHaveBeenCalledWith(`/patients/${patientId}`);
  });

  it('searches patients when the search field changes', async () => {
    const user = userEvent.setup();
    renderList();

    await waitFor(() => {
      expect(screen.getByText('Patient')).toBeInTheDocument();
      expect(screen.getByText('Load')).toBeInTheDocument();
    });

    await user.type(screen.getByLabelText(/search patients/i), 'Load');

    await waitFor(() => {
      expect(fetchPatients).toHaveBeenCalledWith(
        expect.objectContaining({ search: 'Load', page: 1 }),
      );
    });
  });
});
