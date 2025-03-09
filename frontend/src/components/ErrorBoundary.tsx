import React from 'react'
import { Alert, Button, Stack, Text } from '@mantine/core'
import { IconAlertCircle } from '@tabler/icons-react'

interface Props {
  children: React.ReactNode
  onReset?: () => void
}

interface State {
  hasError: boolean
  error?: Error
}

export class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // Log error to error reporting service
    console.error('Error caught by boundary:', error, errorInfo)
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined })
    this.props.onReset?.()
  }

  render() {
    if (this.state.hasError) {
      return (
        <Stack align="center" gap="md" p="xl">
          <Alert
            icon={<IconAlertCircle size={16} />}
            title="Something went wrong!"
            color="red"
            radius="md"
          >
            <Stack gap="sm">
              <Text size="sm">
                We apologize for the inconvenience. Please try refreshing the page or contact support if the problem persists.
              </Text>
              {this.state.error && (
                <Text size="xs" c="dimmed">
                  Error: {this.state.error.message}
                </Text>
              )}
              <Button
                variant="light"
                color="blue"
                size="sm"
                onClick={this.handleReset}
              >
                Try Again
              </Button>
            </Stack>
          </Alert>
        </Stack>
      )
    }

    return this.props.children
  }
} 