import {
  Button,
  Card,
  Group,
  SimpleGrid,
  Stack,
  Switch,
  Text,
  Title,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import {
  fetchAdminPatientFieldRequirements,
  updatePatientFieldRequirements,
} from '../../api/admin';
import { DEFAULT_PATIENT_FIELD_REQUIREMENTS, PatientFieldRequirements } from '../../types/patient';
import { getErrorMessage } from '../../utils/apiError';

type FieldConfig = {
  key: keyof PatientFieldRequirements;
  label: string;
  locked?: boolean;
  description?: string;
};

const FIELD_CONFIG: FieldConfig[] = [
  { key: 'firstName', label: 'First name', locked: true, description: 'Always required for patient identity' },
  { key: 'lastName', label: 'Last name', locked: true, description: 'Always required for patient identity' },
  { key: 'middleName', label: 'Middle name' },
  { key: 'dateOfBirth', label: 'Date of birth', locked: true, description: 'Always required for patient identity' },
  { key: 'gender', label: 'Gender' },
  { key: 'email', label: 'Email' },
  { key: 'phoneNumber', label: 'Phone number' },
  { key: 'addressStreet', label: 'Street address' },
  { key: 'addressCity', label: 'City' },
  { key: 'addressState', label: 'State' },
  { key: 'addressPostalCode', label: 'ZIP code' },
  { key: 'height', label: 'Height' },
  { key: 'weight', label: 'Weight' },
];

export function PatientFieldSettingsPanel() {
  const queryClient = useQueryClient();
  const { data, isLoading, isError } = useQuery({
    queryKey: ['admin-patient-field-requirements'],
    queryFn: fetchAdminPatientFieldRequirements,
  });
  const [draft, setDraft] = useState<PatientFieldRequirements>(DEFAULT_PATIENT_FIELD_REQUIREMENTS);

  useEffect(() => {
    if (data) {
      setDraft(data);
    }
  }, [data]);

  const saveMutation = useMutation({
    mutationFn: updatePatientFieldRequirements,
    onSuccess: (saved) => {
      setDraft(saved);
      queryClient.setQueryData(['admin-patient-field-requirements'], saved);
      queryClient.invalidateQueries({ queryKey: ['patient-field-requirements'] });
      notifications.show({
        title: 'Patient fields updated',
        message: 'Required patient fields were saved for the clinic.',
        color: 'green',
      });
    },
    onError: (error) => {
      notifications.show({
        title: 'Could not save settings',
        message: getErrorMessage(error, 'Failed to update patient field requirements'),
        color: 'red',
      });
    },
  });

  const toggleField = (key: keyof PatientFieldRequirements, checked: boolean) => {
    setDraft((current) => ({
      ...current,
      [key]: checked,
      ...(key === 'firstName' || key === 'lastName' || key === 'dateOfBirth' ? { [key]: true } : {}),
    }));
  };

  const hasChanges = data ? JSON.stringify(data) !== JSON.stringify(draft) : false;

  return (
    <Card withBorder radius="md" padding="lg">
      <Stack gap="md">
        <Stack gap={4}>
          <Title order={3}>Patient chart fields</Title>
          <Text size="sm" c="dimmed">
            Choose which fields staff must complete when creating or editing a patient chart.
            First name, last name, and date of birth always stay required.
          </Text>
        </Stack>

        {isError && (
          <Text size="sm" c="red">
            Could not load current settings. Refresh the page or try again.
          </Text>
        )}

        <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="md">
          {FIELD_CONFIG.map((field) => (
            <Switch
              key={field.key}
              label={field.label}
              description={field.description}
              checked={draft[field.key]}
              disabled={field.locked || isLoading || saveMutation.isPending}
              onChange={(event) => toggleField(field.key, event.currentTarget.checked)}
            />
          ))}
        </SimpleGrid>

        <Group justify="flex-end">
          <Button
            onClick={() => saveMutation.mutate(draft)}
            loading={saveMutation.isPending}
            disabled={!hasChanges || isLoading}
          >
            Save field requirements
          </Button>
        </Group>
      </Stack>
    </Card>
  );
}
