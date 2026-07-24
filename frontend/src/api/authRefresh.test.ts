import { beforeEach, describe, expect, it, vi } from 'vitest';

const { post } = vi.hoisted(() => ({
  post: vi.fn(),
}));

vi.mock('axios', () => ({
  default: {
    create: () => ({
      post,
      interceptors: { response: { use: vi.fn() } },
    }),
  },
}));

vi.mock('./baseUrl', () => ({ API_BASE_URL: 'http://127.0.0.1:8080/api/v1' }));
vi.mock('./token', () => ({ clearLegacyStoredToken: vi.fn() }));

import {
  AUTH_SESSION_EXPIRED_EVENT,
  refreshAccessToken,
  refreshAccessTokenOnce,
} from './authRefresh';

describe('authRefresh', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('refreshAccessToken returns user payload on success', async () => {
    post.mockResolvedValueOnce({
      data: { user: { id: '1', email: 'a@b.c', name: 'A', role: 'staff' } },
    });

    const result = await refreshAccessToken();
    expect(result?.user.email).toBe('a@b.c');
    expect(post).toHaveBeenCalledWith('/refresh');
  });

  it('refreshAccessTokenOnce deduplicates concurrent calls', async () => {
    let resolvePost: (value: unknown) => void = () => {};
    post.mockReturnValue(
      new Promise((resolve) => {
        resolvePost = resolve;
      }),
    );

    const first = refreshAccessTokenOnce();
    const second = refreshAccessTokenOnce();
    resolvePost({ data: { user: { id: '1', email: 'x@y.z', name: 'X', role: 'staff' } } });

    const [a, b] = await Promise.all([first, second]);
    expect(a?.user.email).toBe('x@y.z');
    expect(b?.user.email).toBe('x@y.z');
    expect(post).toHaveBeenCalledTimes(1);
  });

  it('dispatches session expired event when refresh fails', async () => {
    const listener = vi.fn();
    window.addEventListener(AUTH_SESSION_EXPIRED_EVENT, listener);
    post.mockRejectedValueOnce(new Error('401'));

    const result = await refreshAccessToken();
    expect(result).toBeNull();
    expect(listener).toHaveBeenCalled();

    window.removeEventListener(AUTH_SESSION_EXPIRED_EVENT, listener);
  });
});
