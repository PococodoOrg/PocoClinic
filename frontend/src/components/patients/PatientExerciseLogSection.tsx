import {
  Badge,
  Button,
  Group,
  Modal,
  NumberInput,
  Select,
  Stack,
  Text,
  TextInput,
  Textarea,
  Title,
} from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { notifications } from '@mantine/notifications';
import { IconPlus, IconTrash } from '@tabler/icons-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useMemo, useState } from 'react';
import {
  createExerciseEntry,
  createExercisePlan,
  deleteExerciseEntry,
  deleteExercisePlan,
  fetchExerciseEntries,
  fetchExercisePlans,
  updateExercisePlan,
} from '../../api/exerciseLog';
import { TABLET_MEDIA_QUERY } from '../../layout/breakpoints';
import {
  ExerciseLogEntry,
  ExerciseLogEntryInput,
  ExercisePlan,
  ExercisePlanStatus,
} from '../../types/exerciseLog';
import { getErrorMessage } from '../../utils/apiError';
import { formatSessionSummary, progressionHint } from '../../utils/exerciseLog';
import { parseNumberInput } from '../../utils/numberInput';

interface PatientExerciseLogSectionProps {
  patientId: string;
}

const STATUS_OPTIONS = [
  { value: 'active', label: 'Active' },
  { value: 'completed', label: 'Completed' },
  { value: 'archived', label: 'Archived' },
];

function statusColor(status: ExercisePlanStatus): string {
  switch (status) {
    case 'active':
      return 'green';
    case 'completed':
      return 'blue';
    default:
      return 'gray';
  }
}

