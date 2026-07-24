import { useEffect } from 'react';
import { Button, Checkbox, NumberInput, Select, Stack, TextInput, Textarea } from '@mantine/core';
import { useForm } from '@mantine/form';
import { FormField, FormTemplate } from '../../types/form';

interface DynamicFormRendererProps {
  template: FormTemplate;
  initialValues?: Record<string, string | number | boolean>;
  onSubmit: (answers: Record<string, string | number | boolean>) => Promise<void>;
  isLoading?: boolean;
  submitLabel?: string;
  resetOnSubmit?: boolean;
}

function buildInitialValues(
  fields: FormField[],
  overrides?: Record<string, string | number | boolean>,
): Record<string, string | number | boolean> {
  const values = fields.reduce<Record<string, string | number | boolean>>((acc, field) => {
    if (field.type === 'checkbox') {
      acc[field.id] = false;
    } else if (field.type === 'number') {
      acc[field.id] = '';
    } else {
      acc[field.id] = '';
    }
    return acc;
  }, {});

  if (overrides) {
    for (const [key, value] of Object.entries(overrides)) {
      values[key] = value;
    }
  }
  return values;
}

export function DynamicFormRenderer({
  template,
  initialValues,
  onSubmit,
  isLoading,
  submitLabel = 'Save form',
  resetOnSubmit = false,
}: DynamicFormRendererProps) {
  const sortedFields = [...template.fields].sort((a, b) => a.sortOrder - b.sortOrder);

  const form = useForm({
    initialValues: buildInitialValues(sortedFields, initialValues),
    validate: sortedFields.reduce<Record<string, (value: unknown) => string | null>>((rules, field) => {
      rules[field.id] = (value) => {
        if (!field.required) {
          return null;
        }
        if (field.type === 'checkbox') {
          return value === true ? null : `${field.label} is required`;
        }
        if (value === '' || value === null || value === undefined) {
          return `${field.label} is required`;
        }
        return null;
      };
      return rules;
    }, {}),
  });

  useEffect(() => {
    form.setValues(buildInitialValues(sortedFields, initialValues));
    // Mantine form identity is stable; re-run only when template or saved answers change.
  }, [template.id, initialValues]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleSubmit = form.onSubmit(async (values) => {
    const answers: Record<string, string | number | boolean> = {};
    for (const field of sortedFields) {
      const value = values[field.id];
      if (field.type === 'number') {
        if (value === '' || value === null || value === undefined) {
          if (field.required) {
            answers[field.id] = '';
          }
          continue;
        }
        answers[field.id] = Number(value);
      } else if (field.type === 'checkbox') {
        answers[field.id] = Boolean(value);
      } else if (value !== '' && value !== null && value !== undefined) {
        answers[field.id] = String(value);
      } else if (field.required) {
        answers[field.id] = '';
      }
    }
    await onSubmit(answers);
    if (resetOnSubmit) {
      form.reset();
    }
  });

  return (
    <form onSubmit={handleSubmit}>
      <Stack gap="md">
        {sortedFields.map((field) => {
          if (field.type === 'textarea') {
            return (
              <Textarea
                key={field.id}
                label={field.label}
                required={field.required}
                {...form.getInputProps(field.id)}
              />
            );
          }
          if (field.type === 'number') {
            return (
              <NumberInput
                key={field.id}
                label={field.label}
                required={field.required}
                {...form.getInputProps(field.id)}
              />
            );
          }
          if (field.type === 'select') {
            return (
              <Select
                key={field.id}
                label={field.label}
                required={field.required}
                data={field.options ?? []}
                {...form.getInputProps(field.id)}
              />
            );
          }
          if (field.type === 'checkbox') {
            return (
              <Checkbox
                key={field.id}
                label={field.label}
                {...form.getInputProps(field.id, { type: 'checkbox' })}
              />
            );
          }
          return (
            <TextInput
              key={field.id}
              label={field.label}
              required={field.required}
              {...form.getInputProps(field.id)}
            />
          );
        })}
        <Button type="submit" loading={isLoading}>
          {submitLabel}
        </Button>
      </Stack>
    </form>
  );
}
