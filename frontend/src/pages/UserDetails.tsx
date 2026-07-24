import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Badge,
  Button,
  Container,
  Divider,
  Grid,
  Group,
  LoadingOverlay,
  Modal,
  Paper,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { IconArrowLeft, IconEdit, IconRefresh, IconTrash } from '@tabler/icons-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { deleteUser, getUser, reissueBadge, unlockUser } from '../api/users';
import { useAuth } from '../context/AuthContext';
import { UserAuditLog } from '../components/users/UserAuditLog';
import { BadgeIssueModal } from '../components/users/BadgeIssueModal';
import { StaffUser } from '../types/user';

function formatRole(role: string): string {
  return role.charAt(0).toUpperCase() + role.slice(1);
}

function formatDate(value?: string): string {
  if (!value) {
    return 'Never';
  }
  return new Date(value).toLocaleString();
}

export default function UserDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user: currentUser } = useAuth();

  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [reissueConfirmOpen, setReissueConfirmOpen] = useState(false);
  const [issuedUser, setIssuedUser] = useState<StaffUser | null>(null);
  const [issuedKey, setIssuedKey] = useState('');
  const [badgeModalOpen, setBadgeModalOpen] = useState(false);

  const { data: user, isLoading, error } = useQuery({
    queryKey: ['user', id],
    queryFn: () => getUser(id!),
    enabled: Boolean(id),
  });

  const deleteMutation = useMutation({
    mutationFn: () => deleteUser(id!),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['users'] });
      notifications.show({
        title: 'Staff member deleted',
        message: 'The account was removed successfully.',
        color: 'green',
      });
      navigate('/users');
    },
    onError: (mutationError: Error) => {
      notifications.show({
        title: 'Delete failed',
        message: mutationError.message || 'Could not delete staff member',
        color: 'red',
      });
    },
  });

  const unlockMutation = useMutation({
    mutationFn: () => unlockUser(id!),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['user', id] });
      await queryClient.invalidateQueries({ queryKey: ['users'] });
      notifications.show({ color: 'green', message: 'Account unlocked. Staff can sign in again.' });
    },
    onError: (mutationError: Error) => {
      notifications.show({ color: 'red', message: mutationError.message || 'Could not unlock account' });
    },
  });

  const reissueMutation = useMutation({
    mutationFn: () => reissueBadge(id!),
    onSuccess: (response) => {
      setReissueConfirmOpen(false);
      setIssuedUser(response.user);
      setIssuedKey(response.key);
      setBadgeModalOpen(true);
      queryClient.invalidateQueries({ queryKey: ['user-audit', id] });
    },
    onError: (mutationError: Error) => {
      notifications.show({
        title: 'Reissue failed',
        message: mutationError.message || 'Could not reissue badge',
        color: 'red',
      });
    },
  });

  if (error) {
    return (
      <Container>
        <Button
          variant="light"
          leftSection={<IconArrowLeft size={16} />}
          onClick={() => navigate('/users')}
          mb="md"
        >
          Back to staff
        </Button>
        <Text c="red">Error loading staff member.</Text>
      </Container>
    );
  }

  if (!user && !isLoading) {
    return (
      <Container>
        <Text>Staff member not found.</Text>
      </Container>
    );
  }

  const isSelf = currentUser?.id === user?.id;

  return (
    <Container size="lg">
      <Paper radius="md" p="xl" withBorder pos="relative">
        <LoadingOverlay visible={isLoading} />

        {user && (
          <Stack gap="lg">
            <Group justify="space-between" align="flex-start">
              <Group>
                <Button
                  variant="light"
                  leftSection={<IconArrowLeft size={16} />}
                  onClick={() => navigate('/users')}
                >
                  Back to staff
                </Button>
                <div>
                  <Title order={2}>{user.name}</Title>
                  <Text c="dimmed">{user.email}</Text>
                  <Group gap="xs" mt="xs">
                    {!user.isActive && <Badge color="gray">Inactive</Badge>}
                    {user.isLocked && <Badge color="red">Locked</Badge>}
                    {user.mustChangePin && <Badge color="yellow">Default PIN</Badge>}
                  </Group>
                </div>
              </Group>
              <Group>
                <Button
                  variant="light"
                  leftSection={<IconEdit size={16} />}
                  onClick={() => navigate(`/users/${id}/edit`)}
                >
                  Edit
                </Button>
                {user.isLocked && (
                  <Button
                    variant="light"
                    color="orange"
                    loading={unlockMutation.isPending}
                    onClick={() => unlockMutation.mutate()}
                  >
                    Unlock account
                  </Button>
                )}
                <Button
                  variant="light"
                  leftSection={<IconRefresh size={16} />}
                  onClick={() => setReissueConfirmOpen(true)}
                >
                  Reissue badge
                </Button>
                <Button
                  variant="light"
                  color="red"
                  leftSection={<IconTrash size={16} />}
                  disabled={isSelf}
                  onClick={() => setDeleteModalOpen(true)}
                >
                  Delete
                </Button>
              </Group>
            </Group>

            {isSelf && (
              <Text size="sm" c="dimmed">
                You cannot delete your own account from this screen.
              </Text>
            )}

            <Divider />

            <Grid>
              <Grid.Col span={{ base: 12, md: 6 }}>
                <Stack gap="sm">
                  <Title order={4}>Profile</Title>
                  <Group>
                    <Text fw={500}>Role</Text>
                    <Badge variant="light">{formatRole(user.role)}</Badge>
                  </Group>
                  <Group>
                    <Text fw={500}>Last sign in</Text>
                    <Text>{formatDate(user.lastLogin)}</Text>
                  </Group>
                  {user.isLocked && user.lockedUntil && (
                    <Group>
                      <Text fw={500}>Locked until</Text>
                      <Text>{new Date(user.lockedUntil).toLocaleString()}</Text>
                    </Group>
                  )}
                  <Group>
                    <Text fw={500}>Created</Text>
                    <Text>{new Date(user.createdAt).toLocaleString()}</Text>
                  </Group>
                  <Group>
                    <Text fw={500}>Updated</Text>
                    <Text>{new Date(user.updatedAt).toLocaleString()}</Text>
                  </Group>
                </Stack>
              </Grid.Col>
            </Grid>

            <Divider />

            <UserAuditLog userId={user.id} />
          </Stack>
        )}
      </Paper>

      <Modal opened={deleteModalOpen} onClose={() => setDeleteModalOpen(false)} title="Delete staff member?">
        <Stack gap="md">
          <Text>
            Are you sure you want to delete {user?.name}? This removes their account and invalidates
            their badge. This action cannot be undone.
          </Text>
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setDeleteModalOpen(false)}>
              Cancel
            </Button>
            <Button color="red" loading={deleteMutation.isPending} onClick={() => deleteMutation.mutate()}>
              Delete staff member
            </Button>
          </Group>
        </Stack>
      </Modal>

      <Modal
        opened={reissueConfirmOpen}
        onClose={() => setReissueConfirmOpen(false)}
        title="Reissue employee badge?"
      >
        <Stack gap="md">
          <Text>
            This invalidates {user?.name}&apos;s current badge immediately. You will need to print a
            replacement badge with the new QR code. Their PIN will stay the same.
          </Text>
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setReissueConfirmOpen(false)}>
              Cancel
            </Button>
            <Button color="orange" loading={reissueMutation.isPending} onClick={() => reissueMutation.mutate()}>
              Reissue badge
            </Button>
          </Group>
        </Stack>
      </Modal>

      <BadgeIssueModal
        opened={badgeModalOpen}
        user={issuedUser}
        badgeKey={issuedKey}
        mode="reissue"
        onClose={() => {
          setBadgeModalOpen(false);
          setIssuedUser(null);
          setIssuedKey('');
        }}
      />
    </Container>
  );
}
