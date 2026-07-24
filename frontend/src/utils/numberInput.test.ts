import { describe, expect, it } from 'vitest';
import { parseNumberInput } from './numberInput';

describe('parseNumberInput', () => {
  it('keeps finite numbers', () => {
    expect(parseNumberInput(12)).toBe(12);
    expect(parseNumberInput('8')).toBe(8);
    expect(parseNumberInput(0)).toBe(0);
  });

  it('allows empty while typing', () => {
    expect(parseNumberInput('')).toBe('');
    expect(parseNumberInput(null)).toBe('');
    expect(parseNumberInput(undefined)).toBe('');
  });

  it('rejects NaN-producing intermediate values', () => {
    expect(parseNumberInput('-')).toBe('');
    expect(parseNumberInput('.')).toBe('');
    expect(parseNumberInput('abc')).toBe('');
    expect(parseNumberInput(Number.NaN)).toBe('');
    expect(parseNumberInput(Number.POSITIVE_INFINITY)).toBe('');
  });

  it('parses decimal strings when finite', () => {
    expect(parseNumberInput('3.5')).toBe(3.5);
  });
});
