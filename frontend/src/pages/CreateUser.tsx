import { useState } from 'react';
import { Anchor, Container, Stack, Text, Title } from '@mantine/core';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Link, useNavigate } from 'react-router-dom';
import { notifications } from '@mantine/notifications';
import { createUser } from '../api/users';
import { UserForm } from '../components/users/UserForm';
import { BadgeIssueModal } from '../components/users/BadgeIssueModal';
import { CreateUserFormData, StaffUser, UpdateUserFormData } from '../types/user';
import { getErrorMessage } from '../utils/apiError';

export default function CreateUser() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [issuedUser, setIssuedUser] = useState<StaffUser | null>(null);
  const [issuedKey, setIssuedKey] = useState('');
  const [badgeModalOpen, setBadgeModalOpen] = useState(false);

  const createUserMutation = useMutation({
    mutationFn: createUser,
    onSuccess: (response) => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setIssuedUser(response.user);
      setIssuedKey(response.key);
      setBadgeModalOpen(true);
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Error',
        message: getErrorMessage(error, 'Failed to create staff member'),
        color: 'red',
      });
    },
  });

  const handleSubmit = async (data: UpdateUserFormData) => {
    const payload: CreateUserFormData = {
      email: data.email,
      name: data.name,
      role: data.role,
    };
    await createUserMutation.mutateAsync(payload);
  };

  const handleBadgeModalClose = () => {
    setBadgeModalOpen(false);
    setIssuedUser(null);
    setIssuedKey('');
    navigate('/users');
  };

  return (
    <Container size="sm">
      <Stack gap="md" mb="xl">
        <Title order={2}>Add staff member</Title>
        <Text c="dimmed">
          A badge key and default PIN will be generated. You will print the badge immediately after
          creation.
        </Text>
        <Text size="sm">
          <Anchor component={Link} to="/users">
            Back to staff list
          </Anchor>
        </Text>
      </Stack>

      <UserForm onSubmit={handleSubmit} isLoading={createUserMutation.isPending} />

      <BadgeIssueModal
        opened={badgeModalOpen}
        user={issuedUser}
        badgeKey={issuedKey}
        onClose={handleBadgeModalClose}
      />
    </Container>
  );
}
