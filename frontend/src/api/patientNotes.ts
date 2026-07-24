import { PatientNote } from '../types/patientNote';
import { createApiClient } from './client';

const notesApi = createApiClient();

export async function fetchPatientNotes(patientId: string): Promise<PatientNote[]> {
  const response = await notesApi.get<PatientNote[]>(`/patients/${patientId}/notes`);
  return response.data;
}

export async function createPatientNote(patientId: string, body: string): Promise<PatientNote> {
  const response = await notesApi.post<PatientNote>(`/patients/${patientId}/notes`, { body });
  return response.data;
}

export async function updatePatientNote(patientId: string, noteId: string, body: string): Promise<PatientNote> {
  const response = await notesApi.put<PatientNote>(`/patients/${patientId}/notes/${noteId}`, { body });
  return response.data;
}

export async function deletePatientNote(patientId: string, noteId: string): Promise<void> {
  await notesApi.delete(`/patients/${patientId}/notes/${noteId}`);
}