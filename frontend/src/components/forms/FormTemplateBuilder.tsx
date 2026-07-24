import {
  ActionIcon,
  Button,
  Checkbox,
  Group,
  Paper,
  Select,
  Stack,
  Text,
  TextInput,
  Textarea,
} from '@mantine/core';
import { IconArrowDown, IconArrowUp, IconTrash } from '@tabler/icons-react';
import {
  FIELD_TYPE_OPTIONS,
  FORM_TYPE_OPTIONS,
  FormField,
  FormFieldType,
  FormGroup,
  FormTemplateFormData,
  createEmptyField,
} from '../../types/form';

interface FormTemplateBuilderProps {
  value: FormTemplateFormData;
  groups: FormGroup[];
  onChange: (value: FormTemplateFormData) => void;
  submitLabel?: string;
  onSubmit: () => void;
  isLoading?: boolean;
}

export function FormTemplateBuilder({
  value,
  groups,
  onChange,
  submitLabel = 'Save form',
  onSubmit,
  isLoading,
}: FormTemplateBuilderProps) {
  const updateField = (index: number, patch: Partial<FormField>) => {
    const fields = value.fields.map((field, fieldIndex) =>
      fieldIndex === index ? { ...field, ...patch } : field,
    );
    onChange({ ...value, fields });
  };

  const moveField = (index: number, direction: -1 | 1) => {
    const target = index + direction;
    if (target < 0 || target >= value.fields.length) {
      return;
    }
    const fields = [...value.fields];
    [fields[index], fields[target]] = [fields[target], fields[index]];
    onChange({
      ...value,
      fields: fields.map((field, fieldIndex) => ({ ...field, sortOrder: fieldIndex })),
    });
  };

  const removeField = (index: number) => {
    onChange({
      ...value,
      fields: value.fields
        .filter((_, fieldIndex) => fieldIndex !== index)
        .map((field, fieldIndex) => ({ ...field, sortOrder: fieldIndex })),
    });
  };

  const addField = () => {
    onChange({
      ...value,
      fields: [...value.fields, createEmptyField(value.fields.length)],
    });
  };

  const selectedFormType = FORM_TYPE_OPTIONS.find((option) => option.value === value.formType);

  return (
    <Stack gap="md">
      <TextInput
        label="Form name"
        placeholder="Intake checklist"
        value={value.name}
        onChange={(event) => onChange({ ...value, name: event.currentTarget.value })}
        required
      />

      <Group grow align="flex-start">
        <Select
          label="Group"
          placeholder="Select a group"
          data={groups.map((group) => ({ value: group.id, label: group.name }))}
          value={value.groupId || null}
          onChange={(groupId) => onChange({ ...value, groupId: groupId ?? '' })}
          required
        />
        <Select
          label="Form type"
          data={FORM_TYPE_OPTIONS.map((option) => ({
            value: option.value,
            label: option.label,
          }))}
          value={value.formType}
          onChange={(formType) =>
            onChange({ ...value, formType: (formType as FormTemplateFormData['formType']) ?? 'singleton' })
          }
          required
        />
      </Group>

      {selectedFormType && (
        <Text size="sm" c="dimmed">
          {selectedFormType.description}
        </Text>
      )}

      <Stack gap="sm">
        <Group justify="space-between">
          <Text fw={500}>Fields</Text>
          <Button variant="light" size="xs" onClick={addField}>
            Add field
          </Button>
        </Group>

        {value.fields.length === 0 && (
          <Text c="dimmed" size="sm">
            Add at least one field to this form.
          </Text>
        )}

        {value.fields.map((field, index) => (
          <Paper key={field.id} withBorder p="md">
            <Stack gap="sm">
              <Group align="flex-end" grow>
                <TextInput
                  label="Label"
                  value={field.label}
                  onChange={(event) => updateField(index, { label: event.currentTarget.value })}
                  required
                />
                <Select
                  label="Type"
                  data={FIELD_TYPE_OPTIONS}
                  value={field.type}
                  onChange={(nextType) =>
                    updateField(index, { type: (nextType as FormFieldType) ?? 'text' })
                  }
                />
              </Group>

              {field.type === 'select' && (
                <Textarea
                  label="Options"
                  description="One option per line"
                  value={(field.options ?? []).join('\n')}
                  onChange={(event) =>
                    updateField(index, {
                      options: event.currentTarget.value
                        .split('\n')
                        .map((option) => option.trim())
                        .filter(Boolean),
                    })
                  }
                  minRows={2}
                />
              )}

              <Group justify="space-between">
                <Checkbox
                  label="Required"
                  checked={field.required}
                  onChange={(event) => updateField(index, { required: event.currentTarget.checked })}
                />
                <Group gap="xs">
                  <ActionIcon variant="light" onClick={() => moveField(index, -1)} aria-label="Move up">
                    <IconArrowUp size={16} />
                  </ActionIcon>
                  <ActionIcon variant="light" onClick={() => moveField(index, 1)} aria-label="Move down">
                    <IconArrowDown size={16} />
                  </ActionIcon>
                  <ActionIcon
                    color="red"
                    variant="light"
                    onClick={() => removeField(index)}
                    aria-label="Remove field"
                  >
                    <IconTrash size={16} />
                  </ActionIcon>
                </Group>
              </Group>
            </Stack>
          </Paper>
        ))}
      </Stack>

      <Button
        onClick={onSubmit}
        loading={isLoading}
        disabled={!value.name.trim() || !value.groupId || value.fields.length === 0}
      >
        {submitLabel}
      </Button>
    </Stack>
  );
}
