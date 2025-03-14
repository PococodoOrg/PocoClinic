import { Container, Title } from '@mantine/core';
import { PatientForm } from '../components/patients/PatientForm';
import { createPatient } from '../api/patients';
import { PatientFormData } from '../types/patient';

export function CreatePatient() {
  return (
    <Container size="lg">
      <Title order={2} mb="xl">Create New Patient</Title>
      <PatientForm onSubmit={createPatient} />
    </Container>
  );
} 