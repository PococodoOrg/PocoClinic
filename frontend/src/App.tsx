import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { MantineProvider, Container } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Notifications } from '@mantine/notifications';
import { PatientList } from './components/patients/PatientList';
import { CreatePatient } from './pages/CreatePatient';
import { AppLayout } from './components/layout/AppLayout';
import HelpAndSupport from './pages/HelpAndSupport';
import PatientDetails from './pages/PatientDetails';

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
              <Route path="/help" element={<HelpAndSupport />} />
              <Route path="/patients/:id" element={<PatientDetails />} />
              {/* Add more routes as needed */}
            </Routes>
          </AppLayout>
        </Router>
      </MantineProvider>
    </QueryClientProvider>
  );
} 