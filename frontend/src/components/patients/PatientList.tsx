import { useState } from 'react';
import {
  Badge,
  Button,
  Group,
  Paper,
  Select,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchPatients, PatientListParams } from '../../api/patients';
import { IconChevronRight, IconEye, IconSearch, IconUserPlus, IconX } from '@tabler/icons-react';
import { TABLET_MEDIA_QUERY } from '../../layout/breakpoints';
import { Gender, Patient } from '../../types/patient';

const GENDER_OPTIONS = [
  { value: '', label: 'Any gender' },
  { value: 'male', label: 'Male' },
  { value: 'female', label: 'Female' },
  { value: 'other', label: 'Other' },
  { value: 'unknown', label: 'Unknown' },
];

const EMPTY_FILTERS: PatientListParams = {
  search: '',
  gender: '',
  dobFrom: '',
  dobTo: '',
  registeredSince: '',
};

function formatDob(value: string): string {
  return new Date(value).toLocaleDateString();
}

export function PatientList() {
  const navigate = useNavigate();
  const isTablet = useMediaQuery(TABLET_MEDIA_QUERY);
  const [filters, setFilters] = useState<PatientListParams>(EMPTY_FILTERS);
  const [page, setPage] = useState(1);
  const pageSize = isTablet ? 12 : 10;

  const activeFilterCount = [
    filters.gender,
    filters.dobFrom,
    filters.dobTo,
    filters.registeredSince,
  ].filter(Boolean).length;

  const { data, isLoading, error } = useQuery({
    queryKey: ['patients', filters, page, pageSize],
    queryFn: () => fetchPatients({ ...filters, page, pageSize }),
  });

  const updateFilter = (patch: Partial<PatientListParams>) => {
    setFilters((current) => ({ ...current, ...patch }));
    setPage(1);
  };

  const clearFilters = () => {
    setFilters(EMPTY_FILTERS);
    setPage(1);
  };

  const openPatient = (patientId: string) => {
    navigate(`/patients/${patientId}`);
  };

  if (error) {
    return (
      <Text c="red" role="alert">
        Error loading patients. Please try again later.
      </Text>
    );
  }

  const patients = data?.patients ?? [];

  return (
    <Stack gap="md" className="patient-list-page">
      <Paper withBorder p="md" className="workstation-page-header">
        <Stack gap="md">
          <Group justify="space-between" align="flex-start" wrap="wrap" gap="md">
            <div>
              <Title order={2}>Patients</Title>
              <Text size="sm" c="dimmed" mt={4}>
                Find a chart, then tap to open.
              </Text>
            </div>
            <Button
              leftSection={<IconUserPlus size={18} />}
              onClick={() => navigate('/patients/new')}
              size={isTablet ? 'md' : 'sm'}
              className="touch-control"
            >
              New patient
            </Button>
          </Group>

          <TextInput
            placeholder="Search by name or email…"
            leftSection={<IconSearch size="1.1rem" aria-hidden />}
            value={filters.search ?? ''}
            onChange={(event) => updateFilter({ search: event.currentTarget.value })}
            aria-label="Search patients by name or email"
            size={isTablet ? 'md' : 'sm'}
            className="touch-control"
          />

          {isTablet ? (
            <div className="patient-list-filters">
              <Select
                label="Gender"
                data={GENDER_OPTIONS}
                value={filters.gender ?? ''}
                onChange={(value) => updateFilter({ gender: (value ?? '') as Gender | '' })}
                allowDeselect={false}
                size="md"
                className="touch-control"
              />
              <TextInput
                label="DOB from"
                type="date"
                value={filters.dobFrom ?? ''}
                onChange={(event) => updateFilter({ dobFrom: event.currentTarget.value })}
                size="md"
                className="touch-control"
              />
              <TextInput
                label="DOB to"
                type="date"
                value={filters.dobTo ?? ''}
                onChange={(event) => updateFilter({ dobTo: event.currentTarget.value })}
                size="md"
                className="touch-control"
              />
              <TextInput
                className="filter-span-2 touch-control"
                label="Registered since"
                type="date"
                value={filters.registeredSince ?? ''}
                onChange={(event) => updateFilter({ registeredSince: event.currentTarget.value })}
                size="md"
              />
              {activeFilterCount > 0 && (
                <Button
                  variant="subtle"
                  leftSection={<IconX size={16} />}
                  onClick={clearFilters}
                  size="md"
                  className="touch-control filter-span-2"
                >
                  Clear filters
                </Button>
              )}
            </div>
          ) : (
            <Group align="flex-end" wrap="wrap" gap="md">
              <Select
                label="Gender"
                data={GENDER_OPTIONS}
                value={filters.gender ?? ''}
                onChange={(value) => updateFilter({ gender: (value ?? '') as Gender | '' })}
                w={160}
                allowDeselect={false}
              />
              <TextInput
                label="DOB from"
                type="date"
                value={filters.dobFrom ?? ''}
                onChange={(event) => updateFilter({ dobFrom: event.currentTarget.value })}
                w={160}
              />
              <TextInput
                label="DOB to"
                type="date"
                value={filters.dobTo ?? ''}
                onChange={(event) => updateFilter({ dobTo: event.currentTarget.value })}
                w={160}
              />
              <TextInput
                label="Registered since"
                type="date"
                value={filters.registeredSince ?? ''}
                onChange={(event) => updateFilter({ registeredSince: event.currentTarget.value })}
                w={180}
              />
              {activeFilterCount > 0 && (
                <Button variant="subtle" leftSection={<IconX size={14} />} onClick={clearFilters}>
                  Clear filters
                </Button>
              )}
            </Group>
          )}
        </Stack>
      </Paper>

      {isLoading ? (
        <Text c="dimmed" ta="center" py="xl" aria-live="polite">
          Loading patients…
        </Text>
      ) : patients.length === 0 ? (
        <Paper withBorder p="xl" ta="center">
          <Text c="dimmed" mb="md">
            No patients match these filters.
          </Text>
          <Button variant="light" onClick={() => navigate('/patients/new')} className="touch-control">
            Register a new patient
          </Button>
        </Paper>
      ) : isTablet ? (
        <Stack gap="sm" role="list" aria-label="Patient results">
          {patients.map((patient: Patient) => (
            <button
              key={patient.id}
              type="button"
              className="patient-list-card"
              role="listitem"
              onClick={() => openPatient(patient.id)}
              aria-label={`Open chart for ${patient.firstName} ${patient.lastName}`}
            >
              <Group justify="space-between" wrap="nowrap" align="center">
                <Stack gap={4} style={{ minWidth: 0 }}>
                  <Text fw={600} size="lg" lineClamp={1}>
                    {patient.lastName}, {patient.firstName}
                  </Text>
                  <Text size="sm" c="dimmed">
                    DOB {formatDob(patient.dateOfBirth)}
                    {patient.email ? ` · ${patient.email}` : ''}
                  </Text>
                  <Badge variant="light" size="sm" w="fit-content" tt="capitalize">
                    {patient.gender}
                  </Badge>
                </Stack>
                <IconChevronRight size={22} aria-hidden style={{ flexShrink: 0, opacity: 0.55 }} />
              </Group>
            </button>
          ))}
        </Stack>
      ) : (
        <Table.ScrollContainer minWidth={860}>
          <Table striped highlightOnHover withTableBorder stickyHeader>
            <Table.Thead>
              <Table.Tr>
                <Table.Th scope="col">Last name</Table.Th>
                <Table.Th scope="col">First name</Table.Th>
                <Table.Th scope="col">Date of birth</Table.Th>
                <Table.Th scope="col">Gender</Table.Th>
                <Table.Th scope="col" className="patient-list-email">
                  Email
                </Table.Th>
                <Table.Th scope="col">Registered</Table.Th>
                <Table.Th scope="col">
                  <span className="sr-only">Actions</span>
                </Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {patients.map((patient: Patient) => (
                <Table.Tr
                  key={patient.id}
                  tabIndex={0}
                  aria-label={`Open chart for ${patient.firstName} ${patient.lastName}`}
                  onDoubleClick={() => openPatient(patient.id)}
                  onKeyDown={(event) => {
                    if (event.key === 'Enter') {
                      openPatient(patient.id);
                    }
                  }}
                  className="patient-list-row"
                >
                  <Table.Td>{patient.lastName}</Table.Td>
                  <Table.Td>{patient.firstName}</Table.Td>
                  <Table.Td>{formatDob(patient.dateOfBirth)}</Table.Td>
                  <Table.Td>
                    <Badge variant="light" size="sm">
                      {patient.gender}
                    </Badge>
                  </Table.Td>
                  <Table.Td className="patient-list-email">{patient.email || '—'}</Table.Td>
                  <Table.Td>{formatDob(patient.createdAt)}</Table.Td>
                  <Table.Td>
                    <Button
                      variant="light"
                      size="compact-sm"
                      leftSection={<IconEye size={16} />}
                      onClick={() => openPatient(patient.id)}
                      aria-label={`View ${patient.firstName} ${patient.lastName}`}
                    >
                      Open
                    </Button>
                  </Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        </Table.ScrollContainer>
      )}

      {data && data.totalPages > 1 && (
        <Group justify="center" mt="md" role="navigation" aria-label="Patient list pagination" gap="md">
          <Button
            variant="outline"
            disabled={page <= 1}
            onClick={() => setPage(page - 1)}
            size={isTablet ? 'md' : 'sm'}
            className="touch-control"
          >
            Previous
          </Button>
          <Text size="sm" aria-live="polite">
            Page {page} of {data.totalPages}
          </Text>
          <Button
            variant="outline"
            disabled={page >= data.totalPages}
            onClick={() => setPage(page + 1)}
            size={isTablet ? 'md' : 'sm'}
            className="touch-control"
          >
            Next
          </Button>
        </Group>
      )}

      {data && data.totalCount > 0 && (
        <Text size="xs" c="dimmed" ta="center">
          {data.totalCount} patient{data.totalCount === 1 ? '' : 's'} matching
          {activeFilterCount > 0 ? ' filters' : ''}
          {isTablet ? ' · tap a row to open' : ' · double-click or press Enter on a row to open'}
        </Text>
      )}
    </Stack>
  );
}
