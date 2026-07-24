import { Anchor, Breadcrumbs, Text } from '@mantine/core';
import { Link, useLocation, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { fetchPatient } from '../../api/patients';
import { getArticleById } from '../../help/content';

const staticLabels: Record<string, string> = {
  patients: 'Patients',
  admin: 'Admin',
  audit: 'Audit log',
  users: 'Staff',
  forms: 'Forms',
  reports: 'Reports',
  help: 'Help center',
  account: 'Account',
  pin: 'Change PIN',
  new: 'New',
  edit: 'Edit',
};

export function PageBreadcrumbs() {
  const location = useLocation();
  const params = useParams();
  const segments = location.pathname.split('/').filter(Boolean);

  const patientId =
    params.id && segments[0] === 'patients' && segments[1] !== 'new' ? params.id : undefined;
  const helpArticleId = params.articleId && segments[0] === 'help' ? params.articleId : undefined;

  const { data: patient } = useQuery({
    queryKey: ['patient', patientId],
    queryFn: () => fetchPatient(patientId!),
    enabled: Boolean(patientId),
  });

  const helpArticle = helpArticleId ? getArticleById(helpArticleId) : undefined;

  if (segments.length === 0) {
    return null;
  }

  const crumbs: { label: string; to?: string }[] = [];
  let path = '';

  segments.forEach((segment, index) => {
    path += `/${segment}`;
    const isLast = index === segments.length - 1;

    if (segment === params.id && patient) {
      crumbs.push({
        label: `${patient.firstName} ${patient.lastName}`,
        to: isLast ? undefined : path,
      });
      return;
    }

    if (segment === params.articleId && helpArticle) {
      crumbs.push({
        label: helpArticle.title,
        to: isLast ? undefined : path,
      });
      return;
    }

    crumbs.push({
      label: staticLabels[segment] ?? segment,
      to: isLast ? undefined : path,
    });
  });

  return (
    <Breadcrumbs mb="md">
      {crumbs.map((crumb) =>
        crumb.to ? (
          <Anchor component={Link} to={crumb.to} key={crumb.to} size="sm">
            {crumb.label}
          </Anchor>
        ) : (
          <Text key={crumb.label} size="sm">
            {crumb.label}
          </Text>
        ),
      )}
    </Breadcrumbs>
  );
}
