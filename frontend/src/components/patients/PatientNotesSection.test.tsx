import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { PatientNotesSection } from './PatientNotesSection';

vi.mock('../../api/patientNotes', () => ({
  fetchPatientNotes: vi.fn(),
  createPatientNote: vi.fn(),
  updatePatientNote: vi.fn(),
  deletePatientNote: vi.fn(),
}));

vi.mock('../../context/AuthContext', () => ({
  useAuth: () => ({ user: { id: 'user-1', role: 'doctor' } }),
}));

vi.mock('@mantine/notifications', () => ({
  notifications: { show: vi.fn() },
}));

import {
  createPatientNote,
  deletePatientNote,
  fetchPatientNotes,
  updatePatientNote,
} from '../../api/patientNotes';

const patientId = '11111111-1111-1111-1111-111111111111';

function renderSection() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <PatientNotesSection patientId={patientId} />
      </MantineProvider>
    </QueryClientProvider>,
  );
}

describe('PatientNotesSection', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(fetchPatientNotes).mockResolvedValue([]);
  });

  it('creates a clinical note', async () => {
    const user = userEvent.setup();
    vi.mocked(createPatientNote).mockResolvedValue({
      id: 'note-1',
      patientId,
      body: 'Follow-up complete',
      authorId: 'user-1',
      authorName: 'Doctor',
      createdAt: '2026-01-01T00:00:00.000Z',
      updatedAt: '2026-01-01T00:00:00.000Z',
    });

    renderSection();

    await user.type(screen.getByLabelText(/^new note$/i), 'Follow-up complete');
    await user.click(screen.getByRole('button', { name: /add note/i }));

    await waitFor(() => {
      expect(createPatientNote).toHaveBeenCalledWith(patientId, 'Follow-up complete');
    });
  });

  it('edits an existing note', async () => {
    const user = userEvent.setup();
    vi.mocked(fetchPatientNotes).mockResolvedValue([
      {
        id: 'note-1',
        patientId,
        body: 'Initial note',
        authorId: 'user-1',
        authorName: 'Doctor',
        createdAt: '2026-01-01T00:00:00.000Z',
        updatedAt: '2026-01-01T00:00:00.000Z',
      },
    ]);
    vi.mocked(updatePatientNote).mockResolvedValue({
      id: 'note-1',
      patientId,
      body: 'Updated note',
      authorId: 'user-1',
      authorName: 'Doctor',
      createdAt: '2026-01-01T00:00:00.000Z',
      updatedAt: '2026-01-02T00:00:00.000Z',
    });

    renderSection();

    await waitFor(() => {
      expect(screen.getByText('Initial note')).toBeInTheDocument();
    });

    await user.click(screen.getByLabelText(/edit note/i));
    const textarea = await screen.findByDisplayValue('Initial note');
    await user.clear(textarea);
    await user.type(textarea, 'Updated note');
    await user.click(screen.getByRole('button', { name: /save changes/i }));

    await waitFor(() => {
      expect(updatePatientNote).toHaveBeenCalledWith(patientId, 'note-1', 'Updated note');
    });
  });

  it('deletes a note after confirmation', async () => {
    const user = userEvent.setup();
    vi.mocked(fetchPatientNotes).mockResolvedValue([
      {
        id: 'note-1',
        patientId,
        body: 'Remove me',
        authorId: 'user-1',
        authorName: 'Doctor',
        createdAt: '2026-01-01T00:00:00.000Z',
        updatedAt: '2026-01-01T00:00:00.000Z',
      },
    ]);
    vi.mocked(deletePatientNote).mockResolvedValue(undefined);

    renderSection();

    await waitFor(() => {
      expect(screen.getByText('Remove me')).toBeInTheDocument();
    });

    await user.click(screen.getByLabelText(/delete note/i));
    const dialog = await screen.findByRole('dialog');
    await user.click(within(dialog).getByRole('button', { name: /^delete note$/i }));

    await waitFor(() => {
      expect(deletePatientNote).toHaveBeenCalledWith(patientId, 'note-1');
    });
  });
});
