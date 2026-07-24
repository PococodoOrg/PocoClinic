import { useState } from 'react';
import {
  Badge,
  Button,
  Container,
  Group,
  Modal,
  NumberInput,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import {
  createFormGroup,
  deleteFormGroup,
  deleteFormTemplate,
  fetchFormGroups,
  fetchFormTemplates,
  updateFormGroup,
} from '../api/forms';
import { FormGroup, FormGroupFormData, FormTemplate, groupTemplatesByGroup } from '../types/form';

export default function FormTemplates() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [groupModalOpen, setGroupModalOpen] = useState(false);
  const [editingGroup, setEditingGroup] = useState<FormGroup | null>(null);
  const [groupForm, setGroupForm] = useState<FormGroupFormData>({ name: '', sortOrder: 0 });

  const { data: groups = [], isLoading: groupsLoading } = useQuery({
    queryKey: ['form-groups'],
    queryFn: fetchFormGroups,
  });

  const { data: templates = [], isLoading: templatesLoading, error } = useQuery({
    queryKey: ['form-templates'],
    queryFn: fetchFormTemplates,
  });

  const grouped = groupTemplatesByGroup(groups, templates);

  const deleteMutation = useMutation({
    mutationFn: deleteFormTemplate,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['form-templates'] });
      notifications.show({
        title: 'Form deleted',
        message: 'The form template was removed.',
        color: 'green',
      });
    },
    onError: (mutationError: Error) => {
      notifications.show({
        title: 'Delete failed',
        message: mutationError.message,
        color: 'red',
      });
    },
  });

  const saveGroupMutation = useMutation({
    mutationFn: async () => {
      if (editingGroup) {
        return updateFormGroup(editingGroup.id, groupForm);
      }
      return createFormGroup(groupForm);
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['form-groups'] });
      notifications.show({
        title: editingGroup ? 'Group updated' : 'Group created',
        message: 'Form groups organize how staff see forms on patient charts.',
        color: 'green',
      });
      setGroupModalOpen(false);
      setEditingGroup(null);
      setGroupForm({ name: '', sortOrder: 0 });
    },
    onError: (mutationError: Error) => {
      notifications.show({
        title: 'Group save failed',
        message: mutationError.message,
        color: 'red',
      });
    },
  });

  const deleteGroupMutation = useMutation({
    mutationFn: deleteFormGroup,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['form-groups'] });
      await queryClient.invalidateQueries({ queryKey: ['form-templates'] });
      notifications.show({
        title: 'Group deleted',
        message: 'Templates in that group were moved to General.',
        color: 'green',
      });
    },
    onError: (mutationError: Error) => {
      notifications.show({
        title: 'Delete failed',
        message: mutationError.message,
        color: 'red',
      });
    },
  });

  const openCreateGroup = () => {
    setEditingGroup(null);
    setGroupForm({ name: '', sortOrder: groups.length });
    setGroupModalOpen(true);
  };

  const openEditGroup = (group: FormGroup) => {
    setEditingGroup(group);
    setGroupForm({ name: group.name, sortOrder: group.sortOrder });
    setGroupModalOpen(true);
  };

  if (error) {
    return <Text c="red">Could not load forms.</Text>;
  }

  const isLoading = groupsLoading || templatesLoading;

  return (
    <Container size="lg">
      <Stack gap="md" mb="xl">
        <Group justify="space-between">
          <div>
            <Title order={2}>Forms</Title>
            <Text c="dimmed">
              Create form templates, organize them into groups, and choose singleton or log behavior.
            </Text>
          </div>
          <Group>
            <Button variant="default" onClick={() => navigate('/forms/reports')}>
              Reports
            </Button>
            <Button variant="light" onClick={openCreateGroup}>
              New group
            </Button>
            <Button onClick={() => navigate('/forms/new')}>Create form</Button>
          </Group>
        </Group>
      </Stack>

      {isLoading ? (
        <Text c="dimmed">Loading forms...</Text>
      ) : (
        grouped.map(({ group, templates: groupTemplates }) => (
          <Stack key={group.id} gap="sm" mb="xl">
            <Group justify="space-between">
              <Group gap="xs">
                <Text fw={600}>{group.name}</Text>
                <Badge variant="light">{groupTemplates.length} forms</Badge>
              </Group>
              <Group gap="xs">
                <Button size="xs" variant="light" onClick={() => openEditGroup(group)}>
                  Edit group
                </Button>
                {group.id !== '00000000-0000-0000-0000-000000000001' && (
                  <Button
                    size="xs"
                    variant="light"
                    color="red"
                    loading={deleteGroupMutation.isPending}
                    onClick={() => deleteGroupMutation.mutate(group.id)}
                  >
                    Delete group
                  </Button>
                )}
              </Group>
            </Group>

            {groupTemplates.length === 0 ? (
              <Text c="dimmed" size="sm">
                No forms in this group yet.
              </Text>
            ) : (
              <Table>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th>Name</Table.Th>
                    <Table.Th>Type</Table.Th>
                    <Table.Th>Fields</Table.Th>
                    <Table.Th>Updated</Table.Th>
                    <Table.Th>Actions</Table.Th>
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                  {groupTemplates.map((template: FormTemplate) => (
                    <Table.Tr key={template.id}>
                      <Table.Td>{template.name}</Table.Td>
                      <Table.Td>
                        <Badge variant="light" color={template.formType === 'log' ? 'teal' : 'blue'}>
                          {template.formType}
                        </Badge>
                      </Table.Td>
                      <Table.Td>{template.fields.length}</Table.Td>
                      <Table.Td>{new Date(template.updatedAt).toLocaleDateString()}</Table.Td>
                      <Table.Td>
                        <Group gap="xs">
                          <Button
                            size="xs"
                            variant="light"
                            onClick={() => navigate(`/forms/${template.id}/edit`)}
                          >
                            Edit
                          </Button>
                          <Button
                            size="xs"
                            variant="light"
                            color="red"
                            loading={deleteMutation.isPending}
                            onClick={() => deleteMutation.mutate(template.id)}
                          >
                            Delete
                          </Button>
                        </Group>
                      </Table.Td>
                    </Table.Tr>
                  ))}
                </Table.Tbody>
              </Table>
            )}
          </Stack>
        ))
      )}

      <Modal
        opened={groupModalOpen}
        onClose={() => setGroupModalOpen(false)}
        title={editingGroup ? 'Edit form group' : 'Create form group'}
      >
        <Stack gap="md">
          <TextInput
            label="Group name"
            value={groupForm.name}
            onChange={(event) => setGroupForm({ ...groupForm, name: event.currentTarget.value })}
            required
          />
          <NumberInput
            label="Sort order"
            description="Lower numbers appear first on patient charts."
            value={groupForm.sortOrder}
            onChange={(value) => setGroupForm({ ...groupForm, sortOrder: Number(value) || 0 })}
          />
          <Button
            onClick={() => saveGroupMutation.mutate()}
            loading={saveGroupMutation.isPending}
            disabled={!groupForm.name.trim()}
          >
            {editingGroup ? 'Save group' : 'Create group'}
          </Button>
        </Stack>
      </Modal>
    </Container>
  );
}
