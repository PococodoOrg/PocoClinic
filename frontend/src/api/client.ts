import axios, { AxiosInstance } from 'axios';
import { API_BASE_URL } from './baseUrl';
import { attachAuthInterceptor } from './authRefresh';
import { attachSanitizedErrorInterceptor } from '../utils/apiError';

export { API_BASE_URL } from './baseUrl';

export interface ApiClientOptions {
  /** Path appended to the API base URL, e.g. `/auth`. */
  path?: string;
  timeout?: number;
  /** When false, omit the default JSON Content-Type header (e.g. multipart uploads). */
  json?: boolean;
}

export function createApiClient(options: ApiClientOptions = {}): AxiosInstance {
  const { path = '', timeout = 10_000, json = true } = options;

  const client = axios.create({
    baseURL: `${API_BASE_URL}${path}`,
    ...(json ? { headers: { 'Content-Type': 'application/json' } } : {}),
    withCredentials: true,
    timeout,
  });

  attachAuthInterceptor(client);
  attachSanitizedErrorInterceptor(client);

  return client;
}
