-- Admin-configurable required patient fields (clinic-wide setting).

CREATE TABLE IF NOT EXISTS clinic_settings (
    key TEXT PRIMARY KEY,
    value_json TEXT NOT NULL,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

INSERT OR IGNORE INTO clinic_settings (key, value_json) VALUES (
    'patient_field_requirements',
    '{"firstName":true,"lastName":true,"middleName":false,"dateOfBirth":true,"gender":true,"email":false,"phoneNumber":false,"addressStreet":false,"addressCity":false,"addressState":false,"addressPostalCode":false,"height":false,"weight":false}'
);
