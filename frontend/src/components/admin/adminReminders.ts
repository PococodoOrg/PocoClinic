const MONTHLY_TESTING_KEY = 'poco-monthly-testing';
const COMPLIANCE_CHECK_KEY = 'poco-compliance-check';
const HEALTH_CHECK_KEY = 'poco-health-check';
const AUDIT_ARCHIVE_KEY = 'poco-audit-archive';

export function currentMonthKey(now = new Date()): string {
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
}

function loadMonthKey(key: string): string | null {
  try {
    return localStorage.getItem(key);
  } catch {
    return null;
  }
}

function saveMonthKey(key: string, now = new Date()): void {
  localStorage.setItem(key, currentMonthKey(now));
}

export function loadLastMonthlyTesting(): string | null {
  return loadMonthKey(MONTHLY_TESTING_KEY);
}

export function saveMonthlyTestingComplete(now = new Date()): void {
  saveMonthKey(MONTHLY_TESTING_KEY, now);
}

export function monthlyTestingDue(now = new Date()): boolean {
  return loadLastMonthlyTesting() !== currentMonthKey(now);
}

export function loadLastComplianceCheck(): string | null {
  return loadMonthKey(COMPLIANCE_CHECK_KEY);
}

export function saveComplianceCheckComplete(now = new Date()): void {
  saveMonthKey(COMPLIANCE_CHECK_KEY, now);
}

export function complianceCheckDue(now = new Date()): boolean {
  return loadLastComplianceCheck() !== currentMonthKey(now);
}

export function loadLastHealthCheck(): string | null {
  return loadMonthKey(HEALTH_CHECK_KEY);
}

export function saveHealthCheckComplete(now = new Date()): void {
  saveMonthKey(HEALTH_CHECK_KEY, now);
}

export function healthCheckDue(now = new Date()): boolean {
  return loadLastHealthCheck() !== currentMonthKey(now);
}

export function loadLastAuditArchive(): string | null {
  return loadMonthKey(AUDIT_ARCHIVE_KEY);
}

export function saveAuditArchiveComplete(now = new Date()): void {
  saveMonthKey(AUDIT_ARCHIVE_KEY, now);
}

export function auditArchiveDue(now = new Date()): boolean {
  return loadLastAuditArchive() !== currentMonthKey(now);
}
