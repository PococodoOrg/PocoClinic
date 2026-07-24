const STORAGE_KEY = 'poco-admin-setup-checklist';

export function loadSetupChecklist(): Record<string, boolean> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return {};
    }
    const parsed = JSON.parse(raw) as Record<string, boolean>;
    return typeof parsed === 'object' && parsed !== null ? parsed : {};
  } catch {
    return {};
  }
}

export function saveSetupChecklist(state: Record<string, boolean>): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
}

export function toggleSetupStep(stepId: string, completed: boolean): Record<string, boolean> {
  const next = { ...loadSetupChecklist(), [stepId]: completed };
  saveSetupChecklist(next);
  return next;
}
