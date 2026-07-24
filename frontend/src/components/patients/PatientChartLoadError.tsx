import { Button, Container, Stack, Text } from '@mantine/core';
import { IconArrowLeft } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';

interface PatientChartLoadErrorProps {
  patientId?: string;
  message: string;
}

export function PatientChartLoadError({ patientId, message }: PatientChartLoadErrorProps) {
  const navigate = useNavigate();

  return (
    <Container>
      <Stack gap="md">
        <Button
          variant="light"
          leftSection={<IconArrowLeft size={16} />}
          onClick={() => navigate(patientId ? `/patients/${patientId}` : '/patients')}
        >
          Back
        </Button>
        <Text c="red">{message}</Text>
      </Stack>
    </Container>
  );
}
