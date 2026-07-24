import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  Text, 
  Container, 
  Title, 
  Paper, 
  Grid, 
  Group, 
  Button, 
  Stack,
  Badge,
  Divider,
  LoadingOverlay,
  Modal,
} from '@mantine/core';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { deletePatient, fetchPatient } from '../api/patients';
import { IconArrowLeft, IconEdit, IconMail, IconPhone, IconMapPin, IconUser, IconTrash } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { PatientFormsSection } from '../components/forms/PatientFormsSection';
import { PatientNotesSection } from '../components/patients/PatientNotesSection';
import { PatientDocumentsSection } from '../components/patients/PatientDocumentsSection';
import { PatientExerciseLogSection } from '../components/patients/PatientExerciseLogSection';
import { PATIENT_CHART_PREVIEW_LIMIT } from '../constants/patientChart';

export default function PatientDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);

  const { data: patient, isLoading, error } = useQuery({
    queryKey: ['patient', id],
    queryFn: () => fetchPatient(id!),
    enabled: Boolean(id),
  });

  const deleteMutation = useMutation({
    mutationFn: () => deletePatient(id!),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['patients'] });
      notifications.show({
        title: 'Patient deleted',
        message: 'The patient record was removed successfully.',
        color: 'green',
      });
      navigate('/patients');
    },
    onError: () => {
      notifications.show({
        title: 'Error',
        message: 'Failed to delete patient. Please try again.',
        color: 'red',
      });
    },
  });

  if (error) {
    notifications.show({
      title: 'Error',
      message: 'Failed to load patient details. Please try again later.',
      color: 'red'
    });
    return (
      <Container>
        <Button 
          variant="light" 
          leftSection={<IconArrowLeft size={16} />}
          onClick={() => navigate('/patients')}
          mb="md"
        >
          Back to Patients
        </Button>
        <Text c="red">Error loading patient details. Please try again later.</Text>
      </Container>
    );
  }

  if (!patient) {
    return (
      <Container>
        <Button 
          variant="light" 
          leftSection={<IconArrowLeft size={16} />}
          onClick={() => navigate('/patients')}
          mb="md"
        >
          Back to Patients
        </Button>
        <Text>Patient not found.</Text>
      </Container>
    );
  }

  return (
    <Stack gap="lg">
      <Paper radius="md" p="xl" withBorder>
        <LoadingOverlay visible={isLoading} />

        <Stack gap="lg">
          <Group justify="space-between" align="flex-start" className="workstation-page-header" wrap="wrap">
            <Group>
              <Button 
                variant="light" 
                leftSection={<IconArrowLeft size={16} />}
                onClick={() => navigate('/patients')}
              >
                Back to Patients
              </Button>
              <Title order={2}>{patient.firstName} {patient.lastName}</Title>
            </Group>
            <Group>
              <Button 
                variant="light" 
                leftSection={<IconEdit size={16} />}
                onClick={() => navigate(`/patients/${id}/edit`)}
              >
                Edit Patient
              </Button>
              <Button
                variant="light"
                color="red"
                leftSection={<IconTrash size={16} />}
                onClick={() => setDeleteModalOpen(true)}
              >
                Delete
              </Button>
            </Group>
          </Group>

          <Modal
            opened={deleteModalOpen}
            onClose={() => setDeleteModalOpen(false)}
            title="Delete patient"
          >
            <Stack>
              <Text>
                Are you sure you want to delete {patient.firstName} {patient.lastName}? This action cannot be undone.
              </Text>
              <Group justify="flex-end">
                <Button variant="default" onClick={() => setDeleteModalOpen(false)}>
                  Cancel
                </Button>
                <Button
                  color="red"
                  loading={deleteMutation.isPending}
                  onClick={() => deleteMutation.mutate()}
                >
                  Delete Patient
                </Button>
              </Group>
            </Stack>
          </Modal>

          <Divider />

          {/* Patient Information */}
          <Grid>
            {/* Basic Information */}
            <Grid.Col span={{ base: 12, md: 6 }}>
              <Stack gap="md">
                <Title order={3}>Basic Information</Title>
                <Group>
                  <IconUser size={20} />
                  <Text fw={500}>Gender:</Text>
                  <Badge variant="light">{patient.gender}</Badge>
                </Group>
                <Group>
                  <IconUser size={20} />
                  <Text fw={500}>Date of Birth:</Text>
                  <Text>{new Date(patient.dateOfBirth).toLocaleDateString()}</Text>
                </Group>
                {patient.height !== null && patient.height !== undefined && (
                  <Group>
                    <IconUser size={20} />
                    <Text fw={500}>Height:</Text>
                    <Text>{patient.height} cm</Text>
                  </Group>
                )}
                {patient.weight !== null && patient.weight !== undefined && (
                  <Group>
                    <IconUser size={20} />
                    <Text fw={500}>Weight:</Text>
                    <Text>{patient.weight} kg</Text>
                  </Group>
                )}
              </Stack>
            </Grid.Col>

            {/* Contact Information */}
            <Grid.Col span={{ base: 12, md: 6 }}>
              <Stack gap="md">
                <Title order={3}>Contact Information</Title>
                <Group>
                  <IconMail size={20} />
                  <Text fw={500}>Email:</Text>
                  <Text>{patient.email}</Text>
                </Group>
                <Group>
                  <IconPhone size={20} />
                  <Text fw={500}>Phone:</Text>
                  <Text>{patient.phoneNumber}</Text>
                </Group>
                {patient.address && (
                  <Group>
                    <IconMapPin size={20} />
                    <Text fw={500}>Address:</Text>
                    <Text>
                      {[
                        patient.address.street,
                        patient.address.city,
                        patient.address.state,
                        patient.address.postalCode,
                        patient.address.country
                      ].filter(Boolean).join(', ')}
                    </Text>
                  </Group>
                )}
              </Stack>
            </Grid.Col>
          </Grid>

          <Divider className="patient-chart-full-width" />

          <div className="patient-chart-clinical">
            <PatientNotesSection patientId={patient.id} previewLimit={PATIENT_CHART_PREVIEW_LIMIT} />
            <PatientDocumentsSection patientId={patient.id} />
          </div>

          <Divider className="patient-chart-full-width" />

          <div className="patient-chart-full-width">
            <PatientExerciseLogSection patientId={patient.id} previewLimit={PATIENT_CHART_PREVIEW_LIMIT} />
          </div>

          <Divider className="patient-chart-full-width" />

          <div className="patient-chart-full-width">
            <PatientFormsSection patientId={patient.id} previewLimit={PATIENT_CHART_PREVIEW_LIMIT} />
          </div>
        </Stack>
      </Paper>
    </Stack>
  );
} 