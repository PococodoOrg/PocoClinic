import { API_BASE_URL, createApiClient } from './client';
import { AuditEntry } from '../types/audit';
import { PatientFieldRequirements } from '../types/patient';
import { safeDownloadFilename } from '../utils/safeFilename';

const adminApi = createApiClient();

export interface BackupStatus {
  lastBackupAt?: string;
  backupAgeHours?: number;
  status: 'ok' | 'warning' | 'critical' | 'unknown';
  lastBackupFile?: string;
}

export interface ResourceCounts {
  patients: number;
  users: number;
  formTemplates: number;
  formSubmissions: number;
  activeSessions: number;
  lockedAccounts: number;
  defaultPinAccounts: number;
  patientDocuments: number;
}

export interface SecurityOverview {
  documentsStorageReady: boolean;
}

export interface SystemStatus {
  storageMode: 'database' | 'memory';
  databaseConnected: boolean;
  pendingMigrations: string[];
  latestMigration?: string;
  uptimeSeconds: number;
  environment: string;
  appVersion: string;
  backup: BackupStatus;
  counts: ResourceCounts;
  security: SecurityOverview;
}

export interface GenderCount {
  gender: string;
  count: number;
}

export interface PatientCensus {
  totalPatients: number;
  byGender: GenderCount[];
  addedLast30Days: number;
}

export interface BackupEntry {
  filename: string;
  path: string;
  createdAt: string;
  appVersion?: string;
}

export interface BackupResult {
  filename: string;
  path: string;
  createdAt: string;
}

export const fetchSystemStatus = async (): Promise<SystemStatus> => {
  const response = await adminApi.get<SystemStatus>('/admin/system-status');
  return response.data;
};

export const fetchPatientCensus = async (): Promise<PatientCensus> => {
  const response = await adminApi.get<PatientCensus>('/admin/reports/patient-census');
  return response.data;
};

export const fetchBackups = async (): Promise<BackupEntry[]> => {
  const response = await adminApi.get<{ backups: BackupEntry[] }>('/admin/backups');
  return response.data.backups ?? [];
};

export const createBackup = async (): Promise<BackupResult> => {
  const response = await adminApi.post<BackupResult>('/admin/backups', {}, { timeout: 120_000 });
  return response.data;
};

export interface RestoreBackupRequest {
  filename: string;
  confirm: boolean;
}

export const restoreBackup = async (request: RestoreBackupRequest): Promise<void> => {
  await adminApi.post('/admin/restore', request, { timeout: 120_000 });
};

export interface BackupVerifyResult {
  filename: string;
  valid: boolean;
  checkedAt: string;
  message?: string;
  summary?: {
    tableRows: Record<string, number>;
    documentFiles: number;
  };
  integrityIssues?: string[];
}

export interface HealthCheckItem {
  id: string;
  status: 'ok' | 'warning' | 'critical';
  label: string;
  message: string;
}

export interface HealthCheckResult {
  checkedAt: string;
  status: 'ok' | 'warning' | 'critical';
  checks: HealthCheckItem[];
}

export interface ActivitySummary {
  periodDays: number;
  signIns: number;
  failedSignIns: number;
  patientsCreated: number;
  chartNotesAdded: number;
  patientViews: number;
  documentsUploaded: number;
  backupsCreated: number;
}

export const fetchActivitySummary = async (days = 7): Promise<ActivitySummary> => {
  const response = await adminApi.get<ActivitySummary>('/admin/reports/activity-summary', {
    params: { days },
  });
  return response.data;
};

export interface StaffActivityRow {
  userId: string;
  name: string;
  email: string;
  signIns: number;
  failedSignIns: number;
  patientViews: number;
  lastSignIn?: string;
}

export interface StaffActivityReport {
  periodDays: number;
  staff: StaffActivityRow[];
}

