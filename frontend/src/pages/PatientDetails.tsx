import React from 'react';
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
  LoadingOverlay
} from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { fetchPatient } from '../api/patients';
import { IconArrowLeft, IconEdit, IconMail, IconPhone, IconMapPin, IconUser } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';

export default function PatientDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: patient, isLoading, error } = useQuery({
    queryKey: ['patient', id],
    queryFn: () => fetchPatient(id!),
  });

  // Add console logging to debug
  React.useEffect(() => {
    if (patient) {
      console.log('Patient data:', patient);
    }
  }, [patient]);

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
    <Container size="lg">
      <Paper radius="md" p="xl" withBorder>
        <LoadingOverlay visible={isLoading} />
        
        <Stack>
          {/* Header Section */}
          <Group justify="space-between" align="flex-start">
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
            <Button 
              variant="light" 
              leftSection={<IconEdit size={16} />}
              onClick={() => navigate(`/patients/${id}/edit`)}
            >
              Edit Patient
            </Button>
          </Group>

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
        </Stack>
      </Paper>
    </Container>
  );
} 