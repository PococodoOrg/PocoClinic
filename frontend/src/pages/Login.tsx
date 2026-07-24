import React, { useRef, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import {
  Alert,
  Anchor,
  Button,
  Container,
  Paper,
  PinInput,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useAuth } from '../context/AuthContext';
import { safeInternalPath } from '../utils/safeRedirect';

export default function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const { loginStaff, authError, clearAuthError } = useAuth();
  const pinContainerRef = useRef<HTMLDivElement>(null);
  const submittingRef = useRef(false);

  const [badgeKey, setBadgeKey] = useState('');
  const [pin, setPin] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const from = safeInternalPath((location.state as { from?: string; reason?: string } | null)?.from);
  const sessionExpired = (location.state as { reason?: string } | null)?.reason === 'session-expired';

  const focusPin = () => {
    const input = pinContainerRef.current?.querySelector('input');
    input?.focus();
  };

  const handleBadgeChange = (value: string) => {
    const trimmed = value.trim();
    const wasEmpty = badgeKey.length === 0;
    setBadgeKey(trimmed);
    // Badge scanners dump the full code at once — move to PIN without fighting manual typing.
    if (wasEmpty && trimmed.length >= 8) {
      requestAnimationFrame(focusPin);
    }
  };

  const signIn = async (pinValue: string) => {
    if (!badgeKey || pinValue.length !== 4 || submittingRef.current) {
      return;
    }

    submittingRef.current = true;
    setError(null);
    setIsSubmitting(true);

    try {
      await loginStaff({ key: badgeKey, pin: pinValue });
      navigate(from, { replace: true });
    } catch {
      setError('Invalid badge or PIN. Please try again.');
      setPin('');
      focusPin();
    } finally {
      submittingRef.current = false;
      setIsSubmitting(false);
    }
  };

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    void signIn(pin);
  };

  return (
    <div className="workstation-auth-page">
      <Container size={440} className="workstation-auth-card">
        <Paper radius="md" p={{ base: 'lg', sm: 'xl' }} withBorder component="section" aria-labelledby="staff-login-title">
          <Stack gap="lg">
            <div>
              <Text size="sm" c="dimmed" fw={600} tt="uppercase" mb={4}>
                PocoClinic
              </Text>
              <Title order={2} id="staff-login-title">
                Staff sign in
              </Title>
              <Text c="dimmed" size="md" mt="xs">
                Scan your employee badge, then enter your 4-digit PIN.
              </Text>
            </div>

            {sessionExpired && (
              <Alert color="blue" role="status">
                Your session ended. Sign in again to continue.
              </Alert>
            )}

            {authError && (
              <Alert color="yellow" onClose={clearAuthError} withCloseButton role="status">
                {authError}
              </Alert>
            )}

            {error && (
              <Alert color="red" role="alert">
                {error}
              </Alert>
            )}

            <form onSubmit={handleSubmit}>
              <Stack gap="lg">
                <TextInput
                  label="Badge"
                  description="Scan the QR badge, or paste the code if testing"
                  value={badgeKey}
                  onChange={(event) => handleBadgeChange(event.currentTarget.value)}
                  onKeyDown={(event) => {
                    if (event.key === 'Enter') {
                      event.preventDefault();
                      focusPin();
                    }
                  }}
                  size="lg"
                  className="touch-control"
                  autoFocus
                  autoComplete="username"
                  enterKeyHint="next"
                  required
                />

                <Stack gap="xs" ref={pinContainerRef}>
                  <Text component="label" size="sm" fw={500} htmlFor="staff-pin">
                    PIN
                  </Text>
                  <PinInput
                    id="staff-pin"
                    length={4}
                    type="number"
                    mask
                    value={pin}
                    onChange={setPin}
                    onComplete={(value) => {
                      void signIn(value);
                    }}
                    size="lg"
                    oneTimeCode
                    aria-label="Four-digit PIN"
                    disabled={isSubmitting}
                    styles={{
                      root: { justifyContent: 'space-between' },
                      input: { minHeight: 52, flex: 1 },
                    }}
                  />
                </Stack>

                <Button
                  type="submit"
                  loading={isSubmitting}
                  fullWidth
                  size="lg"
                  className="touch-control"
                  disabled={!badgeKey || pin.length !== 4}
                >
                  Sign in
                </Button>
              </Stack>
            </form>

            <Text size="sm" ta="center">
              Setting up the clinic for the first time?{' '}
              <Anchor component={Link} to="/login/admin" size="sm">
                Administrator sign in
              </Anchor>
            </Text>
          </Stack>
        </Paper>
      </Container>
    </div>
  );
}
