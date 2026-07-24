import { useMemo, useState } from 'react';
import {
  Accordion,
  Badge,
  Button,
  Collapse,
  Group,
  Modal,
  Paper,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import {
  fetchFormEntryHistory,
  fetchFormGroups,
  fetchFormTemplates,
  fetchPatientFormSubmissions,
  getFormTemplate,
  savePatientForm,
} from '../../api/forms';
import { DynamicFormRenderer } from './DynamicFormRenderer';
import {
  FormSubmission,
  FormTemplate,
  formatAnswerValue,
  getFieldLabel,
  groupTemplatesByGroup,
} from '../../types/form';
import { getErrorMessage } from '../../utils/apiError';
import { ViewAllChartLink } from '../patients/ViewAllChartLink';

interface PatientFormsSectionProps {
  patientId: string;
  previewLimit?: number;
}

interface FormModalState {
  templateId: string;
  entryId?: string;
  initialAnswers?: Record<string, string | number | boolean>;
  title: string;
  submitLabel: string;
}

function SubmissionAnswers({ submission }: { submission: FormSubmission }) {
  return (
    <Stack gap="xs">
      {Object.entries(submission.answers).map(([fieldId, value]) => (
        <Group key={fieldId} justify="space-between" align="flex-start">
          <Text size="sm" c="dimmed">
            {getFieldLabel(submission, fieldId)}
          </Text>
          <Text size="sm">{formatAnswerValue(value)}</Text>
        </Group>
      ))}
    </Stack>
  );
}

function EntryHistoryPanel({ patientId, entryId }: { patientId: string; entryId: string }) {
  const { data: history = [], isLoading } = useQuery({
    queryKey: ['form-entry-history', patientId, entryId],
    queryFn: () => fetchFormEntryHistory(patientId, entryId),
  });

  const priorVersions = history.filter((item) => !item.isCurrent);

  if (isLoading) {
    return (
      <Text c="dimmed" size="sm">
        Loading history...
      </Text>
    );
  }

  if (priorVersions.length === 0) {
    return (
      <Text c="dimmed" size="sm">
        No prior versions.
      </Text>
    );
  }

  return (
    <Accordion variant="separated">
      {priorVersions.map((version) => (
        <Accordion.Item key={version.id} value={version.id}>
          <Accordion.Control>
            <Group justify="space-between">
              <Text size="sm">Version {version.version}</Text>
              <Text size="xs" c="dimmed">
                {new Date(version.updatedAt).toLocaleString()}
              </Text>
            </Group>
          </Accordion.Control>
          <Accordion.Panel>
            <SubmissionAnswers submission={version} />
          </Accordion.Panel>
        </Accordion.Item>
      ))}
    </Accordion>
  );
}

function TemplateCard({
  patientId,
  template,
  submissions,
  onOpenForm,
}: {
  patientId: string;
  template: FormTemplate;
  submissions: FormSubmission[];
  onOpenForm: (state: FormModalState) => void;
}) {
  const templateSubmissions = submissions.filter((item) => item.templateId === template.id);

  if (template.formType === 'singleton') {
    const current = templateSubmissions[0];

    return (
      <Paper withBorder p="md">
        <Stack gap="sm">
          <Group justify="space-between">
            <Group gap="xs">
              <Text fw={500}>{template.name}</Text>
              <Badge variant="light">Singleton</Badge>
            </Group>
            <Button
              size="xs"
              variant="light"
              onClick={() =>
                onOpenForm({
                  templateId: template.id,
                  entryId: current?.entryId,
                  initialAnswers: current?.answers,
                  title: current ? `Edit ${template.name}` : template.name,
                  submitLabel: current ? 'Save changes' : 'Create form',
                })
              }
            >
              {current ? 'Edit' : 'Start'}
            </Button>
          </Group>

          {current ? (
            <>
              <Text size="xs" c="dimmed">
                Last updated {new Date(current.updatedAt).toLocaleString()} · Version {current.version}
              </Text>
              <SubmissionAnswers submission={current} />
              <HistoryToggle patientId={patientId} entryId={current.entryId} />
            </>
          ) : (
            <Text c="dimmed" size="sm">
              Not started for this patient.
            </Text>
          )}
        </Stack>
      </Paper>
    );
  }

  return (
    <Paper withBorder p="md">
      <Stack gap="sm">
        <Group justify="space-between">
          <Group gap="xs">
            <Text fw={500}>{template.name}</Text>
            <Badge variant="light" color="teal">
              Log
            </Badge>
          </Group>
          <Button
            size="xs"
            onClick={() =>
              onOpenForm({
                templateId: template.id,
                title: `New ${template.name} entry`,
                submitLabel: 'Add entry',
              })
            }
          >
            Add entry
          </Button>
        </Group>

        {templateSubmissions.length === 0 ? (
          <Text c="dimmed" size="sm">
            No entries yet.
          </Text>
        ) : (
          <Accordion variant="contained">
            {templateSubmissions.map((entry) => (
              <Accordion.Item key={entry.entryId} value={entry.entryId}>
                <Accordion.Control>
                  <Group justify="space-between">
                    <Text size="sm">
                      Entry · Version {entry.version}
                    </Text>
                    <Text size="xs" c="dimmed">
                      {new Date(entry.updatedAt).toLocaleString()}
                    </Text>
                  </Group>
                </Accordion.Control>
                <Accordion.Panel>
                  <Stack gap="sm">
                    <Group justify="flex-end">
                      <Button
                        size="xs"
                        variant="light"
                        onClick={() =>
                          onOpenForm({
                            templateId: template.id,
                            entryId: entry.entryId,
                            initialAnswers: entry.answers,
                            title: `Edit ${template.name} entry`,
                            submitLabel: 'Save changes',
                          })
                        }
                      >
                        Edit entry
                      </Button>
                    </Group>
                    <SubmissionAnswers submission={entry} />
                    <HistoryToggle patientId={patientId} entryId={entry.entryId} />
                  </Stack>
                </Accordion.Panel>
              </Accordion.Item>
            ))}
          </Accordion>
        )}
      </Stack>
    </Paper>
  );
}

function HistoryToggle({ patientId, entryId }: { patientId: string; entryId: string }) {
  const [opened, { toggle }] = useDisclosure(false);

  return (
    <Stack gap="xs">
      <Button variant="subtle" size="xs" onClick={toggle} px={0}>
        {opened ? 'Hide version history' : 'Show version history'}
      </Button>
      <Collapse expanded={opened}>
        <EntryHistoryPanel patientId={patientId} entryId={entryId} />
      </Collapse>
    </Stack>
  );
}

export function PatientFormsSection({ patientId, previewLimit }: PatientFormsSectionProps) {
  const queryClient = useQueryClient();
  const [modalState, setModalState] = useState<FormModalState | null>(null);

  const { data: groups = [], isLoading: groupsLoading } = useQuery({
    queryKey: ['form-groups'],
    queryFn: fetchFormGroups,
  });

  const { data: templates = [], isLoading: templatesLoading } = useQuery({
    queryKey: ['form-templates'],
    queryFn: fetchFormTemplates,
  });

  const { data: submissions = [], isLoading: submissionsLoading } = useQuery({
    queryKey: ['patient-form-submissions', patientId],
    queryFn: () => fetchPatientFormSubmissions(patientId),
  });

  const { data: activeTemplate, isLoading: activeTemplateLoading } = useQuery({
    queryKey: ['form-template', modalState?.templateId],
    queryFn: () => getFormTemplate(modalState!.templateId),
    enabled: Boolean(modalState?.templateId),
  });

  const groupedTemplates = useMemo(
    () => groupTemplatesByGroup(groups, templates),
    [groups, templates],
  );

  const templateNameById = useMemo(() => {
    const map = new Map<string, string>();
    templates.forEach((template) => map.set(template.id, template.name));
    return map;
  }, [templates]);

  const recentSubmissions = useMemo(() => {
    return [...submissions]
      .sort(
        (a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime(),
      )
      .slice(0, previewLimit ?? submissions.length);
  }, [submissions, previewLimit]);

  const isPreview = previewLimit !== undefined;
  const showViewAll = isPreview && submissions.length > previewLimit;

  const saveMutation = useMutation({
    mutationFn: ({
      templateId,
      entryId,
      answers,
    }: {
      templateId: string;
      entryId?: string;
      answers: Record<string, string | number | boolean>;
    }) => savePatientForm(patientId, templateId, answers, entryId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['patient-form-submissions', patientId] });
      await queryClient.invalidateQueries({ queryKey: ['form-entry-history', patientId] });
      notifications.show({
        title: 'Form saved',
        message: 'The form was saved to this patient record.',
        color: 'green',
      });
      setModalState(null);
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Could not save form',
        message: getErrorMessage(error, 'Could not save form'),
        color: 'red',
      });
    },
  });

  const isLoading = groupsLoading || templatesLoading || submissionsLoading;

  return (
    <Stack gap="md">
      <Title order={3}>Patient forms</Title>
      <Text c="dimmed" size="sm">
        Forms are organized by clinic-defined groups. Singleton forms keep one active copy; log forms
        accumulate entries over time.
      </Text>

      {isLoading ? (
        <Text c="dimmed" size="sm">
          Loading forms...
        </Text>
      ) : templates.length === 0 ? (
        <Text c="dimmed" size="sm">
          No form templates are available yet. Admins can create them under Forms in the header.
        </Text>
      ) : isPreview ? (
        <>
          {submissions.length === 0 ? (
            <Text c="dimmed" size="sm">
              No form activity yet for this patient.
            </Text>
          ) : (
            <Stack gap="sm" role="list" aria-label="Recent form activity">
              {recentSubmissions.map((submission) => {
                const name =
                  submission.templateName ?? templateNameById.get(submission.templateId) ?? 'Form';
                const answerPreview = Object.values(submission.answers)[0];
                return (
                  <Paper key={submission.entryId} withBorder p="md" role="listitem">
                    <Stack gap={4}>
                      <Group justify="space-between" wrap="nowrap">
                        <Text fw={600} size="sm" lineClamp={1}>
                          {name}
                        </Text>
                        <Text size="xs" c="dimmed">
                          {new Date(submission.updatedAt).toLocaleString()}
                        </Text>
                      </Group>
                      {answerPreview !== undefined && (
                        <Text size="sm" c="dimmed" lineClamp={2}>
                          {formatAnswerValue(answerPreview)}
                        </Text>
                      )}
                    </Stack>
                  </Paper>
                );
              })}
            </Stack>
          )}
          {showViewAll && (
            <ViewAllChartLink
              to={`/patients/${patientId}/forms`}
              total={submissions.length}
              label="form entries"
            />
          )}
        </>
      ) : (
        groupedTemplates.map(({ group, templates: groupTemplates }) => (
          <Stack key={group.id} gap="sm">
            <Text fw={600}>{group.name}</Text>
            {groupTemplates.length === 0 ? (
              <Text c="dimmed" size="sm">
                No forms in this group.
              </Text>
            ) : (
              groupTemplates.map((template) => (
                <TemplateCard
                  key={template.id}
                  patientId={patientId}
                  template={template}
                  submissions={submissions}
                  onOpenForm={setModalState}
                />
              ))
            )}
          </Stack>
        ))
      )}

      <Modal
        opened={Boolean(modalState)}
        onClose={() => setModalState(null)}
        title={modalState?.title ?? 'Form'}
        size="lg"
      >
        {activeTemplateLoading || !activeTemplate || !modalState ? (
          <Text c="dimmed">Loading form...</Text>
        ) : (
          <DynamicFormRenderer
            template={activeTemplate}
            initialValues={modalState.initialAnswers}
            onSubmit={async (answers) => {
              await saveMutation.mutateAsync({
                templateId: modalState.templateId,
                entryId: modalState.entryId,
                answers,
              });
            }}
            isLoading={saveMutation.isPending}
            submitLabel={modalState.submitLabel}
          />
        )}
      </Modal>
    </Stack>
  );
}
