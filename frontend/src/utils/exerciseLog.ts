import { ExerciseLogEntry } from '../types/exerciseLog';

export function formatSessionSummary(entry: Pick<ExerciseLogEntry, 'sets' | 'reps' | 'resistance' | 'difficulty'>): string {
  const parts = [`${entry.sets}×${entry.reps}`];
  if (entry.resistance) {
    parts.push(entry.resistance);
  }
  if (entry.difficulty != null) {
    parts.push(`RPE ${entry.difficulty}`);
  }
  return parts.join(' · ');
}

/** Compare current session reps to the most recent prior session of the same exercise. */
export function progressionHint(
  entries: Array<Pick<ExerciseLogEntry, 'id' | 'exerciseName' | 'performedAt' | 'reps'>>,
  current: Pick<ExerciseLogEntry, 'id' | 'exerciseName' | 'performedAt' | 'reps'>,
): string | null {
  const prior = entries.find(
    (entry) =>
      entry.id !== current.id &&
      entry.exerciseName.toLowerCase() === current.exerciseName.toLowerCase() &&
      new Date(entry.performedAt) < new Date(current.performedAt),
  );
  if (!prior) {
    return null;
  }
  const delta = current.reps - prior.reps;
  if (delta === 0) {
    return `Same reps as last time (${prior.reps})`;
  }
  if (delta > 0) {
    return `+${delta} reps vs last time (${prior.reps})`;
  }
  return `${delta} reps vs last time (${prior.reps})`;
}
