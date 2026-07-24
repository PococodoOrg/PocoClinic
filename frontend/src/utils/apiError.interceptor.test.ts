import axios, { AxiosError, AxiosHeaders } from 'axios';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { ValidationError, attachSanitizedErrorInterceptor, getErrorMessage } from './apiError';

describe('attachSanitizedErrorInterceptor', () => {
  afterEach(() => {
    vi.unstubAllEnvs();
  });

  function clientWithInterceptor() {
    const client = axios.create();
    attachSanitizedErrorInterceptor(client);
    return client;
  }

  it('converts validation payloads into ValidationError', async () => {
    const client = clientWithInterceptor();
    client.defaults.adapter = async () => {
      throw new AxiosError(
        'Request failed',
        'ERR_BAD_REQUEST',
        { headers: new AxiosHeaders() },
        {},
        {
          status: 400,
          statusText: 'Bad Request',
          headers: {},
          config: { headers: new AxiosHeaders() },
          data: {
            code: 'VALIDATION_ERROR',
            message: 'Invalid fields',
            errors: { name: ['required'] },
          },
        },
      );
    };

    await expect(client.get('/x')).rejects.toMatchObject({
      name: 'ValidationError',
      message: 'Invalid fields',
      code: 'VALIDATION_ERROR',
      errors: { name: ['required'] },
    });
  });

  it('hides raw API messages in production builds', async () => {
    vi.stubEnv('DEV', false);
    const client = clientWithInterceptor();
    client.defaults.adapter = async () => {
      throw new AxiosError(
        'Request failed',
        'ERR_BAD_RESPONSE',
        { headers: new AxiosHeaders() },
        {},
        {
          status: 500,
          statusText: 'Server Error',
          headers: {},
          config: { headers: new AxiosHeaders() },
          data: {
            code: 'INTERNAL_ERROR',
            message: 'SQLSTATE 42703 column group_id does not exist',
          },
        },
      );
    };

    try {
      await client.get('/x');
      expect.fail('expected rejection');
    } catch (error) {
      expect(error).toBeInstanceOf(Error);
      expect(error).not.toBeInstanceOf(ValidationError);
      expect((error as Error).message).toBe('Something went wrong. Please try again.');
      expect(getErrorMessage(error)).toBe('Something went wrong. Please try again.');
    }
  });

  it('maps transport failures to a network error', async () => {
    const client = clientWithInterceptor();
    client.defaults.adapter = async () => {
      const error = new AxiosError('Network Error');
      error.request = {};
      throw error;
    };

    await expect(client.get('/x')).rejects.toThrow('Network error - no response received');
  });
});
