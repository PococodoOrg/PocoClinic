import {
  Alert,
  Card,
  Stack,
  Stepper,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { IconAlertTriangle } from '@tabler/icons-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchBackups, restoreBackup } from '../api';
import { BackupPicker } from '../components/BackupPicker';
import { TouchPrimaryButton } from '../components/TouchPrimaryButton';
import { usePiTouch, useTouchUi } from '../context/PiTouchContext';

export function RestoreWizardPage() {
  const [active, setActive] = useState(0);
  const [selected, setSelected] = useState<string | null>(null);
  const [confirmText, setConfirmText] = useState('');
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { isPiTouch } = usePiTouch();
  const { stepperOrientation } = useTouchUi();

  const { data, isLoading } = useQuery({
    queryKey: ['helper-backups'],
    queryFn: fetchBackups,
  });

  const restoreMutation = useMutation({
    mutationFn: () => restoreBackup(selected!, confirmText),
    onSuccess: async (result) => {
      await queryClient.invalidateQueries({ queryKey: ['helper-status'] });
      notifications.show({ title: 'Restore complete', message: result.message, color: 'teal' });
      setActive(3);
    },
    onError: (error: Error) => {
      notifications.show({ title: 'Restore failed', message: error.message, color: 'red' });
    },
  });

  const backups = data?.backups ?? [];
  const canConfirm = confirmText.trim().toUpperCase() === 'RESTORE' && selected;

  return (
    <Stack gap="md">
      <Alert color="red" variant="filled" title="Careful" icon={<IconAlertTriangle size={22} />}>
        Restore replaces all clinic data. Everyone must sign out first.
      </Alert>

      <Card withBorder padding={isPiTouch ? 'md' : 'xl'} radius="md">
        {!isPiTouch && (
          <Title order={2} mb="lg">
            Restore guide
          </Title>
        )}
        {isPiTouch && (
          <Title order={3} mb="md">
            Restore clinic
          </Title>
        )}

        <Stepper
          active={active}
          onStepClick={setActive}
          allowNextStepsSelect={false}
          color="orange"
          orientation={stepperOrientation}
          size={isPiTouch ? 'sm' : 'md'}
        >
          <Stepper.Step label="Understand" description="Important">
            <Stack gap="md" mt="md">
              <Text size={isPiTouch ? 'md' : 'sm'}>
                This turns back time. Records added after the backup date will be lost.
              </Text>
              <Stack gap={8}>
                <Text size={isPiTouch ? 'md' : 'sm'}>• All staff signed out</Text>
                <Text size={isPiTouch ? 'md' : 'sm'}>• Correct USB backup selected</Text>
                <Text size={isPiTouch ? 'md' : 'sm'}>• Clinic leader approved</Text>
              </Stack>
              <TouchPrimaryButton color="orange" onClick={() => setActive(1)}>
                I understand
              </TouchPrimaryButton>
            </Stack>
          </Stepper.Step>

          <Stepper.Step label="Choose" description="Pick date">
            <Stack gap="md" mt="md">
              {isLoading && <Text c="dimmed">Loading backups…</Text>}
              {!isLoading && backups.length === 0 && (
                <Text c="dimmed">No backups found. Create a backup first or copy a file to the backups folder.</Text>
              )}
              {backups.length > 0 && (
                <BackupPicker backups={backups} selected={selected} onSelect={setSelected} />
              )}
              <TouchPrimaryButton color="orange" disabled={!selected} onClick={() => setActive(2)}>
                Use this backup
              </TouchPrimaryButton>
            </Stack>
          </Stepper.Step>

          <Stepper.Step label="Confirm" description="Type RESTORE">
            <Stack gap="md" mt="md">
              <Text size={isPiTouch ? 'md' : 'sm'}>
                Type <strong>RESTORE</strong> to continue.
              </Text>
              <TextInput
                label="Confirmation"
                placeholder="RESTORE"
                size={isPiTouch ? 'lg' : 'md'}
                value={confirmText}
                onChange={(event) => setConfirmText(event.currentTarget.value)}
                styles={isPiTouch ? { input: { minHeight: 56, fontSize: 18 } } : undefined}
              />
              <TouchPrimaryButton
                color="red"
                loading={restoreMutation.isPending}
                disabled={!canConfirm}
                onClick={() => restoreMutation.mutate()}
              >
                Restore clinic data
              </TouchPrimaryButton>
            </Stack>
          </Stepper.Step>

          <Stepper.Step label="Done" description="Verify">
            <Stack gap="md" mt="md">
              <Alert color="green" variant="light" title="Restore finished">
                Have someone sign in and check a patient chart. Take a new backup right away.
              </Alert>
              <TouchPrimaryButton color="teal" onClick={() => navigate('/backup')}>
                Create new backup
              </TouchPrimaryButton>
            </Stack>
          </Stepper.Step>
        </Stepper>
      </Card>
    </Stack>
  );
}
