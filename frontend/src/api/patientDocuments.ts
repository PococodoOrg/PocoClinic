import { PatientDocument } from '../types/patientDocument';
import { createApiClient } from './client';
import { safeDownloadFilename } from '../utils/safeFilename';

const documentsApi = createApiClient({ timeout: 60_000, json: false });

export async function fetchPatientDocuments(patientId: string): Promise<PatientDocument[]> {
  const response = await documentsApi.get<PatientDocument[]>(`/patients/${patientId}/documents`);
  return response.data;
}

export async function uploadPatientDocument(patientId: string, file: File): Promise<PatientDocument> {
  const formData = new FormData();
  formData.append('file', file);
  const response = await documentsApi.post<PatientDocument>(`/patients/${patientId}/documents`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  return response.data;
}

export async function deletePatientDocument(patientId: string, documentId: string): Promise<void> {
  await documentsApi.delete(`/patients/${patientId}/documents/${documentId}`);
}

export async function downloadPatientDocument(patientId: string, documentId: string, fileName: string): Promise<void> {
  const response = await documentsApi.get(`/patients/${patientId}/documents/${documentId}/download`, {
    responseType: 'blob',
  });
  const url = window.URL.createObjectURL(response.data);
  const link = document.createElement('a');
  link.href = url;
  link.download = safeDownloadFilename(fileName);
  link.click();
  window.URL.revokeObjectURL(url);
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) {
    return `${bytes} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export { formatBytes };
