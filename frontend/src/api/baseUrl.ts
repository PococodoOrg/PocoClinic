import { config } from '../config';

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || `${config.apiUrl}/api/v1`;
