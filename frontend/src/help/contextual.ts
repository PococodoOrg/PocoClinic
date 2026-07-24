export function getContextualArticleIds(pathname: string): string[] {
  if (pathname === '/account/pin') {
    return ['change-pin', 'sign-in'];
  }
  if (pathname === '/patients/new') {
    return ['register-patient', 'patient-list'];
  }
  if (/^\/patients\/[^/]+\/edit$/.test(pathname)) {
    return ['register-patient', 'patient-chart'];
  }
  if (/^\/patients\/[^/]+$/.test(pathname)) {
    return ['patient-chart', 'clinical-notes', 'patient-documents', 'patient-exercise', 'patient-forms'];
  }
  if (pathname === '/patients') {
    return ['patient-list', 'register-patient'];
  }
  if (pathname === '/admin') {
    return ['admin-dashboard', 'admin-setup', 'daily-backup', 'system-health-check', 'troubleshooting-common'];
  }
  if (pathname === '/admin/audit') {
    return ['audit-log'];
  }
  if (pathname.startsWith('/users')) {
    return ['staff-management', 'sign-in'];
  }
  if (pathname.startsWith('/forms/reports')) {
    return ['form-reports', 'form-templates'];
  }
  if (pathname.startsWith('/forms')) {
    return ['form-templates', 'patient-forms'];
  }
  if (pathname.startsWith('/help')) {
    return ['how-it-works'];
  }
  return ['how-it-works', 'sign-in', 'patient-list'];
}
