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
  dateOfBirth: string;
  gender: Gender;
  email?: string;
  phoneNumber?: string;
  address?: Address;
  createdAt: string;
  updatedAt: string;
  medicalNumber: string;
}

export interface PaginatedPatients {
  patients: Patient[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
  totalPages: number;
}

export interface PatientFormData {
  firstName: string;
  lastName: string;
  middleName?: string;
  dateOfBirth: Date;
  gender: Gender;
  email?: string;
  phoneNumber?: string;
  address?: Address;
  medicalNumber: string;
} 