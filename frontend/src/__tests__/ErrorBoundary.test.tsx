import React, { useState } from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import { vi, describe, it, expect, beforeAll, afterAll } from 'vitest'
import { MantineProvider } from '@mantine/core'
import { ErrorBoundary } from '../components/ErrorBoundary'
import '@testing-library/jest-dom/vitest'

// Component that throws an error
const BuggyComponent = ({ shouldThrow = true }: { shouldThrow?: boolean }) => {
  if (shouldThrow) {
    throw new Error('Test error')
  }
  return <div>Working component</div>
}

// Test wrapper with providers
const renderWithProviders = (ui: React.ReactElement) => {
  return render(
    <MantineProvider>
      {ui}
    </MantineProvider>
  )
}

// Recovery test wrapper
const RecoveryWrapper = () => {
  const [shouldThrow, setShouldThrow] = useState(true)

  return (
    <ErrorBoundary onReset={() => setShouldThrow(false)}>
      <BuggyComponent shouldThrow={shouldThrow} />
    </ErrorBoundary>
  )
}

describe('ErrorBoundary', () => {
  // Prevent console.error from cluttering test output
  const originalError = console.error
  beforeAll(() => {
    console.error = vi.fn()
  })
  afterAll(() => {
    console.error = originalError
  })

  it('renders children when there is no error', () => {
    renderWithProviders(
      <ErrorBoundary>
        <div>Test content</div>
      </ErrorBoundary>
    )

    expect(screen.getByText('Test content')).toBeInTheDocument()
  })

  it('renders error UI when there is an error', () => {
    renderWithProviders(
      <ErrorBoundary>
        <BuggyComponent />
      </ErrorBoundary>
    )

    expect(screen.getByText('Something went wrong!')).toBeInTheDocument()
    expect(screen.getByText(/Test error/)).toBeInTheDocument()
  })

  it('allows recovery via try again button', () => {
    renderWithProviders(<RecoveryWrapper />)

    // Verify error state
    expect(screen.getByText('Something went wrong!')).toBeInTheDocument()
    expect(screen.getByText(/Test error/)).toBeInTheDocument()

    // Click try again
    fireEvent.click(screen.getByText('Try Again'))

    // Verify recovery
    expect(screen.getByText('Working component')).toBeInTheDocument()
  })
}) 