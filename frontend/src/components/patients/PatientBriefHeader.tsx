import { Badge, Button, Group, Paper, Stack, Text, Title } from '@mantine/core';
import { IconArrowLeft, IconPhone, IconUser } from '@tabler/icons-react';
import { Link } from 'react-router-dom';
import { Patient } from '../../types/patient';

interface PatientBriefHeaderProps {
  patient: Patient;
  sectionTitle: string;
}

export function PatientBriefHeader({ patient, sectionTitle }: PatientBriefHeaderProps) {
  const dob = new Date(patient.dateOfBirth).toLocaleDateString();

  return (
    <Paper withBorder p="md" radius="md">
      <Stack gap="sm">
        <Button
          component={Link}
          to={`/patients/${patient.id}`}
          variant="subtle"
          size="compact-sm"
          leftSection={<IconArrowLeft size={16} />}
          px={0}
        >
          Back to chart
        </Button>
        <Group justify="space-between" align="flex-start" wrap="wrap" gap="md">
          <Stack gap={4}>
            <Title order={2}>
              {patient.firstName} {patient.lastName}
            </Title>
            <Group gap="md" wrap="wrap">
              <Group gap={6}>
                <IconUser size={16} />
                <Text size="sm" c="dimmed">
                  DOB {dob}
                </Text>
              </Group>
              <Badge variant="light" tt="capitalize">
                {patient.gender}
              </Badge>
              {patient.phoneNumber && (
                <Group gap={6}>
                  <IconPhone size={16} />
                  <Text size="sm">{patient.phoneNumber}</Text>
                </Group>
              )}
            </Group>
          </Stack>
          <Title order={4} c="dimmed">
            {sectionTitle}
          </Title>
        </Group>
      </Stack>
    </Paper>
  );
}
