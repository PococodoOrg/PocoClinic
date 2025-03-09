import React from 'react';
import { useForm } from '@mantine/form';
import {
  TextInput,
  Select,
  Group,
  Button,
  Box,
  Stack,
  Grid,
} from '@mantine/core';
import { DateInput } from '@mantine/dates';
import { Patient, PatientFormData, Gender } from '../../types/patient';

interface PatientFormProps {
  initialValues?: Patient;
  onSubmit: (values: PatientFormData) => void;
  isLoading?: boolean;
}

export function PatientForm({ initialValues, onSubmit, isLoading }: PatientFormProps) {
  const form = useForm<PatientFormData>({
    initialValues: initialValues ? {
      ...initialValues,
      dateOfBirth: new Date(initialValues.dateOfBirth),
    } : {
      firstName: '',
      lastName: '',
      middleName: '',
      dateOfBirth: new Date(),
      gender: 'unknown' as Gender,
      email: '',
      phoneNumber: '',
      medicalNumber: '',
      address: {
        street: '',
        city: '',
        state: '',
        postalCode: '',
        country: '',
      },
    },

    validate: {
      firstName: (value) => (value.trim() ? null : 'First name is required'),
      lastName: (value) => (value.trim() ? null : 'Last name is required'),
      dateOfBirth: (value) => (value ? null : 'Date of birth is required'),
      gender: (value) => (value ? null : 'Gender is required'),
      medicalNumber: (value) => (value.trim() ? null : 'Medical number is required'),
      email: (value) => (value ? /^\S+@\S+$/.test(value) ? null : 'Invalid email' : null),
    },
  });

  const handleSubmit = form.onSubmit((values) => {
    onSubmit(values);
  });

  return (
    <Box component="form" onSubmit={handleSubmit}>
      <Stack spacing="md">
        <Grid>
          <Grid.Col span={4}>
            <TextInput
              required
              label="First Name"
              {...form.getInputProps('firstName')}
            />
          </Grid.Col>
          <Grid.Col span={4}>
            <TextInput
              label="Middle Name"
              {...form.getInputProps('middleName')}
            />
          </Grid.Col>
          <Grid.Col span={4}>
            <TextInput
              required
              label="Last Name"
              {...form.getInputProps('lastName')}
            />
          </Grid.Col>
        </Grid>

        <Grid>
          <Grid.Col span={4}>
            <DateInput
              required
              label="Date of Birth"
              {...form.getInputProps('dateOfBirth')}
            />
          </Grid.Col>
          <Grid.Col span={4}>
            <Select
              required
              label="Gender"
              data={[
                { value: 'male', label: 'Male' },
                { value: 'female', label: 'Female' },
                { value: 'other', label: 'Other' },
                { value: 'unknown', label: 'Unknown' },
              ]}
              {...form.getInputProps('gender')}
            />
          </Grid.Col>
          <Grid.Col span={4}>
            <TextInput
              required
              label="Medical Number"
              {...form.getInputProps('medicalNumber')}
            />
          </Grid.Col>
        </Grid>

        <Grid>
          <Grid.Col span={6}>
            <TextInput
              label="Email"
              type="email"
              {...form.getInputProps('email')}
            />
          </Grid.Col>
          <Grid.Col span={6}>
            <TextInput
              label="Phone Number"
              {...form.getInputProps('phoneNumber')}
            />
          </Grid.Col>
        </Grid>

        <Box>
          <TextInput
            label="Street Address"
            {...form.getInputProps('address.street')}
          />
        </Box>

        <Grid>
          <Grid.Col span={4}>
            <TextInput
              label="City"
              {...form.getInputProps('address.city')}
            />
          </Grid.Col>
          <Grid.Col span={3}>
            <TextInput
              label="State"
              {...form.getInputProps('address.state')}
            />
          </Grid.Col>
          <Grid.Col span={2}>
            <TextInput
              label="Postal Code"
              {...form.getInputProps('address.postalCode')}
            />
          </Grid.Col>
          <Grid.Col span={3}>
            <TextInput
              label="Country"
              {...form.getInputProps('address.country')}
            />
          </Grid.Col>
        </Grid>

        <Group position="right" mt="xl">
          <Button type="submit" loading={isLoading}>
            {initialValues ? 'Update Patient' : 'Create Patient'}
          </Button>
        </Group>
      </Stack>
    </Box>
  );
} 