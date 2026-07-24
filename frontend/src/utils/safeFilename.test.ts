import { describe, expect, it } from 'vitest';
import { safeDownloadFilename } from './safeFilename';

describe('safeDownloadFilename', () => {
  it('returns basename only', () => {
    expect(safeDownloadFilename('../../etc/passwd')).toBe('passwd');
  });

  it('strips unsafe characters', () => {
    expect(safeDownloadFilename('report";evil.csv')).toBe('reportevil.csv');
  });

  it('falls back for empty names', () => {
    expect(safeDownloadFilename('')).toBe('download');
  });
});
