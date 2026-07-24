export type UserRole = 'admin' | 'doctor' | 'nurse' | 'staff' | 'patient';

export interface StaffUser {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  lastLogin?: string;
  mustChangePin?: boolean;
  isActive?: boolean;
  isLocked?: boolean;
  lockedUntil?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PaginatedUsers {
  users: StaffUser[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
  totalPages: number;
}

export interface CreateUserFormData {
  email: string;
  name: string;
  role: UserRole;
}

export type UpdateUserFormData = CreateUserFormData & {
  isActive: boolean;
};

export interface CreateUserResponse {
  user: StaffUser;
  key: string;
}

export interface ReissueBadgeResponse {
  user: StaffUser;
  key: string;
}

export const STAFF_ROLES: { value: UserRole; label: string }[] = [
  { value: 'doctor', label: 'Doctor' },
  { value: 'nurse', label: 'Nurse' },
  { value: 'staff', label: 'Staff' },
  { value: 'admin', label: 'Administrator' },
];

export const DEFAULT_PIN = '0000';
