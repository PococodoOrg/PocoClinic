import React from 'react';
import { useParams } from 'react-router-dom';
import { Text, Container, Title } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { fetchPatient } from '../api/patients';

const PatientDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { data: patient, isLoading, error } = useQuery({
    queryKey: ['patient', id],
    queryFn: () => fetchPatient(id),
  });

  if (isLoading) {
    return <Text>Loading patient details...</Text>;
  }

  if (error) {
    return <Text c="red">Error loading patient details. Please try again later.</Text>;
  }

  if (!patient) {
    return <Text>Patient not found.</Text>;
  }

  return (
    <Container>
      <Title order={2}>Patient Details</Title>
      <Text>First Name: {patient.firstName}</Text>
      <Text>Last Name: {patient.lastName}</Text>
      <Text>Date of Birth: {new Date(patient.dateOfBirth).toLocaleDateString()}</Text>
      <Text>Gender: {patient.gender}</Text>
      <Text>Email: {patient.email}</Text>
      {/* Add more patient details as needed */}
    </Container>
  );
};

export default PatientDetails; 