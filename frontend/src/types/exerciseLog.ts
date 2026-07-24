export type ExercisePlanStatus = 'active' | 'completed' | 'archived';

export interface ExercisePlan {
  id: string;
  patientId: string;
  name: string;
  description: string;
  status: ExercisePlanStatus;
  createdBy: string;
  createdByName: string;
  createdAt: string;
  updatedAt?: string;
}

export interface ExerciseLogEntry {
  id: string;
  planId: string;
  patientId: string;
  performedAt: string;
  exerciseName: string;
  sets: number;
  reps: number;
  durationSeconds?: number | null;
  resistance: string;
  difficulty?: number | null;
  notes: string;
  recordedBy: string;
  recordedByName: string;
  createdAt: string;
  updatedAt?: string;
}

export interface CreateExercisePlanInput {
  name: string;
  description?: string;
}

export interface UpdateExercisePlanInput {
  name: string;
  description?: string;
  status: ExercisePlanStatus;
}

export interface ExerciseLogEntryInput {
  exerciseName: string;
  sets: number;
  reps: number;
  durationSeconds?: number | null;
  resistance?: string;
  difficulty?: number | null;
  notes?: string;
  performedAt?: string;
}
