import { AxiosError, AxiosInstance } from 'axios';

const GENERIC_MESSAGE = 'Something went wrong. Please try again.';

interface ApiErrorBody {
  message?: string;
  code?: string;
  errors?: Record<string, string[]>;
}

export class ValidationError extends Error {
  code: string;
  errors?: Record<string, string[]>;

  constructor(message: string, code: string, errors?: Record<string, string[]>) {
    super(message);
    this.name = 'ValidationError';
    this.code = code;
    this.errors = errors;
  }
}

function parseApiErrorBody(data: unknown): ApiErrorBody {
  if (typeof data !== 'object' || data === null) {
    return {};
  }
  const body = data as Record<string, unknown>;
  const message =
    typeof body.message === 'string'
      ? body.message
      : typeof body.error === 'string'
        ? body.error
        : undefined;
  return {
    message,
    code: typeof body.code === 'string' ? body.code : undefined,
    errors: body.errors as Record<string, string[]> | undefined,
  };
}

/** Returns a user-safe error message; detailed text is dev-only. */
export function userFacingApiError(error: unknown, fallback = GENERIC_MESSAGE): string {
  if (import.meta.env.DEV && error instanceof Error && error.message) {
    return error.message;
  }

  const body = parseApiErrorBody(error);
  if (import.meta.env.DEV && body.message?.trim()) {
    return body.message;
  }

  return fallback;
}

/** Reads a message from a caught error (sanitized when thrown by the API client). */
export function getErrorMessage(error: unknown, fallback = GENERIC_MESSAGE): string {
  if (error instanceof ValidationError) {
    return error.message;
  }
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return fallback;
}

export function logApiError(context: string, error: unknown): void {
  if (import.meta.env.DEV) {
    console.error(context, error);
  }
}

/** Sanitizes axios error responses before they reach UI code. */
export function attachSanitizedErrorInterceptor(client: AxiosInstance): void {
  client.interceptors.response.use(
    (response) => response,
    (error: AxiosError) => {
      if (error.response) {
        logApiError('API Error', error.response.data);
        const body = parseApiErrorBody(error.response.data);
        if (body.code === 'VALIDATION_ERROR' && body.errors) {
          return Promise.reject(
            new ValidationError(body.message || 'Validation failed', body.code, body.errors),
          );
        }
        return Promise.reject(new Error(userFacingApiError(body)));
      }
      if (error.request) {
        logApiError('Network Error', error.request);
        return Promise.reject(new Error('Network error - no response received'));
      }
      logApiError('Request Error', error.message);
      return Promise.reject(new Error('Error setting up the request'));
    },
  );
}
