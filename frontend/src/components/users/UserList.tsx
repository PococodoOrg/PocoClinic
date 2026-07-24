import { useState } from 'react';
import { Badge, Button, Group, Table, Text, TextInput } from '@mantine/core';
import { IconEye, IconSearch } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchUsers } from '../../api/users';
import { StaffUser } from '../../types/user';

function formatRole(role: string): string {
  return role.charAt(0).toUpperCase() + role.slice(1);
}

function formatDate(value?: string): string {
  if (!value) {
    return 'Never';
  }
  return new Date(value).toLocaleString();
}

export function UserList() {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const { data, isLoading, error } = useQuery({
    queryKey: ['users', search, page],
    queryFn: () => fetchUsers({ search, page, pageSize }),
  });

  if (error) {
    return <Text c="red">Error loading staff. Please try again later.</Text>;
  }

  const handleViewUser = (userId: string) => {
    navigate(`/users/${userId}`);
  };

  const rows =
    data?.users?.map((user: StaffUser) => (
      <Table.Tr key={user.id} onDoubleClick={() => handleViewUser(user.id)}>
        <Table.Td>{user.name}</Table.Td>
        <Table.Td>{user.email}</Table.Td>
        <Table.Td>
          <Group gap="xs">
            {formatRole(user.role)}
            {!user.isActive && <Badge size="xs" color="gray">Inactive</Badge>}
            {user.isLocked && <Badge size="xs" color="red">Locked</Badge>}
          </Group>
        </Table.Td>
        <Table.Td>{formatDate(user.lastLogin)}</Table.Td>
        <Table.Td>{new Date(user.createdAt).toLocaleDateString()}</Table.Td>
        <Table.Td>
          <Button variant="light" size="xs" leftSection={<IconEye size={14} />} onClick={() => handleViewUser(user.id)}>
            Manage
          </Button>
        </Table.Td>
      </Table.Tr>
    )) ?? [];

  return (
    <div>
      <Group justify="space-between" mb="md">
        <TextInput
          placeholder="Search staff..."
          leftSection={<IconSearch size="1rem" />}
          value={search}
          onChange={(event) => {
            setSearch(event.currentTarget.value);
            setPage(1);
          }}
        />
        <Button onClick={() => navigate('/users/new')}>Add staff member</Button>
      </Group>

      {isLoading ? (
        <Text c="dimmed" ta="center" mt="xl">
          Loading staff...
        </Text>
      ) : rows.length === 0 ? (
        <Text c="dimmed" ta="center" mt="xl">
          No staff members found. Create the first account to print badges.
        </Text>
      ) : (
        <Table>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Name</Table.Th>
              <Table.Th>Email</Table.Th>
              <Table.Th>Role</Table.Th>
              <Table.Th>Last login</Table.Th>
              <Table.Th>Created</Table.Th>
              <Table.Th>Actions</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>{rows}</Table.Tbody>
        </Table>
      )}

      {data && data.totalPages > 1 && (
        <Group justify="center" mt="xl">
          <Button variant="outline" disabled={page <= 1} onClick={() => setPage(page - 1)}>
            Previous
          </Button>
          <Text>
            Page {page} of {data.totalPages}
          </Text>
          <Button variant="outline" disabled={page >= data.totalPages} onClick={() => setPage(page + 1)}>
            Next
          </Button>
        </Group>
      )}
    </div>
  );
}
