import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';

import { API_BASE_URL } from './baseUrl';

import { clearLegacyStoredToken } from './token';



export interface RefreshResponse {

  user: {

    id: string;

    name: string;

    email: string;

    role: string;

    lastLogin?: string;

    mustChangePin?: boolean;

  };

}



export const AUTH_SESSION_EXPIRED_EVENT = 'auth:session-expired';



const notifySessionExpired = (): void => {

  window.dispatchEvent(new CustomEvent(AUTH_SESSION_EXPIRED_EVENT));

};



const refreshClient = axios.create({

  baseURL: `${API_BASE_URL}/auth`,

  headers: { 'Content-Type': 'application/json' },

  withCredentials: true,

  timeout: 10_000,

});



export const refreshAccessToken = async (): Promise<RefreshResponse | null> => {

  try {

    const response = await refreshClient.post<RefreshResponse>('/refresh');

    return response.data;

  } catch {

    clearLegacyStoredToken();

    notifySessionExpired();

    return null;

  }

};



let refreshPromise: Promise<RefreshResponse | null> | null = null;



export const refreshAccessTokenOnce = async (): Promise<RefreshResponse | null> => {

  if (!refreshPromise) {

    refreshPromise = refreshAccessToken().finally(() => {

      refreshPromise = null;

    });

  }

  return refreshPromise;

};



type RetriableConfig = InternalAxiosRequestConfig & { _retry?: boolean };



export const attachAuthInterceptor = (client: ReturnType<typeof axios.create>): void => {

  client.interceptors.response.use(

    (response) => response,

    async (error: AxiosError) => {

      const originalRequest = error.config as RetriableConfig | undefined;

      if (!originalRequest || error.response?.status !== 401 || originalRequest._retry) {

        if (error.response?.status === 401) {

          clearLegacyStoredToken();

          notifySessionExpired();

        }

        return Promise.reject(error);

      }



      originalRequest._retry = true;

      const refreshed = await refreshAccessTokenOnce();

      if (!refreshed) {

        return Promise.reject(error);

      }



      return client(originalRequest);

    },

  );

};


