export interface AdminGuideTask {
  id: string;
  title: string;
  description: string;
  articleId?: string;
  route?: string;
  setupStep?: boolean;
}

export interface AdminGuideSection {
  id: string;
  title: string;
  description: string;
  tasks: AdminGuideTask[];
}

export const adminGuideSections: AdminGuideSection[] = [
  {
    id: 'initial-setup',
    title: 'First-time clinic setup',
    description: 'Do these once when PocoClinic is installed on the clinic server.',
    tasks: [
      {
        id: 'setup-database',
        title: 'Configure persistent database',
        description: 'Set DATABASE_URL to the SQLite file path on the server and run migrations so patient data survives restarts.',
        articleId: 'admin-setup',
        setupStep: true,
      },
      {
        id: 'setup-admin-login',
        title: 'Sign in as bootstrap administrator',
        description: 'Use /login/admin with the admin email, key, and PIN to access staff management and system settings.',
        articleId: 'admin-bootstrap-login',
        route: '/login/admin',
        setupStep: true,
      },
      {
        id: 'setup-safe-vault',
        title: 'Print safe vault credentials sheet',
        description:
          'Print the blank vault sheet, handwrite secrets (admin key, JWT, document encryption key), lock the binder in a safe, then delete bootstrap-admin-once.txt and any digital copies.',
        articleId: 'safe-vault-credentials',
        route: '/help/safe-vault-credentials',
        setupStep: true,
      },
      {
        id: 'setup-ops-binder',
        title: 'Assemble the printed ops binder',
        description:
          'Print sections A (emergency), B (weekly/monthly), C (system help). Store closed; secrets stay in the safe. No daily paper logs.',
        articleId: 'ops-binder',
        route: '/help/ops-binder',
        setupStep: true,
      },
      {
        id: 'setup-staff',
        title: 'Create staff accounts',
        description: 'Add each employee with name, email, and initial PIN. They will choose a private PIN on first sign-in.',
        articleId: 'staff-management',
        route: '/users',
        setupStep: true,
      },
      {
        id: 'setup-badges',
        title: 'Print employee badges',
        description: 'Open each staff member and print their QR badge. Hand badges in person — reissue immediately if lost.',
        articleId: 'staff-management',
        route: '/users',
        setupStep: true,
      },
      {
        id: 'setup-contacts',
        title: 'Fill clinic emergency contacts',
        description: 'Record who to call when the server or network needs help. Stored locally in the browser.',
        articleId: 'clinic-contacts',
        route: '/help/clinic-contacts',
        setupStep: true,
      },
      {
        id: 'setup-usb-labels',
        title: 'Label USB backup drives',
        description: 'Configure Mon/Wed/Fri rotation labels and print stickers for your ops binder.',
        articleId: 'daily-backup',
        setupStep: true,
      },
      {
        id: 'setup-first-backup',
        title: 'Run and verify first backup',
        description: 'Create a backup from the dashboard, click Verify, and copy the file to a labeled USB drive.',
        articleId: 'backup-verification',
        setupStep: true,
      },
      {
        id: 'setup-health-check',
        title: 'Run system health check',
        description: 'Confirm database, document storage, backup directory, and document integrity are all green.',
        articleId: 'system-health-check',
        setupStep: true,
      },
      {
        id: 'setup-compliance',
        title: 'Run compliance checklist',
        description: 'Verify backups, default PINs, audit logging, and storage controls before go-live.',
        articleId: 'compliance-checklist',
        setupStep: true,
      },
    ],
  },
  {
    id: 'daily-ops',
    title: 'Daily operations',
    description: 'Quick tasks administrators handle during a normal clinic week.',
    tasks: [
      {
        id: 'daily-dashboard',
        title: 'Review admin dashboard',
        description: 'Check Action needed alerts, backup age, and resource counts each morning.',
        articleId: 'admin-dashboard',
        route: '/admin',
      },
      {
        id: 'daily-backup',
        title: 'USB backup rotation',
        description: 'Insert today\'s labeled drive, run Backup now, verify success, and store the drive securely.',
        articleId: 'daily-backup',
      },
      {
        id: 'daily-locked-accounts',
        title: 'Help locked-out staff',
        description: 'Accounts lock for 15 minutes after 5 failed PIN attempts. Review Staff if someone needs access sooner.',
        articleId: 'staff-management',
        route: '/users',
      },
    ],
  },
  {
    id: 'staff-security',
    title: 'Staff & security',
    description: 'Manage accounts, badges, and review who did what.',
    tasks: [
      {
        id: 'staff-create',
        title: 'Onboard new employees',
        description: 'Create account → print badge → staff changes default PIN on first sign-in.',
        articleId: 'staff-management',
        route: '/users/new',
      },
      {
        id: 'staff-reissue',
        title: 'Reissue lost badges',
        description: 'Reissue from the staff detail page. The old QR code stops working immediately.',
        articleId: 'staff-management',
        route: '/users',
      },
      {
        id: 'audit-review',
        title: 'Review audit log',
        description: 'Filter sign-ins, patient access, notes, backups, and admin actions when investigating a question.',
        articleId: 'audit-log',
        route: '/admin/audit',
      },
      {
        id: 'audit-export',
        title: 'Archive audit log CSV',
        description: 'Download the last 30 days (admin only), save offline, and mark archived before monthly purge.',
        articleId: 'audit-log',
        route: '/admin/audit',
      },
      {
        id: 'security-overview',
        title: 'Monitor security overview',
        description: 'Watch active sessions, locked accounts, default PIN accounts, and document storage health.',
        articleId: 'security-overview',
        route: '/admin',
      },
    ],
  },
  {
    id: 'forms-reports',
    title: 'Forms & reports',
    description: 'Build templates staff use on patient charts and export submission data.',
    tasks: [
      {
        id: 'forms-templates',
        title: 'Create form templates',
        description: 'Build the forms staff fill out on patient charts — intake, consent, vitals, etc.',
        articleId: 'form-templates',
        route: '/forms',
      },
      {
        id: 'forms-reports',
        title: 'Export form submissions',
        description: 'Filter by date range and download CSV for reporting or quality review.',
        articleId: 'form-reports',
        route: '/forms/reports',
      },
      {
        id: 'census',
        title: 'Patient census',
        description: 'View total patients, gender breakdown, and registrations in the last 30 days on the admin dashboard.',
        articleId: 'admin-dashboard',
        route: '/admin',
      },
      {
        id: 'staff-activity',
        title: 'Review staff activity',
        description: 'See sign-ins, failed attempts, and chart views per staff member on the Reports tab.',
        articleId: 'admin-dashboard',
        route: '/admin',
      },
    ],
  },
  {
    id: 'backup-recovery',
    title: 'Backup & recovery',
    description: 'Protect clinic data and recover after hardware failure.',
    tasks: [
      {
        id: 'backup-helper',
        title: 'Use Backup Helper on server',
        description: 'Guided backup and restore wizards at http://127.0.0.1:9090 — localhost only.',
        articleId: 'backup-helper',
      },
      {
        id: 'backup-verify',
        title: 'Verify backup integrity',
        description: 'Checksums plus document metadata checks before relying on a USB copy.',
        articleId: 'backup-verification',
      },
      {
        id: 'restore-disaster',
        title: 'Restore after disaster',
        description: 'Replace all database and document files from a backup — sign out all staff first.',
        articleId: 'restore-disaster',
      },
      {
        id: 'quarterly-drill',
        title: 'Quarterly restore drill',
        description: 'Test a USB backup on a spare machine when the clinic is closed. Mark complete on the dashboard.',
        articleId: 'quarterly-drill',
      },
    ],
  },
];

export function setupChecklistTasks(): AdminGuideTask[] {
  return adminGuideSections.flatMap((section) => section.tasks.filter((task) => task.setupStep));
}