export function PatientExerciseLogSection({ patientId }: PatientExerciseLogSectionProps) {
  const queryClient = useQueryClient();
  const isTablet = useMediaQuery(TABLET_MEDIA_QUERY);
  const [selectedPlanId, setSelectedPlanId] = useState<string | null>(null);
  const [planModalOpen, setPlanModalOpen] = useState(false);
  const [entryModalOpen, setEntryModalOpen] = useState(false);
  const [deletePlanTarget, setDeletePlanTarget] = useState<ExercisePlan | null>(null);
  const [deleteEntryTarget, setDeleteEntryTarget] = useState<ExerciseLogEntry | null>(null);

  const [planName, setPlanName] = useState('');
  const [planDescription, setPlanDescription] = useState('');
  const [entryExerciseName, setEntryExerciseName] = useState('');
  const [entrySets, setEntrySets] = useState<number | ''>(3);
  const [entryReps, setEntryReps] = useState<number | ''>(10);
  const [entryResistance, setEntryResistance] = useState('');
  const [entryDifficulty, setEntryDifficulty] = useState<number | ''>('');
  const [entryNotes, setEntryNotes] = useState('');

  const resetEntryForm = () => {
    setEntryExerciseName('');
    setEntrySets(3);
    setEntryReps(10);
    setEntryResistance('');
    setEntryDifficulty('');
    setEntryNotes('');
  };

  const { data: plans = [], isLoading: plansLoading } = useQuery({
    queryKey: ['exercise-plans', patientId],
    queryFn: () => fetchExercisePlans(patientId),
  });

  const activePlanId = selectedPlanId ?? plans[0]?.id ?? null;
  const selectedPlan = plans.find((plan) => plan.id === activePlanId) ?? null;

  const { data: entries = [], isLoading: entriesLoading } = useQuery({
    queryKey: ['exercise-entries', patientId, activePlanId],
    queryFn: () => fetchExerciseEntries(patientId, activePlanId!),
    enabled: Boolean(activePlanId),
  });

  const invalidate = async () => {
    await queryClient.invalidateQueries({ queryKey: ['exercise-plans', patientId] });
    if (activePlanId) {
      await queryClient.invalidateQueries({ queryKey: ['exercise-entries', patientId, activePlanId] });
    }
  };

  const createPlanMutation = useMutation({
    mutationFn: () =>
      createExercisePlan(patientId, { name: planName, description: planDescription }),
    onSuccess: async (plan) => {
      setPlanModalOpen(false);
      setPlanName('');
      setPlanDescription('');
      setSelectedPlanId(plan.id);
      await invalidate();
      notifications.show({
        title: 'Plan created',
        message: 'You can start logging exercise sessions.',
        color: 'green',
      });
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Could not create plan',
        message: getErrorMessage(error, 'Could not create exercise plan'),
        color: 'red',
      });
    },
  });

  const updatePlanMutation = useMutation({
    mutationFn: (status: ExercisePlanStatus) =>
      updateExercisePlan(patientId, selectedPlan!.id, {
        name: selectedPlan!.name,
        description: selectedPlan!.description,
        status,
      }),
    onSuccess: async () => {
      await invalidate();
      notifications.show({
        title: 'Plan updated',
        message: 'Plan status saved.',
        color: 'green',
      });
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Could not update plan',
        message: getErrorMessage(error, 'Could not update exercise plan'),
        color: 'red',
      });
    },
  });

  const deletePlanMutation = useMutation({
    mutationFn: (planId: string) => deleteExercisePlan(patientId, planId),
    onSuccess: async () => {
      setDeletePlanTarget(null);
      setSelectedPlanId(null);
      await invalidate();
      notifications.show({
        title: 'Plan deleted',
        message: 'The plan and its session history were removed.',
        color: 'green',
      });
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Could not delete plan',
        message: getErrorMessage(error, 'Could not delete exercise plan'),
        color: 'red',
      });
    },
  });

  const createEntryMutation = useMutation({
    mutationFn: () => {
      const payload: ExerciseLogEntryInput = {
        exerciseName: entryExerciseName.trim(),
        sets: entrySets === '' ? 0 : entrySets,
        reps: entryReps === '' ? 0 : entryReps,
        resistance: entryResistance.trim(),
        difficulty: entryDifficulty === '' ? null : entryDifficulty,
        notes: entryNotes.trim(),
      };
      return createExerciseEntry(patientId, activePlanId!, payload);
    },
    onSuccess: async () => {
      setEntryModalOpen(false);
      resetEntryForm();
      await invalidate();
      notifications.show({
        title: 'Session logged',
        message: 'Exercise progress recorded on the chart.',
        color: 'green',
      });
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Could not log session',
        message: getErrorMessage(error, 'Could not save exercise entry'),
        color: 'red',
      });
    },
  });

  const deleteEntryMutation = useMutation({
    mutationFn: (entry: ExerciseLogEntry) =>
      deleteExerciseEntry(patientId, entry.planId, entry.id),
    onSuccess: async () => {
      setDeleteEntryTarget(null);
      await invalidate();
      notifications.show({
        title: 'Entry deleted',
        message: 'Session removed from the log.',
        color: 'green',
      });
    },
    onError: (error: unknown) => {
      notifications.show({
        title: 'Could not delete entry',
        message: getErrorMessage(error, 'Could not delete exercise entry'),
        color: 'red',
      });
    },
  });

  const recentExerciseNames = useMemo(() => {
    const names = new Set<string>();
    entries.forEach((entry) => names.add(entry.exerciseName));
    return Array.from(names).slice(0, 8);
  }, [entries]);

  return (
    <Stack gap="md" className="patient-chart-full-width">
      <Group justify="space-between" align="flex-start" wrap="wrap">
        <div>
          <Title order={3}>Exercise log</Title>
          <Text size="sm" c="dimmed" mt={4}>
            Physical therapy plans and session progress (reps, sets, resistance).
          </Text>
        </div>
        <Button
          leftSection={<IconPlus size={18} />}
          onClick={() => setPlanModalOpen(true)}
          size={isTablet ? 'md' : 'sm'}
          className="touch-control"
        >
          New plan
        </Button>
      </Group>

      {plansLoading ? (
        <Text c="dimmed" aria-live="polite">
          Loading exercise plans…
        </Text>
      ) : plans.length === 0 ? (
        <Text c="dimmed">
          No exercise plans yet. Create a plan for home or in-clinic PT, then log sessions to track
          improvement.
        </Text>
      ) : (
        <Stack gap="sm" role="list" aria-label="Exercise plans">
          {plans.map((plan) => {
            const selected = plan.id === activePlanId;
            return (
              <button
                key={plan.id}
                type="button"
                className="patient-list-card"
                role="listitem"
                aria-pressed={selected}
                onClick={() => setSelectedPlanId(plan.id)}
                style={
                  selected
                    ? { borderColor: 'var(--mantine-color-blue-5)', boxShadow: 'inset 0 0 0 1px var(--mantine-color-blue-5)' }
                    : undefined
                }
              >
                <Group justify="space-between" wrap="nowrap" align="flex-start">
                  <Stack gap={4} style={{ minWidth: 0 }}>
                    <Text fw={600} size="md" lineClamp={1}>
                      {plan.name}
                    </Text>
                    {plan.description && (
                      <Text size="sm" c="dimmed" lineClamp={2}>
                        {plan.description}
                      </Text>
                    )}
                  </Stack>
                  <Badge color={statusColor(plan.status)} variant="light" tt="capitalize">
                    {plan.status}
                  </Badge>
                </Group>
              </button>
            );
          })}
        </Stack>
      )}

      {selectedPlan && (
        <Stack gap="md">
          <Group justify="space-between" wrap="wrap" gap="sm">
            <Select
              label="Plan status"
              data={STATUS_OPTIONS}
              value={selectedPlan.status}
              onChange={(value) => {
                if (value) {
                  updatePlanMutation.mutate(value as ExercisePlanStatus);
                }
              }}
              w={isTablet ? '100%' : 180}
              size={isTablet ? 'md' : 'sm'}
              className="touch-control"
              allowDeselect={false}
            />
            <Group gap="sm">
              <Button
                leftSection={<IconPlus size={18} />}
                onClick={() => {
                  resetEntryForm();
                  setEntryModalOpen(true);
                }}
                size={isTablet ? 'md' : 'sm'}
                className="touch-control"
              >
                Log session
              </Button>
              <Button
                variant="light"
                color="red"
                leftSection={<IconTrash size={16} />}
                onClick={() => setDeletePlanTarget(selectedPlan)}
                size={isTablet ? 'md' : 'sm'}
                className="touch-control"
              >
                Delete plan
              </Button>
            </Group>
          </Group>

          {entriesLoading ? (
            <Text c="dimmed" aria-live="polite">
              Loading sessions…
            </Text>
          ) : entries.length === 0 ? (
            <Text c="dimmed">No sessions logged yet. Tap Log session after each exercise set.</Text>
          ) : (
            <Stack gap="sm" role="list" aria-label="Exercise sessions">
              {entries.map((entry) => {
                const hint = progressionHint(entries, entry);
                return (
                  <div key={entry.id} className="patient-list-card" role="listitem">
                    <Group justify="space-between" align="flex-start" wrap="nowrap">
                      <Stack gap={4} style={{ minWidth: 0 }}>
                        <Text fw={600}>{entry.exerciseName}</Text>
                        <Text size="sm">{formatSessionSummary(entry)}</Text>
                        <Text size="xs" c="dimmed">
                          {new Date(entry.performedAt).toLocaleString()} · {entry.recordedByName}
                        </Text>
                        {hint && (
                          <Text size="sm" c={hint.startsWith('+') ? 'green' : 'dimmed'}>
                            {hint}
                          </Text>
                        )}
                        {entry.notes && (
                          <Text size="sm" c="dimmed">
                            {entry.notes}
                          </Text>
                        )}
                      </Stack>
                      <Button
                        variant="subtle"
                        color="red"
                        size="compact-md"
                        aria-label={`Delete session ${entry.exerciseName}`}
                        onClick={() => setDeleteEntryTarget(entry)}
                      >
                        <IconTrash size={16} />
                      </Button>
                    </Group>
                  </div>
                );
              })}
            </Stack>
          )}
        </Stack>
      )}

      <Modal
        opened={planModalOpen}
        onClose={() => setPlanModalOpen(false)}
        title="New exercise plan"
        centered
        size="md"
        transitionProps={{ duration: 0 }}
      >
        <Stack gap="md">
          <TextInput
            label="Plan name"
            placeholder="e.g. Knee strengthening"
            value={planName}
            onChange={(event) => setPlanName(event.currentTarget.value)}
            size="md"
            className="touch-control"
            required
          />
          <Textarea
            label="Description / goals"
            placeholder="Optional: frequency, precautions, week goals"
            value={planDescription}
            onChange={(event) => setPlanDescription(event.currentTarget.value)}
            minRows={3}
            size="md"
          />
          <Button
            fullWidth
            size="md"
            className="touch-control"
            loading={createPlanMutation.isPending}
            disabled={!planName.trim()}
            onClick={() => createPlanMutation.mutate()}
          >
            Create plan
          </Button>
        </Stack>
      </Modal>

      <Modal
        opened={entryModalOpen}
        onClose={() => setEntryModalOpen(false)}
        title="Log exercise session"
        centered
        size="md"
        transitionProps={{ duration: 0 }}
      >
        <Stack gap="md">
          <TextInput
            label="Exercise"
            placeholder="e.g. Quad sets"
            value={entryExerciseName}
            onChange={(event) => setEntryExerciseName(event.currentTarget.value)}
            size="md"
            className="touch-control"
            required
          />
          {recentExerciseNames.length > 0 && (
            <Group gap="xs">
              {recentExerciseNames.map((name) => (
                <Button
                  key={name}
                  size="compact-md"
                  variant="light"
                  onClick={() => setEntryExerciseName(name)}
                >
                  {name}
                </Button>
              ))}
            </Group>
          )}
          <Group grow align="flex-start">
            <NumberInput
              label="Sets"
              value={entrySets}
              onChange={(value) => setEntrySets(parseNumberInput(value))}
              min={0}
              allowDecimal={false}
              clampBehavior="blur"
              size="md"
              className="touch-control"
            />
            <NumberInput
              label="Reps"
              value={entryReps}
              onChange={(value) => setEntryReps(parseNumberInput(value))}
              min={0}
              allowDecimal={false}
              clampBehavior="blur"
              size="md"
              className="touch-control"
            />
          </Group>
          <TextInput
            label="Resistance / load"
            placeholder="e.g. yellow band, 5 lb"
            value={entryResistance}
            onChange={(event) => setEntryResistance(event.currentTarget.value)}
            size="md"
            className="touch-control"
          />
          <NumberInput
            label="Difficulty (1-10)"
            value={entryDifficulty}
            onChange={(value) => setEntryDifficulty(parseNumberInput(value))}
            min={1}
            max={10}
            allowDecimal={false}
            clampBehavior="blur"
            size="md"
            className="touch-control"
          />
          <Textarea
            label="Notes"
            placeholder="Form cues, pain, patient report"
            value={entryNotes}
            onChange={(event) => setEntryNotes(event.currentTarget.value)}
            minRows={2}
            size="md"
          />
          <Button
            fullWidth
            size="md"
            className="touch-control"
            loading={createEntryMutation.isPending}
            disabled={!entryExerciseName.trim()}
            onClick={() => createEntryMutation.mutate()}
          >
            Save session
          </Button>
        </Stack>
      </Modal>

      <Modal
        opened={Boolean(deletePlanTarget)}
        onClose={() => setDeletePlanTarget(null)}
        title="Delete exercise plan?"
        centered
      >
        <Stack gap="md">
          <Text>
            Delete <strong>{deletePlanTarget?.name}</strong> and all logged sessions? This cannot be
            undone (audit log retains a record).
          </Text>
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setDeletePlanTarget(null)}>
              Cancel
            </Button>
            <Button
              color="red"
              loading={deletePlanMutation.isPending}
              onClick={() => deletePlanTarget && deletePlanMutation.mutate(deletePlanTarget.id)}
            >
              Delete plan
            </Button>
          </Group>
        </Stack>
      </Modal>

      <Modal
        opened={Boolean(deleteEntryTarget)}
        onClose={() => setDeleteEntryTarget(null)}
        title="Delete session?"
        centered
      >
        <Stack gap="md">
          <Text>
            Remove the {deleteEntryTarget?.exerciseName} session from{' '}
            {deleteEntryTarget
              ? new Date(deleteEntryTarget.performedAt).toLocaleString()
              : 'this plan'}
            ?
          </Text>
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setDeleteEntryTarget(null)}>
              Cancel
            </Button>
            <Button
              color="red"
              loading={deleteEntryMutation.isPending}
              onClick={() => deleteEntryTarget && deleteEntryMutation.mutate(deleteEntryTarget)}
            >
              Delete
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Stack>
  );
}
