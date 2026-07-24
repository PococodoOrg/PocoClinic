import { MantineProvider } from '@mantine/core';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { DynamicFormRenderer } from './DynamicFormRenderer';
import { FormTemplate } from '../../types/form';

const template: FormTemplate = {
  id: 'tpl-1',
  groupId: 'grp-1',
  name: 'Intake',
  formType: 'singleton',
  fields: [
    { id: 'name', label: 'Patient name', type: 'text', required: true, sortOrder: 1 },
    { id: 'pain', label: 'Pain level', type: 'number', required: true, sortOrder: 2 },
    { id: 'consent', label: 'Consent given', type: 'checkbox', required: true, sortOrder: 3 },
  ],
  createdAt: '2026-01-01T00:00:00.000Z',
  updatedAt: '2026-01-01T00:00:00.000Z',
};

function renderForm(onSubmit = vi.fn().mockResolvedValue(undefined)) {
  render(
    <MantineProvider>
      <DynamicFormRenderer template={template} onSubmit={onSubmit} />
    </MantineProvider>,
  );
  return onSubmit;
}

describe('DynamicFormRenderer', () => {
  it('submits typed answers', async () => {
    const user = userEvent.setup();
    const onSubmit = renderForm();

    await user.type(screen.getByLabelText(/patient name/i), 'Ada Lovelace');
    await user.type(screen.getByLabelText(/pain level/i), '4');
    await user.click(screen.getByLabelText(/consent given/i));
    await user.click(screen.getByRole('button', { name: /save form/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: 'Ada Lovelace',
        pain: 4,
        consent: true,
      });
    });
  });
});
