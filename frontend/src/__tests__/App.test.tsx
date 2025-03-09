import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MantineProvider } from '@mantine/core';
import App from '../App';

const queryClient = new QueryClient();

const TestWrapper = ({ children }: { children: React.ReactNode }) => (
  <QueryClientProvider client={queryClient}>
    <MantineProvider>
      {children}
    </MantineProvider>
  </QueryClientProvider>
);

describe('App Component', () => {
  it('renders without crashing', () => {
    render(<App />, { wrapper: TestWrapper });
  });

  it('contains header elements', () => {
    render(<App />, { wrapper: TestWrapper });
    // Check for the app title
    expect(screen.getByText('PocoClinic EMR')).toBeInTheDocument();
    // Check for the Patients button
    expect(screen.getByRole('button', { name: 'Patients' })).toBeInTheDocument();
  });
}); 