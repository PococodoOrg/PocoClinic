import axios, { AxiosError } from 'axios';
import { Patient, PatientFormData, PaginatedPatients, PatientApiData } from '../types/patient';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export class ValidationError extends Error {
  code: string;
  errors?: Record<string, string[]>;

  constructor(message: string, code: string, errors?: Record<string, string[]>) {
    super(message);
    this.name = 'ValidationError';
    this.code = code;
    this.errors = errors;
  }
}

const patientApi = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Helper function to convert form data to API format
const formatPatientData = (data: PatientFormData) => {
  return {
    firstName: data.firstName,
    lastName: data.lastName,
    middleName: data.middleName,
    dateOfBirth: data.dateOfBirth ? data.dateOfBirth.toISOString().split('T')[0] : null,
    gender: data.gender,
    email: data.email,
    phoneNumber: data.phoneNumber,
    address: data.address,
    height: data.height,
    weight: data.weight,
  };
};

// Add request interceptor for authentication
patientApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add response interceptor for error handling
patientApi.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ValidationError>) => {
    if (error.response) {
      console.error('API Error:', error.response.data);
      if (error.response.data.code === 'VALIDATION_ERROR' && error.response.data.errors) {
        throw error.response.data;
      }
      throw new Error(error.response.data.message || 'An error occurred');
    } else if (error.request) {
      console.error('Network Error:', error.request);
      throw new Error('Network error - no response received');
    } else {
      console.error('Request Error:', error.message);
      throw new Error('Error setting up the request');
    }
  }
);

const transformFormDataToApiData = (data: PatientFormData): PatientApiData => {
  return {
    firstName: data.firstName,
    lastName: data.lastName,
    middleName: data.middleName,
    dateOfBirth: data.dateOfBirth ? data.dateOfBirth.toISOString().split('T')[0] : null,
    gender: data.gender,
    email: data.email,
    phoneNumber: data.phoneNumber,
    address: data.address,
    height: data.height,
    weight: data.weight,
  };
};

export const createPatient = async (data: PatientFormData): Promise<Patient> => {
  try {
    const response = await patientApi.post<Patient>('/patients', formatPatientData(data));
    return response.data;
  } catch (error) {
    if (error instanceof AxiosError) {
      const apiError = error.response?.data;
      if (apiError?.code === 'VALIDATION_ERROR') {
        throw new ValidationError(apiError.message, apiError.code, apiError.errors);
      }
    }
    throw new Error('Failed to create patient');
  }
};

export const fetchPatients = async (params?: { page?: number; pageSize?: number; search?: string; }): Promise<PaginatedPatients> => {
  const response = await patientApi.get<PaginatedPatients>('/patients', { 
    params: {
      page: params?.page ?? 1,
      pageSize: params?.pageSize ?? 10,
      search: params?.search
    }
  });
  return response.data;
};

export const getPatient = async (id: string): Promise<Patient> => {
  const response = await patientApi.get<Patient>(`/patients/${id}`);
  return response.data;
};

export const updatePatient = async (id: string, data: PatientFormData): Promise<Patient> => {
  try {
    const response = await patientApi.put<Patient>(`/patients/${id}`, formatPatientData(data));
    return response.data;
  } catch (error) {
    if (error instanceof AxiosError) {
      const apiError = error.response?.data;
      if (apiError?.code === 'VALIDATION_ERROR') {
        throw new ValidationError(apiError.message, apiError.code, apiError.errors);
      }
    }
    throw new Error('Failed to update patient');
  }
};

export const deletePatient = async (id: string): Promise<void> => {
  await patientApi.delete(`/patients/${id}`);
};

// Add fetchPatient function to fetch a single patient by ID
export const fetchPatient = async (id: string): Promise<Patient> => {
  const response = await patientApi.get(`/patients/${id}`);
  return response.data;
}; 