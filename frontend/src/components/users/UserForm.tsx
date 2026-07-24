import { Button, Select, Stack, Switch, TextInput } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useEffect } from 'react';
import { STAFF_ROLES, StaffUser, UpdateUserFormData, UserRole } from '../../types/user';

interface UserFormProps {
  onSubmit: (data: UpdateUserFormData) => Promise<void>;
  isLoading?: boolean;
  initialValues?: StaffUser;
  submitLabel?: string;
  includeActiveToggle?: boolean;
}

export function UserForm({
  onSubmit,
  isLoading,
  initialValues,
  submitLabel = 'Create staff member',
  includeActiveToggle = false,
}: UserFormProps) {
  const isEdit = includeActiveToggle || Boolean(initialValues);

  const form = useForm<UpdateUserFormData>({
    initialValues: {
      email: initialValues?.email ?? '',
      name: initialValues?.name ?? '',
      role: (initialValues?.role ?? 'staff') as UserRole,
      isActive: initialValues?.isActive ?? true,
    },
    validate: {
      email: (value) => (/^\S+@\S+\.\S+$/.test(value) ? null : 'Enter a valid email'),
      name: (value) => (value.trim().length > 0 ? null : 'Name is required'),
      role: (value) => (value ? null : 'Role is required'),
    },
  });

  useEffect(() => {
    if (initialValues) {
      form.setValues({
        email: initialValues.email,
        name: initialValues.name,
        role: initialValues.role,
        isActive: initialValues.isActive ?? true,
      });
    }
  }, [initialValues]); // eslint-disable-line react-hooks/exhaustive-deps -- edit form reset when user record changes

  const handleSubmit = form.onSubmit(async (values) => {
    await onSubmit(values);
  });

  return (
    <form onSubmit={handleSubmit}>
      <Stack gap="md">
        <TextInput
          label="Full name"
          placeholder="Jane Smith"
          required
          {...form.getInputProps('name')}
        />
        <TextInput
          label="Email"
          type="email"
          placeholder="jane.smith@clinic.local"
          required
          {...form.getInputProps('email')}
        />
        <Select label="Role" data={STAFF_ROLES} required {...form.getInputProps('role')} />
        {isEdit && includeActiveToggle && (
          <Switch
            label="Active account"
            description="Inactive staff cannot sign in but their history is kept."
            {...form.getInputProps('isActive', { type: 'checkbox' })}
          />
        )}
        <Button type="submit" loading={isLoading}>
          {submitLabel}
        </Button>
      </Stack>
    </form>
  );
}
