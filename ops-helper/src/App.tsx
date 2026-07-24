import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { AppShell } from './components/AppShell';
import { usePiTouch } from './context/PiTouchContext';
import { ChecklistPage } from './pages/ChecklistPage';
import { HomePage } from './pages/HomePage';
import { BackupWizardPage } from './pages/BackupWizardPage';
import { RestoreWizardPage } from './pages/RestoreWizardPage';

function AppRoutes() {
  const { isPiTouch } = usePiTouch();

  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/backup" element={<BackupWizardPage />} />
      <Route path="/restore" element={<RestoreWizardPage />} />
      {!isPiTouch && <Route path="/checklist" element={<ChecklistPage />} />}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AppShell>
        <AppRoutes />
      </AppShell>
    </BrowserRouter>
  );
}
