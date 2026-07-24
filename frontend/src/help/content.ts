import { HelpArticle, HelpCategory } from './types';

export const helpCategories: HelpCategory[] = [
  {
    id: 'getting-started',
    label: 'Getting started',
    description: 'How PocoClinic works at your clinic on the private LAN.',
  },
  {
    id: 'daily-care',
    label: 'Daily patient care',
    description: 'Patients, chart notes, and forms during a normal shift.',
  },
  {
    id: 'account-security',
    label: 'Sign-in & security',
    description: 'Badges, PINs, sessions, and keeping accounts safe.',
  },
  {
    id: 'administration',
    label: 'Administration',
    description: 'Staff, forms, reports, audit log, and system health.',
  },
  {
    id: 'backup-restore',
    label: 'Backup & restore',
    description: 'USB backups, the local helper, and disaster recovery.',
  },
  {
    id: 'troubleshooting',
    label: 'Troubleshooting',
    description: 'Common problems and what to try first.',
  },
];

export const helpArticles: HelpArticle[] = [
  {
    id: 'how-it-works',
    title: 'How PocoClinic works at your clinic',
    summary: 'Private LAN, clinic workstations, no internet required for daily care.',
    category: 'getting-started',
    audience: 'all',
    tags: ['network', 'lan', 'offline', 'overview'],
    sections: [
      {
        paragraphs: [
          'PocoClinic runs on a clinic-owned server (Raspberry Pi or small PC) on your private local network. Staff use clinic desktop or laptop browsers at the front desk and in exam rooms.',
          'Daily care — sign-in, patient charts, forms, and USB backups — works without internet. Security (badge + PIN, audit log, session timeouts) still applies on the LAN.',
        ],
      },
      {
        heading: 'What you need at the clinic',
        list: [
          'One server running PocoClinic and the database',
          'Clinic-owned PCs on the same private network',
          'Staff badge QR codes and personal 4-digit PINs',
          'Labeled USB drives for backup rotation (admin)',
        ],
      },
      {
        alert: {
          color: 'blue',
          title: 'Not a mobile or cloud system',
          body: 'There are no phone apps and no cloud hosting. Do not expose the server to the public internet.',
        },
      },
    ],
  },
  {
    id: 'sign-in',
    title: 'Sign in with badge and PIN',
    summary: 'Scan your employee badge, enter your PIN, and start your session.',
    category: 'account-security',
    audience: 'all',
    tags: ['login', 'badge', 'pin', 'qr'],
    sections: [
      {
        steps: [
          'Open the clinic browser and go to the PocoClinic sign-in page.',
          'Scan your employee badge QR code (USB scanner or paste the key if prompted).',
          'Enter your 4-digit PIN and press Sign in.',
        ],
      },
      {
        heading: 'First-time staff',
        paragraphs: [
          'Your administrator creates your account and prints your badge. If you still have the default PIN, you will be prompted to choose a private PIN before using patient records.',
        ],
      },
      {
        heading: 'Account lockout',
        list: [
          'After 5 failed PIN attempts, your account locks for 15 minutes.',
          'Wait for the lockout to expire, or ask an administrator for help from Staff management.',
        ],
      },
      {
        heading: 'Session timeout',
        paragraphs: [
          'For privacy, your session ends after 15 minutes of inactivity. Sign in again with badge + PIN when prompted.',
        ],
      },
    ],
  },
  {
    id: 'change-pin',
    title: 'Change your PIN',
    summary: 'Set or update your personal 4-digit PIN from the sidebar.',
    category: 'account-security',
    audience: 'all',
    tags: ['pin', 'password', 'account'],
    sections: [
      {
        steps: [
          'Open Change PIN in the sidebar under Account.',
          'Enter your current PIN, then your new 4-digit PIN twice.',
          'Save. Use the new PIN at your next sign-in.',
        ],
      },
      {
        alert: {
          color: 'yellow',
          title: 'Keep your PIN private',
          body: 'Do not share your PIN or write it on your badge. Treat it like an ATM PIN.',
        },
      },
    ],
  },
  {
    id: 'keyboard-and-accessibility',
    title: 'Keyboard and screen reader tips',
    summary: 'Navigate PocoClinic from the front desk without a mouse when needed.',
    category: 'getting-started',
    audience: 'all',
    tags: ['keyboard', 'accessibility', 'a11y'],
    sections: [
      {
        list: [
          'Press Tab to move between links, buttons, and form fields.',
          'Use Skip to main content (first Tab on a page) to jump past the sidebar.',
          'Press Enter or Space to activate buttons and links.',
          'The ? button in the header opens contextual help for the current page.',
        ],
      },
      {
        paragraphs: [
          'PocoClinic targets clinic desktop browsers at 1024px and wider — not mobile phones. Ask your administrator if you need display adjustments at the front desk.',
        ],
      },
    ],
  },
  {
    id: 'patient-list',
    title: 'Find and open patients',
    summary: 'Search the patient list and open a chart from the front desk.',
    category: 'daily-care',
    audience: 'all',
    tags: ['patients', 'search', 'list'],
    sections: [
      {
        steps: [
          'Open Patients in the sidebar.',
          'Type a name in the search box to filter the list.',
          'Click a patient row to open their chart.',
        ],
      },
      {
        heading: 'Register a new patient',
        paragraphs: ['Use the New patient button on the patient list. Fill in demographics and save.'],
      },
    ],
  },
  {
    id: 'register-patient',
    title: 'Register a new patient',
    summary: 'Create a patient record with demographics and contact information.',
    category: 'daily-care',
    audience: 'all',
    tags: ['patients', 'create', 'registration'],
    sections: [
      {
        steps: [
          'From Patients, click New patient.',
          'Enter first name, last name, date of birth, and gender (required).',
          'Add contact details, address, height, and weight as needed.',
          'Click Save to create the record.',
        ],
      },
      {
        paragraphs: [
          'Creating and viewing patients is recorded in the audit log for compliance.',
        ],
      },
    ],
  },
  {
    id: 'patient-chart',
    title: 'Patient chart overview',
    summary: 'Demographics, notes, documents, forms, and exercise log on the patient detail page.',
    category: 'daily-care',
    audience: 'all',
    tags: ['patients', 'chart', 'detail'],
    sections: [
      {
        paragraphs: [
          'The patient chart shows demographics, clinical notes, documents, forms, and the exercise log (when used).',
        ],
      },
      {
        heading: 'Edit or remove',
        list: [
          'Edit Patient updates demographics.',
          'Delete removes the entire record permanently — use only when sure.',
        ],
      },
    ],
  },
  {
    id: 'clinical-notes',
    title: 'Clinical chart notes',
    summary: 'Add timestamped notes visible to all staff; every note is audit-logged.',
    category: 'daily-care',
    audience: 'all',
    tags: ['notes', 'chart', 'documentation'],
    sections: [
      {
        steps: [
          'Open the patient chart.',
          'Scroll to Clinical notes.',
          'Type your note and click Add note.',
        ],
      },
      {
        paragraphs: [
          'Notes show author name and time in a timeline. You can edit or delete your own notes using the pencil and trash icons.',
          'Administrators can delete any staff member\'s note. Staff cannot edit another person\'s note.',
          'Create, edit, and delete actions are all recorded in the audit log even after a note is removed from the chart.',
        ],
      },
    ],
  },
  {
    id: 'patient-documents',
    title: 'Patient documents',
    summary: 'Upload referrals, consent PDFs, and scans to the patient chart on the clinic server.',
    category: 'daily-care',
    audience: 'all',
    tags: ['documents', 'upload', 'pdf', 'files'],
    sections: [
      {
        steps: [
          'Open the patient chart and scroll to Documents.',
          'Click Upload file and choose PDF, PNG, JPEG, GIF, or plain text (max 10 MB).',
          'Download or remove files as needed. All uploads are audit-logged.',
        ],
      },
      {
        alert: {
          color: 'blue',
          title: 'Stored on the clinic server',
          body: 'New uploads are encrypted (AES-256-GCM) and stored in the clinic database. They are included in USB backups with the rest of the chart. Older installs may still have a few files under DOCUMENTS_DIR on disk — those are included under documents/ in the backup when present.',
        },
      },
    ],
  },
  {
    id: 'patient-exercise',
    title: 'Exercise log (physical therapy)',
    summary: 'Create exercise plans and log sets, reps, resistance, and difficulty on the patient chart.',
    category: 'daily-care',
    audience: 'all',
    tags: ['exercise', 'PT', 'therapy', 'plans'],
    sections: [
      {
        steps: [
          'Open the patient chart and scroll to Exercise log.',
          'Create or select a plan for the patient.',
          'Add session entries (exercise name, sets, reps, optional resistance and difficulty).',
        ],
      },
      {
        paragraphs: [
          'Plans and log entries stay on the clinic server and are included in USB database backups. Create, edit, and delete actions are audit-logged.',
        ],
      },
    ],
  },
  {
    id: 'patient-forms',
    title: 'Patient forms',
    summary: 'Fill in form templates attached to a patient chart.',
    category: 'daily-care',
    audience: 'all',
    tags: ['forms', 'templates', 'submissions'],
    sections: [
      {
        steps: [
          'Open the patient chart and scroll to Forms.',
          'Choose a form template from the list.',
          'Complete the fields and save.',
        ],
      },
      {
        paragraphs: [
          'Administrators create and manage form templates under Forms in the sidebar. Saved submissions stay on the patient chart.',
        ],
      },
    ],
  },
  {
    id: 'admin-setup',
    title: 'First-time clinic setup',
    summary: 'Configure the server, database, staff, badges, and first backup.',
    category: 'administration',
    audience: 'admin',
    tags: ['setup', 'install', 'database', 'migrate', 'admin'],
    printable: true,
    sections: [
      {
        paragraphs: [
          'Complete these steps once when PocoClinic is installed on the clinic server. Use the Administrator guide on the admin dashboard to track progress.',
        ],
      },
      {
        heading: 'Server requirements',
        list: [
          'Clinic-owned server (Raspberry Pi or small PC) on a private LAN',
          'SQLite file at DATABASE_URL (e.g. /var/lib/pococlinic/pococlinic.db on the server)',
          'BACKUP_DIR writable on disk; DOCUMENT_ENCRYPTION_KEY set for production (patient documents encrypt in the database)',
          'Staff tablets/workstations on the same network — browsers only, no internet required for daily care',
        ],
      },
      {
        heading: 'Setup steps',
        steps: [
          'Set DATABASE_URL and other secrets in /etc/pococlinic/env (production) or backend/.env (development).',
          'Set DOCUMENT_ENCRYPTION_KEY (openssl rand -base64 32) and record it on the Safe vault sheet.',
          'Run migrations: migrate.bat or go run ./cmd/migrate (dev); sudo /opt/pococlinic/bin/migrate (production Pi).',
          'Start services: run-backend.bat + run-frontend.bat (dev) or systemctl start pococlinic (production).',
          'Sign in at /login/admin with the bootstrap administrator account.',
          'Print the Safe vault credentials sheet, fill secrets by hand, lock in the safe, then delete digital copies including bootstrap-admin-once.txt.',
          'Assemble the printed ops binder (docs/ops/binder): Section A emergency, B weekly/monthly, C system help. No daily paper logs.',
          'Create staff accounts, print badges, and have each person set a private PIN.',
          'Fill clinic emergency contacts under Help → Clinic emergency contacts.',
          'Label USB drives (Admin dashboard → USB rotation) and run your first backup.',
          'Verify the backup file and run System health check on the admin dashboard.',
        ],
      },
      {
        alert: {
          color: 'blue',
          title: 'Development vs production',
          body: 'Developers use a local SQLite file and .env. Production clinics use /etc/pococlinic/env on the server — see docs/guide/for-administrators/first-time-setup.md.',
        },
      },
    ],
  },
  {
    id: 'admin-bootstrap-login',
    title: 'Administrator bootstrap sign-in',
    summary: 'How admins sign in before staff badges exist, and for server-side setup tasks.',
    category: 'administration',
    audience: 'admin',
    tags: ['admin', 'login', 'bootstrap', 'setup'],
    sections: [
      {
        paragraphs: [
          'Staff use badge + PIN at /login. Administrators also use /login/admin with email, admin key, and PIN for bootstrap access — especially before badges are printed.',
        ],
      },
      {
        steps: [
          'Open /login/admin in the clinic browser.',
          'Enter administrator email, key, and 4-digit PIN.',
          'You land in the EMR with access to Dashboard, Staff, Forms, and Audit log.',
        ],
      },
      {
        heading: 'When to use admin login',
        list: [
          'First-time setup before staff accounts exist',
          'Printing badges from Staff management',
          'Backup, restore, and system health tasks',
          'Creating form templates and reviewing audit logs',
        ],
      },
    ],
  },
  {
    id: 'safe-vault-credentials',
    title: 'Safe vault credentials',
    summary: 'Print once, fill secrets by hand, lock in a safe, then delete every digital copy.',
    category: 'administration',
    audience: 'admin',
    tags: ['safe', 'vault', 'credentials', 'binder', 'secrets', 'print', 'encryption'],
    printable: true,
    sections: [
      {
        paragraphs: [
          'PocoClinic is designed so a non-technical administrator can recover the clinic from a printed binder. Secrets must not live forever on the server hard drive or in email.',
          'Use the printable Safe vault credentials sheet in the repository (docs/ops/safe-credentials-vault.md). Print the blank template, fill it by hand, put it in the ops binder, lock the binder in a safe, then delete digital copies.',
        ],
      },
      {
        heading: 'Print → lock → delete',
        steps: [
          'Print the blank vault sheet (use Print on this article, or print docs/ops/safe-credentials-vault.md).',
          'Handwrite admin email/key/PIN, JWT secrets, DOCUMENT_ENCRYPTION_KEY, and DATABASE_URL from first-time setup.',
          'Add emergency contacts and safe key-holder names.',
          'Place the filled sheet in the ops binder and lock the binder in the clinic safe.',
          'Delete data/bootstrap-admin-once.txt from the server.',
          'Remove any Notepad/.env printouts, screenshots, or chat messages that contain the filled values.',
          'Have a second administrator witness and sign the sheet.',
        ],
      },
      {
        heading: 'What belongs on the vault sheet',
        list: [
          'Bootstrap / break-glass administrator email, key, and PIN policy',
          'JWT_ACCESS_SECRET and JWT_REFRESH_SECRET',
          'DOCUMENT_ENCRYPTION_KEY (required to open encrypted patient documents after restore)',
          'DATABASE_URL and server OS passwords if used',
          'Break-glass steps for lockout, rebuild, and suspected leak',
        ],
      },
      {
        alert: {
          color: 'yellow',
          title: 'Never store filled secrets in the EMR',
          body: 'Do not paste live keys into Help notes, forms, or browser storage. This article only explains the process. The lasting copy is paper in the safe.',
        },
      },
      {
        heading: 'Related binder inserts (no secrets)',
        list: [
          'Administrator runbook — daily USB backup and restore',
          'Physical security binder — room, network, badge policy',
          'Emergency procedures — outage and incident steps',
        ],
      },
    ],
  },
  {
    id: 'ops-binder',
    title: 'Printed ops binder',
    summary: 'Three uses: emergency when down, weekly/monthly processes, static system help. No daily logs.',
    category: 'administration',
    audience: 'admin',
    tags: ['binder', 'print', 'ops', 'runbook', 'howto', 'emergency', 'monthly'],
    printable: true,
    sections: [
      {
        paragraphs: [
          'Keep the ops binder closed in a cabinet. Open it for emergencies, scheduled weekly/monthly work, or printed system help — not for everyday charting.',
          'Print the packet from docs/ops/binder/. Day-to-day backup status stays on the Admin dashboard; do not keep a daily paper backup log in the binder.',
        ],
      },
      {
        heading: 'Section A — Emergency (system down)',
        list: [
          'Emergency procedures',
          'Paper care when EMR is down',
          'Incident and restore log + emergency contacts',
          'Vault credentials stay in the safe for rebuild',
        ],
      },
      {
        heading: 'Section B — Weekly / monthly',
        list: [
          'Weekly and monthly process sheet',
          'Monthly testing checklist',
          'Security audit checklist',
          'Backup and restore runbook (instructions only)',
          'Physical security insert',
        ],
      },
      {
        heading: 'Section C — System help (device-independent)',
        list: [
          'What is PocoClinic',
          'Staff how-to',
          'Administrator how-to',
          'Same content on every PC, laptop, or tablet — no per-device notes',
        ],
      },
      {
        alert: {
          color: 'blue',
          title: 'Safe vs binder',
          body: 'Filled Safe vault credentials stay locked in the safe. The binder holds procedures and checklists only.',
        },
      },
    ],
  },
  {
    id: 'admin-dashboard',
    title: 'Admin dashboard',
    summary: 'System health, backup status, patient census, administrator guide, and action-needed alerts.',
    category: 'administration',
    audience: 'admin',
    tags: ['admin', 'health', 'status'],
    sections: [
      {
        steps: [
          'Sign in as an administrator.',
          'Open Dashboard in the sidebar under Administration.',
          'Start with the Administrator guide and Action needed alerts.',
          'Review storage mode, backup age, security overview, and resource counts.',
        ],
      },
      {
        heading: 'Dashboard sections',
        list: [
          'Administrator guide — setup checklist and links to every admin task',
          'Action needed — backup overdue, database offline, migrations, locked accounts, restore drill',
          'System health check — run on demand to verify database, storage, and document integrity',
          'USB rotation — today\'s backup drive and printable labels',
          'Maintenance reminders — restore drill tracker and document storage alerts',
          'Backups — create, verify, and restore backup files',
        ],
      },
      {
        heading: 'Action-needed alerts',
        list: [
          'Backup overdue or aging — run a USB backup today.',
          'Database offline or in-memory mode — configure DATABASE_URL and migrations.',
          'Pending migrations — run migrate.bat (dev) or sudo /opt/pococlinic/bin/migrate (production).',
          'Locked staff accounts — review Staff management.',
          'Restore drill overdue — test a backup on a spare machine.',
        ],
      },
    ],
  },
  {
    id: 'system-health-check',
    title: 'System health check',
    summary: 'Verify database, backup directory, encrypted documents, and integrity.',
    category: 'administration',
    audience: 'admin',
    tags: ['health', 'integrity', 'admin', 'documents'],
    sections: [
      {
        steps: [
          'Open Admin dashboard.',
          'Scroll to System health check.',
          'Click Run check and review each line item.',
        ],
      },
      {
        heading: 'What is checked',
        table: {
          headers: ['Check', 'Meaning'],
          rows: [
            ['Database', 'Persistent storage connected (not in-memory mode)'],
            ['Migrations', 'All schema migrations applied'],
            ['Backup directory', 'BACKUP_DIR exists and is writable'],
            ['Document storage', 'Encrypted document blobs in the database (AES-256-GCM)'],
            ['Storage usage', 'Backup (and legacy document directory) sizes in MB; warns when volumes grow large'],
            ['Document integrity', 'Each document has encrypted DB content or a matching legacy on-disk file; flags orphans'],
          ],
        },
      },
      {
        paragraphs: [
          'Run after setup, after restore, and whenever document uploads fail. Pair with backup Verify for offline archive integrity.',
        ],
      },
      {
        heading: 'Audit log retention',
        paragraphs: [
          'Clinics may schedule audit log rotation with go run ./cmd/audit-purge and AUDIT_RETENTION_DAYS. Each purge writes a system.audit.purged event to the audit log. See the security audit checklist in the ops binder.',
        ],
      },
    ],
  },
  {
    id: 'compliance-checklist',
    title: 'HIPAA compliance checklist',
    summary: 'Automated operational controls for backups, authentication, audit logging, and storage.',
    category: 'administration',
    audience: 'admin',
    tags: ['compliance', 'hipaa', 'admin', 'audit'],
    sections: [
      {
        steps: [
          'Open Admin dashboard → Overview tab.',
          'Find HIPAA compliance checklist and click Run checklist.',
          'Resolve any Fail items immediately; document Warn items in the ops binder.',
        ],
      },
      {
        heading: 'Controls checked',
        list: [
          'Persistent database storage (not in-memory mode)',
          'Document storage encrypted in database',
          'Backup within 24–48 hours',
          'Schema migrations current',
          'No staff on default PIN',
          'Locked accounts and failed sign-ins reviewed',
          'Audit events recorded in the last 7 days',
        ],
      },
      {
        paragraphs: [
          'Run monthly with the ops binder monthly testing checklist. Download audit CSV from Admin → Audit log and archive offline before purge.',
        ],
      },
    ],
  },
  {
    id: 'security-overview',
    title: 'Security overview',
    summary: 'Sessions, locked accounts, default PINs, and document storage health at a glance.',
    category: 'administration',
    audience: 'admin',
    tags: ['security', 'sessions', 'pin', 'admin'],
    sections: [
      {
        paragraphs: [
          'The Security overview card on the admin dashboard summarizes auth and storage health. PocoClinic assumes a private LAN but still enforces badge + PIN, audit logging, and session timeouts.',
        ],
      },
      {
        heading: 'Metrics',
        list: [
          'Active sessions — staff currently signed in',
          'Locked accounts — too many failed PIN attempts (15-minute lockout)',
          'Default PIN — staff who must change PIN on next sign-in',
          'Patient documents — count of uploaded files on the server',
          'Document storage — encrypted blobs in the database (legacy DOCUMENTS_DIR only if older files remain)',
        ],
      },
      {
        alert: {
          color: 'yellow',
          title: 'Default PIN accounts',
          body: 'New staff start with an administrator-assigned PIN. They must set a private PIN before routine patient work. Follow up if the count stays high.',
        },
      },
    ],
  },
  {
    id: 'form-reports',
    title: 'Form submission reports',
    summary: 'Filter and export form data as CSV for quality review or reporting.',
    category: 'administration',
    audience: 'admin',
    tags: ['forms', 'reports', 'csv', 'export'],
    sections: [
      {
        steps: [
          'Open Forms in the sidebar.',
          'Click Reports (or go to /forms/reports).',
          'Choose a template and optional date range.',
          'Review the table and use Export CSV if needed.',
        ],
      },
      {
        paragraphs: [
          'Reports include current submissions only. Data stays on the clinic server — export files should be handled per your clinic privacy policy.',
        ],
      },
    ],
  },
  {
    id: 'staff-management',
    title: 'Manage staff accounts',
    summary: 'Create staff, print badges, reissue badges, and edit accounts.',
    category: 'administration',
    audience: 'admin',
    tags: ['staff', 'users', 'badge', 'admin'],
    sections: [
      {
        steps: [
          'Open Staff in the sidebar.',
          'Click New staff member to create an account (name, email, initial PIN).',
          'Open a staff member to print their badge QR or reissue a lost badge.',
        ],
      },
      {
        alert: {
          color: 'red',
          title: 'Reissue invalidates the old badge',
          body: 'When you reissue a badge, the previous QR code stops working immediately. Print and hand the new badge to the employee.',
        },
      },
    ],
  },
  {
    id: 'form-templates',
    title: 'Form templates (admin)',
    summary: 'Create and edit the forms staff fill out on patient charts.',
    category: 'administration',
    audience: 'admin',
    tags: ['forms', 'templates', 'builder'],
    sections: [
      {
        steps: [
          'Open Forms in the sidebar.',
          'Create a new template or edit an existing one.',
          'Add fields, save, and staff will see the template on patient charts.',
        ],
      },
      {
        paragraphs: [
          'Use Reports to export submission data and filter by date range.',
        ],
      },
    ],
  },
  {
    id: 'audit-log',
    title: 'Audit log',
    summary: 'Review sign-ins, patient access, notes, backups, and admin actions.',
    category: 'administration',
    audience: 'admin',
    tags: ['audit', 'compliance', 'security'],
    sections: [
      {
        steps: [
          'Open Dashboard → Audit log, or Audit log in the sidebar.',
          'Filter by event type or user if needed.',
          'Use this to investigate access questions or backup/restore events.',
          'Each month: download the CSV (admin only), save it offline or in your ops binder, then mark archived on the dashboard.',
        ],
      },
      {
        paragraphs: [
          'Patient viewed, updated, deleted, and clinical note events are logged with the staff member and timestamp.',
          'Monthly audit purge removes old rows from the live database. Archive the CSV first so your clinic keeps a long-term copy.',
        ],
      },
    ],
  },
  {
    id: 'daily-backup',
    title: 'Daily USB backup',
    summary: 'Protect clinic data with a labeled USB drive every day.',
    category: 'backup-restore',
    audience: 'admin',
    tags: ['backup', 'usb', 'daily'],
    printable: true,
    sections: [
      {
        heading: 'Daily checklist',
        steps: [
          'Confirm staff can sign in and the system looks normal.',
          'Insert today\'s labeled USB drive into the server.',
          'Run backup: Admin → Backup now, or use the Backup Helper at http://127.0.0.1:9090 on the server.',
          'Confirm success and a new backup file appears.',
          'Copy or store the backup on USB; remove drive and lock it up.',
          'Confirm Admin dashboard backup status is green — no daily paper log in the ops binder.',
        ],
      },
      {
        heading: 'USB rotation',
        paragraphs: [
          'Use at least two drives (e.g. Mon / Wed / Fri). Never rely on a single USB stick.',
          'Admin dashboard → USB rotation shows today\'s drive, editable labels, and a Print labels button for the USB kit.',
        ],
      },
      {
        heading: 'Backup age colors (Admin dashboard)',
        table: {
          headers: ['Color', 'Meaning'],
          rows: [
            ['Green', 'Backup within 24 hours'],
            ['Yellow', '24–48 hours — backup today'],
            ['Red', 'Over 48 hours — backup before end of day'],
          ],
        },
      },
    ],
  },
  {
    id: 'backup-helper',
    title: 'Backup Helper (on the server PC)',
    summary: 'Guided backup and restore at http://127.0.0.1:9090 — localhost only.',
    category: 'backup-restore',
    audience: 'admin',
    tags: ['backup', 'helper', 'wizard', 'localhost'],
    sections: [
      {
        paragraphs: [
          'The Backup Helper runs only on the clinic server computer, not on the wider network. Open http://127.0.0.1:9090 in a browser on that machine.',
        ],
      },
      {
        heading: 'Daily backup wizard',
        steps: [
          'Start the helper (run-ops-helper.bat on Windows dev, or systemd on Pi).',
          'Click Backup and follow each step: USB plugged in → Create backup → Copy to USB → Done.',
          'Print the checklist from the helper for your ops binder if needed.',
        ],
      },
      {
        alert: {
          color: 'blue',
          title: 'Pi touchscreen mode',
          body: 'On a Raspberry Pi with a touchscreen on the server, use the Pi build for large backup/restore buttons. Clinical charting still uses staff workstations on the LAN.',
        },
      },
    ],
  },
  {
    id: 'restore-disaster',
    title: 'Restore after disaster',
    summary: 'Recover all clinic data from a USB backup — use with care.',
    category: 'backup-restore',
    audience: 'admin',
    tags: ['restore', 'disaster', 'recovery'],
    printable: true,
    sections: [
      {
        alert: {
          color: 'red',
          title: 'Restore replaces everything',
          body: 'All patients, forms, staff sessions, and audit history will be replaced by the backup contents. Sign out all staff except one administrator before restoring.',
        },
      },
      {
        heading: 'From the Backup Helper (recommended)',
        steps: [
          'Sign out all staff except one administrator.',
          'On the server PC, open http://127.0.0.1:9090 → Restore.',
          'Select the backup, type RESTORE to confirm, and wait for completion.',
          'Verify sign-in and patient count; take a fresh backup immediately.',
        ],
      },
      {
        heading: 'From the Admin dashboard',
        steps: [
          'Sign out all staff except one administrator.',
          'Admin → Backups → Restore on the chosen file → confirm.',
          'Verify sign-in and patient count.',
        ],
      },
      {
        heading: 'After restore verification',
        list: [
          'Administrator can sign in',
          'Patient list loads with expected count',
          'Staff badges work (reissue if backup was old)',
          'Take a new backup and log the restore in your binder',
        ],
      },
    ],
  },
  {
    id: 'backup-verification',
    title: 'Verify a backup file',
    summary: 'Confirm checksums match the manifest before relying on a USB copy.',
    category: 'backup-restore',
    audience: 'admin',
    tags: ['backup', 'verify', 'checksum'],
    sections: [
      {
        steps: [
          'Open Admin → Backups.',
          'Click Verify on the backup file you want to check.',
          'A green notification means checksums and document integrity passed.',
        ],
      },
      {
        paragraphs: [
          'Verification checks manifest checksums and that each patient_documents row has encrypted content in the DB export or a matching legacy file under documents/ in the archive. Run after copying to USB and during quarterly restore drills.',
        ],
      },
    ],
  },
  {
    id: 'quarterly-drill',
    title: 'Quarterly restore drill',
    summary: 'Prove backups work on a test machine — not during clinic hours.',
    category: 'backup-restore',
    audience: 'admin',
    tags: ['restore', 'drill', 'test'],
    printable: true,
    sections: [
      {
        steps: [
          'Use a spare PC or test machine — not production during clinic hours.',
          'Install PocoClinic and restore the latest USB backup with confirmation.',
          'Verify login and patient count match expectations.',
          'Sign the drill log in your ops binder.',
        ],
      },
    ],
  },
  {
    id: 'troubleshooting-common',
    title: 'Common problems',
    summary: 'Quick fixes for sign-in, data loss, backup, and lockouts.',
    category: 'troubleshooting',
    audience: 'all',
    tags: ['fix', 'error', 'problem'],
    sections: [
      {
        table: {
          headers: ['Symptom', 'Likely cause', 'What to try'],
          rows: [
            ['Cannot sign in', 'Wrong PIN or locked account', 'Wait 15 min after lockout; verify badge scan; ask admin'],
            ['Blank page or spinner', 'Server not running', 'Ask admin to restart the clinic server'],
            ['Patient data gone after restart', 'Database not configured', 'Admin: set DATABASE_URL and run migrations'],
            ['Backup button disabled', 'In-memory mode', 'Admin: configure persistent database'],
            ['Restore fails checksum', 'Corrupted USB file', 'Try an earlier backup drive'],
            ['Session keeps ending', '15 min inactivity timeout', 'Normal — sign in again with badge + PIN'],
          ],
        },
      },
    ],
  },
  {
    id: 'clinic-contacts',
    title: 'Clinic emergency contacts',
    summary: 'Who to call when the system or network needs help.',
    category: 'troubleshooting',
    audience: 'all',
    tags: ['contact', 'emergency', 'phone'],
    sections: [
      {
        paragraphs: [
          'Administrators can fill in contact details below. This information is stored only in this browser on the clinic network — not sent to the cloud.',
        ],
      },
    ],
  },
];

