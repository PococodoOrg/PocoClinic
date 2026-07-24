-- PocoClinic SQLite schema (squashed).
-- UUIDs and timestamps use TEXT. Binary columns use BLOB. Booleans use INTEGER 0/1.

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    key_hash BLOB NOT NULL,
    key_salt BLOB NOT NULL,
    pin_hash BLOB NOT NULL,
    pin_salt BLOB NOT NULL,
    key_lookup TEXT,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TEXT,
    last_login TEXT,
    must_change_pin INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_key_lookup ON users (key_lookup) WHERE key_lookup IS NOT NULL;

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address TEXT,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);

CREATE TABLE IF NOT EXISTS patients (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    middle_name TEXT,
    date_of_birth TEXT NOT NULL,
    gender TEXT NOT NULL,
    email TEXT,
    phone_number TEXT,
    height REAL,
    weight REAL,
    address_street TEXT,
    address_city TEXT,
    address_state TEXT,
    address_postal_code TEXT,
    address_country TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    user_id TEXT,
    resource_type TEXT,
    resource_id TEXT,
    ip_address TEXT,
    user_agent TEXT,
    details TEXT,
    success INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_event_type ON audit_logs (event_type);

CREATE TABLE IF NOT EXISTS form_groups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

INSERT INTO form_groups (id, name, sort_order)
SELECT '00000000-0000-0000-0000-000000000001', 'General', 0
WHERE NOT EXISTS (SELECT 1 FROM form_groups WHERE id = '00000000-0000-0000-0000-000000000001');

CREATE TABLE IF NOT EXISTS form_templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    fields TEXT NOT NULL,
    group_id TEXT NOT NULL DEFAULT '00000000-0000-0000-0000-000000000001'
        REFERENCES form_groups(id),
    form_type TEXT NOT NULL DEFAULT 'singleton',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS form_submissions (
    id TEXT PRIMARY KEY,
    template_id TEXT NOT NULL REFERENCES form_templates(id),
    patient_id TEXT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    submitted_by TEXT NOT NULL REFERENCES users(id),
    answers TEXT NOT NULL,
    entry_id TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    is_current INTEGER NOT NULL DEFAULT 1,
    field_snapshot TEXT,
    updated_at TEXT,
    updated_by TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_form_submissions_patient_id ON form_submissions (patient_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_form_submissions_entry_id ON form_submissions (entry_id, version DESC);
CREATE INDEX IF NOT EXISTS idx_form_submissions_current ON form_submissions (patient_id, template_id, is_current);
CREATE INDEX IF NOT EXISTS idx_form_submissions_report ON form_submissions (template_id, is_current, updated_at);

CREATE TABLE IF NOT EXISTS patient_notes (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    author_id TEXT NOT NULL REFERENCES users(id),
    body TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT
);

CREATE INDEX IF NOT EXISTS idx_patient_notes_patient_id ON patient_notes (patient_id, created_at DESC);

CREATE TABLE IF NOT EXISTS patient_documents (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    uploaded_by TEXT NOT NULL REFERENCES users(id),
    file_name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes INTEGER NOT NULL,
    storage_key TEXT,
    encrypted_content BLOB,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_patient_documents_patient_id ON patient_documents (patient_id, created_at DESC);

CREATE TABLE IF NOT EXISTS exercise_plans (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active',
    created_by TEXT NOT NULL REFERENCES users(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT,
    CONSTRAINT exercise_plans_status_check CHECK (status IN ('active', 'completed', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_exercise_plans_patient_id
    ON exercise_plans (patient_id, created_at DESC);

CREATE TABLE IF NOT EXISTS exercise_log_entries (
    id TEXT PRIMARY KEY,
    plan_id TEXT NOT NULL REFERENCES exercise_plans(id) ON DELETE CASCADE,
    patient_id TEXT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    performed_at TEXT NOT NULL,
    exercise_name TEXT NOT NULL,
    sets INTEGER NOT NULL DEFAULT 0,
    reps INTEGER NOT NULL DEFAULT 0,
    duration_seconds INTEGER,
    resistance TEXT NOT NULL DEFAULT '',
    difficulty INTEGER,
    notes TEXT NOT NULL DEFAULT '',
    recorded_by TEXT NOT NULL REFERENCES users(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT,
    CONSTRAINT exercise_log_entries_sets_check CHECK (sets >= 0),
    CONSTRAINT exercise_log_entries_reps_check CHECK (reps >= 0),
    CONSTRAINT exercise_log_entries_difficulty_check CHECK (difficulty IS NULL OR (difficulty >= 1 AND difficulty <= 10))
);

CREATE INDEX IF NOT EXISTS idx_exercise_log_entries_plan_id
    ON exercise_log_entries (plan_id, performed_at DESC);

CREATE INDEX IF NOT EXISTS idx_exercise_log_entries_patient_id
    ON exercise_log_entries (patient_id, performed_at DESC);
