import { LoadingOverlay, Paper, Stack } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import { fetchPatient } from '../api/patients';
import { PatientBriefHeader } from '../components/patients/PatientBriefHeader';
import { PatientExerciseLogSection } from '../components/patients/PatientExerciseLogSection';
import { PatientChartLoadError } from '../components/patients/PatientChartLoadError';

export default function PatientExerciseLogListPage() {
  const { id } = useParams<{ id: string }>();

  const { data: patient, isLoading, error } = useQuery({
    queryKey: ['patient', id],
    queryFn: () => fetchPatient(id!),
    enabled: Boolean(id),
  });

  if (error || (!isLoading && !patient)) {
    return <PatientChartLoadError patientId={id} message="Could not load exercise log." />;
  }

  return (
    <Stack gap="lg">
      {patient && <PatientBriefHeader patient={patient} sectionTitle="Exercise log" />}
      <Paper radius="md" p="xl" withBorder pos="relative">
        <LoadingOverlay visible={isLoading} />
        {patient && <PatientExerciseLogSection patientId={patient.id} />}
      </Paper>
    </Stack>
  );
}
