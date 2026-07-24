const DRILL_KEY = 'poco-last-restore-drill';
export const RESTORE_DRILL_INTERVAL_DAYS = 90;

export function loadLastRestoreDrill(): Date | null {
  try {
    const raw = localStorage.getItem(DRILL_KEY);
    if (!raw) {
      return null;
    }
    const parsed = new Date(raw);
    return Number.isNaN(parsed.getTime()) ? null : parsed;
  } catch {
    return null;
  }
}

export function saveLastRestoreDrill(date = new Date()): void {
  localStorage.setItem(DRILL_KEY, date.toISOString());
}

export function daysSinceLastDrill(last: Date | null, now = new Date()): number | null {
  if (!last) {
    return null;
  }
  const ms = now.getTime() - last.getTime();
  return Math.floor(ms / (1000 * 60 * 60 * 24));
}

export function restoreDrillDue(last: Date | null, now = new Date()): boolean {
  const days = daysSinceLastDrill(last, now);
  if (days === null) {
    return true;
  }
  return days >= RESTORE_DRILL_INTERVAL_DAYS;
}
