/** Returns a basename safe for browser download attributes. */
export function safeDownloadFilename(name: string): string {
  const base = name.replace(/\\/g, '/').split('/').pop()?.trim() ?? '';
  // Strip control chars and path/download metacharacters.
  // eslint-disable-next-line no-control-regex -- intentional sanitization of user-supplied names
  const cleaned = base.replace(/[\x00-\x1f"';\\]/g, '').trim();
  if (!cleaned || cleaned === '.') {
    return 'download';
  }
  return cleaned;
}
