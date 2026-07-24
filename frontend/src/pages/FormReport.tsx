import {
  Button,
  Container,
  Group,
  Pagination,
  Select,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { fetchAllTemplateSubmissions, fetchFormTemplates, fetchTemplateSubmissions } from '../api/forms';
import { safeDownloadFilename } from '../utils/safeFilename';
import { getErrorMessage } from '../utils/apiError';
import { formatAnswerValue, FormField, FormSubmissionReport } from '../types/form';

const PAGE_SIZE = 25;

function formatDate(value: string): string {
  return new Date(value).toLocaleString();
}

function buildColumns(items: FormSubmissionReport[]): FormField[] {
  const snapshot = items.find((item) => item.fieldSnapshot?.length)?.fieldSnapshot;
  if (snapshot?.length) {
    return [...snapshot].sort((a, b) => a.sortOrder - b.sortOrder);
  }
  const ids = new Set<string>();
  items.forEach((item) => {
    Object.keys(item.answers ?? {}).forEach((key) => ids.add(key));
  });
  return [...ids].map((id, index) => ({
    id,
    label: id,
    type: 'text' as const,
    required: false,
    sortOrder: index,
  }));
}

function downloadCsv(filename: string, headers: string[], rows: string[][]) {
  const escape = (value: string) => `"${value.replace(/"/g, '""')}"`;
  const lines = [headers.map(escape).join(','), ...rows.map((row) => row.map(escape).join(','))];
  const blob = new Blob([lines.join('\n')], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = safeDownloadFilename(filename);
  link.click();
  URL.revokeObjectURL(url);
}

export default function FormReport() {
  const navigate = useNavigate();
  const [templateId, setTemplateId] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');

  const [exporting, setExporting] = useState(false);

  const { data: templates = [], isLoading: templatesLoading } = useQuery({
    queryKey: ['form-templates'],
    queryFn: fetchFormTemplates,
  });

  const selectedTemplate = templates.find((template) => template.id === templateId);

  const { data: report, isLoading: reportLoading, error } = useQuery({
    queryKey: ['form-report', templateId, page, from, to],
    queryFn: () =>
      fetchTemplateSubmissions(templateId!, {
        page,
        pageSize: PAGE_SIZE,
        from: from || undefined,
        to: to || undefined,
      }),
    enabled: Boolean(templateId),
  });

  const columns = useMemo(() => buildColumns(report?.items ?? []), [report?.items]);
  const totalPages = report ? Math.max(1, Math.ceil(report.total / report.pageSize)) : 1;

  const handleExport = async () => {
    if (!selectedTemplate || !report?.total) {
      return;
    }

    setExporting(true);
    try {
      const allItems = await fetchAllTemplateSubmissions(templateId!, {
        from: from || undefined,
        to: to || undefined,
      });
      const exportColumns = buildColumns(allItems);
      const headers = [
        'Patient',
        'Updated',
        'Submitted by',
        ...exportColumns.map((column) => column.label),
      ];
      const rows = allItems.map((item) => [
        item.patientName || item.patientId,
        formatDate(item.updatedAt),
        item.submittedByName || item.submittedBy,
        ...exportColumns.map((column) => formatAnswerValue(item.answers?.[column.id])),
      ]);
      downloadCsv(`${selectedTemplate.name.replace(/\s+/g, '-').toLowerCase()}-report.csv`, headers, rows);
      notifications.show({
        color: 'green',
        message: `Exported ${allItems.length} submission${allItems.length === 1 ? '' : 's'}.`,
      });
    } catch {
      notifications.show({ color: 'red', message: 'Failed to export report.' });
    } finally {
      setExporting(false);
    }
  };

  return (
    <Container size="xl" py="md">
      <Group justify="space-between" mb="lg">
        <Stack gap={4}>
          <Title order={2}>Form reports</Title>
          <Text c="dimmed" size="sm">
            View current submissions across patients for a selected form template.
          </Text>
        </Stack>
        <Group>
          <Button variant="default" onClick={() => navigate('/forms')}>
            Manage forms
          </Button>
          <Button variant="light" loading={exporting} disabled={!report?.total} onClick={handleExport}>
            Export all CSV ({report?.total ?? 0})
          </Button>
        </Group>
      </Group>

      <Group align="flex-end" mb="md" wrap="wrap">
        <Select
          label="Form template"
          placeholder={templatesLoading ? 'Loading…' : 'Choose a form'}
          data={templates.map((template) => ({ value: template.id, label: template.name }))}
          value={templateId}
          onChange={(value) => {
            setTemplateId(value);
            setPage(1);
          }}
          searchable
          nothingFoundMessage="No forms found"
          w={320}
        />
        <TextInput
          label="Updated from"
          type="date"
          value={from}
          onChange={(event) => {
            setFrom(event.currentTarget.value);
            setPage(1);
          }}
        />
        <TextInput
          label="Updated to"
          type="date"
          value={to}
          onChange={(event) => {
            setTo(event.currentTarget.value);
            setPage(1);
          }}
        />
      </Group>

      {!templateId && (
        <Text c="dimmed">Select a form template to load the report.</Text>
      )}

      {error && (
        <Text c="red" mb="md">
          {getErrorMessage(error, 'Failed to load report')}
        </Text>
      )}

      {templateId && (
        <Stack gap="md">
          <Group justify="space-between">
            <Text size="sm" c="dimmed">
              {reportLoading
                ? 'Loading submissions…'
                : `${report?.total ?? 0} current submission${report?.total === 1 ? '' : 's'}`}
            </Text>
            {report && report.total > report.pageSize && (
              <Pagination value={page} onChange={setPage} total={totalPages} />
            )}
          </Group>

          <Table striped highlightOnHover withTableBorder>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Patient</Table.Th>
                <Table.Th>Updated</Table.Th>
                <Table.Th>Submitted by</Table.Th>
                {columns.map((column) => (
                  <Table.Th key={column.id}>{column.label}</Table.Th>
                ))}
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {!reportLoading && report?.items.length === 0 && (
                <Table.Tr>
                  <Table.Td colSpan={3 + columns.length}>
                    <Text c="dimmed">No submissions match this filter.</Text>
                  </Table.Td>
                </Table.Tr>
              )}
              {report?.items.map((item) => (
                <Table.Tr key={item.id}>
                  <Table.Td>{item.patientName || item.patientId}</Table.Td>
                  <Table.Td>{formatDate(item.updatedAt)}</Table.Td>
                  <Table.Td>{item.submittedByName || item.submittedBy}</Table.Td>
                  {columns.map((column) => (
                    <Table.Td key={column.id}>{formatAnswerValue(item.answers?.[column.id])}</Table.Td>
                  ))}
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>

          {report && report.total > report.pageSize && (
            <Group justify="center">
              <Pagination value={page} onChange={setPage} total={totalPages} />
            </Group>
          )}
        </Stack>
      )}
    </Container>
  );
}
