import React, { useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import {
  Alert,
  Anchor,
  Button,
  Container,
  Paper,
  PasswordInput,
  PinInput,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useAuth } from '../context/AuthContext';
import { safeInternalPath } from '../utils/safeRedirect';

export default function AdminLogin() {
  const navigate = useNavigate();
  const location = useLocation();
  const { loginAdmin } = useAuth();

  const [email, setEmail] = useState(import.meta.env.DEV ? 'admin@pococlinic.local' : '');
  const [key, setKey] = useState('');
  const [pin, setPin] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const from = safeInternalPath((location.state as { from?: string } | null)?.from);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (pin.length !== 4) {
      return;
    }

    setError(null);
    setIsSubmitting(true);

    try {
      await loginAdmin({ email, key, pin });
      navigate(from, { replace: true });
    } catch {
      setError('Invalid email, access key, or PIN. Please try again.');
      setPin('');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="workstation-auth-page">
      <Container size={440} className="workstation-auth-card">
        <Paper radius="md" p={{ base: 'lg', sm: 'xl' }} withBorder component="section" aria-labelledby="admin-login-title">
          <Stack gap="lg">
            <div>
              <Text size="sm" c="dimmed" fw={600} tt="uppercase" mb={4}>
                PocoClinic
              </Text>
              <Title order={2} id="admin-login-title">
                Administrator sign in
              </Title>
              <Text c="dimmed" size="md" mt="xs">
                Use this path for first-time setup, creating staff accounts, and printing badges.
                Daily staff should use badge sign-in instead.
              </Text>
            </div>

            {error && (
              <Alert color="red" role="alert">
                {error}
              </Alert>
            )}

            <form onSubmit={handleSubmit}>
              <Stack gap="lg">
                <TextInput
                  label="Email"
                  type="email"
                  value={email}
                  onChange={(event) => setEmail(event.currentTarget.value)}
                  size="lg"
                  className="touch-control"
                  autoComplete="username"
                  required
                />
                <PasswordInput
                  label="Access key"
                  description="From bootstrap-admin-once.txt on the server (Pi: /var/lib/pococlinic/), or your safe vault sheet"
                  value={key}
                  onChange={(event) => setKey(event.currentTarget.value)}
                  size="lg"
                  className="touch-control"
                  autoComplete="current-password"
                  required
                />
                <Stack gap="xs">
                  <Text component="label" size="sm" fw={500}>
                    PIN
                  </Text>
                  <PinInput
                    length={4}
                    type="number"
                    mask
                    value={pin}
                    onChange={setPin}
                    size="lg"
                    oneTimeCode
                    aria-label="Four-digit administrator PIN"
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
                  disabled={!email || !key || pin.length !== 4}
                >
                  Sign in
                </Button>
              </Stack>
            </form>

            {import.meta.env.DEV && (
              <Text size="xs" c="dimmed">
                Default dev account: admin@pococlinic.local with PIN 0000. The access key is printed in
                the backend server logs on first startup.
              </Text>
            )}

            <Text size="sm" ta="center">
              <Anchor component={Link} to="/login">
                Back to staff badge sign in
              </Anchor>
            </Text>
          </Stack>
        </Paper>
      </Container>
    </div>
  );
}
