import { useState } from 'react';
import { Anchor, Container, Stack, Text, Title } from '@mantine/core';
import { Link, useNavigate } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { createFormTemplate, fetchFormGroups } from '../api/forms';
import { FormTemplateBuilder } from '../components/forms/FormTemplateBuilder';
import { FormTemplateFormData, createEmptyField } from '../types/form';
import { getErrorMessage } from '../utils/apiError';

export default function CreateFormTemplate() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: groups = [], isLoading: groupsLoading } = useQuery({
    queryKey: ['form-groups'],
    queryFn: fetchFormGroups,
  });

  const [formData, setFormData] = useState<FormTemplateFormData>({
    groupId: '',
    name: '',
    formType: 'singleton',
    fields: [createEmptyField(0)],
  });

  const createMutation = useMutation({
    mutationFn: createFormTemplate,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['form-templates'] });
      notifications.show({
        title: 'Form created',
        message: 'Staff can now use this form on patient records.',
        color: 'green',
      });
      navigate('/forms');
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Create failed',
        message: getErrorMessage(error, 'Could not create form template'),
        color: 'red',
      });
    },
  });

  const handleSubmit = () => {
    createMutation.mutate(formData);
  };

  if (groupsLoading) {
    return (
      <Container size="md">
        <Text c="dimmed">Loading...</Text>
      </Container>
    );
  }

  const resolvedFormData =
    formData.groupId || groups.length === 0
      ? formData
      : { ...formData, groupId: groups[0].id };

  return (
    <Container size="md">
      <Stack gap="md" mb="xl">
        <Title order={2}>Create form</Title>
        <Text c="dimmed">
          Build a form and choose whether it is a singleton (one active copy) or a log (many entries).
        </Text>
        <Text size="sm">
          <Anchor component={Link} to="/forms">
            Back to forms
          </Anchor>
        </Text>
      </Stack>

      <FormTemplateBuilder
        value={resolvedFormData}
        groups={groups}
        onChange={setFormData}
        onSubmit={handleSubmit}
        isLoading={createMutation.isPending}
        submitLabel="Create form"
      />
    </Container>
  );
}
