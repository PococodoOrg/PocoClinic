import { useEffect, useState } from 'react';
import { Anchor, Container, Stack, Text, Title } from '@mantine/core';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { fetchFormGroups, getFormTemplate, updateFormTemplate } from '../api/forms';
import { FormTemplateBuilder } from '../components/forms/FormTemplateBuilder';
import { FormTemplateFormData } from '../types/form';

export default function EditFormTemplate() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState<FormTemplateFormData | null>(null);

  const { data: groups = [], isLoading: groupsLoading } = useQuery({
    queryKey: ['form-groups'],
    queryFn: fetchFormGroups,
  });

  const { data: template, isLoading, error } = useQuery({
    queryKey: ['form-template', id],
    queryFn: () => getFormTemplate(id!),
    enabled: Boolean(id),
  });

  useEffect(() => {
    if (template) {
      setFormData({
        groupId: template.groupId,
        name: template.name,
        formType: template.formType,
        fields: template.fields,
      });
    }
  }, [template]);

  const updateMutation = useMutation({
    mutationFn: (data: FormTemplateFormData) => updateFormTemplate(id!, data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['form-templates'] });
      await queryClient.invalidateQueries({ queryKey: ['form-template', id] });
      notifications.show({
        title: 'Form updated',
        message: 'Template changes were saved.',
        color: 'green',
      });
      navigate('/forms');
    },
    onError: (mutationError: Error) => {
      notifications.show({
        title: 'Update failed',
        message: mutationError.message,
        color: 'red',
      });
    },
  });

  if (isLoading || groupsLoading || !formData) {
    return (
      <Container size="md">
        <Text>{error ? 'Could not load form.' : 'Loading...'}</Text>
      </Container>
    );
  }

  return (
    <Container size="md">
      <Stack gap="md" mb="xl">
        <Title order={2}>Edit form</Title>
        <Text size="sm">
          <Anchor component={Link} to="/forms">
            Back to forms
          </Anchor>
        </Text>
      </Stack>

      <FormTemplateBuilder
        value={formData}
        groups={groups}
        onChange={setFormData}
        onSubmit={() => updateMutation.mutate(formData)}
        isLoading={updateMutation.isPending}
        submitLabel="Save changes"
      />
    </Container>
  );
}
