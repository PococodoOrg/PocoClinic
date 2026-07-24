import { createApiClient } from './client';
import {
  FormGroup,
  FormGroupFormData,
  FormSubmission,
  FormSubmissionReport,
  FormTemplate,
  FormTemplateFormData,
  SubmissionReportPage,
} from '../types/form';

const formsApi = createApiClient();

export const fetchFormGroups = async (): Promise<FormGroup[]> => {
  const response = await formsApi.get<{ groups: FormGroup[] }>('/form-groups');
  return response.data.groups ?? [];
};

export const createFormGroup = async (data: FormGroupFormData): Promise<FormGroup> => {
  const response = await formsApi.post<FormGroup>('/form-groups', data);
  return response.data;
};

export const updateFormGroup = async (id: string, data: FormGroupFormData): Promise<FormGroup> => {
  const response = await formsApi.put<FormGroup>(`/form-groups/${id}`, data);
  return response.data;
};

export const deleteFormGroup = async (id: string): Promise<void> => {
  await formsApi.delete(`/form-groups/${id}`);
};

export const fetchFormTemplates = async (): Promise<FormTemplate[]> => {
  const response = await formsApi.get<{ templates: FormTemplate[] }>('/form-templates');
  return response.data.templates ?? [];
};

export interface TemplateSubmissionsQuery {
  page?: number;
  pageSize?: number;
  from?: string;
  to?: string;
}

export const fetchTemplateSubmissions = async (
  templateId: string,
  query: TemplateSubmissionsQuery = {},
): Promise<SubmissionReportPage> => {
  const response = await formsApi.get<SubmissionReportPage>(`/form-templates/${templateId}/submissions`, {
    params: query,
  });
  return response.data;
};

export const fetchAllTemplateSubmissions = async (
  templateId: string,
  query: Omit<TemplateSubmissionsQuery, 'page' | 'pageSize'> = {},
): Promise<FormSubmissionReport[]> => {
  const pageSize = 100;
  let page = 1;
  const items: FormSubmissionReport[] = [];

  while (true) {
    const report = await fetchTemplateSubmissions(templateId, { ...query, page, pageSize });
    items.push(...report.items);
    if (items.length >= report.total || report.items.length === 0) {
      break;
    }
    page += 1;
  }

  return items;
};

export const getFormTemplate = async (id: string): Promise<FormTemplate> => {
  const response = await formsApi.get<FormTemplate>(`/form-templates/${id}`);
  return response.data;
};

export const createFormTemplate = async (data: FormTemplateFormData): Promise<FormTemplate> => {
  const response = await formsApi.post<FormTemplate>('/form-templates', data);
  return response.data;
};

export const updateFormTemplate = async (id: string, data: FormTemplateFormData): Promise<FormTemplate> => {
  const response = await formsApi.put<FormTemplate>(`/form-templates/${id}`, data);
  return response.data;
};

export const deleteFormTemplate = async (id: string): Promise<void> => {
  await formsApi.delete(`/form-templates/${id}`);
};

export const fetchPatientFormSubmissions = async (patientId: string): Promise<FormSubmission[]> => {
  const response = await formsApi.get<{ submissions: FormSubmission[] }>(
    `/patients/${patientId}/form-submissions`,
  );
  return response.data.submissions ?? [];
};

export const fetchFormEntryHistory = async (
  patientId: string,
  entryId: string,
): Promise<FormSubmission[]> => {
  const response = await formsApi.get<{ history: FormSubmission[] }>(
    `/patients/${patientId}/form-entries/${entryId}/history`,
  );
  return response.data.history ?? [];
};

export const savePatientForm = async (
  patientId: string,
  templateId: string,
  answers: Record<string, string | number | boolean>,
  entryId?: string,
): Promise<FormSubmission> => {
  const response = await formsApi.post<FormSubmission>(`/patients/${patientId}/form-submissions`, {
    templateId,
    entryId,
    answers,
  });
  return response.data;
};
