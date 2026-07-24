import { beforeEach, describe, expect, it, vi } from 'vitest';

const { get, post, put, del } = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  del: vi.fn(),
}));

vi.mock('./client', () => ({
  createApiClient: () => ({
    get,
    post,
    put,
    delete: del,
  }),
}));

import {
  createExerciseEntry,
  createExercisePlan,
  deleteExerciseEntry,
  deleteExercisePlan,
  fetchExerciseEntries,
  fetchExercisePlans,
  updateExerciseEntry,
  updateExercisePlan,
} from './exerciseLog';

describe('exerciseLog API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('fetches plans and entries', async () => {
    get.mockResolvedValueOnce({ data: [{ id: 'p1', name: 'Knee' }] });
    get.mockResolvedValueOnce({ data: [{ id: 'e1', reps: 10 }] });

    await expect(fetchExercisePlans('patient-1')).resolves.toEqual([{ id: 'p1', name: 'Knee' }]);
    expect(get).toHaveBeenCalledWith('/patients/patient-1/exercise-plans');

    await expect(fetchExerciseEntries('patient-1', 'plan-1')).resolves.toEqual([{ id: 'e1', reps: 10 }]);
    expect(get).toHaveBeenCalledWith('/patients/patient-1/exercise-plans/plan-1/entries');
  });

  it('creates and updates plans', async () => {
    post.mockResolvedValueOnce({ data: { id: 'p1', name: 'Knee' } });
    put.mockResolvedValueOnce({ data: { id: 'p1', status: 'completed' } });

    await createExercisePlan('patient-1', { name: 'Knee', description: 'goals' });
    expect(post).toHaveBeenCalledWith('/patients/patient-1/exercise-plans', {
      name: 'Knee',
      description: 'goals',
    });

    await updateExercisePlan('patient-1', 'p1', {
      name: 'Knee',
      description: '',
      status: 'completed',
    });
    expect(put).toHaveBeenCalledWith('/patients/patient-1/exercise-plans/p1', {
      name: 'Knee',
      description: '',
      status: 'completed',
    });
  });

  it('creates updates and deletes entries', async () => {
    post.mockResolvedValueOnce({ data: { id: 'e1' } });
    put.mockResolvedValueOnce({ data: { id: 'e1', reps: 12 } });
    del.mockResolvedValueOnce({});

    await createExerciseEntry('patient-1', 'p1', {
      exerciseName: 'Quad sets',
      sets: 3,
      reps: 10,
    });
    expect(post).toHaveBeenCalledWith(
      '/patients/patient-1/exercise-plans/p1/entries',
      expect.objectContaining({ exerciseName: 'Quad sets' }),
    );

    await updateExerciseEntry('patient-1', 'p1', 'e1', {
      exerciseName: 'Quad sets',
      sets: 3,
      reps: 12,
    });
    expect(put).toHaveBeenCalledWith(
      '/patients/patient-1/exercise-plans/p1/entries/e1',
      expect.objectContaining({ reps: 12 }),
    );

    await deleteExerciseEntry('patient-1', 'p1', 'e1');
    expect(del).toHaveBeenCalledWith('/patients/patient-1/exercise-plans/p1/entries/e1');

    await deleteExercisePlan('patient-1', 'p1');
    expect(del).toHaveBeenCalledWith('/patients/patient-1/exercise-plans/p1');
  });
});
