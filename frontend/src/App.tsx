import { Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { MantineProvider, LoadingOverlay, localStorageColorSchemeManager } from '@mantine/core';
import { QueryClientProvider } from '@tanstack/react-query';
import { Notifications } from '@mantine/notifications';
import { AuthProvider } from './context/AuthContext';
import { SessionTimeoutManager } from './components/auth/SessionTimeoutManager';
import { AuthSessionExpiredHandler } from './components/auth/AuthSessionExpiredHandler';
import { HelpProvider } from './context/HelpContext';
import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { AdminRoute } from './components/auth/AdminRoute';
import { AppLayout } from './components/layout/AppLayout';
import { queryClient } from './queryClient';

const Login = lazy(() => import('./pages/Login'));
const AdminLogin = lazy(() => import('./pages/AdminLogin'));
const PatientList = lazy(() => import('./components/patients/PatientList').then((module) => ({ default: module.PatientList })));
const CreatePatient = lazy(() => import('./pages/CreatePatient'));
const HelpAndSupport = lazy(() => import('./pages/HelpAndSupport'));
const PatientDetails = lazy(() => import('./pages/PatientDetails'));
const PatientNotesListPage = lazy(() => import('./pages/PatientNotesListPage'));
const PatientExerciseLogListPage = lazy(() => import('./pages/PatientExerciseLogListPage'));
const PatientFormsListPage = lazy(() => import('./pages/PatientFormsListPage'));
const EditPatient = lazy(() => import('./pages/EditPatient'));
const Users = lazy(() => import('./pages/Users'));
const CreateUser = lazy(() => import('./pages/CreateUser'));
const UserDetails = lazy(() => import('./pages/UserDetails'));
const EditUser = lazy(() => import('./pages/EditUser'));
const ChangePin = lazy(() => import('./pages/ChangePin'));
const FormTemplates = lazy(() => import('./pages/FormTemplates'));
const FormReport = lazy(() => import('./pages/FormReport'));
const AdminDashboard = lazy(() => import('./pages/AdminDashboard'));
const AdminAuditLog = lazy(() => import('./pages/AdminAuditLog'));
const CreateFormTemplate = lazy(() => import('./pages/CreateFormTemplate'));
const EditFormTemplate = lazy(() => import('./pages/EditFormTemplate'));

const LoadingFallback = () => (  <div style={{ position: 'relative', minHeight: '200px' }}>
    <LoadingOverlay visible={true} />
  </div>
);

const colorSchemeManager = localStorageColorSchemeManager({
  key: 'poco-color-scheme',
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider defaultColorScheme="auto" colorSchemeManager={colorSchemeManager}>
        <Notifications />
        <AuthProvider>
          <Router>
            <AuthSessionExpiredHandler />
            <SessionTimeoutManager />
            <Suspense fallback={<LoadingFallback />}>
              <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/login/admin" element={<AdminLogin />} />
                <Route
                  path="/*"
                  element={
                    <ProtectedRoute>
                      <HelpProvider>
                        <AppLayout>
                        <Suspense fallback={<LoadingFallback />}>
                          <Routes>
                            <Route path="/" element={<Navigate to="/patients" replace />} />
                            <Route path="/patients" element={<PatientList />} />
                            <Route path="/patients/new" element={<CreatePatient />} />
                            <Route path="/help" element={<HelpAndSupport />} />
                            <Route path="/help/:articleId" element={<HelpAndSupport />} />
                            <Route path="/account/pin" element={<ChangePin />} />
                            <Route path="/patients/:id" element={<PatientDetails />} />
                            <Route path="/patients/:id/notes" element={<PatientNotesListPage />} />
                            <Route path="/patients/:id/exercise-log" element={<PatientExerciseLogListPage />} />
                            <Route path="/patients/:id/forms" element={<PatientFormsListPage />} />
                            <Route path="/patients/:id/edit" element={<EditPatient />} />
                            <Route
                              path="/forms"
                              element={
                                <AdminRoute>
                                  <FormTemplates />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/admin"
                              element={
                                <AdminRoute>
                                  <AdminDashboard />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/admin/audit"
                              element={
                                <AdminRoute>
                                  <AdminAuditLog />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/forms/reports"
                              element={
                                <AdminRoute>
                                  <FormReport />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/forms/new"
                              element={
                                <AdminRoute>
                                  <CreateFormTemplate />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/forms/:id/edit"
                              element={
                                <AdminRoute>
                                  <EditFormTemplate />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/users"
                              element={
                                <AdminRoute>
                                  <Users />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/users/new"
                              element={
                                <AdminRoute>
                                  <CreateUser />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/users/:id/edit"
                              element={
                                <AdminRoute>
                                  <EditUser />
                                </AdminRoute>
                              }
                            />
                            <Route
                              path="/users/:id"
                              element={
                                <AdminRoute>
                                  <UserDetails />
                                </AdminRoute>
                              }
                            />
                          </Routes>
                        </Suspense>
                        </AppLayout>
                      </HelpProvider>
                    </ProtectedRoute>
                  }
                />
              </Routes>
            </Suspense>
          </Router>
        </AuthProvider>
      </MantineProvider>
    </QueryClientProvider>
  );
}
