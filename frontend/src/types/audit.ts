export interface AuditEntry {
  id: string;
  eventType: string;
  userId?: string;
  resourceType?: string;
  resourceId?: string;
  ipAddress?: string;
  userAgent?: string;
  details?: Record<string, string>;
  success: boolean;
  createdAt: string;
}

export interface PaginatedAuditEntries {
  events: AuditEntry[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
  totalPages: number;
}

export const AUDIT_EVENT_LABELS: Record<string, string> = {
  'auth.login.success': 'Signed in',
  'auth.login.failure': 'Failed sign in',
  'auth.login.locked': 'Account locked',
  'auth.token.refresh': 'Session refreshed',
  'auth.logout': 'Signed out',
  'auth.badge.reissued': 'Badge reissued',
  'auth.pin.changed': 'PIN changed',
  'user.created': 'Account created',
  'user.updated': 'Profile updated',
  'user.deleted': 'Account deleted',
  'user.unlocked': 'Account unlocked',
  'patient.created': 'Patient created',
  'patient.updated': 'Patient updated',
  'patient.deleted': 'Patient deleted',
  'patient.viewed': 'Patient viewed',
  'patient.note.created': 'Clinical note added',
  'patient.note.updated': 'Clinical note edited',
  'patient.note.deleted': 'Clinical note deleted',
  'patient.document.uploaded': 'Document uploaded',
  'patient.document.viewed': 'Document downloaded',
  'patient.document.deleted': 'Document removed',
  'system.backup.created': 'Backup created',
  'system.backup.restored': 'Database restored',
  'system.healthcheck.run': 'Health check run',
  'system.compliance.run': 'Compliance check run',
  'system.audit.purged': 'Audit log rotation',
};

export function formatAuditEvent(entry: AuditEntry): string {
  const base = AUDIT_EVENT_LABELS[entry.eventType] ?? entry.eventType;
  const details = entry.details ?? {};

  if (entry.eventType === 'auth.login.success' && details.loginMode) {
    return `${base} (${details.loginMode})`;
  }
  if (entry.eventType === 'auth.login.failure' && details.reason) {
    return `${base}: ${details.reason.replace(/_/g, ' ')}`;
  }
  if (entry.eventType === 'auth.badge.reissued' && details.actorId) {
    return `${base} by administrator`;
  }
  if (entry.eventType.startsWith('patient.') && entry.resourceId) {
    return `${base}`;
  }

  return base;
}
