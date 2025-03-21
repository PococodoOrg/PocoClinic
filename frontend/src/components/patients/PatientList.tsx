import { useState } from 'react';
import { Table, Button, Group, Text, TextInput } from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchPatients } from '../../api/patients';
import { IconSearch, IconEye } from '@tabler/icons-react';
import { Patient } from '../../types/patient';

export function PatientList() {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const { data, isLoading, error } = useQuery({
    queryKey: ['patients', search, page],
    queryFn: () => fetchPatients({ search, page, pageSize })
  });

  if (error) {
    return (
      <Text c="red">Error loading patients. Please try again later.</Text>
    );
  }

  const handleViewPatient = (patientId: string) => {
    navigate(`/patients/${patientId}`);
  };

  const rows = data?.patients?.map((patient: Patient) => (
    <Table.Tr key={patient.id} onDoubleClick={() => handleViewPatient(patient.id)}>
      <Table.Td>{patient.firstName}</Table.Td>
      <Table.Td>{patient.lastName}</Table.Td>
      <Table.Td>{new Date(patient.dateOfBirth).toLocaleDateString()}</Table.Td>
      <Table.Td>{patient.gender}</Table.Td>
      <Table.Td>{patient.email}</Table.Td>
      <Table.Td>
        <Button variant="light" onClick={() => handleViewPatient(patient.id)}>
          <IconEye size={16} />
        </Button>
      </Table.Td>
    </Table.Tr>
  )) ?? [];

  const handlePreviousPage = () => {
    if (page > 1) {
      setPage(page - 1);
    }
  };

  const handleNextPage = () => {
    if (data && page < data.totalPages) {
      setPage(page + 1);
    }
  };

  return (
    <div>
      <Group justify="space-between" mb="md">
        <TextInput
          placeholder="Search patients..."
          leftSection={<IconSearch size="1rem" />}
          value={search}
          onChange={(event) => setSearch(event.currentTarget.value)}
        />
        <Button onClick={() => navigate('/patients/new')}>
          Add Patient
        </Button>
      </Group>

      {isLoading ? (
        <Text c="dimmed" ta="center" mt="xl">
          Loading patients...
        </Text>
      ) : rows.length === 0 ? (
        <Text c="dimmed" ta="center" mt="xl">
          No patients found.
        </Text>
      ) : (
        <Table>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>First Name</Table.Th>
              <Table.Th>Last Name</Table.Th>
              <Table.Th>Date of Birth</Table.Th>
              <Table.Th>Gender</Table.Th>
              <Table.Th>Email</Table.Th>
              <Table.Th>Actions</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>{rows}</Table.Tbody>
        </Table>
      )}

      {data && data.totalPages > 1 && (
        <Group justify="center" mt="xl">
          <Button 
            variant="outline" 
            disabled={page <= 1}
            onClick={handlePreviousPage}
          >
            Previous
          </Button>
          <Text>
            Page {page} of {data.totalPages}
          </Text>
          <Button 
            variant="outline" 
            disabled={page >= data.totalPages}
            onClick={handleNextPage}
          >
            Next
          </Button>
        </Group>
      )}
    </div>
  );
} 