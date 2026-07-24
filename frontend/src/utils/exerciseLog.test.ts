import { describe, expect, it } from 'vitest';
import { formatSessionSummary, progressionHint } from './exerciseLog';

describe('formatSessionSummary', () => {
  it('formats sets and reps', () => {
    expect(formatSessionSummary({ sets: 3, reps: 10, resistance: '', difficulty: null })).toBe('3×10');
  });

  it('includes resistance and difficulty when present', () => {
    expect(
      formatSessionSummary({ sets: 2, reps: 8, resistance: 'yellow band', difficulty: 4 }),
    ).toBe('2×8 · yellow band · RPE 4');
  });
});

describe('progressionHint', () => {
  const older = {
    id: '1',
    exerciseName: 'Quad sets',
    performedAt: '2026-01-01T10:00:00.000Z',
    reps: 10,
  };
  const newer = {
    id: '2',
    exerciseName: 'Quad sets',
    performedAt: '2026-01-08T10:00:00.000Z',
    reps: 12,
  };

  it('returns null when there is no prior session', () => {
    expect(progressionHint([newer], newer)).toBeNull();
  });

  it('reports improvement vs last time', () => {
    expect(progressionHint([newer, older], newer)).toBe('+2 reps vs last time (10)');
  });

  it('reports decrease and same reps', () => {
    const down = { ...newer, id: '3', reps: 8, performedAt: '2026-01-15T10:00:00.000Z' };
    expect(progressionHint([down, newer, older], down)).toBe('-4 reps vs last time (12)');

    const same = { ...newer, id: '4', reps: 12, performedAt: '2026-01-22T10:00:00.000Z' };
    expect(progressionHint([same, newer, older], same)).toBe('Same reps as last time (12)');
  });

  it('matches exercise names case-insensitively', () => {
    const prior = { ...older, exerciseName: 'quad sets' };
    const current = { ...newer, exerciseName: 'QUAD SETS' };
    expect(progressionHint([current, prior], current)).toBe('+2 reps vs last time (10)');
  });

  it('ignores later sessions when computing prior', () => {
    const future = {
      id: '9',
      exerciseName: 'Quad sets',
      performedAt: '2026-12-01T10:00:00.000Z',
      reps: 20,
    };
    expect(progressionHint([future, newer, older], newer)).toBe('+2 reps vs last time (10)');
  });

  it('ignores different exercises', () => {
    const other = {
      id: '9',
      exerciseName: 'Hamstring curl',
      performedAt: '2026-01-01T10:00:00.000Z',
      reps: 10,
    };
    expect(progressionHint([newer, other], newer)).toBeNull();
  });
});
