import { AxiosError } from 'axios';
import { Patient, PatientFormData, PaginatedPatients, Gender, PatientFieldRequirements } from '../types/patient';
import { createApiClient } from './client';
import { ValidationError } from '../utils/apiError';

export { ValidationError } from '../utils/apiError';

export interface PatientListParams {
  page?: number;
  pageSize?: number;
  search?: string;
  gender?: Gender | '';
  dobFrom?: string;
  dobTo?: string;
  registeredSince?: string;
}

const patientApi = createApiClient();

const formatPatientData = (data: PatientFormData) => ({
  firstName: data.firstName,
  lastName: data.lastName,
  middleName: data.middleName,
  dateOfBirth: data.dateOfBirth,
  gender: data.gender,
  email: data.email,
  phoneNumber: data.phoneNumber,
  address: data.address,
  height: data.height,
  weight: data.weight,
});

function rethrowValidation(error: unknown, fallback: string): never {
  if (error instanceof ValidationError) {
    throw error;
  }
  if (error instanceof AxiosError && error.response?.data?.code === 'VALIDATION_ERROR') {
    throw error;
  }
  throw error instanceof Error ? error : new Error(fallback);
}

export const createPatient = async (data: PatientFormData): Promise<Patient> => {
  try {
    const response = await patientApi.post<Patient>('/patients', formatPatientData(data));
    return response.data;
  } catch (error) {
    rethrowValidation(error, 'Failed to create patient');
  }
};

export const fetchPatients = async (params: PatientListParams = {}): Promise<PaginatedPatients> => {
  const response = await patientApi.get<PaginatedPatients>('/patients', {
    params: {
      page: params.page ?? 1,
      pageSize: params.pageSize ?? 10,
      search: params.search,
      gender: params.gender || undefined,
      dobFrom: params.dobFrom || undefined,
      dobTo: params.dobTo || undefined,
      registeredSince: params.registeredSince || undefined,
    },
  });
  return response.data;
};

export const getPatient = async (id: string): Promise<Patient> => {
  const response = await patientApi.get<Patient>(`/patients/${id}`);
  return response.data;
};

export const fetchPatient = getPatient;

export const updatePatient = async (id: string, data: PatientFormData): Promise<Patient> => {
  try {
    const response = await patientApi.put<Patient>(`/patients/${id}`, formatPatientData(data));
    return response.data;
  } catch (error) {
    rethrowValidation(error, 'Failed to update patient');
  }
};

export const deletePatient = async (id: string): Promise<void> => {
  await patientApi.delete(`/patients/${id}`);
};

export const fetchPatientFieldRequirements = async (): Promise<PatientFieldRequirements> => {
  const response = await patientApi.get<{ requirements: PatientFieldRequirements }>('/patients/field-requirements');
  return response.data.requirements;
};
