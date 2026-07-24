import {
  Button,
  FileButton,
  Group,
  Stack,
  Table,
  Text,
  Title,
} from '@mantine/core';
import { IconDownload, IconTrash, IconUpload } from '@tabler/icons-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import {
  deletePatientDocument,
  downloadPatientDocument,
  fetchPatientDocuments,
  formatBytes,
  uploadPatientDocument,
} from '../../api/patientDocuments';
import { getErrorMessage } from '../../utils/apiError';

const ALLOWED_EXTENSIONS = ['.pdf', '.png', '.jpg', '.jpeg', '.gif', '.txt'];

interface PatientDocumentsSectionProps {
  patientId: string;
}

export function PatientDocumentsSection({ patientId }: PatientDocumentsSectionProps) {
  const queryClient = useQueryClient();

  const { data: documents = [], isLoading } = useQuery({
    queryKey: ['patient-documents', patientId],
    queryFn: () => fetchPatientDocuments(patientId),
  });

  const uploadMutation = useMutation({
    mutationFn: (file: File) => uploadPatientDocument(patientId, file),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['patient-documents', patientId] });
      notifications.show({
        title: 'Document uploaded',
        message: 'File saved on the clinic server.',
        color: 'green',
      });
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Upload failed',
        message: getErrorMessage(error, 'Could not upload document'),
        color: 'red',
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (documentId: string) => deletePatientDocument(patientId, documentId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['patient-documents', patientId] });
      notifications.show({
        title: 'Document removed',
        message: 'The file was deleted from the clinic server.',
        color: 'green',
      });
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Delete failed',
        message: getErrorMessage(error, 'Could not delete document'),
        color: 'red',
      });
    },
  });

  const handleFile = (file: File | null) => {
    if (!file) {
      return;
    }
    const lower = file.name.toLowerCase();
    if (!ALLOWED_EXTENSIONS.some((ext) => lower.endsWith(ext))) {
      notifications.show({
        title: 'File type not allowed',
        message: 'Use PDF, PNG, JPEG, GIF, or plain text (.txt). Max 10 MB.',
        color: 'red',
      });
      return;
    }
    if (file.size > 10 * 1024 * 1024) {
      notifications.show({
        title: 'File too large',
        message: 'Maximum size is 10 MB.',
        color: 'red',
      });
      return;
    }
    uploadMutation.mutate(file);
  };

  return (
    <Stack gap="md">
      <Group justify="space-between">
        <Title order={3}>Documents</Title>
        <FileButton onChange={handleFile} accept=".pdf,.png,.jpg,.jpeg,.gif,.txt">
          {(props) => (
            <Button
              {...props}
              leftSection={<IconUpload size={16} />}
              loading={uploadMutation.isPending}
            >
              Upload file
            </Button>
          )}
        </FileButton>
      </Group>
      <Text size="sm" c="dimmed">
        Referrals, consent forms, and scans stored on the clinic server (PDF, images, or text — max 10 MB).
      </Text>

      {isLoading ? (
        <Text size="sm" c="dimmed">Loading documents…</Text>
      ) : documents.length === 0 ? (
        <Text size="sm" c="dimmed">No documents attached yet.</Text>
      ) : (
        <Table withTableBorder>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>File</Table.Th>
              <Table.Th>Size</Table.Th>
              <Table.Th>Uploaded</Table.Th>
              <Table.Th>By</Table.Th>
              <Table.Th />
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {documents.map((doc) => (
              <Table.Tr key={doc.id}>
                <Table.Td>{doc.fileName}</Table.Td>
                <Table.Td>{formatBytes(doc.sizeBytes)}</Table.Td>
                <Table.Td>{new Date(doc.createdAt).toLocaleString()}</Table.Td>
                <Table.Td>{doc.uploadedByName || 'Staff'}</Table.Td>
                <Table.Td>
                  <Group gap="xs" justify="flex-end">
                    <Button
                      size="xs"
                      variant="light"
                      leftSection={<IconDownload size={14} />}
                      onClick={() => downloadPatientDocument(patientId, doc.id, doc.fileName)}
                    >
                      Download
                    </Button>
                    <Button
                      size="xs"
                      variant="light"
                      color="red"
                      leftSection={<IconTrash size={14} />}
                      loading={deleteMutation.isPending}
                      onClick={() => deleteMutation.mutate(doc.id)}
                    >
                      Remove
                    </Button>
                  </Group>
                </Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      )}
    </Stack>
  );
}
