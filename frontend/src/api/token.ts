import { config } from '../config';

/** Removes legacy tokens stored before cookie-based auth. */
export const clearLegacyStoredToken = (): void => {
  localStorage.removeItem(config.auth.tokenKey);
  localStorage.removeItem(config.auth.refreshTokenKey);
};
