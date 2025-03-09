// Auth Types
export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface User {
  id: string
  name: string
  email: string
  role: UserRole
  lastLogin?: string
  createdAt?: string
  updatedAt?: string
}

export type UserRole = 'admin' | 'doctor' | 'nurse' | 'staff' | 'patient'

// Patient Types
export interface Patient {
  id: string
  firstName: string
  lastName: string
  dateOfBirth: string
  gender: PatientGender
  email?: string
  phone?: string
  address?: string
  createdAt?: string
  updatedAt?: string
}

export type PatientGender = 'male' | 'female' | 'other' | 'unknown'

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface PaginationParams {
  page?: number
  pageSize?: number
  search?: string
}

// Error Types
export interface ApiError {
  message: string
  code?: string
  details?: Record<string, string[]>
} 