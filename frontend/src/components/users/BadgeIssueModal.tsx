import { useState } from 'react';
import { Alert, Button, Code, CopyButton, Group, Modal, Stack, Text } from '@mantine/core';
import { IconCheck, IconCopy, IconPrinter } from '@tabler/icons-react';
import { BadgePrintCard } from './BadgePrintCard';
import { DEFAULT_PIN, StaffUser } from '../../types/user';

interface BadgeIssueModalProps {
  opened: boolean;
  user: StaffUser | null;
  badgeKey: string;
  mode?: 'create' | 'reissue';
  onClose: () => void;
}

export function BadgeIssueModal({
  opened,
  user,
  badgeKey,
  mode = 'create',
  onClose,
}: BadgeIssueModalProps) {
  const [acknowledged, setAcknowledged] = useState(false);

  const handlePrint = () => {
    window.print();
  };

  const handleClose = () => {
    setAcknowledged(false);
    onClose();
  };

  if (!user) {
    return null;
  }

  return (
    <>
      <style>
        {`
          @media print {
            body * {
              visibility: hidden;
            }
            #badge-print-card,
            #badge-print-card * {
              visibility: visible;
            }
            #badge-print-card {
              position: absolute;
              left: 0;
              top: 0;
              border: none;
            }
          }
        `}
      </style>
      <Modal
        opened={opened}
        onClose={handleClose}
        title={mode === 'reissue' ? 'Print replacement badge' : 'Print employee badge'}
        size="md"
        closeOnClickOutside={false}
        closeOnEscape={false}
        withCloseButton={acknowledged}
      >
        <Stack gap="md">
          <Alert color="yellow" title="Save this badge key now">
            {mode === 'reissue'
              ? 'The previous badge no longer works. Print the replacement badge or copy the new key before closing.'
              : 'The access key is shown only once. Print the badge or copy the key before closing this dialog.'}
          </Alert>

          <BadgePrintCard
            name={user.name}
            email={user.email}
            badgeKey={badgeKey}
            pinReminder={
              mode === 'reissue'
                ? 'Scan this replacement badge at staff sign-in, then enter your existing PIN.'
                : undefined
            }
          />

          <Stack gap="xs">
            <Text size="sm" fw={500}>
              Badge key
            </Text>
            <Group gap="xs">
              <Code block style={{ flex: 1 }}>
                {badgeKey}
              </Code>
              <CopyButton value={badgeKey}>
                {({ copied, copy }) => (
                  <Button
                    variant="light"
                    leftSection={copied ? <IconCheck size={16} /> : <IconCopy size={16} />}
                    onClick={copy}
                  >
                    {copied ? 'Copied' : 'Copy'}
                  </Button>
                )}
              </CopyButton>
            </Group>
          </Stack>

          {mode === 'create' && (
            <Alert color="blue" title={`Default PIN: ${DEFAULT_PIN}`}>
              Tell {user.name} to sign in with this PIN at the staff login screen. They should change
              it after first login from Account → Change PIN.
            </Alert>
          )}

          {mode === 'reissue' && (
            <Alert color="blue" title="PIN unchanged">
              {user.name} keeps the same PIN. Only the badge key changed.
            </Alert>
          )}

          <Group justify="space-between">
            <Button variant="light" leftSection={<IconPrinter size={16} />} onClick={handlePrint}>
              Print badge
            </Button>
            <Button
              onClick={() => {
                setAcknowledged(true);
                handleClose();
              }}
            >
              I saved the badge key
            </Button>
          </Group>
        </Stack>
      </Modal>
    </>
  );
}
