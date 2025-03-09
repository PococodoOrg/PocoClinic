import axios from 'axios';
import { Patient, PatientFormData, PaginatedPatients } from '../types/patient';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface GetPatientsParams {
  page?: number;
  pageSize?: number;
  search?: string;
}

export const patientApi = {
  // Create a new patient
  createPatient: async (data: PatientFormData): Promise<Patient> => {
    const response = await api.post<Patient>('/patients', data);
    return response.data;
  },

  // Get paginated list of patients
  getPatients: async (params: GetPatientsParams): Promise<PaginatedPatients> => {
    const response = await api.get<PaginatedPatients>('/patients', { params });
    return response.data;
  },

  // Get a single patient by ID
  getPatient: async (id: string): Promise<Patient> => {
    const response = await api.get<Patient>(`/patients/${id}`);
    return response.data;
  },

  // Update a patient
  updatePatient: async (id: string, data: PatientFormData): Promise<Patient> => {
    const response = await api.put<Patient>(`/patients/${id}`, data);
    return response.data;
  },

  // Delete a patient
  deletePatient: async (id: string): Promise<void> => {
    await api.delete(`/patients/${id}`);
  },
};

// Add request interceptor for authentication
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      console.error('API Error:', error.response.data);
      throw new Error(error.response.data.error || 'An error occurred');
    } else if (error.request) {
      // The request was made but no response was received
      console.error('Network Error:', error.request);
      throw new Error('Network error - no response received');
    } else {
      // Something happened in setting up the request that triggered an Error
      console.error('Request Error:', error.message);
      throw new Error('Error setting up the request');
    }
  }
); 