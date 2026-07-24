export interface ClinicContact {
  role: string;
  name: string;
  phone: string;
}

export const defaultClinicContacts: ClinicContact[] = [
  { role: 'Primary administrator', name: '', phone: '' },
  { role: 'Backup administrator', name: '', phone: '' },
  { role: 'IT / vendor support', name: '', phone: '' },
];

const STORAGE_KEY = 'poco-clinic-contacts';

export function loadClinicContacts(): ClinicContact[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return defaultClinicContacts;
    }
    const parsed = JSON.parse(raw) as ClinicContact[];
    if (!Array.isArray(parsed) || parsed.length === 0) {
      return defaultClinicContacts;
    }
    return parsed;
  } catch {
    return defaultClinicContacts;
  }
}

export function saveClinicContacts(contacts: ClinicContact[]): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(contacts));
}
