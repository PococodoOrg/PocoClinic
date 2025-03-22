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
  NumberInput,
} from '@mantine/core';
import { DateInput } from '@mantine/dates';
import { Patient, PatientFormData, Gender } from '../../types/patient';
import { notifications } from '@mantine/notifications';
import { useNavigate } from 'react-router-dom';
import { ValidationError } from '../../api/patients';

interface PatientFormProps {
  initialValues?: Patient;
  onSubmit: (values: PatientFormData) => Promise<any>;
  isLoading?: boolean;
}

export function PatientForm({ initialValues, onSubmit, isLoading }: PatientFormProps) {
  const navigate = useNavigate();
  const form = useForm<PatientFormData>({
    initialValues: initialValues ? {
      ...initialValues,
      dateOfBirth: new Date(initialValues.dateOfBirth),
      height: initialValues.height ?? null,
      weight: initialValues.weight ?? null,
    } : {
      firstName: '',
      lastName: '',
      dateOfBirth: null,
      gender: 'unknown' as Gender,
      email: '',
      phone: '',
      address: '',
      city: '',
      state: '',
      zipCode: '',
      height: null,
      weight: null,
    },

    validate: {
      firstName: (value) => (!value || value.trim().length === 0 ? 'First name is required' : null),
      lastName: (value) => (!value || value.trim().length === 0 ? 'Last name is required' : null),
      dateOfBirth: (value) => (!value ? 'Date of birth is required' : null),
      gender: (value) => (!value ? 'Gender is required' : null),
      email: (value) => {
        if (!value || value.trim().length === 0) return 'Email is required';
        if (!/^\S+@\S+\.\S+$/.test(value)) return 'Invalid email format';
        return null;
      },
      phone: (value) => {
        if (!value || value.trim().length === 0) return 'Phone number is required';
        if (!/^\+?[\d\s-()]+$/.test(value)) return 'Invalid phone number format';
        return null;
      },
      zipCode: (value?: string) => {
        if (!value) return null;
        if (!/^\d{5}(-\d{4})?$/.test(value)) return 'Invalid ZIP code format (e.g., 12345 or 12345-6789)';
        return null;
      },
      height: (value?: number | null) => {
        if (value !== null && value !== undefined && (value <= 0 || value > 300)) return 'Height must be between 1 and 300 cm';
        return null;
      },
      weight: (value?: number | null) => {
        if (value !== null && value !== undefined && (value <= 0 || value > 500)) return 'Weight must be between 1 and 500 kg';
        return null;
      },
    },

    // No need for transformValues as we handle data transformation in the API layer
  });

  const handleSubmit = async (values: PatientFormData) => {
    try {
      await onSubmit(values);
      navigate('/patients');
    } catch (error) {
      if ((error as ValidationError).code === 'VALIDATION_ERROR') {
        const validationError = error as ValidationError;
        if (validationError.errors) {
          // Set field-specific errors
          Object.entries(validationError.errors).forEach(([field, messages]) => {
            form.setFieldError(field, messages[0]);
          });
        }
      } else {
        notifications.show({
          title: 'Error',
          message: error instanceof Error ? error.message : 'Failed to save patient',
          color: 'red'
        });
      }
    }
  };

  return (
    <Box component="form" onSubmit={form.onSubmit(handleSubmit)}>
      <Stack gap="md">
        <Grid>
          <Grid.Col span={6}>
            <TextInput
              required
              label="First Name"
              placeholder="Enter first name"
              {...form.getInputProps('firstName')}
            />
          </Grid.Col>
          <Grid.Col span={6}>
            <TextInput
              required
              label="Last Name"
              placeholder="Enter last name"
              {...form.getInputProps('lastName')}
            />
          </Grid.Col>
        </Grid>

        <Grid>
          <Grid.Col span={6}>
            <DateInput
              required
              label="Date of Birth"
              placeholder="Select date"
              maxDate={new Date()}
              {...form.getInputProps('dateOfBirth')}
            />
          </Grid.Col>
          <Grid.Col span={6}>
            <Select
              required
              label="Gender"
              placeholder="Select gender"
              data={[
                { value: 'male', label: 'Male' },
                { value: 'female', label: 'Female' },
                { value: 'other', label: 'Other' },
                { value: 'unknown', label: 'Unknown' },
              ]}
              {...form.getInputProps('gender')}
            />
          </Grid.Col>
        </Grid>

        <Grid>
          <Grid.Col span={6}>
            <TextInput
              required
              label="Email"
              placeholder="Enter email"
              type="email"
              {...form.getInputProps('email')}
            />
          </Grid.Col>
          <Grid.Col span={6}>
            <TextInput
              required
              label="Phone"
              placeholder="Enter phone number"
              {...form.getInputProps('phone')}
            />
          </Grid.Col>
        </Grid>

        <TextInput
          label="Address"
          placeholder="Enter street address"
          {...form.getInputProps('address')}
        />

        <Grid>
          <Grid.Col span={4}>
            <TextInput
              label="City"
              placeholder="Enter city"
              {...form.getInputProps('city')}
            />
          </Grid.Col>
          <Grid.Col span={4}>
            <TextInput
              label="State"
              placeholder="Enter state"
              {...form.getInputProps('state')}
            />
          </Grid.Col>
          <Grid.Col span={4}>
            <TextInput
              label="ZIP Code"
              placeholder="Enter ZIP code"
              {...form.getInputProps('zipCode')}
            />
          </Grid.Col>
        </Grid>

        <Grid>
          <Grid.Col span={6}>
            <NumberInput
              label="Height (cm)"
              placeholder="Enter height"
              min={1}
              max={300}
              {...form.getInputProps('height')}
            />
          </Grid.Col>
          <Grid.Col span={6}>
            <NumberInput
              label="Weight (kg)"
              placeholder="Enter weight"
              min={1}
              max={500}
              {...form.getInputProps('weight')}
            />
          </Grid.Col>
        </Grid>

        <Group justify="flex-end" mt="xl">
          <Button type="submit" loading={isLoading}>
            {initialValues ? 'Update Patient' : 'Create Patient'}
          </Button>
        </Group>
      </Stack>
    </Box>
  );
} 