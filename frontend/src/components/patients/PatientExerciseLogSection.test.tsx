import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { PatientExerciseLogSection } from './PatientExerciseLogSection';

vi.mock('../../api/exerciseLog', () => ({
  fetchExercisePlans: vi.fn(),
  fetchExerciseEntries: vi.fn(),
  createExercisePlan: vi.fn(),
  updateExercisePlan: vi.fn(),
  deleteExercisePlan: vi.fn(),
  createExerciseEntry: vi.fn(),
  updateExerciseEntry: vi.fn(),
  deleteExerciseEntry: vi.fn(),
}));

vi.mock('@mantine/notifications', () => ({
  notifications: { show: vi.fn() },
}));

import { notifications } from '@mantine/notifications';
import {
  createExerciseEntry,
  createExercisePlan,
  deleteExercisePlan,
  fetchExerciseEntries,
  fetchExercisePlans,
  updateExercisePlan,
} from '../../api/exerciseLog';

const patientId = '11111111-1111-1111-1111-111111111111';
const planId = '22222222-2222-2222-2222-222222222222';

function renderSection() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <PatientExerciseLogSection patientId={patientId} />
      </MantineProvider>
    </QueryClientProvider>,
  );
}

describe('PatientExerciseLogSection', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(fetchExercisePlans).mockResolvedValue([]);
    vi.mocked(fetchExerciseEntries).mockResolvedValue([]);
  });

  it('shows empty state and creates a plan', async () => {
    const user = userEvent.setup();
    vi.mocked(createExercisePlan).mockResolvedValue({
      id: planId,
      patientId,
      name: 'Knee strengthening',
      description: 'Week 1',
      status: 'active',
      createdBy: 'user-1',
      createdByName: 'Clinician',
      createdAt: '2026-01-01T00:00:00.000Z',
    });
    vi.mocked(fetchExercisePlans)
      .mockResolvedValueOnce([])
      .mockResolvedValue([
        {
          id: planId,
          patientId,
          name: 'Knee strengthening',
          description: 'Week 1',
          status: 'active',
          createdBy: 'user-1',
          createdByName: 'Clinician',
          createdAt: '2026-01-01T00:00:00.000Z',
        },
      ]);

    renderSection();

    expect(await screen.findByText(/No exercise plans yet/i)).toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: /New plan/i }));
    const nameInput = await screen.findByPlaceholderText(/Knee strengthening/i);
    await user.type(nameInput, 'Knee strengthening');
    await user.click(screen.getByRole('button', { name: /Create plan/i }));

    await waitFor(() => {
      expect(createExercisePlan).toHaveBeenCalledWith(patientId, {
        name: 'Knee strengthening',
        description: '',
      });
    });
  });

  it('logs a session for the selected plan', async () => {
    const user = userEvent.setup();
    vi.mocked(fetchExercisePlans).mockResolvedValue([
      {
        id: planId,
        patientId,
        name: 'Knee strengthening',
        description: '',
        status: 'active',
        createdBy: 'user-1',
        createdByName: 'Clinician',
        createdAt: '2026-01-01T00:00:00.000Z',
      },
    ]);
    vi.mocked(fetchExerciseEntries).mockResolvedValue([]);
    vi.mocked(createExerciseEntry).mockResolvedValue({
      id: '33333333-3333-3333-3333-333333333333',
      planId,
      patientId,
      performedAt: '2026-01-08T12:00:00.000Z',
      exerciseName: 'Quad sets',
      sets: 3,
      reps: 10,
      resistance: '',
      notes: '',
      recordedBy: 'user-1',
      recordedByName: 'Clinician',
      createdAt: '2026-01-08T12:00:00.000Z',
    });

    renderSection();

    expect(await screen.findByText('Knee strengthening')).toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: /Log session/i }));
    const exerciseInput = await screen.findByPlaceholderText(/Quad sets/i);
    await user.type(exerciseInput, 'Quad sets');
    await user.click(screen.getByRole('button', { name: /Save session/i }));

    await waitFor(() => {
      expect(createExerciseEntry).toHaveBeenCalledWith(
        patientId,
        planId,
        expect.objectContaining({
          exerciseName: 'Quad sets',
          sets: 3,
          reps: 10,
        }),
      );
    });
  });

  it('deletes a plan after confirmation', async () => {
    const user = userEvent.setup();
    const plan = {
      id: planId,
      patientId,
      name: 'Knee strengthening',
      description: '',
      status: 'active' as const,
      createdBy: 'user-1',
      createdByName: 'Clinician',
      createdAt: '2026-01-01T00:00:00.000Z',
    };
    vi.mocked(fetchExercisePlans).mockResolvedValue([plan]);
    vi.mocked(fetchExerciseEntries).mockResolvedValue([]);
    vi.mocked(deleteExercisePlan).mockResolvedValue(undefined);

    renderSection();

    expect(await screen.findByText('Knee strengthening')).toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: /Delete plan/i }));
    const dialog = await screen.findByRole('dialog', { name: /Delete exercise plan/i });
    await user.click(within(dialog).getByRole('button', { name: /^Delete plan$/i }));

    await waitFor(() => {
      expect(deleteExercisePlan).toHaveBeenCalledWith(patientId, planId);
    });
  });

  it('surfaces API errors when creating a plan fails', async () => {
    const user = userEvent.setup();
    vi.mocked(createExercisePlan).mockRejectedValue(new Error('Plan name is required'));

    renderSection();
    expect(await screen.findByText(/No exercise plans yet/i)).toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: /New plan/i }));
    const nameInput = await screen.findByPlaceholderText(/Knee strengthening/i);
    await user.type(nameInput, 'Broken plan');
    await user.click(screen.getByRole('button', { name: /Create plan/i }));

    await waitFor(() => {
      expect(notifications.show).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Could not create plan',
          message: 'Plan name is required',
          color: 'red',
        }),
      );
    });
  });

  it('updates plan status from the select control', async () => {
    const user = userEvent.setup();
    const plan = {
      id: planId,
      patientId,
      name: 'Knee strengthening',
      description: '',
      status: 'active' as const,
      createdBy: 'user-1',
      createdByName: 'Clinician',
      createdAt: '2026-01-01T00:00:00.000Z',
    };
    vi.mocked(fetchExercisePlans).mockResolvedValue([plan]);
    vi.mocked(fetchExerciseEntries).mockResolvedValue([]);
    vi.mocked(updateExercisePlan).mockResolvedValue({ ...plan, status: 'completed' });

    renderSection();
    expect(await screen.findByText('Knee strengthening')).toBeInTheDocument();

    const statusInputs = screen.getAllByLabelText(/Plan status/i);
    await user.click(statusInputs[0]);
    const completedOption = await screen.findByRole('option', { name: 'Completed', hidden: true });
    await user.click(completedOption);

    await waitFor(() => {
      expect(updateExercisePlan).toHaveBeenCalledWith(patientId, planId, {
        name: 'Knee strengthening',
        description: '',
        status: 'completed',
      });
    });
  });

  it('shows progression hint for improving reps', async () => {
    vi.mocked(fetchExercisePlans).mockResolvedValue([
      {
        id: planId,
        patientId,
        name: 'Knee strengthening',
        description: '',
        status: 'active',
        createdBy: 'user-1',
        createdByName: 'Clinician',
        createdAt: '2026-01-01T00:00:00.000Z',
      },
    ]);
    vi.mocked(fetchExerciseEntries).mockResolvedValue([
      {
        id: 'e2',
        planId,
        patientId,
        performedAt: '2026-01-08T12:00:00.000Z',
        exerciseName: 'Quad sets',
        sets: 3,
        reps: 12,
        resistance: '',
        notes: '',
        recordedBy: 'user-1',
        recordedByName: 'Clinician',
        createdAt: '2026-01-08T12:00:00.000Z',
      },
      {
        id: 'e1',
        planId,
        patientId,
        performedAt: '2026-01-01T12:00:00.000Z',
        exerciseName: 'Quad sets',
        sets: 3,
        reps: 10,
        resistance: '',
        notes: '',
        recordedBy: 'user-1',
        recordedByName: 'Clinician',
        createdAt: '2026-01-01T12:00:00.000Z',
      },
    ]);

    renderSection();

    expect(await screen.findByText('+2 reps vs last time (10)')).toBeInTheDocument();
  });
});
