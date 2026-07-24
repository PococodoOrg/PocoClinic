/** Safe NumberInput onChange helper — never returns NaN (crashes Mantine). */
export function parseNumberInput(value: string | number | null | undefined): number | '' {
  if (value === '' || value == null) {
    return '';
  }
  const parsed = typeof value === 'number' ? value : Number(value);
  return Number.isFinite(parsed) ? parsed : '';
}
