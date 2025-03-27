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
  middleName?: string;
  dateOfBirth: string;  // YYYY-MM-DD format
  gender: Gender;
  email: string;
  phoneNumber: string;
  address?: Address;
  height?: number | null;  // Height in centimeters
  weight?: number | null;  // Weight in kilograms
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
  middleName?: string;
  dateOfBirth: Date | null;  // Use Date for form handling
  gender: Gender;
  email: string;
  phoneNumber: string;  // Changed from phone to phoneNumber for consistency
  address?: Address;
  street?: string;
  city?: string;
  state?: string;
  zipCode?: string;
  height?: number | null;  // Height in centimeters
  weight?: number | null;  // Weight in kilograms
}

// API data interface with string date for the API
export interface PatientApiData {
  firstName: string;
  lastName: string;
  middleName?: string;
  dateOfBirth: string | null;
  gender: Gender;
  email: string;
  phoneNumber: string;
  address?: Address;
  height?: number | null;  // Height in centimeters
  weight?: number | null;  // Weight in kilograms
} 