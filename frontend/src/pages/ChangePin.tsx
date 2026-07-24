import { Alert, Container, Stack, Text, Title } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { useMutation } from '@tanstack/react-query';
import { useLocation, useNavigate } from 'react-router-dom';
import { changePin } from '../api/auth';
import { ChangePinForm } from '../components/account/ChangePinForm';
import { useAuth } from '../context/AuthContext';
import { safeInternalPath } from '../utils/safeRedirect';
import { getErrorMessage } from '../utils/apiError';

export default function ChangePin() {
  const { user, setUser } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const redirectTo = safeInternalPath((location.state as { from?: string } | null)?.from);

  const changePinMutation = useMutation({
    mutationFn: changePin,
    onSuccess: (updatedUser) => {
      setUser(updatedUser);
      notifications.show({
        title: 'PIN updated',
        message: 'Your PIN has been changed successfully.',
        color: 'green',
      });
      if (!updatedUser.mustChangePin) {
        navigate(redirectTo, { replace: true });
      }
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Could not change PIN',
        message: getErrorMessage(error, 'Please check your current PIN and try again.'),
        color: 'red',
      });
    },
  });

  const handleSubmit = async (data: { currentPin: string; newPin: string }) => {
    await changePinMutation.mutateAsync(data);
  };

  return (
    <Container size="sm">
      <Stack gap="md" mb="xl">
        <Title order={2}>Change PIN</Title>
        {user?.mustChangePin && (
          <Alert color="yellow" title="Choose your personal PIN">
            New staff accounts start with the default PIN. Set a private 4-digit PIN before opening patient charts.
          </Alert>
        )}
        <Text c="dimmed">          Any employee can update their own 4-digit PIN here. You will enter this PIN after scanning
          your badge at staff sign-in.
        </Text>
      </Stack>

      <ChangePinForm onSubmit={handleSubmit} isLoading={changePinMutation.isPending} />
    </Container>
  );
}
