export type FormFieldType = 'text' | 'textarea' | 'number' | 'select' | 'checkbox';

export type FormType = 'singleton' | 'log';

export interface FormField {
  id: string;
  label: string;
  type: FormFieldType;
  required: boolean;
  options?: string[];
  sortOrder: number;
}

export interface FormGroup {
  id: string;
  name: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface FormTemplate {
  id: string;
  groupId: string;
  name: string;
  formType: FormType;
  fields: FormField[];
  createdAt: string;
  updatedAt: string;
}

export interface FormSubmission {
  id: string;
  entryId: string;
  version: number;
  isCurrent: boolean;
  templateId: string;
  templateName?: string;
  formType?: FormType;
  patientId: string;
  submittedBy: string;
  updatedBy: string;
  answers: Record<string, string | number | boolean>;
  fieldSnapshot?: FormField[];
  createdAt: string;
  updatedAt: string;
}

export interface FormSubmissionReport extends FormSubmission {
  patientName: string;
  submittedByName: string;
}

export interface SubmissionReportPage {
  items: FormSubmissionReport[];
  total: number;
  page: number;
  pageSize: number;
}

export interface FormTemplateFormData {
  groupId: string;
  name: string;
  formType: FormType;
  fields: FormField[];
}

export interface FormGroupFormData {
  name: string;
  sortOrder: number;
}

export const FIELD_TYPE_OPTIONS = [
  { value: 'text', label: 'Short text' },
  { value: 'textarea', label: 'Long text' },
  { value: 'number', label: 'Number' },
  { value: 'select', label: 'Dropdown' },
  { value: 'checkbox', label: 'Checkbox' },
];

export const FORM_TYPE_OPTIONS = [
  {
    value: 'singleton',
    label: 'Singleton',
    description: 'One active copy per patient. Edits create historical versions.',
  },
  {
    value: 'log',
    label: 'Log',
    description: 'Ongoing log with multiple entries. Each entry can be edited and versioned.',
  },
];

export function createEmptyField(sortOrder: number): FormField {
  return {
    id: `field_${crypto.randomUUID().slice(0, 8)}`,
    label: '',
    type: 'text',
    required: false,
    options: [],
    sortOrder,
  };
}

export function formatAnswerValue(value: unknown): string {
  if (typeof value === 'boolean') {
    return value ? 'Yes' : 'No';
  }
  if (value === null || value === undefined || value === '') {
    return '—';
  }
  return String(value);
}

export function getFieldLabel(submission: FormSubmission, fieldId: string): string {
  const fromSnapshot = submission.fieldSnapshot?.find((field) => field.id === fieldId)?.label;
  return fromSnapshot ?? fieldId;
}

export function groupTemplatesByGroup(
  groups: FormGroup[],
  templates: FormTemplate[],
): Array<{ group: FormGroup; templates: FormTemplate[] }> {
  const sortedGroups = [...groups].sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name));
  return sortedGroups.map((group) => ({
    group,
    templates: templates
      .filter((template) => template.groupId === group.id)
      .sort((a, b) => a.name.localeCompare(b.name)),
  }));
}
