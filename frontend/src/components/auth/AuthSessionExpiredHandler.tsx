import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { AUTH_SESSION_EXPIRED_EVENT } from '../../api/authRefresh';
import { useAuth } from '../../context/AuthContext';

/** Signs out and redirects to login when refresh fails or the API returns 401. */
export function AuthSessionExpiredHandler() {
  const { logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const onSessionExpired = () => {
      void logout().finally(() => {
        navigate('/login', { replace: true, state: { reason: 'session-expired' } });
      });
    };

    window.addEventListener(AUTH_SESSION_EXPIRED_EVENT, onSessionExpired);
    return () => {
      window.removeEventListener(AUTH_SESSION_EXPIRED_EVENT, onSessionExpired);
    };
  }, [logout, navigate]);

  return null;
}
