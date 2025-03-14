export type Gender = 'male' | 'female' | 'other' | 'unknown';

export interface Address {
  street: string;
  city: string;
  state: string;
  postalCode: string;
  country: string;
}

export interface Patient {
  id: string;
  firstName: string;
  lastName: string;
  dateOfBirth: string;  // YYYY-MM-DD format
  gender: Gender;
  email: string;
  phone: string;
  address?: string;
  city?: string;
  state?: string;
  zipCode?: string;
  height?: number;
  weight?: number;
  createdAt: string;
  updatedAt: string;
}

export interface PaginatedPatients {
  patients: Patient[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
  totalPages: number;
}

// Form data interface with Date object for the form
export interface PatientFormData {
  firstName: string;
  lastName: string;
  dateOfBirth: Date | null;  // Use Date for form handling
  gender: Gender;
  email: string;
  phone: string;
  address?: string;
  city?: string;
  state?: string;
  zipCode?: string;
  height?: number | null;
  weight?: number | null;
}

// API data interface with string date for the API
export interface PatientApiData {
  firstName: string;
  lastName: string;
  dateOfBirth: string | null;
  gender: Gender;
  email: string;
  phone: string;
  address?: string;
  city?: string;
  state?: string;
  zipCode?: string;
  height?: number | null;
  weight?: number | null;
} 