import {
  Alert,
  Card,
  Group,
  List,
  Stack,
  Stepper,
  Text,
  ThemeIcon,
  Title,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { IconCheck, IconCircleDashed, IconDeviceUsb, IconDownload } from '@tabler/icons-react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createBackup } from '../api';
import { TouchPrimaryButton } from '../components/TouchPrimaryButton';
import { usePiTouch, useTouchUi } from '../context/PiTouchContext';

export function BackupWizardPage() {
  const [active, setActive] = useState(0);
  const [resultFile, setResultFile] = useState<string | null>(null);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { isPiTouch } = usePiTouch();
  const { stepperOrientation, textSize } = useTouchUi();

  const backupMutation = useMutation({
    mutationFn: createBackup,
    onSuccess: async (result) => {
      setResultFile(result.filename);
      setActive(2);
      await queryClient.invalidateQueries({ queryKey: ['helper-status'] });
      await queryClient.invalidateQueries({ queryKey: ['helper-backups'] });
      notifications.show({
        title: 'Backup finished',
        message: result.message,
        color: 'teal',
      });
    },
    onError: (error: Error) => {
      notifications.show({ title: 'Backup failed', message: error.message, color: 'red' });
    },
  });

  return (
    <Stack gap="lg">
      <Card withBorder padding={isPiTouch ? 'md' : 'xl'} radius="md">
        {!isPiTouch && (
          <>
            <Title order={2} mb="xs">
              Daily backup guide
            </Title>
            <Text c="dimmed" mb="lg">
              Follow these four easy steps. You can pause anytime — we will not skip ahead.
            </Text>
          </>
        )}
        {isPiTouch && (
          <Title order={3} mb="md">
            Daily backup
          </Title>
        )}

        <Stepper
          active={active}
          onStepClick={setActive}
          allowNextStepsSelect={false}
          color="teal"
          orientation={stepperOrientation}
          size={isPiTouch ? 'sm' : 'md'}
        >
          <Stepper.Step label="Prepare" description="USB ready">
            <Stack gap="md" mt="md">
              <List spacing="sm" size={textSize}>
                <List.Item>Ask staff to save open patient forms.</List.Item>
                <List.Item>Plug in today&apos;s labeled USB drive.</List.Item>
              </List>
              <TouchPrimaryButton onClick={() => setActive(1)}>USB is plugged in</TouchPrimaryButton>
            </Stack>
          </Stepper.Step>

          <Stepper.Step label="Create" description="Save file">
            <Stack gap="md" mt="md">
              <Alert color="blue" variant="light" title="Saving a protected copy">
                Patients, forms, and staff accounts will be saved on this device.
              </Alert>
              <TouchPrimaryButton
                leftSection={<IconDownload size={22} />}
                loading={backupMutation.isPending}
                onClick={() => backupMutation.mutate()}
              >
                Create backup now
              </TouchPrimaryButton>
            </Stack>
          </Stepper.Step>

          <Stepper.Step label="Copy" description="To USB">
            <Stack gap="md" mt="md">
              {resultFile && (
                <Alert color="teal" variant="light" icon={<IconCheck size={18} />}>
                  Saved: <strong>{resultFile}</strong>
                </Alert>
              )}
              <List
                spacing="sm"
                size={textSize}
                icon={
                  <ThemeIcon color="teal" size={28} radius="xl">
                    <IconDeviceUsb size={16} />
                  </ThemeIcon>
                }
              >
                <List.Item>Open the USB drive.</List.Item>
                <List.Item>Copy the backup file onto the USB.</List.Item>
                <List.Item>Wait until finished — do not unplug early.</List.Item>
              </List>
              <TouchPrimaryButton onClick={() => setActive(3)} disabled={!resultFile}>
                USB copy finished
              </TouchPrimaryButton>
            </Stack>
          </Stepper.Step>

          <Stepper.Step label="Done" description="All set">
            <Stack gap="md" mt="md">
              <Alert color="green" variant="light" title="Great job!" icon={<IconCheck size={18} />}>
                Remove the USB and store it in the locked drawer.
              </Alert>
              {!isPiTouch && (
                <Group>
                  <TouchPrimaryButton variant="light" onClick={() => navigate('/checklist')}>
                    Print checklist
                  </TouchPrimaryButton>
                  <TouchPrimaryButton variant="default" onClick={() => navigate('/')}>
                    Back to home
                  </TouchPrimaryButton>
                </Group>
              )}
              {isPiTouch && (
                <TouchPrimaryButton variant="light" onClick={() => navigate('/')}>
                  Back to status
                </TouchPrimaryButton>
              )}
            </Stack>
          </Stepper.Step>
        </Stepper>

        {!resultFile && active >= 2 && (
          <Text size="sm" c="dimmed" mt="md">
            <IconCircleDashed size={14} style={{ verticalAlign: 'middle' }} /> Create the backup file first (step 2).
          </Text>
        )}
      </Card>
    </Stack>
  );
}
