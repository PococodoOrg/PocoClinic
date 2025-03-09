import { useEffect, useState } from 'react';
import { Table, TextInput, Group, Button, Pagination, Text, LoadingOverlay } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { patientApi } from '../../api/patients';
import { Patient } from '../../types/patient';
import { useNavigate } from 'react-router-dom';

export function PatientList() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
      setPage(1);
    }, 300);
    return () => clearTimeout(timer);
  }, [search]);

  const { data, isLoading } = useQuery({
    queryKey: ['patients', page, debouncedSearch],
    queryFn: () => patientApi.getPatients({ page, pageSize: 10, search: debouncedSearch }),
  });

  const handleRowClick = (patient: Patient) => {
    navigate(`/patients/${patient.id}`);
  };

  return (
    <div style={{ position: 'relative' }}>
      <LoadingOverlay visible={isLoading} />
      
      <Group position="apart" mb="md">
        <TextInput
          placeholder="Search patients..."
          value={search}
          onChange={(event) => setSearch(event.currentTarget.value)}
          style={{ width: '300px' }}
        />
        <Button onClick={() => navigate('/patients/new')}>Add New Patient</Button>
      </Group>

      <Table striped highlightOnHover>
        <thead>
          <tr>
            <th>Medical Number</th>
            <th>Name</th>
            <th>Date of Birth</th>
            <th>Gender</th>
            <th>Contact</th>
          </tr>
        </thead>
        <tbody>
          {data?.patients.map((patient) => (
            <tr key={patient.id} onClick={() => handleRowClick(patient)} style={{ cursor: 'pointer' }}>
              <td>{patient.medicalNumber}</td>
              <td>{`${patient.lastName}, ${patient.firstName}${patient.middleName ? ` ${patient.middleName}` : ''}`}</td>
              <td>{new Date(patient.dateOfBirth).toLocaleDateString()}</td>
              <td style={{ textTransform: 'capitalize' }}>{patient.gender}</td>
              <td>{patient.email || patient.phoneNumber || '-'}</td>
            </tr>
          ))}
        </tbody>
      </Table>

      {data?.patients.length === 0 && (
        <Text color="dimmed" align="center" mt="xl">
          No patients found
        </Text>
      )}

      {data && data.totalPages > 1 && (
        <Group position="center" mt="xl">
          <Pagination
            total={data.totalPages}
            value={page}
            onChange={setPage}
          />
        </Group>
      )}
    </div>
  );
} 