import React from 'react';
import ReactDOM from 'react-dom/client';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from './App';
import { theme } from './theme';
import { PiTouchProvider } from './context/PiTouchContext';
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import './styles.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={theme} defaultColorScheme="light">
        <Notifications position="top-center" />
        <PiTouchProvider>
          <App />
        </PiTouchProvider>
      </MantineProvider>
    </QueryClientProvider>
  </React.StrictMode>,
);
