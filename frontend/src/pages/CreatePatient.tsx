import { useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { PatientForm } from '../components/patients/PatientForm';
import { patientApi } from '../api/patients';
import { PatientFormData } from '../types/patient';
import { Title, Container, Paper } from '@mantine/core';

export function CreatePatient() {
  const navigate = useNavigate();

  const createPatientMutation = useMutation({
    mutationFn: (data: PatientFormData) => patientApi.createPatient(data),
    onSuccess: () => {
      notifications.show({
        title: 'Success',
        message: 'Patient created successfully',
        color: 'green',
      });
      navigate('/patients');
    },
    onError: (error) => {
      notifications.show({
        title: 'Error',
        message: error instanceof Error ? error.message : 'Failed to create patient',
        color: 'red',
      });
    },
  });

  const handleSubmit = (values: PatientFormData) => {
    createPatientMutation.mutate(values);
  };

  return (
    <Container size="lg">
      <Title order={2} mb="lg">Create New Patient</Title>
      <Paper shadow="xs" p="md">
        <PatientForm 
          onSubmit={handleSubmit}
          isLoading={createPatientMutation.isPending}
        />
      </Paper>
    </Container>
  );
} 