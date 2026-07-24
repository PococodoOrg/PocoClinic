import { createApiClient } from './client';
import { clearLegacyStoredToken } from './token';

export { attachAuthInterceptor, refreshAccessToken, refreshAccessTokenOnce } from './authRefresh';
export type { RefreshResponse } from './authRefresh';
export { clearLegacyStoredToken } from './token';

export interface AuthUser {
  id: string;
  name: string;
  email: string;
  role: string;
  lastLogin?: string;
  mustChangePin?: boolean;
}

export interface StaffLoginRequest {
  key: string;
  pin: string;
}

export interface AdminLoginRequest {
  email: string;
  key: string;
  pin: string;
}

export interface LoginResponse {
  user: AuthUser;
}

const authApi = createApiClient({ path: '/auth' });

export const loginStaff = async (credentials: StaffLoginRequest): Promise<LoginResponse> => {
  const response = await authApi.post<LoginResponse>('/login/staff', credentials);
  return response.data;
};

export const loginAdmin = async (credentials: AdminLoginRequest): Promise<LoginResponse> => {
  const response = await authApi.post<LoginResponse>('/login/admin', credentials);
  return response.data;
};

export const getCurrentUser = async (): Promise<AuthUser> => {
  const response = await authApi.get<AuthUser>('/me');
  return response.data;
};

export const logoutRequest = async (): Promise<void> => {
  try {
    await authApi.post('/logout');
  } finally {
    clearLegacyStoredToken();
  }
};

export interface ChangePINRequest {
  currentPin: string;
  newPin: string;
}

export const changePin = async (payload: ChangePINRequest): Promise<AuthUser> => {
  const response = await authApi.post<AuthUser>('/change-pin', payload);
  return response.data;
};
