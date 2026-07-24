import { Button } from '@mantine/core';
import { Link } from 'react-router-dom';

interface ViewAllChartLinkProps {
  to: string;
  total: number;
  label: string;
}

export function ViewAllChartLink({ to, total, label }: ViewAllChartLinkProps) {
  return (
    <Button component={Link} to={to} variant="light" fullWidth>
      View all {label} ({total})
    </Button>
  );
}
