export interface PatientNote {
  id: string;
  patientId: string;
  authorId: string;
  authorName: string;
  body: string;
  createdAt: string;
  updatedAt?: string;
}