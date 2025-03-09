import { formatDate } from '../../utils/date';

describe('Date Utils', () => {
  it('formats date correctly', () => {
    const testDate = new Date(2024, 2, 8); // March 8, 2024
    expect(formatDate(testDate)).toBe('March 8, 2024');
  });

  it('handles different months', () => {
    const testDate = new Date(2024, 11, 25); // December 25, 2024
    expect(formatDate(testDate)).toBe('December 25, 2024');
  });
}); 