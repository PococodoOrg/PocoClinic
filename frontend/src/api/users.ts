import { createApiClient } from './client';
import {
  CreateUserFormData,
  CreateUserResponse,
  PaginatedUsers,
  ReissueBadgeResponse,
  StaffUser,
  UpdateUserFormData,
} from '../types/user';
import { PaginatedAuditEntries } from '../types/audit';

const usersApi = createApiClient({ path: '/auth' });

export const fetchUsers = async (params?: {
  page?: number;
  pageSize?: number;
  search?: string;
}): Promise<PaginatedUsers> => {
  const response = await usersApi.get<PaginatedUsers>('/users', {
    params: {
      page: params?.page ?? 1,
      pageSize: params?.pageSize ?? 10,
      search: params?.search,
    },
  });
  return response.data;
};

export const getUser = async (id: string): Promise<StaffUser> => {
  const response = await usersApi.get<StaffUser>(`/users/${id}`);
  return response.data;
};

export const createUser = async (data: CreateUserFormData): Promise<CreateUserResponse> => {
  const response = await usersApi.post<CreateUserResponse>('/register', data);
  return response.data;
};

export const reissueBadge = async (userId: string): Promise<ReissueBadgeResponse> => {
  const response = await usersApi.post<ReissueBadgeResponse>(`/users/${userId}/reissue-badge`);
  return response.data;
};

export const unlockUser = async (userId: string): Promise<StaffUser> => {
  const response = await usersApi.post<StaffUser>(`/users/${userId}/unlock`);
  return response.data;
};

export const updateUser = async (userId: string, data: UpdateUserFormData): Promise<StaffUser> => {
  const response = await usersApi.put<StaffUser>(`/users/${userId}`, data);
  return response.data;
};

export const deleteUser = async (userId: string): Promise<void> => {
  await usersApi.delete(`/users/${userId}`);
};

export const fetchUserAuditLogs = async (
  userId: string,
  params?: { page?: number; pageSize?: number },
): Promise<PaginatedAuditEntries> => {
  const response = await usersApi.get<PaginatedAuditEntries>(`/users/${userId}/audit-logs`, {
    params: {
      page: params?.page ?? 1,
      pageSize: params?.pageSize ?? 20,
    },
  });
  return response.data;
};
