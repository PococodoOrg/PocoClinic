import React, { Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { MantineProvider, Container, LoadingOverlay } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Notifications } from '@mantine/notifications';
import { PatientList } from './components/patients/PatientList';
import { AppLayout } from './components/layout/AppLayout';

// Lazy load components
const CreatePatient = lazy(() => import('./pages/CreatePatient'));
const HelpAndSupport = lazy(() => import('./pages/HelpAndSupport'));
const PatientDetails = lazy(() => import('./pages/PatientDetails'));
const EditPatient = lazy(() => import('./pages/EditPatient'));

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      gcTime: 10 * 60 * 1000, // 10 minutes (replaces cacheTime in v5)
      refetchOnWindowFocus: false, // Don't refetch when window regains focus
      retry: 1, // Only retry failed requests once
    },
  },
});

// Loading component
const LoadingFallback = () => (
  <div style={{ position: 'relative', minHeight: '200px' }}>
    <LoadingOverlay visible={true} />
  </div>
);

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <Notifications />
        <Router>
          <AppLayout>
            <Suspense fallback={<LoadingFallback />}>
              <Routes>
                <Route path="/" element={<Navigate to="/patients" replace />} />
                <Route path="/patients" element={<PatientList />} />
                <Route path="/patients/new" element={<CreatePatient />} />
                <Route path="/help" element={<HelpAndSupport />} />
                <Route path="/patients/:id" element={<PatientDetails />} />
                <Route path="/patients/:id/edit" element={<EditPatient />} />
              </Routes>
            </Suspense>
          </AppLayout>
        </Router>
      </MantineProvider>
    </QueryClientProvider>
  );
} 