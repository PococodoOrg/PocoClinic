import { Button, Group, Modal, Stack, Text } from '@mantine/core';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { refreshAccessTokenOnce } from '../../api/auth';
import { useAuth } from '../../context/AuthContext';

const SESSION_MS = 15 * 60 * 1000;
const WARNING_MS = 13 * 60 * 1000;

export function SessionTimeoutManager() {
  const { isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [warningOpen, setWarningOpen] = useState(false);
  const lastActivityRef = useRef(Date.now());
  const warningShownRef = useRef(false);

  const markActivity = useCallback(() => {
    lastActivityRef.current = Date.now();
    if (warningShownRef.current) {
      warningShownRef.current = false;
      setWarningOpen(false);
    }
  }, []);

  const expireSession = useCallback(async () => {
    setWarningOpen(false);
    await logout();
    navigate('/login', { replace: true, state: { reason: 'session-expired' } });
  }, [logout, navigate]);

  const extendSession = useCallback(async () => {
    markActivity();
    const refreshed = await refreshAccessTokenOnce();
    if (!refreshed) {
      await expireSession();
    }
  }, [expireSession, markActivity]);

  useEffect(() => {
    if (!isAuthenticated) {
      setWarningOpen(false);
      warningShownRef.current = false;
      return;
    }

    lastActivityRef.current = Date.now();

    const activityEvents = ['mousedown', 'keydown', 'touchstart', 'scroll'] as const;
    activityEvents.forEach((event) => window.addEventListener(event, markActivity, { passive: true }));

    const interval = window.setInterval(() => {
      const idleMs = Date.now() - lastActivityRef.current;

      if (idleMs >= SESSION_MS) {
        void expireSession();
        return;
      }

      if (idleMs >= WARNING_MS && !warningShownRef.current) {
        warningShownRef.current = true;
        setWarningOpen(true);
      }
    }, 15_000);

    return () => {
      activityEvents.forEach((event) => window.removeEventListener(event, markActivity));
      window.clearInterval(interval);
    };
  }, [expireSession, isAuthenticated, markActivity]);

  if (!isAuthenticated) {
    return null;
  }

  return (
    <Modal
      opened={warningOpen}
      onClose={() => setWarningOpen(false)}
      title="Still there?"
      centered
      closeOnClickOutside={false}
      withCloseButton={false}
    >
      <Stack gap="md">
        <Text size="sm">
          Your session will sign out soon after 15 minutes of inactivity. Click stay signed in to keep working,
          or sign out now if you are finished.
        </Text>
        <Group justify="flex-end">
          <Button variant="default" onClick={() => void expireSession()}>
            Sign out
          </Button>
          <Button onClick={() => void extendSession()}>
            Stay signed in
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
