import {
  CreateExercisePlanInput,
  ExerciseLogEntry,
  ExerciseLogEntryInput,
  ExercisePlan,
  UpdateExercisePlanInput,
} from '../types/exerciseLog';
import { createApiClient } from './client';

const api = createApiClient();

export async function fetchExercisePlans(patientId: string): Promise<ExercisePlan[]> {
  const response = await api.get<ExercisePlan[]>(`/patients/${patientId}/exercise-plans`);
  return response.data;
}

export async function createExercisePlan(
  patientId: string,
  input: CreateExercisePlanInput,
): Promise<ExercisePlan> {
  const response = await api.post<ExercisePlan>(`/patients/${patientId}/exercise-plans`, input);
  return response.data;
}

export async function updateExercisePlan(
  patientId: string,
  planId: string,
  input: UpdateExercisePlanInput,
): Promise<ExercisePlan> {
  const response = await api.put<ExercisePlan>(
    `/patients/${patientId}/exercise-plans/${planId}`,
    input,
  );
  return response.data;
}

export async function deleteExercisePlan(patientId: string, planId: string): Promise<void> {
  await api.delete(`/patients/${patientId}/exercise-plans/${planId}`);
}

export async function fetchExerciseEntries(
  patientId: string,
  planId: string,
): Promise<ExerciseLogEntry[]> {
  const response = await api.get<ExerciseLogEntry[]>(
    `/patients/${patientId}/exercise-plans/${planId}/entries`,
  );
  return response.data;
}

export async function createExerciseEntry(
  patientId: string,
  planId: string,
  input: ExerciseLogEntryInput,
): Promise<ExerciseLogEntry> {
  const response = await api.post<ExerciseLogEntry>(
    `/patients/${patientId}/exercise-plans/${planId}/entries`,
    input,
  );
  return response.data;
}

export async function updateExerciseEntry(
  patientId: string,
  planId: string,
  entryId: string,
  input: ExerciseLogEntryInput,
): Promise<ExerciseLogEntry> {
  const response = await api.put<ExerciseLogEntry>(
    `/patients/${patientId}/exercise-plans/${planId}/entries/${entryId}`,
    input,
  );
  return response.data;
}

export async function deleteExerciseEntry(
  patientId: string,
  planId: string,
  entryId: string,
): Promise<void> {
  await api.delete(`/patients/${patientId}/exercise-plans/${planId}/entries/${entryId}`);
}
