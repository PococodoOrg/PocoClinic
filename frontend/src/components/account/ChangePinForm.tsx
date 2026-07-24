import { Button, PasswordInput, Stack } from '@mantine/core';
import { useForm } from '@mantine/form';

export interface ChangePinFormData {
  currentPin: string;
  newPin: string;
  confirmPin: string;
}

interface ChangePinFormProps {
  onSubmit: (data: Pick<ChangePinFormData, 'currentPin' | 'newPin'>) => Promise<void>;
  isLoading?: boolean;
}

export function ChangePinForm({ onSubmit, isLoading }: ChangePinFormProps) {
  const form = useForm<ChangePinFormData>({
    initialValues: {
      currentPin: '',
      newPin: '',
      confirmPin: '',
    },
    validate: {
      currentPin: (value) => (/^\d{4}$/.test(value) ? null : 'Enter your current 4-digit PIN'),
      newPin: (value) => (/^\d{4}$/.test(value) ? null : 'Enter a 4-digit PIN'),
      confirmPin: (value, values) =>
        value === values.newPin ? null : 'PIN confirmation does not match',
    },
  });

  const handleSubmit = form.onSubmit(async (values) => {
    await onSubmit({
      currentPin: values.currentPin,
      newPin: values.newPin,
    });
    form.reset();
  });

  return (
    <form onSubmit={handleSubmit}>
      <Stack gap="md" maw={360}>
        <PasswordInput
          label="Current PIN"
          maxLength={4}
          required
          {...form.getInputProps('currentPin')}
        />
        <PasswordInput
          label="New PIN"
          description="Choose a 4-digit PIN only you know"
          maxLength={4}
          required
          {...form.getInputProps('newPin')}
        />
        <PasswordInput
          label="Confirm new PIN"
          maxLength={4}
          required
          {...form.getInputProps('confirmPin')}
        />
        <Button type="submit" loading={isLoading}>
          Update PIN
        </Button>
      </Stack>
    </form>
  );
}
