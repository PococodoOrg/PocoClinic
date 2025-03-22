import { Container, Title } from '@mantine/core';
import { PatientForm } from '../components/patients/PatientForm';
import { createPatient } from '../api/patients';
import { PatientFormData } from '../types/patient';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { useNavigate } from 'react-router-dom';

export default function CreatePatient() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const createPatientMutation = useMutation({
    mutationFn: createPatient,
    onSuccess: () => {
      // Invalidate and refetch patients query
      queryClient.invalidateQueries({ queryKey: ['patients'] });
      notifications.show({
        title: 'Success',
        message: 'Patient created successfully',
        color: 'green'
      });
      navigate('/patients');
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Error',
        message: error.message || 'Failed to create patient',
        color: 'red'
      });
    }
  });

  const handleSubmit = async (data: PatientFormData) => {
    await createPatientMutation.mutateAsync(data);
  };

  return (
    <Container size="lg">
      <Title order={2} mb="xl">Create New Patient</Title>
      <PatientForm onSubmit={handleSubmit} isLoading={createPatientMutation.isPending} />
    </Container>
  );
} 