export const fetchStaffActivity = async (days = 30): Promise<StaffActivityReport> => {
  const response = await adminApi.get<StaffActivityReport>('/admin/reports/staff-activity', {
    params: { days },
  });
  return response.data;
};

export interface ComplianceItem {
  id: string;
  status: 'pass' | 'warn' | 'fail';
  label: string;
  message: string;
  control?: string;
}

export interface ComplianceCheckResult {
  checkedAt: string;
  status: 'ok' | 'warn' | 'fail';
  checks: ComplianceItem[];
}

export const fetchComplianceCheck = async (): Promise<ComplianceCheckResult> => {
  const response = await adminApi.get<ComplianceCheckResult>('/admin/compliance-check');
  return response.data;
};

export const buildAuditExportUrl = (params: { days?: number; eventType?: string; userId?: string } = {}): string => {
  const search = new URLSearchParams();
  if (params.days) {
    search.set('days', String(params.days));
  }
  if (params.eventType) {
    search.set('eventType', params.eventType);
  }
  if (params.userId) {
    search.set('userId', params.userId);
  }
  const query = search.toString();
  return `${API_BASE_URL}/admin/audit-logs/export.csv${query ? `?${query}` : ''}`;
};

export interface PerformanceSnapshot {
  collectedAt: string;
  goroutines: number;
  memoryAllocMB: number;
  memorySysMB: number;
  databasePool?: {
    totalConns: number;
    idleConns: number;
    inUseConns: number;
    maxConns: number;
  };
}

export const fetchPerformance = async (): Promise<PerformanceSnapshot> => {
  const response = await adminApi.get<PerformanceSnapshot>('/admin/performance');
  return response.data;
};

export const downloadAuthenticatedCsv = async (url: string, filename: string): Promise<void> => {
  const parsed = new URL(url, window.location.origin);
  const apiRoot = new URL(API_BASE_URL, window.location.origin);
  const apiPath = apiRoot.pathname.replace(/\/$/, '');

  if (parsed.origin !== window.location.origin && parsed.origin !== apiRoot.origin) {
    throw new Error('Invalid export URL');
  }
  if (!parsed.pathname.startsWith(`${apiPath}/admin/`)) {
    throw new Error('Invalid export URL');
  }

  const response = await fetch(url, {
    credentials: 'include',
  });
  if (!response.ok) {
    throw new Error('Export failed');
  }
  const blob = await response.blob();
  const objectUrl = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = objectUrl;
  link.download = safeDownloadFilename(filename);
  link.click();
  URL.revokeObjectURL(objectUrl);
};

export const fetchHealthCheck = async (): Promise<HealthCheckResult> => {
  const response = await adminApi.get<HealthCheckResult>('/admin/health-check');
  return response.data;
};

export const verifyBackup = async (filename: string): Promise<BackupVerifyResult> => {
  const response = await adminApi.post<BackupVerifyResult>('/admin/backups/verify', { filename });
  return response.data;
};

export interface PaginatedAuditEntries {
  events: AuditEntry[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
  totalPages: number;
}

export interface AuditLogsQuery {
  page?: number;
  pageSize?: number;
  eventType?: string;
  userId?: string;
}

export const fetchAuditLogs = async (query: AuditLogsQuery = {}): Promise<PaginatedAuditEntries> => {
  const response = await adminApi.get<PaginatedAuditEntries>('/admin/audit-logs', { params: query });
  return response.data;
};

export const fetchAdminPatientFieldRequirements = async (): Promise<PatientFieldRequirements> => {
  const response = await adminApi.get<{ requirements: PatientFieldRequirements }>('/admin/patient-field-requirements');
  return response.data.requirements;
};

export const updatePatientFieldRequirements = async (
  requirements: PatientFieldRequirements,
): Promise<PatientFieldRequirements> => {
  const response = await adminApi.put<{ requirements: PatientFieldRequirements }>(
    '/admin/patient-field-requirements',
    { requirements },
  );
  return response.data.requirements;
};
