import {
  IconClipboardList,
  IconFileReport,
  IconForms,
  IconHelp,
  IconHome,
  IconKey,
  IconShieldLock,
  IconUsers,
} from '@tabler/icons-react';
import type { ComponentType } from 'react';

export interface NavItem {
  label: string;
  path: string;
  icon: ComponentType<{ size?: number; stroke?: number }>;
  adminOnly?: boolean;
}

export interface NavSection {
  title: string;
  items: NavItem[];
}

export const navSections: NavSection[] = [
  {
    title: 'Clinical',
    items: [
      { label: 'Patients', path: '/patients', icon: IconHome },
    ],
  },
  {
    title: 'Administration',
    items: [
      { label: 'Dashboard', path: '/admin', icon: IconShieldLock, adminOnly: true },
      { label: 'Staff', path: '/users', icon: IconUsers, adminOnly: true },
      { label: 'Forms', path: '/forms', icon: IconForms, adminOnly: true },
      { label: 'Reports', path: '/forms/reports', icon: IconFileReport, adminOnly: true },
      { label: 'Audit log', path: '/admin/audit', icon: IconClipboardList, adminOnly: true },
    ],
  },
  {
    title: 'Account',
    items: [
      { label: 'Change PIN', path: '/account/pin', icon: IconKey },
      { label: 'Help & Support', path: '/help', icon: IconHelp },
    ],
  },
];

export function isNavActive(pathname: string, path: string): boolean {
  if (path === '/patients') {
    return pathname === '/patients' || pathname.startsWith('/patients/');
  }
  if (path === '/forms') {
    return pathname.startsWith('/forms') && !pathname.startsWith('/forms/reports');
  }
  return pathname === path || pathname.startsWith(`${path}/`);
}
