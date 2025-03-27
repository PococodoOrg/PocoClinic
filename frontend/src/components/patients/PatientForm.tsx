import React, { useState } from 'react';
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
  SegmentedControl,
  Text,
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

type MeasurementUnit = 'metric' | 'standard';

// Conversion functions
const inchesToCm = (inches: number) => inches * 2.54;
const cmToInches = (cm: number) => cm / 2.54;
const lbsToKg = (lbs: number) => lbs * 0.453592;
const kgToLbs = (kg: number) => kg / 0.453592;

export function PatientForm({ initialValues, onSubmit, isLoading }: PatientFormProps) {
  const navigate = useNavigate();
  const [heightUnit, setHeightUnit] = useState<MeasurementUnit>('metric');
  const [weightUnit, setWeightUnit] = useState<MeasurementUnit>('metric');

  const form = useForm<PatientFormData>({
    initialValues: initialValues ? {
      ...initialValues,
      dateOfBirth: new Date(initialValues.dateOfBirth),
      height: initialValues.height ?? null,
      weight: initialValues.weight ?? null,
      phoneNumber: initialValues.phoneNumber,
      address: initialValues.address,
      street: initialValues.address?.street || '',
      city: initialValues.address?.city || '',
      state: initialValues.address?.state || '',
      zipCode: initialValues.address?.postalCode || '',
    } : {
      firstName: '',
      lastName: '',
      dateOfBirth: null,
      gender: 'unknown' as Gender,
      email: '',
      phoneNumber: '',
      address: undefined,
      street: '',
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
      phoneNumber: (value) => {
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
        if (value !== null && value !== undefined) {
          const cmValue = heightUnit === 'standard' ? inchesToCm(value) : value;
          if (cmValue <= 0 || cmValue > 300) return 'Invalid height';
        }
        return null;
      },
      weight: (value?: number | null) => {
        if (value !== null && value !== undefined) {
          const kgValue = weightUnit === 'standard' ? lbsToKg(value) : value;
          if (kgValue <= 0 || kgValue > 500) return 'Invalid weight';
        }
        return null;
      },
    },
  });

  const handleSubmit = async (values: PatientFormData) => {
    try {
      // Convert measurements to metric before submitting
      const submitValues: PatientFormData = {
        ...values,
        phoneNumber: values.phoneNumber.trim(), // Ensure phoneNumber is trimmed
        // Structure address fields into an Address object
        address: values.street?.trim() || values.city?.trim() || values.state?.trim() || values.zipCode?.trim() ? {
          street: values.street?.trim() || '',
          city: values.city?.trim() || '',
          state: values.state?.trim() || '',
          postalCode: values.zipCode?.trim() || '',
          country: 'US' // Default to US for now
        } : undefined,
        // Remove individual address fields
        street: undefined,
        city: undefined,
        state: undefined,
        zipCode: undefined,
        // Convert measurements to metric
        height: values.height ? (heightUnit === 'standard' ? inchesToCm(values.height) : values.height) : null,
        weight: values.weight ? (weightUnit === 'standard' ? lbsToKg(values.weight) : values.weight) : null,
      };
      await onSubmit(submitValues);
      navigate('/patients');
    } catch (error) {
      if ((error as ValidationError).code === 'VALIDATION_ERROR') {
        const validationError = error as ValidationError;
        if (validationError.errors) {
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
              label="Phone Number"
              placeholder="Enter phone number"
              {...form.getInputProps('phoneNumber')}
            />
          </Grid.Col>
        </Grid>

        <TextInput
          label="Address"
          placeholder="Enter street address"
          {...form.getInputProps('street')}
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
            <Stack gap="xs">
              <Group justify="space-between">
                <Text size="sm" fw={500}>Height</Text>
                <SegmentedControl
                  size="xs"
                  value={heightUnit}
                  onChange={(value) => setHeightUnit(value as MeasurementUnit)}
                  data={[
                    { label: 'cm', value: 'metric' },
                    { label: 'in', value: 'standard' },
                  ]}
                />
              </Group>
              <NumberInput
                placeholder={`Enter height (${heightUnit === 'metric' ? 'cm' : 'inches'})`}
                min={1}
                max={heightUnit === 'metric' ? 300 : 120}
                allowDecimal={false}
                value={form.values.height ?? ''}
                onChange={(value) => form.setFieldValue('height', value === '' ? null : Number(value))}
                error={form.errors.height}
              />
            </Stack>
          </Grid.Col>
          <Grid.Col span={6}>
            <Stack gap="xs">
              <Group justify="space-between">
                <Text size="sm" fw={500}>Weight</Text>
                <SegmentedControl
                  size="xs"
                  value={weightUnit}
                  onChange={(value) => setWeightUnit(value as MeasurementUnit)}
                  data={[
                    { label: 'kg', value: 'metric' },
                    { label: 'lbs', value: 'standard' },
                  ]}
                />
              </Group>
              <NumberInput
                placeholder={`Enter weight (${weightUnit === 'metric' ? 'kg' : 'lbs'})`}
                min={1}
                max={weightUnit === 'metric' ? 500 : 1100}
                allowDecimal={false}
                value={form.values.weight ?? ''}
                onChange={(value) => form.setFieldValue('weight', value === '' ? null : Number(value))}
                error={form.errors.weight}
              />
            </Stack>
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