import { useMemo, useState } from 'react';
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
  Loader,
} from '@mantine/core';
import { DateInput } from '@mantine/dates';
import { useQuery } from '@tanstack/react-query';
import {
  Patient,
  PatientFormData,
  Gender,
  PatientFieldRequirements,
  DEFAULT_PATIENT_FIELD_REQUIREMENTS,
} from '../../types/patient';
import { notifications } from '@mantine/notifications';
import { useNavigate } from 'react-router-dom';
import { ValidationError } from '../../api/patients';
import { fetchPatientFieldRequirements } from '../../api/patients';
import { getErrorMessage } from '../../utils/apiError';

interface PatientFormProps {
  initialValues?: Patient;
  onSubmit: (values: PatientFormData) => Promise<any>;
  isLoading?: boolean;
}

type MeasurementUnit = 'metric' | 'standard';

const inchesToCm = (inches: number) => inches * 2.54;
const lbsToKg = (lbs: number) => lbs * 0.453592;

function buildValidators(
  requirements: PatientFieldRequirements,
  heightUnit: MeasurementUnit,
  weightUnit: MeasurementUnit,
) {
  const requiredText = (required: boolean, label: string) => (value?: string | null) => {
    if (!required) {
      return null;
    }
    if (!value || value.trim().length === 0) {
      return `${label} is required`;
    }
    return null;
  };

  return {
    firstName: requiredText(requirements.firstName, 'First name'),
    lastName: requiredText(requirements.lastName, 'Last name'),
    middleName: requiredText(requirements.middleName, 'Middle name'),
    dateOfBirth: (value: string | null) => {
      if (!requirements.dateOfBirth) {
        return null;
      }
      return !value ? 'Date of birth is required' : null;
    },
    gender: (value: Gender | null | undefined) => {
      if (!requirements.gender) {
        return null;
      }
      return !value ? 'Gender is required' : null;
    },
    email: (value?: string) => {
      if (requirements.email && (!value || value.trim().length === 0)) {
        return 'Email is required';
      }
      if (value && value.trim().length > 0 && !/^\S+@\S+\.\S+$/.test(value)) {
        return 'Invalid email format';
      }
      return null;
    },
    phoneNumber: (value?: string) => {
      if (requirements.phoneNumber && (!value || value.trim().length === 0)) {
        return 'Phone number is required';
      }
      if (value && value.trim().length > 0 && !/^\+?[\d\s-()]+$/.test(value)) {
        return 'Invalid phone number format';
      }
      return null;
    },
    street: requiredText(requirements.addressStreet, 'Street address'),
    city: requiredText(requirements.addressCity, 'City'),
    state: requiredText(requirements.addressState, 'State'),
    zipCode: (value?: string) => {
      if (requirements.addressPostalCode && (!value || value.trim().length === 0)) {
        return 'ZIP code is required';
      }
      if (value && value.trim().length > 0 && !/^\d{5}(-\d{4})?$/.test(value)) {
        return 'Invalid ZIP code format (e.g., 12345 or 12345-6789)';
      }
      return null;
    },
    height: (value?: number | null) => {
      if (requirements.height && (value === null || value === undefined)) {
        return 'Height is required';
      }
      if (value !== null && value !== undefined) {
        const cmValue = heightUnit === 'standard' ? inchesToCm(value) : value;
        if (cmValue <= 0 || cmValue > 300) {
          return 'Invalid height';
        }
      }
      return null;
    },
    weight: (value?: number | null) => {
      if (requirements.weight && (value === null || value === undefined)) {
        return 'Weight is required';
      }
      if (value !== null && value !== undefined) {
        const kgValue = weightUnit === 'standard' ? lbsToKg(value) : value;
        if (kgValue <= 0 || kgValue > 500) {
          return 'Invalid weight';
        }
      }
      return null;
    },
  };
}

