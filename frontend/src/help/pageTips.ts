export interface PageHelpTipConfig {
  tipId: string;
  title: string;
  message: string;
  articleId: string;
}

export const pageHelpTips: Record<string, PageHelpTipConfig> = {
  '/patients': {
    tipId: 'patients-list',
    title: 'Finding patients',
    message: 'Search by first or last name, then double-click a row or use the view button to open a chart.',
    articleId: 'patient-list',
  },
  '/patients/new': {
    tipId: 'patients-new',
    title: 'New registration',
    message: 'Enter legal name and date of birth at minimum. Contact details help with follow-up but are optional.',
    articleId: 'register-patient',
  },
  '/admin': {
    tipId: 'admin-dashboard',
    title: 'Administrator dashboard',
    message: 'Use Overview for alerts and status, Administrator guide for setup tasks, and Backup & maintenance for USB rotation and health checks.',
    articleId: 'admin-dashboard',
  },
  '/account/pin': {
    tipId: 'change-pin',
    title: 'Personal PIN',
    message: 'Choose a PIN only you know. You will use it with your badge every sign-in.',
    articleId: 'change-pin',
  },
};

export function getPageHelpTip(pathname: string): PageHelpTipConfig | null {
  if (pageHelpTips[pathname]) {
    return pageHelpTips[pathname];
  }
  if (/^\/patients\/[^/]+$/.test(pathname)) {
    return {
      tipId: 'patient-chart',
      title: 'Patient chart',
      message: 'Add clinical notes for visit summaries. Upload PDFs or scans under Documents. Complete forms at the bottom of the chart.',
      articleId: 'patient-chart',
    };
  }
  if (pathname.startsWith('/users')) {
    return {
      tipId: 'staff-management',
      title: 'Staff accounts',
      message: 'Print badge QR codes for new employees. Reissue a badge immediately if one is lost — the old QR stops working.',
      articleId: 'staff-management',
    };
  }
  return null;
}
