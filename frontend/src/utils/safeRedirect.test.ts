import { describe, expect, it } from 'vitest';
import { safeInternalPath } from './safeRedirect';

describe('safeInternalPath', () => {
  it('accepts internal paths', () => {
    expect(safeInternalPath('/patients')).toBe('/patients');
    expect(safeInternalPath('/admin/audit')).toBe('/admin/audit');
  });

  it('rejects external and protocol-relative paths', () => {
    expect(safeInternalPath('//evil.example')).toBe('/patients');
    expect(safeInternalPath('https://evil.example')).toBe('/patients');
    expect(safeInternalPath('/\\evil')).toBe('/patients');
  });

  it('rejects URL-encoded protocol-relative paths', () => {
    expect(safeInternalPath('/%2f%2fevil.example')).toBe('/patients');
  });

  it('falls back for invalid values', () => {
    expect(safeInternalPath(null)).toBe('/patients');
    expect(safeInternalPath(undefined, '/login')).toBe('/login');
  });
});
