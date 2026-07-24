export interface PatientDocument {
  id: string;
  patientId: string;
  uploadedBy: string;
  uploadedByName: string;
  fileName: string;
  contentType: string;
  sizeBytes: number;
  createdAt: string;
}
