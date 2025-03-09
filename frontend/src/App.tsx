import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { MantineProvider, Container } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Notifications } from '@mantine/notifications';
import { PatientList } from './components/patients/PatientList';
import { CreatePatient } from './pages/CreatePatient';
import { AppLayout } from './components/layout/AppLayout';

const queryClient = new QueryClient();

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <Notifications />
        <Router>
          <AppLayout>
            <Routes>
              <Route path="/" element={<Navigate to="/patients" replace />} />
              <Route path="/patients" element={<PatientList />} />
              <Route path="/patients/new" element={<CreatePatient />} />
              {/* Add more routes as needed */}
            </Routes>
          </AppLayout>
        </Router>
      </MantineProvider>
    </QueryClientProvider>
  );
} 