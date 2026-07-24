import { useParams, useNavigate } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Anchor, Container, Group, Paper, Stack, Text, Title } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { Link } from 'react-router-dom';
import { getUser, updateUser } from '../api/users';
import { UserForm } from '../components/users/UserForm';
import { UpdateUserFormData } from '../types/user';

export default function EditUser() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: user, isLoading, error } = useQuery({
    queryKey: ['user', id],
    queryFn: () => getUser(id!),
    enabled: Boolean(id),
  });

  const updateMutation = useMutation({
    mutationFn: (data: UpdateUserFormData) => updateUser(id!, data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['user', id] });
      await queryClient.invalidateQueries({ queryKey: ['users'] });
      await queryClient.invalidateQueries({ queryKey: ['user-audit', id] });
      notifications.show({
        title: 'Staff updated',
        message: 'Profile changes were saved.',
        color: 'green',
      });
      navigate(`/users/${id}`);
    },
    onError: (mutationError: Error) => {
      notifications.show({
        title: 'Update failed',
        message: mutationError.message || 'Could not update staff member',
        color: 'red',
      });
    },
  });

  const handleSubmit = async (data: UpdateUserFormData) => {
    await updateMutation.mutateAsync(data);
  };

  if (isLoading) {
    return (
      <Container size="sm">
        <Text>Loading...</Text>
      </Container>
    );
  }

  if (error || !user) {
    return (
      <Container size="sm">
        <Text c="red">Could not load staff member.</Text>
      </Container>
    );
  }

  return (
    <Container size="sm">
      <Stack gap="md" mb="xl">
        <Group justify="space-between">
          <Title order={2}>Edit staff member</Title>
        </Group>
        <Text size="sm">
          <Anchor component={Link} to={`/users/${id}`}>
            Back to {user.name}
          </Anchor>
        </Text>
      </Stack>

      <Paper p="md" withBorder>
        <UserForm
          initialValues={user}
          includeActiveToggle
          onSubmit={handleSubmit}
          isLoading={updateMutation.isPending}
          submitLabel="Save changes"
        />
      </Paper>
    </Container>
  );
}