export function getArticleById(id: string): HelpArticle | undefined {
  return helpArticles.find((article) => article.id === id);
}

export function getCategoryById(id: string): HelpCategory | undefined {
  return helpCategories.find((category) => category.id === id);
}

export function getRelatedArticles(articleId: string, isAdmin: boolean, limit = 4): HelpArticle[] {
  const article = getArticleById(articleId);
  if (!article) {
    return [];
  }
  return helpArticles
    .filter((candidate) => {
      if (candidate.id === articleId) {
        return false;
      }
      if (candidate.audience === 'admin' && !isAdmin) {
        return false;
      }
      return candidate.category === article.category || candidate.tags.some((tag) => article.tags.includes(tag));
    })
    .slice(0, limit);
}

export function getFeaturedArticleIds(isAdmin: boolean): string[] {
  if (isAdmin) {
    return ['how-it-works', 'admin-setup', 'sign-in', 'daily-backup'];
  }
  return ['how-it-works', 'sign-in', 'patient-list', 'troubleshooting-common'];
}

export function groupArticlesByCategory(
  articles: HelpArticle[],
): { category: HelpCategory; articles: HelpArticle[] }[] {
  return helpCategories
    .map((category) => ({
      category,
      articles: articles.filter((article) => article.category === category.id),
    }))
    .filter((group) => group.articles.length > 0);
}
