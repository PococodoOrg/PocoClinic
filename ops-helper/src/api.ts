export interface HelperStatus {
  databaseConnected: boolean;
  storageMode: 'database' | 'memory';
  lastBackupAt?: string;
  backupAgeHours?: number;
  backupStatus: 'ok' | 'warning' | 'critical' | 'unknown';
  lastBackupFile?: string;
  backupCount: number;
  patientCount: number;
  appVersion: string;
  backupDir: string;
  mainAppUrl: string;
}

export interface BackupEntry {
  filename: string;
  createdAt: string;
  appVersion?: string;
  friendlyAt: string;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    headers: { 'Content-Type': 'application/json' },
    ...init,
  });
  const body = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(body.error ?? 'Something went wrong. Please try again.');
  }
  return body as T;
}

export function fetchStatus(): Promise<HelperStatus> {
  return request<HelperStatus>('/api/status');
}

export function fetchBackups(): Promise<{ backups: BackupEntry[] }> {
  return request<{ backups: BackupEntry[] }>('/api/backups');
}

export function createBackup(): Promise<{ filename: string; message: string }> {
  return request('/api/backups', { method: 'POST' });
}

export function restoreBackup(filename: string, confirmText: string): Promise<{ message: string }> {
  return request('/api/restore', {
    method: 'POST',
    body: JSON.stringify({ filename, confirmText }),
  });
}

export function statusHeadline(status: HelperStatus): { title: string; detail: string; color: string } {
  if (status.storageMode !== 'database') {
    return {
      title: 'Database not connected',
      detail: 'Set DATABASE_URL on this computer, then restart the backup helper.',
      color: 'red',
    };
  }
  switch (status.backupStatus) {
    case 'ok':
      return {
        title: "You're protected",
        detail: 'Your latest backup is less than 24 hours old. Nice work!',
        color: 'teal',
      };
    case 'warning':
      return {
        title: 'Time for a backup',
        detail: 'Your last backup is more than a day old. A quick backup now keeps everyone safe.',
        color: 'yellow',
      };
    case 'critical':
      return {
        title: 'Backup needed today',
        detail: 'It has been more than two days since the last backup. Please back up before clinic closes.',
        color: 'red',
      };
    default:
      return {
        title: 'No backup yet',
        detail: 'Create your first backup — it only takes a minute.',
        color: 'gray',
      };
  }
}
