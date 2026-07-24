import { Alert, Anchor, Group, Text } from '@mantine/core';
import { IconHelp, IconX } from '@tabler/icons-react';
import { useState } from 'react';
import { Link } from 'react-router-dom';

interface PageHelpTipProps {
  tipId: string;
  title: string;
  message: string;
  articleId: string;
  learnMoreLabel?: string;
}

function isDismissed(tipId: string): boolean {
  return localStorage.getItem(`poco-help-dismiss-${tipId}`) === '1';
}

export function PageHelpTip({ tipId, title, message, articleId, learnMoreLabel = 'Open guide' }: PageHelpTipProps) {
  const [visible, setVisible] = useState(() => !isDismissed(tipId));

  if (!visible) {
    return null;
  }

  const dismiss = () => {
    localStorage.setItem(`poco-help-dismiss-${tipId}`, '1');
    setVisible(false);
  };

  return (
    <Alert
      color="blue"
      variant="light"
      title={
        <Group gap="xs">
          <IconHelp size={16} />
          <Text fw={600} size="sm">{title}</Text>
        </Group>
      }
      mb="md"
      withCloseButton
      closeButtonLabel="Dismiss tip"
      icon={null}
      onClose={dismiss}
    >
      <Group justify="space-between" align="flex-end" wrap="wrap" gap="xs">
        <Text size="sm">{message}</Text>
        <Group gap="xs">
          <Anchor component={Link} to={`/help/${articleId}`} size="sm" fw={600}>
            {learnMoreLabel}
          </Anchor>
          <Anchor component="button" type="button" size="xs" c="dimmed" onClick={dismiss} style={{ border: 'none', background: 'none', cursor: 'pointer' }}>
            <Group gap={4}>
              <IconX size={12} />
              Don&apos;t show again
            </Group>
          </Anchor>
        </Group>
      </Group>
    </Alert>
  );
}
