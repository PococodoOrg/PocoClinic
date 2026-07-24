import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Container, Title, Paper, Text, Button, Group, Stack } from '@mantine/core';
import { PatientForm } from '../components/patients/PatientForm';
import { PatientFormData } from '../types/patient';
import { getPatient, updatePatient } from '../api/patients';
import { notifications } from '@mantine/notifications';
import { getErrorMessage } from '../utils/apiError';

export default function EditPatient() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // Fetch patient data
  const { data: patient, isLoading, error } = useQuery({
    queryKey: ['patient', id],
    queryFn: () => getPatient(id!),
  });

  // Update patient mutation
  const updateMutation = useMutation({
    mutationFn: (data: PatientFormData) => updatePatient(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['patient', id] });
      queryClient.invalidateQueries({ queryKey: ['patients'] });
      notifications.show({
        title: 'Success',
        message: 'Patient updated successfully',
        color: 'green',
      });
      navigate(`/patients/${id}`);
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Error',
        message: getErrorMessage(error, 'Failed to update patient'),
        color: 'red',
      });
    },
  });

  const handleSubmit = async (data: PatientFormData) => {
    updateMutation.mutate(data);
  };

  if (isLoading) {
    return (
      <Container size="md">
        <Text>Loading...</Text>
      </Container>
    );
  }

  if (error) {
    return (
      <Container size="md">
        <Text c="red">Error loading patient data</Text>
      </Container>
    );
  }

  if (!patient) {
    return (
      <Container size="md">
        <Text>Patient not found</Text>
      </Container>
    );
  }

  return (
    <Container size="md">
      <Stack gap="md">
        <Group justify="space-between">
          <Title order={2}>Edit Patient</Title>
          <Button variant="subtle" onClick={() => navigate(`/patients/${id}`)}>
            Cancel
          </Button>
        </Group>

        <Paper p="md" withBorder>
          <PatientForm
            initialValues={patient}
            onSubmit={handleSubmit}
            isLoading={updateMutation.isPending}
          />
        </Paper>
      </Stack>
    </Container>
  );
} 