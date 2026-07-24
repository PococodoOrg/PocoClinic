import { MantineProvider } from '@mantine/core';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { ChangePinForm } from './ChangePinForm';

function renderForm(onSubmit = vi.fn().mockResolvedValue(undefined)) {
  render(
    <MantineProvider>
      <ChangePinForm onSubmit={onSubmit} />
    </MantineProvider>,
  );
  return onSubmit;
}

describe('ChangePinForm', () => {
  it('requires 4-digit PINs and matching confirmation', async () => {
    const user = userEvent.setup();
    renderForm();

    await user.click(screen.getByRole('button', { name: /update pin/i }));

    expect(await screen.findByText(/4-digit pin/i)).toBeInTheDocument();
  });

  it('submits current and new PIN when valid', async () => {
    const user = userEvent.setup();
    const onSubmit = renderForm();

    const inputs = screen.getAllByLabelText(/pin/i);
    await user.type(inputs[0], '1234');
    await user.type(inputs[1], '5678');
    await user.type(inputs[2], '5678');
    await user.click(screen.getByRole('button', { name: /update pin/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({ currentPin: '1234', newPin: '5678' });
    });
  });
});
