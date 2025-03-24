import axios, { AxiosError } from 'axios';
import { Patient, PatientFormData, PaginatedPatients } from '../types/patient';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export interface ValidationError {
  code: string;
  message: string;
  errors?: Record<string, string[]>;
}

export const patientApi = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Helper function to convert form data to API format
const formatPatientData = (data: PatientFormData) => {
  const { phone, ...rest } = data;
  return {
    ...rest,
    // Format date as full ISO timestamp at midnight UTC
    dateOfBirth: data.dateOfBirth 
      ? new Date(data.dateOfBirth.getFullYear(), data.dateOfBirth.getMonth(), data.dateOfBirth.getDate()).toISOString()
      : null,
    // Format measurements as numbers
    height: data.height !== null && data.height !== undefined ? Number(data.height) : null,
    weight: data.weight !== null && data.weight !== undefined ? Number(data.weight) : null,
    // Map phone to phoneNumber
    phoneNumber: phone,
    // Keep the address object as is
    address: data.address,
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

export const createPatient = async (data: PatientFormData): Promise<Patient> => {
  const response = await patientApi.post<Patient>('/patients', formatPatientData(data));
  return response.data;
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
  const response = await patientApi.put<Patient>(`/patients/${id}`, formatPatientData(data));
  return response.data;
};

export const deletePatient = async (id: string): Promise<void> => {
  await patientApi.delete(`/patients/${id}`);
};

// Add fetchPatient function to fetch a single patient by ID
export const fetchPatient = async (id: string): Promise<Patient> => {
  const response = await patientApi.get(`/patients/${id}`);
  return response.data;
}; 