export function PatientForm({ initialValues, onSubmit, isLoading }: PatientFormProps) {
  const navigate = useNavigate();
  const [heightUnit, setHeightUnit] = useState<MeasurementUnit>('metric');
  const [weightUnit, setWeightUnit] = useState<MeasurementUnit>('metric');

  const { data: requirements = DEFAULT_PATIENT_FIELD_REQUIREMENTS, isLoading: requirementsLoading } = useQuery({
    queryKey: ['patient-field-requirements'],
    queryFn: fetchPatientFieldRequirements,
    staleTime: 5 * 60 * 1000,
  });

  const validators = useMemo(
    () => buildValidators(requirements, heightUnit, weightUnit),
    [requirements, heightUnit, weightUnit],
  );

  const form = useForm<PatientFormData>({
    initialValues: initialValues ? {
      ...initialValues,
      dateOfBirth: initialValues.dateOfBirth,
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
      middleName: '',
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
    validate: validators,
  });

  const handleSubmit = async (values: PatientFormData) => {
    try {
      const submitValues: PatientFormData = {
        ...values,
        phoneNumber: values.phoneNumber?.trim() ?? '',
        address: values.street?.trim() || values.city?.trim() || values.state?.trim() || values.zipCode?.trim() ? {
          street: values.street?.trim() || '',
          city: values.city?.trim() || '',
          state: values.state?.trim() || '',
          postalCode: values.zipCode?.trim() || '',
          country: 'US',
        } : undefined,
        street: undefined,
        city: undefined,
        state: undefined,
        zipCode: undefined,
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
          message: getErrorMessage(error, 'Failed to save patient'),
          color: 'red',
        });
      }
    }
  };

  if (requirementsLoading) {
    return (
      <Group justify="center" py="xl">
        <Loader size="sm" />
        <Text size="sm" c="dimmed">Loading form settings…</Text>
      </Group>
    );
  }

  return (
    <Box component="form" onSubmit={form.onSubmit(handleSubmit)}>
      <Stack gap="md">
        <Grid>
          <Grid.Col span={6}>
            <TextInput
              required={requirements.firstName}
              label="First Name"
              placeholder="Enter first name"
              {...form.getInputProps('firstName')}
            />
          </Grid.Col>
          <Grid.Col span={6}>
            <TextInput
              required={requirements.lastName}
              label="Last Name"
              placeholder="Enter last name"
              {...form.getInputProps('lastName')}
            />
          </Grid.Col>
        </Grid>

        <TextInput
          required={requirements.middleName}
          label="Middle Name"
          placeholder="Enter middle name"
          {...form.getInputProps('middleName')}
        />

        <Grid>
          <Grid.Col span={6}>
            <DateInput
              required={requirements.dateOfBirth}
              label="Date of Birth"
              placeholder="Select date"
              maxDate={new Date()}
              {...form.getInputProps('dateOfBirth')}
            />
          </Grid.Col>
          <Grid.Col span={6}>
            <Select
              required={requirements.gender}
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
              required={requirements.email}
              label="Email"
              placeholder="Enter email"
              type="email"
              {...form.getInputProps('email')}
            />
          </Grid.Col>
          <Grid.Col span={6}>
            <TextInput
              required={requirements.phoneNumber}
              label="Phone Number"
              placeholder="Enter phone number"
              {...form.getInputProps('phoneNumber')}
            />
          </Grid.Col>
        </Grid>

        <TextInput
          required={requirements.addressStreet}
          label="Address"
          placeholder="Enter street address"
          {...form.getInputProps('street')}
        />

        <Grid>
          <Grid.Col span={4}>
            <TextInput
              required={requirements.addressCity}
              label="City"
              placeholder="Enter city"
              {...form.getInputProps('city')}
            />
          </Grid.Col>
          <Grid.Col span={4}>
            <TextInput
              required={requirements.addressState}
              label="State"
              placeholder="Enter state"
              {...form.getInputProps('state')}
            />
          </Grid.Col>
          <Grid.Col span={4}>
            <TextInput
              required={requirements.addressPostalCode}
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
                <Text size="sm" fw={500}>
                  Height{requirements.height ? ' *' : ''}
                </Text>
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
                <Text size="sm" fw={500}>
                  Weight{requirements.weight ? ' *' : ''}
                </Text>
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
