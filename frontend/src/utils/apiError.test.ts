import { describe, expect, it } from 'vitest';
import { ValidationError, getErrorMessage } from './apiError';

describe('getErrorMessage', () => {
  it('reads ValidationError and Error messages', () => {
    expect(getErrorMessage(new ValidationError('Bad field', 'VALIDATION_ERROR'))).toBe('Bad field');
    expect(getErrorMessage(new Error('Network error - no response received'))).toBe(
      'Network error - no response received',
    );
  });

  it('falls back for unknown values', () => {
    expect(getErrorMessage(null, 'Could not save')).toBe('Could not save');
    expect(getErrorMessage(42)).toBe('Something went wrong. Please try again.');
  });
});
