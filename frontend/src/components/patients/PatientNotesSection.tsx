import {
  ActionIcon,
  Button,
  Card,
  Group,
  Modal,
  Stack,
  Text,
  Textarea,
  Timeline,
  Title,
} from '@mantine/core';
import { IconPencil, IconTrash } from '@tabler/icons-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { notifications } from '@mantine/notifications';
import {
  createPatientNote,
  deletePatientNote,
  fetchPatientNotes,
  updatePatientNote,
} from '../../api/patientNotes';
import { useAuth } from '../../context/AuthContext';
import { getErrorMessage } from '../../utils/apiError';
import { PatientNote } from '../../types/patientNote';
import { ViewAllChartLink } from './ViewAllChartLink';

interface PatientNotesSectionProps {
  patientId: string;
  previewLimit?: number;
}

export function PatientNotesSection({ patientId, previewLimit }: PatientNotesSectionProps) {
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const [draft, setDraft] = useState('');
  const [editingNote, setEditingNote] = useState<PatientNote | null>(null);
  const [editBody, setEditBody] = useState('');
  const [deleteTarget, setDeleteTarget] = useState<PatientNote | null>(null);

  const { data: notes = [], isLoading } = useQuery({
    queryKey: ['patient-notes', patientId],
    queryFn: () => fetchPatientNotes(patientId),
  });

  const invalidateNotes = async () => {
    await queryClient.invalidateQueries({ queryKey: ['patient-notes', patientId] });
  };

  const createMutation = useMutation({
    mutationFn: (body: string) => createPatientNote(patientId, body),
    onSuccess: async () => {
      setDraft('');
      await invalidateNotes();
      notifications.show({
        title: 'Note saved',
        message: 'Clinical note added to the chart.',
        color: 'green',
      });
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Could not save note',
        message: getErrorMessage(error, 'Could not save note'),
        color: 'red',
      });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ noteId, body }: { noteId: string; body: string }) =>
      updatePatientNote(patientId, noteId, body),
    onSuccess: async () => {
      setEditingNote(null);
      setEditBody('');
      await invalidateNotes();
      notifications.show({
        title: 'Note updated',
        message: 'Your changes were saved.',
        color: 'green',
      });
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Could not update note',
        message: getErrorMessage(error, 'Could not update note'),
        color: 'red',
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (noteId: string) => deletePatientNote(patientId, noteId),
    onSuccess: async () => {
      setDeleteTarget(null);
      await invalidateNotes();
      notifications.show({
        title: 'Note deleted',
        message: 'The note was removed from the chart. The audit log retains a record.',
        color: 'green',
      });
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Could not delete note',
        message: getErrorMessage(error, 'Could not delete note'),
        color: 'red',
      });
    },
  });

  const openEdit = (note: PatientNote) => {
    setEditingNote(note);
    setEditBody(note.body);
  };

  const isAuthor = (note: PatientNote) => user?.id === note.authorId;
  const isAdmin = user?.role === 'admin';
  const canEdit = (note: PatientNote) => isAuthor(note);
  const canDelete = (note: PatientNote) => isAuthor(note) || isAdmin;
  const visibleNotes = previewLimit ? notes.slice(0, previewLimit) : notes;
  const showViewAll = previewLimit !== undefined && notes.length > previewLimit;

  return (
    <Stack gap="md">
      <Title order={3}>Clinical notes</Title>
      <Text size="sm" c="dimmed">
        Chart notes are visible to all staff. Edit or delete your own notes; administrators can delete any note. All changes are recorded in the audit log.
      </Text>

      <Card withBorder padding="md">
        <Stack gap="sm">
          <Textarea
            label="New note"
            placeholder="Enter visit summary, observations, or follow-up instructions…"
            minRows={3}
            maxLength={10000}
            value={draft}
            onChange={(event) => setDraft(event.currentTarget.value)}
          />
          <Group justify="flex-end">
            <Button
              loading={createMutation.isPending}
              disabled={!draft.trim()}
              onClick={() => createMutation.mutate(draft.trim())}
            >
              Add note
            </Button>
          </Group>
        </Stack>
      </Card>

      {isLoading ? (
        <Text c="dimmed" size="sm">Loading notes…</Text>
      ) : notes.length === 0 ? (
        <Text c="dimmed" size="sm">No clinical notes yet.</Text>
      ) : (
        <Timeline active={visibleNotes.length} bulletSize={24} lineWidth={2}>
          {visibleNotes.map((note) => (
            <Timeline.Item
              key={note.id}
              title={
                <Group gap="xs" wrap="nowrap">
                  <Text fw={600} size="sm" component="span">
                    {note.authorName || 'Staff'}
                  </Text>
                  {canEdit(note) && (
                    <ActionIcon
                      variant="subtle"
                      size="sm"
                      aria-label="Edit note"
                      onClick={() => openEdit(note)}
                    >
                      <IconPencil size={14} />
                    </ActionIcon>
                  )}
                  {canDelete(note) && (
                    <ActionIcon
                      variant="subtle"
                      size="sm"
                      color="red"
                      aria-label="Delete note"
                      onClick={() => setDeleteTarget(note)}
                    >
                      <IconTrash size={14} />
                    </ActionIcon>
                  )}
                </Group>
              }
              lineVariant="dashed"
            >
              <Text size="sm" style={{ whiteSpace: 'pre-wrap' }}>
                {note.body}
              </Text>
              <Text size="xs" c="dimmed" mt={4}>
                {new Date(note.createdAt).toLocaleString()}
                {note.updatedAt && (
                  <> · edited {new Date(note.updatedAt).toLocaleString()}</>
                )}
              </Text>
            </Timeline.Item>
          ))}
        </Timeline>
      )}

      {showViewAll && (
        <ViewAllChartLink
          to={`/patients/${patientId}/notes`}
          total={notes.length}
          label="notes"
        />
      )}

      <Modal
        opened={editingNote !== null}
        onClose={() => setEditingNote(null)}
        title="Edit note"
      >
        <Stack gap="md">
          <Textarea
            label="Note"
            minRows={4}
            maxLength={10000}
            value={editBody}
            onChange={(event) => setEditBody(event.currentTarget.value)}
          />
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setEditingNote(null)}>Cancel</Button>
            <Button
              loading={updateMutation.isPending}
              disabled={!editBody.trim()}
              onClick={() => {
                if (editingNote) {
                  updateMutation.mutate({ noteId: editingNote.id, body: editBody.trim() });
                }
              }}
            >
              Save changes
            </Button>
          </Group>
        </Stack>
      </Modal>

      <Modal
        opened={deleteTarget !== null}
        onClose={() => setDeleteTarget(null)}
        title="Delete note?"
      >
        <Stack gap="md">
          <Text size="sm">
            This removes the note from the patient chart. The audit log keeps a record that it was deleted.
          </Text>
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setDeleteTarget(null)}>Cancel</Button>
            <Button
              color="red"
              loading={deleteMutation.isPending}
              onClick={() => {
                if (deleteTarget) {
                  deleteMutation.mutate(deleteTarget.id);
                }
              }}
            >
              Delete note
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Stack>
  );
}
