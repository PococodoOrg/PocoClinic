import { LoadingOverlay, Paper, Stack } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import { fetchPatient } from '../api/patients';
import { PatientBriefHeader } from '../components/patients/PatientBriefHeader';
import { PatientFormsSection } from '../components/forms/PatientFormsSection';
import { PatientChartLoadError } from '../components/patients/PatientChartLoadError';

export default function PatientFormsListPage() {
  const { id } = useParams<{ id: string }>();

  const { data: patient, isLoading, error } = useQuery({
    queryKey: ['patient', id],
    queryFn: () => fetchPatient(id!),
    enabled: Boolean(id),
  });

  if (error || (!isLoading && !patient)) {
    return <PatientChartLoadError patientId={id} message="Could not load patient forms." />;
  }

  return (
    <Stack gap="lg">
      {patient && <PatientBriefHeader patient={patient} sectionTitle="Patient forms" />}
      <Paper radius="md" p="xl" withBorder pos="relative">
        <LoadingOverlay visible={isLoading} />
        {patient && <PatientFormsSection patientId={patient.id} />}
      </Paper>
    </Stack>
  );
}
