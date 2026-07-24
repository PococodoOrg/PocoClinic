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
import { BackupVerifyResult, fetchBackups, fetchStatus, restoreBackup, verifyBackup } from '../api';
import { BackupPicker } from '../components/BackupPicker';
import { TouchPrimaryButton } from '../components/TouchPrimaryButton';
import { usePiTouch, useTouchUi } from '../context/PiTouchContext';

export function RestoreWizardPage() {
  const [active, setActive] = useState(0);
  const [selected, setSelected] = useState<string | null>(null);
  const [confirmText, setConfirmText] = useState('');
  const [verifyResult, setVerifyResult] = useState<BackupVerifyResult | null>(null);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { isPiTouch } = usePiTouch();
  const { stepperOrientation } = useTouchUi();

  const { data: status } = useQuery({
    queryKey: ['helper-status'],
    queryFn: fetchStatus,
  });

  const { data, isLoading } = useQuery({
    queryKey: ['helper-backups'],
    queryFn: fetchBackups,
  });

  const verifyMutation = useMutation({
    mutationFn: verifyBackup,
    onSuccess: (result) => {
      setVerifyResult(result);
      notifications.show({
        title: result.valid ? 'Backup verified' : 'Verification failed',
        message: result.message ?? result.filename,
        color: result.valid ? 'green' : 'red',
      });
    },
    onError: (error: Error) => {
      notifications.show({ title: 'Verification failed', message: error.message, color: 'red' });
    },
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
  const mainAppRunning = status?.mainAppOnline ?? false;
  const restoreBlocked = mainAppRunning || (verifyResult?.valid === false);

  return (
    <Stack gap="md">
      <Alert color="red" variant="filled" title="Careful" icon={<IconAlertTriangle size={22} />}>
        Restore replaces all clinic data. Everyone must sign out first.
      </Alert>

      {mainAppRunning && (
        <Alert color="orange" variant="light" title="Stop the main app first">
          PocoClinic is still running. Stop it before restore (
          <strong>sudo systemctl stop pococlinic</strong> on the server, or close the dev backend on this PC).
        </Alert>
      )}

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
                <Text size={isPiTouch ? 'md' : 'sm'}>• Main PocoClinic app stopped</Text>
                <Text size={isPiTouch ? 'md' : 'sm'}>• Correct USB backup selected</Text>
                <Text size={isPiTouch ? 'md' : 'sm'}>• Clinic leader approved</Text>
              </Stack>
              <TouchPrimaryButton color="orange" disabled={mainAppRunning} onClick={() => setActive(1)}>
                {mainAppRunning ? 'Stop the main app first' : 'I understand'}
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
                <BackupPicker
                  backups={backups}
                  selected={selected}
                  onSelect={(filename) => {
                    setSelected(filename);
                    setVerifyResult(null);
                  }}
                />
              )}
              {selected && (
                <TouchPrimaryButton
                  color="orange"
                  variant="light"
                  loading={verifyMutation.isPending}
                  onClick={() => verifyMutation.mutate(selected)}
                >
                  Verify this backup
                </TouchPrimaryButton>
              )}
              {verifyResult && selected === verifyResult.filename && (
                <Alert color={verifyResult.valid ? 'green' : 'red'} variant="light" title="Verification result">
                  {verifyResult.message}
                </Alert>
              )}
              <TouchPrimaryButton
                color="orange"
                disabled={!selected || verifyResult?.valid === false}
                onClick={() => setActive(2)}
              >
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
                disabled={!canConfirm || restoreBlocked}
                onClick={() => restoreMutation.mutate()}
              >
                Restore clinic data
              </TouchPrimaryButton>
            </Stack>
          </Stepper.Step>

          <Stepper.Step label="Done" description="Verify">
            <Stack gap="md" mt="md">
              <Alert color="green" variant="light" title="Restore finished">
                Restart the main PocoClinic service if you stopped it, then have someone sign in and check a patient
                chart. Take a new backup right away.
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
