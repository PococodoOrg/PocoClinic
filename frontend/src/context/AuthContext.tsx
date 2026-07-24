import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';

import {

  AdminLoginRequest,

  AuthUser,

  StaffLoginRequest,

  clearLegacyStoredToken,

  getCurrentUser,

  loginAdmin,

  loginStaff,

  logoutRequest,

  refreshAccessToken,

} from '../api/auth';
import { queryClient } from '../queryClient';



interface AuthContextValue {

  user: AuthUser | null;

  isAuthenticated: boolean;

  isLoading: boolean;

  authError: string | null;

  clearAuthError: () => void;

  loginStaff: (credentials: StaffLoginRequest) => Promise<void>;

  loginAdmin: (credentials: AdminLoginRequest) => Promise<void>;

  logout: () => Promise<void>;

  refreshUser: () => Promise<void>;

  setUser: (user: AuthUser | null) => void;

}



const AuthContext = createContext<AuthContextValue | undefined>(undefined);



async function restoreAuthenticatedUser(): Promise<AuthUser | null> {

  try {

    return await getCurrentUser();

  } catch {

    const refreshed = await refreshAccessToken();

    if (!refreshed) {

      return null;

    }

    return refreshed.user;

  }

}



export function AuthProvider({ children }: { children: React.ReactNode }) {

  const [user, setUser] = useState<AuthUser | null>(null);

  const [isLoading, setIsLoading] = useState(true);

  const [authError, setAuthError] = useState<string | null>(null);



  useEffect(() => {

    clearLegacyStoredToken();



    const restoreSession = async () => {

      try {

        const currentUser = await Promise.race([

          restoreAuthenticatedUser(),

          new Promise<never>((_, reject) => {

            window.setTimeout(() => reject(new Error('session restore timed out')), 8_000);

          }),

        ]);

        setUser(currentUser);

        setAuthError(null);

      } catch {

        clearLegacyStoredToken();

        setUser(null);

        setAuthError(
          import.meta.env.DEV
            ? 'Could not reach the clinic server. Start the backend on port 8080, then refresh.'
            : 'Could not reach the clinic server. Check your connection and refresh the page.',
        );

      } finally {

        setIsLoading(false);

      }

    };



    restoreSession();

  }, []);



  const handleLoginStaff = useCallback(async (credentials: StaffLoginRequest) => {

    const response = await loginStaff(credentials);

    setUser(response.user);

  }, []);



  const handleLoginAdmin = useCallback(async (credentials: AdminLoginRequest) => {

    const response = await loginAdmin(credentials);

    setUser(response.user);

  }, []);



  const logout = useCallback(async () => {

    await logoutRequest();

    queryClient.clear();

    setUser(null);

  }, []);



  const refreshUser = useCallback(async () => {

    const currentUser = await getCurrentUser();

    setUser(currentUser);

  }, []);



  const value = useMemo(

    () => ({

      user,

      isAuthenticated: user !== null,

      isLoading,

      authError,

      clearAuthError: () => setAuthError(null),

      loginStaff: handleLoginStaff,

      loginAdmin: handleLoginAdmin,

      logout,

      refreshUser,

      setUser,

    }),

    [user, isLoading, authError, handleLoginStaff, handleLoginAdmin, logout, refreshUser],

  );



  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;

}



export function useAuth(): AuthContextValue {

  const context = useContext(AuthContext);

  if (!context) {

    throw new Error('useAuth must be used within an AuthProvider');

  }

  return context;

}

